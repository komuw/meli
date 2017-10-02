package main

import (
	"archive/tar"
	"bytes"
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"sync"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/docker/api/types/strslice"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"

	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"
)

/* DOCS:
1. https://godoc.org/github.com/moby/moby/client
2. https://docs.docker.com/engine/api/v1.31/
*/

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
	//Volumes     []string `yaml:"volumes,omitempty"`
	//VolumesFrom []string `yaml:"volumes_from,omitempty"`
	//Links          yaml.MaporColonSlice `yaml:"links,omitempty"`
}

type dockerComposeConfig struct {
	Version  string                   `yaml:"version,omitempty"`
	Services map[string]serviceConfig `yaml:"services"`
	//networks map[string]     `yaml:"networks,omitempty"`
	//volumes map[string]                  `yaml:"volumes,omitempty"`
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
		log.Fatal(err, "unable to create/get network")
	}

	var dockerCyaml dockerComposeConfig
	err = yaml.Unmarshal([]byte(data), &dockerCyaml)
	if err != nil {
		log.Fatal(errors.Wrap(err, "unable to parse docker-compose file contents"))
	}

	var wg sync.WaitGroup
	for _, v := range dockerCyaml.Services {
		wg.Add(1)
		fmt.Println("docker service", v)
		go fakepullImage(v, networkID, networkName, &wg)
		//go pullImage(v, networkID, networkName, &wg)
	}
	wg.Wait()
}

func fakepullImage(s serviceConfig, networkName, networkID string, wg *sync.WaitGroup) {
	defer wg.Done()
	fmt.Println()

	if s.Build != (buildstruct{}) {
		fmt.Printf("%+v\n", s)

		ctx := context.Background()
		cli, err := client.NewEnvClient()
		if err != nil {
			log.Fatal(errors.Wrap(err, "unable to intialize docker client"))
		}
		defer cli.Close()
		buf := new(bytes.Buffer)
		tw := tar.NewWriter(buf)
		defer tw.Close()

		dockerFile := s.Build.Dockerfile
		if s.Build.Dockerfile == "" {
			dockerFile = "Dockerfile"
		}
		dockerFileReader, err := os.Open(dockerFile)
		if err != nil {
			log.Fatal(err, " :unable to open Dockerfile")
		}
		readDockerFile, err := ioutil.ReadAll(dockerFileReader)
		if err != nil {
			log.Fatal(err, " :unable to read dockerfile")
		}

		tarHeader := &tar.Header{
			Name: dockerFile,
			Size: int64(len(readDockerFile)),
		}
		err = tw.WriteHeader(tarHeader)
		if err != nil {
			log.Fatal(err, " :unable to write tar header")
		}
		_, err = tw.Write(readDockerFile)
		if err != nil {
			log.Fatal(err, " :unable to write tar body")
		}
		dockerFileTarReader := bytes.NewReader(buf.Bytes())
		imageBuildResponse, err := cli.ImageBuild(
			ctx,
			dockerFileTarReader,
			types.ImageBuildOptions{
				//PullParent:     true,
				//Squash:     true, currently only supported in experimenta mode
				Tags:           []string{"meli_" + strings.ToLower(dockerFile)},
				Remove:         true, //remove intermediary containers after build
				ForceRemove:    true,
				SuppressOutput: false,
				Dockerfile:     dockerFile,
				Context:        dockerFileTarReader})
		if err != nil {
			log.Fatal(err, " :unable to build docker image")
		}
		defer imageBuildResponse.Body.Close()
		_, err = io.Copy(os.Stdout, imageBuildResponse.Body)
		if err != nil {
			log.Fatal(err, " :unable to read image build response")
		}

	}

}

func pullImage(s serviceConfig, networkID, networkName string, wg *sync.WaitGroup) {
	defer wg.Done()
	fmt.Println()
	fmt.Println("docker servie:", s)
	fmt.Println()
	formattedImageName := fomatImageName(s.Image)
	ctx := context.Background()
	cli, err := client.NewEnvClient()
	if err != nil {
		log.Fatal(errors.Wrap(err, "unable to intialize docker client"))
	}
	defer cli.Close()

	// 1. Pull Image
	imagePullResp, err := cli.ImagePull(
		ctx,
		s.Image,
		types.ImagePullOptions{})
	if err != nil {
		log.Println(errors.Wrap(err, "unable to pull image"))
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
	//2.3 create command
	cmd := strslice.StrSlice{}
	if s.Command != "" {
		sliceCommand := strings.Fields(s.Command)
		cmd = strslice.StrSlice(sliceCommand)
	}
	//2.4 create restart policy
	restartPolicy := container.RestartPolicy{}
	if s.Restart != "" {
		// you cannot set MaximumRetryCount for the following restart policies;
		// always, no, unless-stopped
		if s.Restart == "on-failure" {
			restartPolicy = container.RestartPolicy{Name: s.Restart, MaximumRetryCount: 3}
		} else {
			restartPolicy = container.RestartPolicy{Name: s.Restart}
		}

	}

	// TODO: we should skip creating the container again if already exists
	// instead of creating a uniquely named container name
	containerCreateResp, err := cli.ContainerCreate(
		ctx,
		&container.Config{
			Image:        s.Image,
			Labels:       labelsMap,
			Env:          s.Environment,
			ExposedPorts: portsMap,
			Cmd:          cmd},
		&container.HostConfig{
			PublishAllPorts: false,
			PortBindings:    portBindingMap,
			NetworkMode:     container.NetworkMode(networkName),
			RestartPolicy:   restartPolicy},
		nil,
		formattedImageName)
	if err != nil {
		log.Println(errors.Wrap(err, "unable to create container"))
	}

	// 3. Connect container to network
	err = cli.NetworkConnect(
		ctx,
		networkID,
		containerCreateResp.ID,
		&network.EndpointSettings{})
	if err != nil {
		log.Println(errors.Wrap(err, "unable to connect container to network"))
	}

	// 4. Start container
	err = cli.ContainerStart(
		ctx,
		containerCreateResp.ID,
		types.ContainerStartOptions{})
	if err != nil {
		log.Println(errors.Wrap(err, "unable to start container"))
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
		log.Println(errors.Wrap(err, "unable to get container logs"))
	}
	defer containerLogResp.Close()
	_, err = io.Copy(os.Stdout, containerLogResp)
	if err != nil {
		log.Println(errors.Wrap(err, "unable to write to stdout"))
	}
}
