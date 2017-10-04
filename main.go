package main

import (
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
	Volumes     []string    `yaml:"volumes,omitempty"`
	//Links          yaml.MaporColonSlice `yaml:"links,omitempty"`
}

type dockerComposeConfig struct {
	Version  string                   `yaml:"version,omitempty"`
	Services map[string]serviceConfig `yaml:"services"`
	Volumes  map[string]string        `yaml:"volumes,omitempty"`
	//networks map[string]     `yaml:"networks,omitempty"`
}

func main() {
	data, err := ioutil.ReadFile("docker-compose.yml")
	if err != nil {
		log.Fatal(err, "unable to read docker-compose file")
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
		log.Fatal(err, "unable to parse docker-compose file contents")
	}

	ctx := context.Background()

	// Create top level volumes, if any
	if len(dockerCyaml.Volumes) > 0 {
		fmt.Println("len", len(dockerCyaml.Volumes))
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
	for _, v := range dockerCyaml.Services {
		wg.Add(1)
		fmt.Println("docker service", v)
		go fakepullImage(ctx, v, networkID, networkName, &wg)
		//go pullImage(ctx, v, networkID, networkName, &wg)
	}
	wg.Wait()
}

func fakepullImage(ctx context.Context, s serviceConfig, networkName, networkID string, wg *sync.WaitGroup) {
	defer wg.Done()
	fmt.Println()
	if len(s.Volumes) > 0 {
		//Volumes map[string]struct{}
		// fmt.Println("service level volume:", s.Volumes)
		fmt.Printf("service level volume22: %#v", s.Volumes)
		fmt.Println()
		fmt.Printf("service level volume33: %#v", s.Volumes[0])

		x := fomatServiceVolumes(s.Volumes[0])
		fmt.Println()
		fmt.Printf("x %+v:", x[1])
		fmt.Println()
		// "Volumes": {
		//         "/home": {}
		//     }
	}
}

func pullImage(ctx context.Context, s serviceConfig, networkID, networkName string, wg *sync.WaitGroup) {
	defer wg.Done()
	fmt.Println()
	fmt.Println("docker servie:", s)
	fmt.Println()

	cli, err := client.NewEnvClient()
	if err != nil {
		log.Fatal(err, "unable to intialize docker client")
	}
	defer cli.Close()

	// 1. Pull Image
	formattedImageName := fomatImageName("containerFromBuild")
	if len(s.Image) > 0 {
		formattedImageName = fomatImageName(s.Image)
		// TODO move cli.ImagePull into image.go
		PullDockerImage(ctx, s.Image)
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
				log.Println(err, "unable to create a nat.Port")
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
	//2.5 build image
	imageName := s.Image
	if s.Build != (buildstruct{}) {
		imageName = BuildDockerImage(ctx, s.Build.Dockerfile)
	}

	//2.6 add volumes
	if len(s.Volumes) > 0 {
		//Volumes map[string]struct{}
		// fmt.Println("service level volume:", s.Volumes)
		fmt.Printf("service level volume22: %#v", s.Volumes)
		fmt.Printf("service level volume33: %#v", s.Volumes[0])

		x := fomatServiceVolumes(s.Volumes[0])
		fmt.Println()
		fmt.Printf("x %+v:", x[1])
		fmt.Println()
		// "Volumes": {
		//         "/home": {}
		//     }
	}

	// TODO: we should skip creating the container again if already exists
	// instead of creating a uniquely named container name
	containerCreateResp, err := cli.ContainerCreate(
		ctx,
		&container.Config{
			Image:        imageName,
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
		log.Println(err, "unable to create container")
	}

	// 3. Connect container to network
	err = cli.NetworkConnect(
		ctx,
		networkID,
		containerCreateResp.ID,
		&network.EndpointSettings{})
	if err != nil {
		log.Println(err, "unable to connect container to network")
	}

	// 4. Start container
	err = cli.ContainerStart(
		ctx,
		containerCreateResp.ID,
		types.ContainerStartOptions{})
	if err != nil {
		log.Println(err, "unable to start container")
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
		log.Println(err, "unable to get container logs")
	}
	defer containerLogResp.Close()
	_, err = io.Copy(os.Stdout, containerLogResp)
	if err != nil {
		log.Println(err, "unable to write to stdout")
	}
}
