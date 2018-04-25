package meli

import (
	"bytes"
	"context"
	"io"
	"io/ioutil"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/network"
	volumetypes "github.com/docker/docker/api/types/volume"
)

type emptyStruct struct{}

// Buildstruct represents a docker-compose' build section
type Buildstruct struct {
	// remember to use caps so that they can be exported
	Context    string `yaml:"context,omitempty"`
	Dockerfile string `yaml:"dockerfile,omitempty"`
}

// ComposeService represents a docker-compose' service section
type ComposeService struct {
	Image       string      `yaml:"image,omitempty"`
	Ports       []string    `yaml:"ports,omitempty"`
	Labels      []string    `yaml:"labels,omitempty"`
	Environment []string    `yaml:"environment,omitempty"`
	Command     string      `yaml:"command,flow,omitempty"`
	Restart     string      `yaml:"restart,omitempty"`
	Build       Buildstruct `yaml:"build,omitempty"`
	Volumes     []string    `yaml:"volumes,omitempty"`
	Links       []string    `yaml:"links,omitempty"`
}

// DockerComposeConfig represents a docker-compose file
type DockerComposeConfig struct {
	Version  string                    `yaml:"version,omitempty"`
	Services map[string]ComposeService `yaml:"services"`
	Volumes  map[string]string         `yaml:"volumes,omitempty"`
}

// DockerContainer represents a docker container
type DockerContainer struct {
	ServiceName       string
	ComposeService    ComposeService
	NetworkID         string
	NetworkName       string
	FollowLogs        bool
	DockerComposeFile string
	ContainerID       string // this assumes that there can only be one container per docker-compose service
	LogMedium         io.Writer
	CurentDir         string
	Rebuild           bool
}

// UpdateContainerID updates a conatiners ID
func (dc *DockerContainer) UpdateContainerID(containerID string) {
	dc.ContainerID = containerID
}

// MeliAPiClient is meli's client to interact with the docker daemon server
type MeliAPiClient interface {
	// we implement this interface so that we can be able to mock it in tests
	// https://medium.com/@zach_4342/dependency-injection-in-golang-e587c69478a8
	ImagePull(ctx context.Context, ref string, options types.ImagePullOptions) (io.ReadCloser, error)
	ImageBuild(ctx context.Context, buildContext io.Reader, options types.ImageBuildOptions) (types.ImageBuildResponse, error)
	ContainerCreate(ctx context.Context, config *container.Config, hostConfig *container.HostConfig, networkingConfig *network.NetworkingConfig, containerName string) (container.ContainerCreateCreatedBody, error)
	ContainerStart(ctx context.Context, containerID string, options types.ContainerStartOptions) error
	ContainerLogs(ctx context.Context, container string, options types.ContainerLogsOptions) (io.ReadCloser, error)
	NetworkList(ctx context.Context, options types.NetworkListOptions) ([]types.NetworkResource, error)
	NetworkCreate(ctx context.Context, name string, options types.NetworkCreate) (types.NetworkCreateResponse, error)
	NetworkConnect(ctx context.Context, networkID, containerID string, config *network.EndpointSettings) error
	VolumeCreate(ctx context.Context, options volumetypes.VolumesCreateBody) (types.Volume, error)
	ContainerList(ctx context.Context, options types.ContainerListOptions) ([]types.Container, error)
	ContainerRemove(ctx context.Context, containerID string, options types.ContainerRemoveOptions) error
}

type imageProgress struct {
	Status         string `json:"status,omitempty"`
	Stream         string `json:"stream,omitempty"`
	Progress       string `json:"progress,omitempty"`
	ProgressDetail string `json:"progressDetail,omitempty"`
}

type mockDockerClient struct{}

func (m *mockDockerClient) ImagePull(ctx context.Context, refStr string, options types.ImagePullOptions) (io.ReadCloser, error) {
	return ioutil.NopCloser(bytes.NewBuffer([]byte("Pulling from library/testImage"))), nil
}
func (m *mockDockerClient) ImageBuild(ctx context.Context, buildContext io.Reader, options types.ImageBuildOptions) (types.ImageBuildResponse, error) {
	return types.ImageBuildResponse{Body: ioutil.NopCloser(bytes.NewBuffer([]byte("BUILT library/testImage"))), OSType: "linux baby!"}, nil
}
func (m *mockDockerClient) ContainerCreate(ctx context.Context, config *container.Config, hostConfig *container.HostConfig, networkingConfig *network.NetworkingConfig, containerName string) (container.ContainerCreateCreatedBody, error) {
	return container.ContainerCreateCreatedBody{ID: "myContainerId001"}, nil
}

func (m *mockDockerClient) ContainerStart(ctx context.Context, containerID string, options types.ContainerStartOptions) error {
	return nil
}

func (m *mockDockerClient) ContainerLogs(ctx context.Context, container string, options types.ContainerLogsOptions) (io.ReadCloser, error) {
	return ioutil.NopCloser(bytes.NewBuffer([]byte("SHOWING LOGS for library/testImage"))), nil
}

func (m *mockDockerClient) NetworkList(ctx context.Context, options types.NetworkListOptions) ([]types.NetworkResource, error) {
	return []types.NetworkResource{}, nil
}
func (m *mockDockerClient) NetworkCreate(ctx context.Context, name string, options types.NetworkCreate) (types.NetworkCreateResponse, error) {
	return types.NetworkCreateResponse{ID: "myNetworkId002"}, nil
}

func (m *mockDockerClient) NetworkConnect(ctx context.Context, networkID, containerID string, config *network.EndpointSettings) error {
	return nil
}

func (m *mockDockerClient) VolumeCreate(ctx context.Context, options volumetypes.VolumesCreateBody) (types.Volume, error) {
	return types.Volume{Name: "MyVolume007"}, nil
}

func (m *mockDockerClient) ContainerList(ctx context.Context, options types.ContainerListOptions) ([]types.Container, error) {
	return []types.Container{types.Container{ID: "myExistingContainerId00912"}}, nil
}
func (m *mockDockerClient) ContainerRemove(ctx context.Context, containerID string, options types.ContainerRemoveOptions) error {
	return nil
}
