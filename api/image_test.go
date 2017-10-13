package api

import (
	"context"
	"testing"

	"github.com/docker/docker/client"
)

func BenchmarkPullDockerImage(b *testing.B) {
	var ctx = context.Background()
	var cli, _ = client.NewEnvClient()
	defer cli.Close()
	for n := 0; n < b.N; n++ {
		err := PullDockerImage(ctx, "busybox", cli)
		b.Log(err)
	}
}

func TestGetPullDockerImage(t *testing.T) {
	tt := []struct {
		input       string
		expectedErr error
	}{
		{"busybox", nil},
	}
	var ctx = context.Background()
	var cli, _ = client.NewEnvClient()
	defer cli.Close()

	for _, v := range tt {
		err := PullDockerImage(ctx, v.input, cli)
		if err != nil {
			t.Errorf("\nran PullDockerImage(%#+v) \ngot %s \nwanted %#+v", v.input, err, v.expectedErr)
		}
	}
}
