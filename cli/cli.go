/*
Package main provides the command line interface for the meli application.
*/
package main

import (
	"context"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"runtime"
	"runtime/pprof"
	"strings"
	"sync"

	"github.com/docker/docker/client"
	"github.com/komuw/meli"
	"github.com/pkg/errors"
	yaml "gopkg.in/yaml.v2"
)

/* DOCS:
1. https://godoc.org/github.com/moby/moby/client
2. https://docs.docker.com/engine/api/latest
*/

var version string

// Cli parses input from stdin
func Cli() (showVersion, followLogs, rebuild bool, dockerComposeFile string, cpuprofile, memprofile string) {
	// TODO; use a more sensible cli lib.
	var up bool
	var d bool
	var build bool

	flag.BoolVar(
		&showVersion,
		"version",
		false,
		"Show version information.")
	flag.BoolVar(
		&showVersion,
		"v",
		false,
		"Show version information.")
	flag.BoolVar(
		&up,
		"up",
		false,
		"Builds, re/creates, starts, and attaches to containers for a service.")
	flag.BoolVar(
		&d,
		"d",
		false,
		"Run containers in the background")
	flag.BoolVar(
		&build,
		"build",
		false,
		"Rebuild services")
	flag.StringVar(
		&dockerComposeFile,
		"f",
		"docker-compose.yml",
		"path to docker-compose.yml file.")
	flag.StringVar(
		&cpuprofile,
		"cpuprofile",
		"",
		"write cpu profile to `file`. This is only useful for debugging meli.")
	flag.StringVar(
		&memprofile,
		"memprofile",
		"",
		"write memory profile to `file`. This is only useful for debugging meli.")

	flag.Parse()

	if showVersion {
		return true, followLogs, rebuild, "", cpuprofile, memprofile
	}
	if !up {
		fmt.Println("to use Meli, run: \n\n\t meli -up")
		os.Exit(0)
	}
	if d {
		followLogs = false
	} else {
		followLogs = true
	}
	if build {
		rebuild = true
	}

	return false, followLogs, rebuild, dockerComposeFile, cpuprofile, memprofile
}

func main() {
	showVersion, followLogs, rebuild, dockerComposeFile, cpuprofile, memprofile := Cli()

	if cpuprofile != "" {
		f, err := os.Create(cpuprofile)
		if err != nil {
			e := errors.Wrap(err, "could not create CPU profile")
			log.Fatalf("%+v", e)
		}
		if err := pprof.StartCPUProfile(f); err != nil {
			e := errors.Wrap(err, "could not start CPU profile")
			log.Fatalf("%+v", e)
		}
		defer pprof.StopCPUProfile()
	}

	if showVersion {
		fmt.Println("Meli version: ", version)
		os.Exit(0)
	}

	data, err := ioutil.ReadFile(dockerComposeFile)
	if err != nil {
		e := errors.Wrap(err, "unable to read docker-compose file")
		log.Fatalf("%+v", e)
	}

	var dockerCyaml meli.DockerComposeConfig
	err = yaml.Unmarshal([]byte(data), &dockerCyaml)
	if err != nil {
		e := errors.Wrap(err, "unable to unmarshal docker-compose file contents")
		log.Fatalf("%+v", e)
	}

	ctx := context.Background()
	cli, err := client.NewEnvClient()
	if err != nil {
		e := errors.Wrap(err, "unable to intialize docker client")
		log.Fatalf("%+v", e)
	}
	defer cli.Close()
	curentDir, err := os.Getwd()
	if err != nil {
		e := errors.Wrap(err, "unable to get the current working directory")
		log.Fatalf("%+v", e)
	}
	networkName := "meli_network_" + getCwdName(curentDir)
	networkID, err := meli.GetNetwork(ctx, networkName, cli)
	if err != nil {
		log.Fatalf("%+v", err)
	}
	meli.LoadAuth()

	// Create top level volumes, if any
	if len(dockerCyaml.Volumes) > 0 {
		for k := range dockerCyaml.Volumes {
			// TODO we need to synchronise here else we'll get a race
			// but I think we can get away for now because:
			// 1. there are on average a lot more containers in a compose file
			// than volumes, so the sync in the for loop for containers is enough
			// 2. since we intend to stream logs as containers run(see; issues/24);
			// then meli will be up long enough for the volume creation goroutines to have finished.
			go meli.CreateDockerVolume(ctx, cli, "meli_"+k, "local", os.Stdout)
		}
	}

	var wg sync.WaitGroup
	for k, v := range dockerCyaml.Services {
		wg.Add(1)

		// use dotted filepath. make it also work for windows
		r := strings.NewReplacer("/", ".", ":", ".", "\\", ".")
		dotFormattedrCurentDir := r.Replace(curentDir)
		v.Labels = append(v.Labels, fmt.Sprintf("meli_service=meli_%s%s", k, dotFormattedrCurentDir))

		dc := &meli.DockerContainer{
			ServiceName:       k,
			ComposeService:    v,
			NetworkID:         networkID,
			NetworkName:       networkName,
			FollowLogs:        followLogs,
			DockerComposeFile: dockerComposeFile,
			LogMedium:         os.Stdout,
			CurentDir:         dotFormattedrCurentDir,
			Rebuild:           rebuild,
			EnvFile:           v.EnvFile}
		go startComposeServices(ctx, cli, &wg, dc)
	}
	wg.Wait()

	if memprofile != "" {
		f, err := os.Create(memprofile)
		if err != nil {
			e := errors.Wrap(err, "could not create memory profile")
			log.Fatalf("%+v", e)
		}
		runtime.GC() // get up-to-date statistics
		if err := pprof.WriteHeapProfile(f); err != nil {
			e := errors.Wrap(err, "could not write memory profile")
			log.Fatalf("%+v", e)
		}
		f.Close()
	}
}

func startComposeServices(ctx context.Context, cli *client.Client, wg *sync.WaitGroup, dc *meli.DockerContainer) {
	defer wg.Done()

	/*
		1. Pull Image
		2. Create a container
		3. Connect container to network
		4. Start container
		5. Stream container logs
	*/

	if len(dc.ComposeService.Image) > 0 {
		err := meli.PullDockerImage(ctx, cli, dc)
		if err != nil {
			// clean exit since we want other goroutines for fetching other images
			// to continue running
			fmt.Printf("\n\n%+v", err)
			return
		}
	}
	alreadyCreated, _, err := meli.CreateContainer(ctx, cli, dc)
	if err != nil {
		// clean exit since we want other goroutines for fetching other images
		// to continue running
		fmt.Printf("\n\n%+v", err)
		return
	}

	if !alreadyCreated {
		err = meli.ConnectNetwork(ctx, cli, dc)
		if err != nil {
			// create whitespace so that error is visible to human
			fmt.Printf("\n\n%+v", err)
			return
		}
	}

	err = meli.ContainerStart(ctx, cli, dc)
	if err != nil {
		fmt.Printf("\n\n%+v", err)
		return
	}

	err = meli.ContainerLogs(ctx, cli, dc)
	if err != nil {
		fmt.Printf("\n\n%+v", err)
		return
	}
}

func getCwdName(path string) string {
	//TODO: investigate if this will work cross platform
	// it might be  :unable to handle paths in windows OS
	f := func(c rune) bool {
		if c == 47 {
			// 47 is the '/' character
			return true
		}
		return false
	}
	pathSlice := strings.FieldsFunc(path, f)
	return pathSlice[len(pathSlice)-1]
}
