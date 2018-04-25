/*
Package meli provides programming interface to interact with the docker daemon.

Example usage:
	package main

	import (
		"context"
		"log"
		"os"

		"github.com/docker/docker/client"
		"github.com/komuw/meli"
	)

	func main() {
		dc := &meli.DockerContainer{
			ComposeService: meli.ComposeService{Image: "busybox"},
			LogMedium:      os.Stdout,
			FollowLogs:     true}

		ctx := context.Background()
		cli, err := client.NewEnvClient()
		if err != nil {
			log.Fatal(err, " :unable to intialize docker client")
		}
		defer cli.Close()

		meli.GetAuth() // read dockerhub info
		err = meli.PullDockerImage(ctx, cli, dc)
		log.Println(err)
	}

*/
package meli

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/api/types/strslice"
	"github.com/docker/go-connections/nat"
)

// CreateContainer creates a docker container
func CreateContainer(ctx context.Context, cli APIclient, dc *DockerContainer) (bool, string, error) {
	// 1. make labels
	labelsMap := make(map[string]string)
	if len(dc.ComposeService.Labels) > 0 {
		for _, v := range dc.ComposeService.Labels {
			onelabel := formatLabels(v)
			labelsMap[onelabel[0]] = onelabel[1]
		}
	}

	// reuse container if already running
	// only reuse containers if we aren't rebuilding
	meliService := labelsMap["meli_service"]
	filters := filters.NewArgs()
	filters.Add("label", fmt.Sprintf("meli_service=%s", meliService))
	containers, err := cli.ContainerList(ctx, types.ContainerListOptions{
		Quiet:   true,
		All:     true,
		Filters: filters})
	if err != nil {
		fmt.Println(" :unable to list containers")
	}
	if len(containers) > 0 {
		if !dc.Rebuild {
			dc.UpdateContainerID(containers[0].ID)
			return true, containers[0].ID, nil
		}
		shadowErr := cli.ContainerRemove(ctx, containers[0].ID, types.ContainerRemoveOptions{Force: true})
		if err != nil {
			fmt.Println(shadowErr, " :unable to remove existing container, ", containers[0].ID)
		}
	}

	// 2. make ports
	portsMap := make(map[nat.Port]struct{})
	portBindingMap := make(map[nat.Port][]nat.PortBinding)
	if len(dc.ComposeService.Ports) > 0 {
		for _, v := range dc.ComposeService.Ports {
			oneport := formatPorts(v)
			// issues/96
			hostport := ""
			containerport := ""
			if len(oneport) == 1 {
				hostport = ""
				containerport = oneport[0]
			} else {
				hostport = oneport[0]
				containerport = oneport[1]
			}
			myPortBinding := nat.PortBinding{HostPort: hostport}
			port, shadowErr := nat.NewPort("tcp", containerport)
			if shadowErr != nil {
				fmt.Println(shadowErr, " :unable to create a nat.Port")
			}
			portsMap[port] = emptyStruct{}
			portBindingMap[port] = []nat.PortBinding{myPortBinding}
		}
	}
	// 3. create command
	cmd := strslice.StrSlice{}
	if dc.ComposeService.Command != "" {
		sliceCommand := strings.Fields(dc.ComposeService.Command)
		cmd = strslice.StrSlice(sliceCommand)
	}
	// 4. create restart policy
	restartPolicy := container.RestartPolicy{}
	if dc.ComposeService.Restart != "" {
		// You cannot set MaximumRetryCount for the following restart policies;
		// always, no, unless-stopped
		if dc.ComposeService.Restart == "on-failure" {
			restartPolicy = container.RestartPolicy{Name: dc.ComposeService.Restart, MaximumRetryCount: 3}
		} else {
			restartPolicy = container.RestartPolicy{Name: dc.ComposeService.Restart}
		}

	}
	// 5. build image
	imageNamePtr := &dc.ComposeService.Image
	if dc.ComposeService.Build != (Buildstruct{}) {
		imageName, shadowErr := BuildDockerImage(ctx, cli, dc)
		if shadowErr != nil {
			return false, "", &popagateError{originalErr: shadowErr}
		}
		// done this way so that we can manipulate the value of the
		// imageName inside this scope
		imageNamePtr = &imageName
	}
	imageName := *imageNamePtr

	// 6. add volumes
	volume := make(map[string]struct{})
	binds := []string{}
	if len(dc.ComposeService.Volumes) > 0 {
		for _, v := range dc.ComposeService.Volumes {
			vol := formatServiceVolumes(v, dc.DockerComposeFile)
			volume[vol[1]] = emptyStruct{}
			// TODO: handle other read/write modes
			whatToBind := vol[0] + ":" + vol[1] + ":rw"
			binds = append(binds, whatToBind)
		}
	}

	containerCreateResp, err := cli.ContainerCreate(
		ctx,
		&container.Config{
			Image:        imageName,
			Labels:       labelsMap,
			Env:          dc.ComposeService.Environment,
			ExposedPorts: portsMap,
			Cmd:          cmd,
			Volumes:      volume},
		&container.HostConfig{
			DNS: []string{
				"8.8.8.8",
				"8.8.4.4",
				"2001:4860:4860::8888",
				"2001:4860:4860::8844"},
			DNSSearch: []string{
				"8.8.8.8",
				"8.8.4.4",
				"2001:4860:4860::8888",
				"2001:4860:4860::8844"},
			PublishAllPorts: false,
			PortBindings:    portBindingMap,
			NetworkMode:     container.NetworkMode(dc.NetworkName),
			RestartPolicy:   restartPolicy,
			Binds:           binds,
			Links:           dc.ComposeService.Links},
		nil,
		dc.ServiceName)
	if err != nil {
		return false, "", &popagateError{
			originalErr: err,
			newErr:      errors.New(" :unable to create container")}
	}

	dc.UpdateContainerID(containerCreateResp.ID)
	return false, containerCreateResp.ID, nil
}

// ContainerStart starts a docker container via docker daemon server
func ContainerStart(ctx context.Context, cli APIclient, dc *DockerContainer) error {
	err := cli.ContainerStart(
		ctx,
		dc.ContainerID,
		types.ContainerStartOptions{})
	if err != nil {
		return &popagateError{
			originalErr: err,
			newErr:      fmt.Errorf(" :unable to start container %s", dc.ContainerID)}

	}
	return nil
}

// ContainerLogs returns the logs generated by a container in an io.ReadCloser.
func ContainerLogs(ctx context.Context, cli APIclient, dc *DockerContainer) error {
	containerLogResp, err := cli.ContainerLogs(
		ctx,
		dc.ContainerID,
		types.ContainerLogsOptions{
			ShowStdout: true,
			ShowStderr: true,
			Timestamps: true,
			Follow:     dc.FollowLogs,
			Details:    true,
			Tail:       "all"})

	if err != nil {
		if err != nil {
			return &popagateError{
				originalErr: err,
				newErr:      fmt.Errorf(" :unable to get container logs %s", dc.ContainerID)}
		}
	}

	scanner := bufio.NewScanner(containerLogResp)
	for scanner.Scan() {
		fmt.Fprintln(dc.LogMedium, dc.ServiceName, "::", scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		fmt.Println(" :unable to log output for container", dc.ContainerID, err)
	}

	containerLogResp.Close()
	return nil
}
