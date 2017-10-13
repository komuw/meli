package api

import (
	"context"
	"testing"

	"github.com/docker/docker/client"
)

var ctx = context.Background()
var cli, err = client.NewEnvClient()

func BenchmarkPullDockerImage(b *testing.B) {

	for n := 0; n < b.N; n++ {
		err := PullDockerImage(ctx, "busybox", cli)
		b.Log(err)
	}
}
