package api

import (
	"context"
	"io/ioutil"
	"testing"
)

func TestCreateContainer(t *testing.T) {
	tt := []struct {
		dc                      *DockerContainer
		expected                string
		containerAlreadycreated bool
		expectedErr             error
	}{
		{
			&DockerContainer{
				ComposeService:    ComposeService{Image: "busybox", Restart: "unless-stopped"},
				ServiceName:       "myservice",
				NetworkName:       "myNetworkName",
				DockerComposeFile: "DockerFile",
				ContainerID:       "myExistingContainerId00912",
			},
			"myExistingContainerId00912",
			true,
			nil},
		{
			&DockerContainer{
				ComposeService:    ComposeService{Image: "busybox", Restart: "unless-stopped"},
				ServiceName:       "myservice",
				NetworkName:       "myNetworkName",
				DockerComposeFile: "DockerFile",
				ContainerID:       "myContainerId001",
				Rebuild:           true,
			},
			"myContainerId001",
			false,
			nil},
	}
	var ctx = context.Background()
	cli := &MockDockerClient{}
	for _, v := range tt {
		alreadyCreated, actual, err := CreateContainer(ctx, cli, v.dc)
		if err != nil {
			t.Errorf("\nCalled CreateContainer(%#+v) \ngot %s \nwanted %#+v", v.dc, err, v.expectedErr)
		}
		if actual != v.expected {
			t.Errorf("\nCalled CreateContainer(%#+v) \ngot %s \nwanted %#+v", v.dc, actual, v.expected)
		}
		if alreadyCreated != v.containerAlreadycreated {
			t.Errorf("\nCalled CreateContainer(%#+v) \ngot %#+v \nwanted %#+v", v.dc, alreadyCreated, true)
		}
	}
}

func TestContainerStart(t *testing.T) {
	tt := []struct {
		dc          *DockerContainer
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
		dc          *DockerContainer
		expectedErr error
	}{
		{&DockerContainer{ContainerID: "myContainerId", FollowLogs: true, LogMedium: ioutil.Discard}, nil},
		{&DockerContainer{ContainerID: "myContainerId", FollowLogs: false, LogMedium: ioutil.Discard}, nil},
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
		ComposeService:    ComposeService{Image: "busybox", Restart: "unless-stopped"},
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
		_ = ContainerLogs(ctx, cli, &DockerContainer{ContainerID: "containerId", LogMedium: ioutil.Discard})
	}
}
