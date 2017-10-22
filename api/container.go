package api

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/api/types/strslice"
	"github.com/docker/go-connections/nat"
)

//func CreateContainer(ctx context.Context, s ServiceConfig, k, networkName, formattedImageName, dockerComposeFile string, cli MeliAPiClient) (bool, string, error) {

func CreateContainer(ctx context.Context, cli MeliAPiClient, xyz *XYZ) (bool, string, error) {
	formattedImageName := FormatImageName(xyz.ServiceName)

	// 2.1 make labels
	labelsMap := make(map[string]string)
	if len(xyz.ServiceConfig.Labels) > 0 {
		for _, v := range xyz.ServiceConfig.Labels {
			onelabel := FormatLabels(v)
			labelsMap[onelabel[0]] = onelabel[1]
		}
	}

	// reuse container if already running
	meliService := labelsMap["meli_service"]
	filters := filters.NewArgs()
	filters.Add("label", fmt.Sprintf("meli_service=%s", meliService))
	containers, err := cli.ContainerList(ctx, types.ContainerListOptions{
		Quiet:   true,
		All:     true,
		Filters: filters})
	if err != nil {
		log.Println(" :unable to list containers")
	}
	if len(containers) > 0 {
		return true, containers[0].ID, nil
	}

	//2.2 make ports
	portsMap := make(map[nat.Port]struct{})
	portBindingMap := make(map[nat.Port][]nat.PortBinding)
	if len(xyz.ServiceConfig.Ports) > 0 {
		for _, v := range xyz.ServiceConfig.Ports {
			oneport := FormatPorts(v)
			hostport := oneport[0]
			containerport := oneport[1]
			port, err := nat.NewPort("tcp", containerport)
			myPortBinding := nat.PortBinding{HostPort: hostport}
			if err != nil {
				log.Println(err, " :unable to create a nat.Port")
			}
			portsMap[port] = EmptyStruct{}
			portBindingMap[port] = []nat.PortBinding{myPortBinding}
		}
	}
	//2.3 create command
	cmd := strslice.StrSlice{}
	if xyz.ServiceConfig.Command != "" {
		sliceCommand := strings.Fields(xyz.ServiceConfig.Command)
		cmd = strslice.StrSlice(sliceCommand)
	}
	//2.4 create restart policy
	restartPolicy := container.RestartPolicy{}
	if xyz.ServiceConfig.Restart != "" {
		// You cannot set MaximumRetryCount for the following restart policies;
		// always, no, unless-stopped
		if xyz.ServiceConfig.Restart == "on-failure" {
			restartPolicy = container.RestartPolicy{Name: xyz.ServiceConfig.Restart, MaximumRetryCount: 3}
		} else {
			restartPolicy = container.RestartPolicy{Name: xyz.ServiceConfig.Restart}
		}

	}
	//2.5 build image
	imageNamePtr := &xyz.ServiceConfig.Image
	if xyz.ServiceConfig.Build != (Buildstruct{}) {
		imageName, err := BuildDockerImage(ctx, cli, xyz)
		if err != nil {
			return false, "", &popagateError{originalErr: err}
		}
		// done this way so that we can manipulate the value of the
		// imageName inside this scope
		imageNamePtr = &imageName
	}
	imageName := *imageNamePtr

	//2.6 add volumes
	volume := make(map[string]struct{})
	binds := []string{}
	if len(xyz.ServiceConfig.Volumes) > 0 {
		for _, v := range xyz.ServiceConfig.Volumes {
			vol := FormatServiceVolumes(v, xyz.DockerComposeFile)
			volume[vol[1]] = EmptyStruct{}
			// TODO: handle other read/write modes
			whatToBind := vol[0] + ":" + vol[1] + ":rw"
			binds = append(binds, whatToBind)
		}
	}

	// TODO: we should skip creating the container again if already exists
	// instead of creating a uniquely named container name
	containerCreateResp, err := cli.ContainerCreate(
		ctx,
		&container.Config{
			Image:        imageName,
			Labels:       labelsMap,
			Env:          xyz.ServiceConfig.Environment,
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
			NetworkMode:     container.NetworkMode(xyz.NetworkName),
			RestartPolicy:   restartPolicy,
			Binds:           binds,
			Links:           xyz.ServiceConfig.Links},
		nil,
		formattedImageName)
	if err != nil {
		return false, "", &popagateError{
			originalErr: err,
			newErr:      errors.New(" :unable to create container")}
	}

	xyz.UpdateContainerID(containerCreateResp.ID)
	return false, containerCreateResp.ID, nil
}

func ContainerStart(ctx context.Context, cli MeliAPiClient, xyz *XYZ) error {
	err := cli.ContainerStart(
		ctx,
		xyz.ContainerID,
		types.ContainerStartOptions{})
	if err != nil {
		return &popagateError{
			originalErr: err,
			newErr:      fmt.Errorf(" :unable to start container %s", xyz.ContainerID)}

	}
	return nil
}

func ContainerLogs(ctx context.Context, cli MeliAPiClient, xyz *XYZ) error {
	containerLogResp, err := cli.ContainerLogs(
		ctx,
		xyz.ContainerID,
		types.ContainerLogsOptions{
			ShowStdout: true,
			ShowStderr: true,
			Timestamps: true,
			Follow:     xyz.FollowLogs,
			Details:    true,
			Tail:       "all"})

	if err != nil {
		if err != nil {
			return &popagateError{
				originalErr: err,
				newErr:      fmt.Errorf(" :unable to get container logs %s", xyz.ContainerID)}
		}
	}
	defer containerLogResp.Close()

	// supplying your own buffer is perfomant than letting the system do it for you
	buff := make([]byte, 2048)
	io.CopyBuffer(os.Stdout, containerLogResp, buff)

	return nil
}
