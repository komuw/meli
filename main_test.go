package main

import (
	"context"
	"io/ioutil"
	"log"
	"sync"
	"testing"

	"github.com/docker/docker/client"
	"github.com/komuw/meli/api"
)

func BenchmarkStartContainers(b *testing.B) {
	var wg sync.WaitGroup
	var ctx = context.Background()
	cli, err := client.NewEnvClient()
	if err != nil {
		log.Fatal(err, " :unable to intialize docker client")
	}
	defer cli.Close()

	dc := &api.DockerContainer{
		ServiceName:       "myservicename",
		LogMedium:         ioutil.Discard,
		DockerComposeFile: "/testdata/docker-compose.yml",
		ComposeService: api.ComposeService{
			Build: api.Buildstruct{Dockerfile: "Dockerfile"}},
	}
	for n := 0; n < b.N; n++ {
		wg.Add(1)
		startContainers(ctx, cli, &wg, dc)
	}
}
