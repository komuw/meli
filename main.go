package main

import (
	"context"
	"io/ioutil"
	"log"
	"os"
	"sync"

	"gopkg.in/yaml.v2"
)

/* DOCS:
1. https://godoc.org/github.com/moby/moby/client
2. https://docs.docker.com/engine/api/v1.31/
*/

var Version = "0.0.0.1"

type emptyStruct struct{}

type buildstruct struct {
	// remember to use caps so that they can be exported
	Context    string `yaml:"context,omitempty"`
	Dockerfile string `yaml:"dockerfile,omitempty"`
}

type serviceConfig struct {
	Image       string      `yaml:"image,omitempty"`
	Ports       []string    `yaml:"ports,omitempty"`
	Labels      []string    `yaml:"labels,omitempty"`
	Environment []string    `yaml:"environment,omitempty"`
	Command     string      `yaml:"command,flow,omitempty"`
	Restart     string      `yaml:"restart,omitempty"`
	Build       buildstruct `yaml:"build,omitempty"`
	Volumes     []string    `yaml:"volumes,omitempty"`
}

type dockerComposeConfig struct {
	Version  string                   `yaml:"version,omitempty"`
	Services map[string]serviceConfig `yaml:"services"`
	Volumes  map[string]string        `yaml:"volumes,omitempty"`
}

func main() {
	showLogs := Cli()

	data, err := ioutil.ReadFile("docker-compose.yml")
	if err != nil {
		log.Fatal(err, " :unable to read docker-compose file")
	}
	curentDir, err := os.Getwd()
	if err != nil {
		log.Fatal(err, " :unable to get the current working directory")
	}
	networkName := "meli_network_" + getCwdName(curentDir)
	networkID, err := GetNetwork(networkName)
	if err != nil {
		log.Fatal(err, " :unable to create/get network")
	}

	var dockerCyaml dockerComposeConfig
	err = yaml.Unmarshal([]byte(data), &dockerCyaml)
	if err != nil {
		log.Fatal(err, " :unable to parse docker-compose file contents")
	}

	ctx := context.Background()

	// Create top level volumes, if any
	if len(dockerCyaml.Volumes) > 0 {
		for k := range dockerCyaml.Volumes {
			// TODO we need to synchronise here else we'll get a race
			// but I think we can get away for now because:
			// 1. there are on average a lot more containers in a compose file
			// than volumes, so the sync in the for loop for containers is enough
			// 2. since we intend to stream logs as containers run(see; issues/24);
			// then meli will be up long enough for the volume creation goroutines to have finished.
			go CreateDockerVolume(ctx, "meli_"+k, "local")
		}
	}

	var wg sync.WaitGroup
	for k, v := range dockerCyaml.Services {
		wg.Add(1)
		//go fakestartContainers(ctx, k, v, networkID, networkName, &wg, showLogs)
		go startContainers(ctx, k, v, networkID, networkName, &wg, showLogs)
	}
	wg.Wait()
}

func fakestartContainers(ctx context.Context, k string, s serviceConfig, networkName, networkID string, wg *sync.WaitGroup, showLogs bool) {
	defer wg.Done()
}

func startContainers(ctx context.Context, k string, s serviceConfig, networkID, networkName string, wg *sync.WaitGroup, showLogs bool) {
	defer wg.Done()

	/*
		1. Pull Image
		2. Create a container
		3. Connect container to network
		4. Start container
		5. Stream container logs
	*/

	formattedContainerName := formatContainerName(k)
	if len(s.Image) > 0 {
		err := PullDockerImage(ctx, s.Image)
		if err != nil {
			// clean exit since we want other goroutines for fetching other images
			// to continue running
			log.Println("\n", err)
			return
		}
	}
	containerID, err := CreateContainer(
		ctx,
		s,
		networkName,
		formattedContainerName)
	if err != nil {
		// clean exit since we want other goroutines for fetching other images
		// to continue running
		log.Println("\n", err)
		return
	}

	err = ConnectNetwork(
		ctx,
		networkID,
		containerID)
	if err != nil {
		// create whitespace so that error is visible to human
		log.Println("\n", err)
		return
	}

	err = ContainerStart(
		ctx, containerID)
	if err != nil {
		log.Println("\n", err)
		return
	}

	err = ContainerLogs(
		ctx,
		containerID,
		showLogs)
	if err != nil {
		log.Println("\n", err)
		return
	}
}
