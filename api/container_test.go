package api

import (
	"context"
	"log"
	"testing"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/client"
)

func TestCreateContainer(t *testing.T) {
	tt := []struct {
		s                 ServiceConfig
		networkName       string
		imgName           string
		dockerComposeFile string
		expected          string
		expectedErr       error
	}{
		{
			ServiceConfig{Image: "busybox", Restart: "unless-stopped"},
			"myNetworkName",
			"myImageName",
			"DockerFile",
			"myContainerId001",
			nil},
	}
	var ctx = context.Background()
	cli := &MockDockerClient{}
	for _, v := range tt {
		actual, err := CreateContainer(ctx, v.s, v.networkName, v.imgName, v.dockerComposeFile, cli)
		if err != nil {
			t.Errorf("\nran CreateContainer(%#+v) \ngot %s \nwanted %#+v", v.s, err, v.expectedErr)
		}
		if actual != v.expected {
			t.Errorf("\nran CreateContainer(%#+v) \ngot %s \nwanted %#+v", v.s, actual, v.expected)
		}
	}
}

func TestContainerStart(t *testing.T) {
	tt := []struct {
		input       string
		expectedErr error
	}{
		{"myContainerId", nil},
	}
	var ctx = context.Background()
	cli := &MockDockerClient{}
	for _, v := range tt {
		err := ContainerStart(ctx, v.input, cli)
		if err != nil {
			t.Errorf("\nran ContainerStart(%#+v) \ngot %s \nwanted %#+v", v.input, err, v.expectedErr)
		}
	}
}

func TestContainerLogs(t *testing.T) {
	tt := []struct {
		containerID string
		followLogs  bool
		expectedErr error
	}{
		{"myContainerId", true, nil},
		{"myContainerId", false, nil},
	}
	var ctx = context.Background()
	cli := &MockDockerClient{}
	for _, v := range tt {
		err := ContainerLogs(ctx, v.containerID, v.followLogs, cli)
		if err != nil {
			t.Errorf("\nran ContainerLogs(%#+v) \ngot %s \nwanted %#+v", v.containerID, err, v.expectedErr)
		}
	}
}

func BenchmarkCreateContainer(b *testing.B) {
	var ctx = context.Background()
	cli := &MockDockerClient{}
	for n := 0; n < b.N; n++ {
		_, _ = CreateContainer(
			ctx,
			ServiceConfig{Image: "busybox", Restart: "unless-stopped"},
			"mynetwork",
			"myImage",
			"dockerfile",
			cli)
	}
}

func BenchmarkContainerStart(b *testing.B) {
	var ctx = context.Background()
	cli := &MockDockerClient{}
	for n := 0; n < b.N; n++ {
		_ = ContainerStart(ctx, "containerId", cli)
	}
}

func BenchmarkContainerLogs(b *testing.B) {
	var ctx = context.Background()
	cli := &MockDockerClient{}
	for n := 0; n < b.N; n++ {
		_ = ContainerLogs(ctx, "containerID", true, cli)
	}
}

func TestContainerList(t *testing.T) {
	var ctx = context.Background()
	cli, err := client.NewEnvClient()
	if err != nil {
		log.Fatal(err, " :unable to intialize docker client")
	}
	defer cli.Close()

	filters := filters.NewArgs()
	filters.Add("label", "meli_service=meli_buildservice")
	listOpts := types.ContainerListOptions{Quiet: true, All: true, Filters: filters}
	_, err = cli.ContainerList(ctx, listOpts)
	t.Log("err", err)
	//litter.Dump(containers)
}
