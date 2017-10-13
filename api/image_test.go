package api

import (
	"context"
	"testing"
)

func TestGetPullDockerImage(t *testing.T) {
	tt := []struct {
		input       string
		expectedErr error
	}{
		{"busybox", nil},
	}
	var ctx = context.Background()
	cli := &MockDockerClient{}
	for _, v := range tt {
		err := PullDockerImage(ctx, v.input, cli)
		if err != nil {
			t.Errorf("\nran PullDockerImage(%#+v) \ngot %s \nwanted %#+v", v.input, err, v.expectedErr)
		}
	}
}

func TestGetBuildDockerImage(t *testing.T) {
	tt := []struct {
		input       string
		expected    string
		expectedErr error
	}{
		{"../testdata/Dockerfile", "meli_../testdata/dockerfile", nil},
	}
	var ctx = context.Background()
	cli := &MockDockerClient{}
	for _, v := range tt {
		actual, err := BuildDockerImage(ctx, v.input, cli)
		if err != nil {
			t.Errorf("\nran BuildDockerImage(%#+v) \ngot %s \nwanted %#+v", v.input, err, v.expectedErr)
		}
		if actual != v.expected {
			t.Errorf("\nran BuildDockerImage(%#+v) \ngot %s \nwanted %#+v", v.input, actual, v.expected)
		}
	}
}

func BenchmarkPullDockerImage(b *testing.B) {
	var ctx = context.Background()
	cli := &MockDockerClient{}
	for n := 0; n < b.N; n++ {
		_ = PullDockerImage(ctx, "busybox", cli)
	}
}

func BenchmarkBuildDockerImage(b *testing.B) {
	var ctx = context.Background()
	cli := &MockDockerClient{}
	for n := 0; n < b.N; n++ {
		_, _ = BuildDockerImage(ctx, "meli_../testdata/dockerfile", cli)
	}
}
