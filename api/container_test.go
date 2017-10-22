package api

import (
	"context"
	"testing"
)

func TestCreateContainer(t *testing.T) {
	tt := []struct {
		xyz         *XYZ
		expected    string
		expectedErr error
	}{
		{
			&XYZ{
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
		//CreateContainer(ctx context.Context, cli MeliAPiClient, xyz *XYZ)
		alreadyCreated, actual, err := CreateContainer(ctx, cli, v.xyz)
		if err != nil {
			t.Errorf("\nCalled CreateContainer(%#+v) \ngot %s \nwanted %#+v", v.xyz, err, v.expectedErr)
		}
		if actual != v.expected {
			t.Errorf("\nCalled CreateContainer(%#+v) \ngot %s \nwanted %#+v", v.xyz, actual, v.expected)
		}
		if alreadyCreated != true {
			t.Errorf("\nCalled CreateContainer(%#+v) \ngot %#+v \nwanted %#+v", v.xyz, alreadyCreated, true)
		}
	}
}

func TestContainerStart(t *testing.T) {
	tt := []struct {
		xyz         *XYZ
		expectedErr error
	}{
		{&XYZ{ContainerID: "myContainerId"}, nil},
	}
	var ctx = context.Background()
	cli := &MockDockerClient{}
	for _, v := range tt {
		err := ContainerStart(ctx, cli, v.xyz)
		if err != nil {
			t.Errorf("\nCalled ContainerStart(%#+v) \ngot %s \nwanted %#+v", v.xyz, err, v.expectedErr)
		}
	}
}

func TestContainerLogs(t *testing.T) {
	tt := []struct {
		xyz         *XYZ
		expectedErr error
	}{
		{&XYZ{ContainerID: "myContainerId", FollowLogs: true}, nil},
		{&XYZ{ContainerID: "myContainerId", FollowLogs: false}, nil},
	}
	var ctx = context.Background()
	cli := &MockDockerClient{}
	for _, v := range tt {
		err := ContainerLogs(ctx, cli, v.xyz)
		if err != nil {
			t.Errorf("\nCalled ContainerLogs(%#+v) \ngot %s \nwanted %#+v", v.xyz, err, v.expectedErr)
		}
	}
}

func BenchmarkCreateContainer(b *testing.B) {
	var ctx = context.Background()
	cli := &MockDockerClient{}
	xyz := &XYZ{
		ServiceConfig:     ServiceConfig{Image: "busybox", Restart: "unless-stopped"},
		ServiceName:       "myservice",
		NetworkName:       "myNetworkName",
		DockerComposeFile: "DockerFile",
		ContainerID:       "myExistingContainerId00912",
	}
	for n := 0; n < b.N; n++ {
		_, _, _ = CreateContainer(ctx, cli, xyz)
	}
}

func BenchmarkContainerStart(b *testing.B) {
	var ctx = context.Background()
	cli := &MockDockerClient{}
	for n := 0; n < b.N; n++ {
		_ = ContainerStart(ctx, cli, &XYZ{ContainerID: "containerId"})
	}
}

func BenchmarkContainerLogs(b *testing.B) {
	var ctx = context.Background()
	cli := &MockDockerClient{}
	for n := 0; n < b.N; n++ {
		_ = ContainerLogs(ctx, cli, &XYZ{ContainerID: "containerId"})
	}
}
