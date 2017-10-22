package api

import (
	"context"
	"testing"
)

func TestGetPullDockerImage(t *testing.T) {
	tt := []struct {
		dc         *DockerContainer
		expectedErr error
	}{
		{&DockerContainer{ServiceConfig: ServiceConfig{Image: "busybox"}}, nil}}
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
		dc         *DockerContainer
		expected    string
		expectedErr error
	}{&DockerContainer{
		ServiceName: "myservicename",
		ServiceConfig: ServiceConfig{
			Build: Buildstruct{
				Dockerfile: "../testdata/Dockerfile"}}}, "meli_myservicename", nil}

	var ctx = context.Background()
	cli := &MockDockerClient{}
	for _, v := range tt {
		actual, err := BuildDockerImage(ctx, v.serviceName, v.dockerFile, cli)
		if err != nil {
			t.Errorf("\nCalled BuildDockerImage(%#+v) \ngot %s \nwanted %#+v", v.dockerFile, err, v.expectedErr)
		}
		if actual != v.expected {
			t.Errorf("\nCalled BuildDockerImage(%#+v) \ngot %s \nwanted %#+v", v.dockerFile, actual, v.expected)
		}
	}
}

func BenchmarkPullDockerImage(b *testing.B) {
	var ctx = context.Background()
	cli := &MockDockerClient{}
	GetAuth()
	for n := 0; n < b.N; n++ {
		_ = PullDockerImage(ctx, "busybox", cli)
	}
}

func BenchmarkBuildDockerImage(b *testing.B) {
	var ctx = context.Background()
	cli := &MockDockerClient{}
	for n := 0; n < b.N; n++ {
		_, _ = BuildDockerImage(ctx, "meliserice", "meli_../testdata/dockerfile", cli)
	}
}
