package api

import (
	"bytes"
	"context"
	"io"
	"io/ioutil"
	"testing"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/network"
)

type xockDockerClient struct{}

func (m *xockDockerClient) ImagePull(ctx context.Context, refStr string, options types.ImagePullOptions) (io.ReadCloser, error) {
	return ioutil.NopCloser(bytes.NewBuffer([]byte("Pulling from library/testImage"))), nil
}
func (m *xockDockerClient) ImageBuild(ctx context.Context, buildContext io.Reader, options types.ImageBuildOptions) (types.ImageBuildResponse, error) {
	return types.ImageBuildResponse{Body: ioutil.NopCloser(bytes.NewBuffer([]byte("BUILT library/testImage"))), OSType: "linux baby!"}, nil
}
func (m *xockDockerClient) ContainerCreate(ctx context.Context, config *container.Config, hostConfig *container.HostConfig, networkingConfig *network.NetworkingConfig, containerName string) (container.ContainerCreateCreatedBody, error) {
	return container.ContainerCreateCreatedBody{ID: "myContainerId001"}, nil
}

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
	cli := &xockDockerClient{}
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

func BenchmarkCreateContainer(b *testing.B) {
	var ctx = context.Background()
	cli := &xockDockerClient{}
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
