package api

import (
	"context"
	"testing"
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
