package meli

import (
	"context"
	"io/ioutil"
	"strings"
	"testing"
)

func TestPullDockerImage(t *testing.T) {
	tt := []struct {
		dc          *DockerContainer
		expectedErr error
	}{
		{&DockerContainer{ComposeService: ComposeService{Image: "busybox"}, LogMedium: ioutil.Discard}, nil},
	}
	var ctx = context.Background()
	cli := &mockDockerClient{}
	for _, v := range tt {
		err := PullDockerImage(ctx, cli, v.dc)
		if err != nil {
			t.Errorf("\nCalled PullDockerImage(%#+v) \ngot %s \nwanted %#+v", v.dc, err, v.expectedErr)
		}
	}
}

func TestBuildDockerImage(t *testing.T) {
	tt := []struct {
		dc          *DockerContainer
		expected    string
		expectedErr error
	}{
		{
			&DockerContainer{
				ServiceName:       "myservicename",
				DockerComposeFile: "docker-compose.yml",
				ComposeService: ComposeService{
					Build: Buildstruct{Dockerfile: "testdata/Dockerfile"}},
				LogMedium: ioutil.Discard},
			"meli_myservicename",
			nil},
		{
			&DockerContainer{
				ServiceName:       "myservicename",
				DockerComposeFile: "docker-compose.yml",
				ComposeService: ComposeService{
					Build: Buildstruct{Dockerfile: "testdata/Dockerfile"}},
				LogMedium: ioutil.Discard,
				Rebuild:   true,
			},
			"meli_myservicename",
			nil},
	}

	var ctx = context.Background()
	cli := &mockDockerClient{}
	LoadAuth()
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
	cli := &mockDockerClient{}
	dc := &DockerContainer{ComposeService: ComposeService{Image: "busybox"}, LogMedium: ioutil.Discard}
	LoadAuth()
	for n := 0; n < b.N; n++ {
		_ = PullDockerImage(ctx, cli, dc)
	}
}

func BenchmarkBuildDockerImage(b *testing.B) {
	var ctx = context.Background()
	cli := &mockDockerClient{}
	dc := &DockerContainer{
		ServiceName:       "myservicename",
		DockerComposeFile: "docker-compose.yml",
		ComposeService: ComposeService{
			Build: Buildstruct{Dockerfile: "testdata/Dockerfile"}},
		LogMedium: ioutil.Discard,
		Rebuild:   true,
	}
	LoadAuth()
	for n := 0; n < b.N; n++ {
		_, _ = BuildDockerImage(ctx, cli, dc)
	}
}

func BenchmarkPoolReadFrom(b *testing.B) {
	r := strings.NewReader("hello")
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		_, _ = poolReadFrom(r)
	}
}
