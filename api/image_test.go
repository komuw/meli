package api

import (
	"context"
	"testing"
)

func TestGetPullDockerImage(t *testing.T) {
	tt := []struct {
		dc          *DockerContainer
		expectedErr error
	}{
		{&DockerContainer{ServiceConfig: ServiceConfig{Image: "busybox"}}, nil},
	}
	var ctx = context.Background()
	cli := &MockDockerClient{}
	for _, v := range tt {
		err := PullDockerImage(ctx, cli, v.dc)
		if err != nil {
			t.Errorf("\nCalled PullDockerImage(%#+v) \ngot %s \nwanted %#+v", v.dc, err, v.expectedErr)
		}
	}
}

func TestGetBuildDockerImage(t *testing.T) {
	tt := []struct {
		dc          *DockerContainer
		expected    string
		expectedErr error
	}{
		{
			&DockerContainer{
				ServiceName:       "myservicename",
				DockerComposeFile: "docker-compose.yml",
				ServiceConfig: ServiceConfig{
					Build: Buildstruct{Dockerfile: "../testdata/Dockerfile"}}},
			"meli_myservicename",
			nil},
	}

	var ctx = context.Background()
	cli := &MockDockerClient{}
	GetAuth()
	for _, v := range tt {
		actual, err := BuildDockerImage(ctx, cli, v.dc)
		if err != nil {
			t.Errorf("\nCalled BuildDockerImage(%#+v) \ngot %s \nwanted %#+v", v.dc, err, v.expectedErr)
		}
		if actual != v.expected {
			t.Errorf("\nCalled BuildDockerImage(%#+v) \ngot %s \nwanted %#+v", v.dc, actual, v.expected)
		}
	}
}

func BenchmarkPullDockerImage(b *testing.B) {
	var ctx = context.Background()
	cli := &MockDockerClient{}
	dc := &DockerContainer{ServiceConfig: ServiceConfig{Image: "busybox"}}
	GetAuth()
	for n := 0; n < b.N; n++ {
		_ = PullDockerImage(ctx, cli, dc)
	}
}

func BenchmarkBuildDockerImage(b *testing.B) {
	var ctx = context.Background()
	cli := &MockDockerClient{}
	dc := &DockerContainer{
		ServiceName: "myservicename",
		ServiceConfig: ServiceConfig{
			Build: Buildstruct{Dockerfile: "../testdata/Dockerfile"}}}
	for n := 0; n < b.N; n++ {
		_, _ = BuildDockerImage(ctx, cli, dc)
	}
}
