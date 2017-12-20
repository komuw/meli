package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"sync"

	"github.com/docker/docker/client"
	"github.com/komuw/meli/api"
	"github.com/komuw/meli/cli"

	"gopkg.in/yaml.v2"
)

/* DOCS:
1. https://godoc.org/github.com/moby/moby/client
2. https://docs.docker.com/engine/api/v1.31/
*/

var version string

func main() {
	showVersion, followLogs, dockerComposeFile := cli.Cli()
	if showVersion {
		fmt.Println("Meli version: ", version)
		os.Exit(0)
	}

	data, err := ioutil.ReadFile(dockerComposeFile)
	if err != nil {
		log.Fatal(err, " :unable to read docker-compose file")
	}

	var dockerCyaml api.DockerComposeConfig
	err = yaml.Unmarshal([]byte(data), &dockerCyaml)
	if err != nil {
		log.Fatal(err, " :unable to parse docker-compose file contents")
	}

	ctx := context.Background()
	cli, err := client.NewEnvClient()
	if err != nil {
		log.Fatal(err, " :unable to intialize docker client")
	}
	defer cli.Close()
	curentDir, err := os.Getwd()
	if err != nil {
		log.Fatal(err, " :unable to get the current working directory")
	}
	networkName := "meli_network_" + api.GetCwdName(curentDir)
	networkID, err := api.GetNetwork(ctx, networkName, cli)
	if err != nil {
		log.Fatal(err, " :unable to create/get network")
	}
	api.GetAuth()

	// Create top level volumes, if any
	if len(dockerCyaml.Volumes) > 0 {
		for k := range dockerCyaml.Volumes {
			// TODO we need to synchronise here else we'll get a race
			// but I think we can get away for now because:
			// 1. there are on average a lot more containers in a compose file
			// than volumes, so the sync in the for loop for containers is enough
			// 2. since we intend to stream logs as containers run(see; issues/24);
			// then meli will be up long enough for the volume creation goroutines to have finished.
			go api.CreateDockerVolume(ctx, cli, "meli_"+k, "local", os.Stdout)
		}
	}

	var wg sync.WaitGroup
	for k, v := range dockerCyaml.Services {
		wg.Add(1)
		v.Labels = append(v.Labels, fmt.Sprintf("meli_service=meli_%s", curentDir))

		dc := &api.DockerContainer{
			ServiceName:       k,
			ComposeService:    v,
			NetworkID:         networkID,
			NetworkName:       networkName,
			FollowLogs:        followLogs,
			DockerComposeFile: dockerComposeFile,
			LogMedium:         os.Stdout,
			CurentDir:         curentDir}
		go startComposeServices(ctx, cli, &wg, dc)
	}
	wg.Wait()
}

func startComposeServices(ctx context.Context, cli *client.Client, wg *sync.WaitGroup, dc *api.DockerContainer) {
	defer wg.Done()

	/*
		1. Pull Image
		2. Create a container
		3. Connect container to network
		4. Start container
		5. Stream container logs
	*/

	if len(dc.ComposeService.Image) > 0 {
		err := api.PullDockerImage(ctx, cli, dc)
		if err != nil {
			// clean exit since we want other goroutines for fetching other images
			// to continue running
			fmt.Printf("\n\t service=%s error=%s", dc.ServiceName, err)
			return
		}
	}
	alreadyCreated, _, err := api.CreateContainer(ctx, cli, dc)
	if err != nil {
		// clean exit since we want other goroutines for fetching other images
		// to continue running
		fmt.Printf("\n\t service=%s error=%s", dc.ServiceName, err)
		return
	}

	if !alreadyCreated {
		err = api.ConnectNetwork(ctx, cli, dc)
		if err != nil {
			// create whitespace so that error is visible to human
			fmt.Printf("\n\t service=%s error=%s", dc.ServiceName, err)
			return
		}
	}

	err = api.ContainerStart(ctx, cli, dc)
	if err != nil {
		fmt.Printf("\n\t service=%s error=%s", dc.ServiceName, err)
		return
	}

	err = api.ContainerLogs(ctx, cli, dc)
	if err != nil {
		fmt.Printf("\n\t service=%s error=%s", dc.ServiceName, err)
		return
	}
}
