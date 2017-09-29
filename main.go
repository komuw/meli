package main

import (
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"sync"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"

	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"
)

/*
version: '3'
services:
  redis:
    image: 'redis:3.0-alpine'

  busybox:
    image: busybox
*/
type serviceConfig struct {
	Build string `yaml:"build,omitempty"`
	//Command        yaml.Command         `yaml:"command,flow,omitempty"`
	Dockerfile  string   `yaml:"dockerfile,omitempty"`
	Environment []string `yaml:"environment,omitempty"`
	Image       string   `yaml:"image,omitempty"`
	//Links          yaml.MaporColonSlice `yaml:"links,omitempty"`
	Name        string   `yaml:"name,omitempty"`
	Ports       []string `yaml:"ports,omitempty"`
	Restart     string   `yaml:"restart,omitempty"`
	Volumes     []string `yaml:"volumes,omitempty"`
	VolumesFrom []string `yaml:"volumes_from,omitempty"`
	Expose      []string `yaml:"expose,omitempty"`
	Labels      []string `yaml:"labels,omitempty"`
}

type dockerComposeConfig struct {
	Version  string                   `yaml:"version,omitempty"`
	Services map[string]serviceConfig `yaml:"services"`
	//networks map[string]     `yaml:"networks,omitempty"`
	//volumes map[string]                  `yaml:"volumes,omitempty"`
}

func (dcy *dockerComposeConfig) Parse(data []byte) error {
	return yaml.Unmarshal(data, dcy)
}

func main() {
	data, err := ioutil.ReadFile("docker-compose.yml")
	if err != nil {
		log.Fatal(errors.Wrap(err, "unable to read docker-compose file"))
	}
	curentDir, err := os.Getwd()
	if err != nil {
		log.Fatal(err, "unable to get the current working directory")
	}
	networkName := "meli_network_" + getCwdName(curentDir)
	networkID, err := getNetwork(networkName)

	if err != nil {
		log.Fatal(err)
	}

	var dockerCyaml dockerComposeConfig
	if err := dockerCyaml.Parse(data); err != nil {
		log.Fatal(errors.Wrap(err, "unable to parse docker-compose file contents"))
	}

	var wg sync.WaitGroup

	for _, v := range dockerCyaml.Services {
		wg.Add(1)
		fmt.Println("docker service", v)
		//go fakepullImage(v, networkID, networkName, &wg)
		go pullImage(v, networkID, networkName, &wg)
	}
	wg.Wait()

}

func fakepullImage(s serviceConfig, networkName, networkID string, wg *sync.WaitGroup) {
	defer wg.Done()
	fmt.Println()
}

func pullImage(s serviceConfig, networkID, networkName string, wg *sync.WaitGroup) {
	defer wg.Done()
	formattedImageName := fomatImageName(s.Image)
	fmt.Println()
	fmt.Println("dockerImage, networkID, name:", s.Image, networkID, formattedImageName)
	fmt.Println()
	ctx := context.Background()
	cli, err := client.NewEnvClient()
	if err != nil {
		panic(errors.Wrap(err, "unable to intialize docker client"))
	}

	// 1. Pull Image
	imagePullResp, err := cli.ImagePull(
		ctx,
		s.Image,
		types.ImagePullOptions{})
	if err != nil {
		panic(errors.Wrap(err, "unable to pull image"))
	}
	defer imagePullResp.Close()
	_, err = io.Copy(os.Stdout, imagePullResp)
	if err != nil {
		log.Println(errors.Wrap(err, "unable to write to stdout"))
	}

	// 2. Create a container
	// 2.1 make labels
	labelsMap := make(map[string]string)
	if len(s.Labels) > 0 {
		for _, v := range s.Labels {
			onelabel := fomatLabels(v)
			labelsMap[onelabel[0]] = onelabel[1]
			fmt.Println("labelsMap", labelsMap)
		}
	}
	//2.2 make ports
	portsMap := make(map[nat.Port]struct{})
	type emptyStruct struct{}
	portBindingMap := make(map[nat.Port][]nat.PortBinding)
	if len(s.Ports) > 0 {
		for _, v := range s.Ports {
			oneport := fomatPorts(v)
			hostport := oneport[0]
			containerport := oneport[1]
			port, err := nat.NewPort("tcp", containerport)
			myPortBinding := nat.PortBinding{HostPort: hostport}
			if err != nil {
				log.Println(errors.Wrap(err, "unable to create a nat.Port"))
			}
			portsMap[port] = emptyStruct{}
			portBindingMap[port] = []nat.PortBinding{myPortBinding}
		}
	}
	// TODO: we should skip creating the container again if already exists
	// instead of creating a uniquely named container name
	containerCreateResp, err := cli.ContainerCreate(
		ctx,
		&container.Config{Image: s.Image, Labels: labelsMap, Env: s.Environment, ExposedPorts: portsMap},
		&container.HostConfig{PublishAllPorts: false, PortBindings: portBindingMap, NetworkMode: container.NetworkMode(networkName)},
		nil,
		formattedImageName)
	if err != nil {
		log.Fatal(errors.Wrap(err, "unable to create container"))
	}

	// 3. Connect container to network
	err = cli.NetworkConnect(
		ctx,
		networkID,
		containerCreateResp.ID,
		&network.EndpointSettings{})
	if err != nil {
		log.Fatal(errors.Wrap(err, "unable to connect container to network"))
	}

	// 4. Start container
	err = cli.ContainerStart(
		ctx,
		containerCreateResp.ID,
		types.ContainerStartOptions{})
	if err != nil {
		panic(errors.Wrap(err, "unable to start container"))
	}

	// 5. Stream container logs to stdOut
	containerLogResp, err := cli.ContainerLogs(
		ctx,
		containerCreateResp.ID,
		types.ContainerLogsOptions{
			ShowStdout: true,
			ShowStderr: true,
			Timestamps: true})
	if err != nil {
		panic(errors.Wrap(err, "unable to get container logs"))
	}
	defer containerLogResp.Close()
	_, err = io.Copy(os.Stdout, containerLogResp)
	if err != nil {
		log.Println(errors.Wrap(err, "unable to write to stdout"))
	}
}
