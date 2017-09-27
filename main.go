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
	"github.com/docker/docker/client"
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
		log.Fatal(err)
	}

	var dockerCyaml dockerComposeConfig
	if err := dockerCyaml.Parse(data); err != nil {
		log.Fatal(err)
	}

	for _, v := range dockerCyaml.Services {
		fmt.Println()
		fmt.Println(v.Image)
		fmt.Println()
		pullImage(v.Image)
	}

}

func pullImage(imagename string) {
	ctx := context.Background()
	cli, err := client.NewEnvClient()
	if err != nil {
		panic(err)
	}

	imagePullResp, err := cli.ImagePull(ctx, imagename, types.ImagePullOptions{})
	if err != nil {
		panic(err)
	}
	defer imagePullResp.Close()
	_, err = io.Copy(os.Stdout, imagePullResp)
	if err != nil {
		log.Println(err)
	}

	containerCreateResp, err := cli.ContainerCreate(ctx, &container.Config{Image: imagename}, &container.HostConfig{PublishAllPorts: true}, nil, "")
	if err != nil {
		log.Fatal(err)
	}

	err = cli.ContainerStart(ctx, containerCreateResp.ID, types.ContainerStartOptions{})
	if err != nil {
		panic(err)
	}

	containerLogResp, err := cli.ContainerLogs(ctx, containerCreateResp.ID, types.ContainerLogsOptions{ShowStdout: true, ShowStderr: true, Timestamps: true})
	if err != nil {
		panic(err)
	}
	defer containerLogResp.Close()

	_, err = io.Copy(os.Stdout, containerLogResp)
	if err != nil {
		log.Println(err)
	}
}
