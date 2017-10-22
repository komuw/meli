package api

import (
	"context"
	"testing"
)

func TestCreateContainer(t *testing.T) {
	tt := []struct {
		dc         *DockerContainer
		expected    string
		expectedErr error
	}{
		{
			&DockerContainer{
				ServiceConfig:     ServiceConfig{Image: "busybox", Restart: "unless-stopped"},
				ServiceName:       "myservice",
				NetworkName:       "myNetworkName",
				DockerComposeFile: "DockerFile",
				ContainerID:       "myExistingContainerId00912",
			},
			"myExistingContainerId00912",
			nil},
	}
	var ctx = context.Background()
	cli := &MockDockerClient{}
	for _, v := range tt {
		//CreateContainer(ctx context.Context, cli MeliAPiClient, dc *DockerContainer)
		alreadyCreated, actual, err := CreateContainer(ctx, cli, v.dc)
		if err != nil {
			t.Errorf("\nCalled CreateContainer(%#+v) \ngot %s \nwanted %#+v", v.dc, err, v.expectedErr)
		}
		if actual != v.expected {
			t.Errorf("\nCalled CreateContainer(%#+v) \ngot %s \nwanted %#+v", v.dc, actual, v.expected)
		}
		if alreadyCreated != true {
			t.Errorf("\nCalled CreateContainer(%#+v) \ngot %#+v \nwanted %#+v", v.dc, alreadyCreated, true)
		}
	}
}

func TestContainerStart(t *testing.T) {
	tt := []struct {
		dc         *DockerContainer
		expectedErr error
	}{
		{&DockerContainer{ContainerID: "myContainerId"}, nil},
	}
	var ctx = context.Background()
	cli := &MockDockerClient{}
	for _, v := range tt {
		err := ContainerStart(ctx, cli, v.dc)
		if err != nil {
			t.Errorf("\nCalled ContainerStart(%#+v) \ngot %s \nwanted %#+v", v.dc, err, v.expectedErr)
		}
	}
}

func TestContainerLogs(t *testing.T) {
	tt := []struct {
		dc         *DockerContainer
		expectedErr error
	}{
		{&DockerContainer{ContainerID: "myContainerId", FollowLogs: true}, nil},
		{&DockerContainer{ContainerID: "myContainerId", FollowLogs: false}, nil},
	}
	var ctx = context.Background()
	cli := &MockDockerClient{}
	for _, v := range tt {
		err := ContainerLogs(ctx, cli, v.dc)
		if err != nil {
			t.Errorf("\nCalled ContainerLogs(%#+v) \ngot %s \nwanted %#+v", v.dc, err, v.expectedErr)
		}
	}
}

func BenchmarkCreateContainer(b *testing.B) {
	var ctx = context.Background()
	cli := &MockDockerClient{}
	dc := &DockerContainer{
		ServiceConfig:     ServiceConfig{Image: "busybox", Restart: "unless-stopped"},
		ServiceName:       "myservice",
		NetworkName:       "myNetworkName",
		DockerComposeFile: "DockerFile",
		ContainerID:       "myExistingContainerId00912",
	}
	for n := 0; n < b.N; n++ {
		_, _, _ = CreateContainer(ctx, cli, dc)
	}
}

func BenchmarkContainerStart(b *testing.B) {
	var ctx = context.Background()
	cli := &MockDockerClient{}
	for n := 0; n < b.N; n++ {
		_ = ContainerStart(ctx, cli, &DockerContainer{ContainerID: "containerId"})
	}
}

func BenchmarkContainerLogs(b *testing.B) {
	var ctx = context.Background()
	cli := &MockDockerClient{}
	for n := 0; n < b.N; n++ {
		_ = ContainerLogs(ctx, cli, &DockerContainer{ContainerID: "containerId"})
	}
}
