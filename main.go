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

var version = "master"

func main() {
	followLogs, dockerComposeFile := cli.Cli()

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

	// Create top level volumes, if any
	if len(dockerCyaml.Volumes) > 0 {
		for k := range dockerCyaml.Volumes {
			// TODO we need to synchronise here else we'll get a race
			// but I think we can get away for now because:
			// 1. there are on average a lot more containers in a compose file
			// than volumes, so the sync in the for loop for containers is enough
			// 2. since we intend to stream logs as containers run(see; issues/24);
			// then meli will be up long enough for the volume creation goroutines to have finished.
			go api.CreateDockerVolume(ctx, "meli_"+k, "local", cli)
		}
	}

	var wg sync.WaitGroup
	for k, v := range dockerCyaml.Services {
		wg.Add(1)
		v.Labels = append(v.Labels, fmt.Sprintf("meli_service=meli_%s", k))
		//go fakestartContainers(ctx, k, v, networkID, networkName, &wg, followLogs,  dockerComposeFile)
		go startContainers(
			ctx,
			k,
			v,
			networkID,
			networkName,
			&wg,
			followLogs,
			dockerComposeFile,
			cli)
	}
	wg.Wait()
}

func fakestartContainers(
	ctx context.Context,
	k string,
	s api.ServiceConfig,
	networkName, networkID string,
	wg *sync.WaitGroup,
	followLogs bool,
	dockerComposeFile string) {
	defer wg.Done()
}

func startContainers(
	ctx context.Context,
	k string,
	s api.ServiceConfig,
	networkID, networkName string,
	wg *sync.WaitGroup,
	followLogs bool,
	dockerComposeFile string,
	cli *client.Client) {
	defer wg.Done()

	/*
		1. Pull Image
		2. Create a container
		3. Connect container to network
		4. Start container
		5. Stream container logs
	*/

	formattedContainerName := api.FormatContainerName(k)
	if len(s.Image) > 0 {
		err := api.PullDockerImage(ctx, s.Image, cli)
		if err != nil {
			// clean exit since we want other goroutines for fetching other images
			// to continue running
			log.Printf("\n\t service=%s error=%s", k, err)
			return
		}
	}
	alreadyCreated, containerID, err := api.CreateContainer(
		ctx,
		s,
		networkName,
		formattedContainerName,
		dockerComposeFile,
		cli)
	if err != nil {
		// clean exit since we want other goroutines for fetching other images
		// to continue running
		log.Printf("\n\t service=%s error=%s", k, err)
		return
	}

	if !alreadyCreated {
		err = api.ConnectNetwork(
			ctx,
			networkID,
			containerID,
			cli)
		if err != nil {
			// create whitespace so that error is visible to human
			log.Printf("\n\t service=%s error=%s", k, err)
			return
		}
	}

	err = api.ContainerStart(
		ctx,
		containerID,
		cli)
	if err != nil {
		log.Printf("\n\t service=%s error=%s", k, err)
		return
	}

	err = api.ContainerLogs(
		ctx,
		containerID,
		followLogs,
		cli)
	if err != nil {
		log.Printf("\n\t service=%s error=%s", k, err)
		return
	}
}
