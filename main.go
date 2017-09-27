package main

import (
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/docker/client"

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
	Dockerfile string `yaml:"dockerfile,omitempty"`
	//Environment    yaml.MaporEqualSlice `yaml:"environment,omitempty"`
	Image string `yaml:"image,omitempty"`
	//Links          yaml.MaporColonSlice `yaml:"links,omitempty"`
	Name        string   `yaml:"name,omitempty"`
	Ports       []string `yaml:"ports,omitempty"`
	Restart     string   `yaml:"restart,omitempty"`
	Volumes     []string `yaml:"volumes,omitempty"`
	VolumesFrom []string `yaml:"volumes_from,omitempty"`
	Expose      []string `yaml:"expose,omitempty"`
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
	networkID, err := getNetwork()
	if err != nil {
		log.Fatal(err)
	}

	var dockerCyaml dockerComposeConfig
	if err := dockerCyaml.Parse(data); err != nil {
		log.Fatal(errors.Wrap(err, "unable to parse docker-compose file contents"))
	}

	for _, v := range dockerCyaml.Services {
		fmt.Println()
		fmt.Println("image, networkID, name: ", v.Image, networkID, v.Name)
		fmt.Println()
		pullImage(v.Image, networkID)
	}

}

func pullImage(imagename, networkID string) {
	ctx := context.Background()
	cli, err := client.NewEnvClient()
	if err != nil {
		panic(errors.Wrap(err, "unable to intialize docker client"))
	}

	imagePullResp, err := cli.ImagePull(
		ctx,
		imagename,
		types.ImagePullOptions{})
	if err != nil {
		panic(errors.Wrap(err, "unable to pull image"))
	}
	defer imagePullResp.Close()
	_, err = io.Copy(os.Stdout, imagePullResp)
	if err != nil {
		log.Println(errors.Wrap(err, "unable to write to stdout"))
	}

	containerCreateResp, err := cli.ContainerCreate(
		ctx,
		&container.Config{Image: imagename},
		&container.HostConfig{PublishAllPorts: true},
		nil,
		"")
	if err != nil {
		log.Fatal(errors.Wrap(err, "unable to create container"))
	}

	//
	//Connect container to network
	err = cli.NetworkConnect(
		ctx,
		networkID,
		containerCreateResp.ID,
		&network.EndpointSettings{})
	if err != nil {
		log.Fatal(errors.Wrap(err, "unable to connect container to network"))
	}
	//

	err = cli.ContainerStart(
		ctx,
		containerCreateResp.ID,
		types.ContainerStartOptions{})
	if err != nil {
		panic(errors.Wrap(err, "unable to start container"))
	}

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
