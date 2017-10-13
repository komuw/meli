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
