package main

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/strslice"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
)

func CreateContainer(ctx context.Context, s serviceConfig, networkName, formattedImageName string) container.ContainerCreateCreatedBody {
	cli, err := client.NewEnvClient()
	if err != nil {
		log.Println(err, "unable to intialize docker client")
	}
	defer cli.Close()

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
	volume := make(map[string]struct{})
	binds := []string{}
	if len(s.Volumes) > 0 {
		vol := fomatServiceVolumes(s.Volumes[0])
		volume[vol[1]] = emptyStruct{}
		// TODO: handle other read/write modes
		whatToBind := "meli_" + vol[0] + ":" + vol[1] + ":rw"
		binds = append(binds, whatToBind)
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
			Cmd:          cmd,
			Volumes:      volume},
		&container.HostConfig{
			PublishAllPorts: false,
			PortBindings:    portBindingMap,
			NetworkMode:     container.NetworkMode(networkName),
			RestartPolicy:   restartPolicy,
			Binds:           binds},
		nil,
		formattedImageName)
	if err != nil {
		log.Println(err, "unable to create container")
	}

	return containerCreateResp
}
