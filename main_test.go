package main

import (
	"context"
	"io/ioutil"
	"log"
	"os"
	"sync"
	"testing"

	"github.com/docker/docker/client"
	"github.com/komuw/meli/api"
)

func BenchmarkstartComposeServices(b *testing.B) {
	var wg sync.WaitGroup
	var ctx = context.Background()
	cli, err := client.NewEnvClient()
	if err != nil {
		log.Fatal(err, " :unable to intialize docker client")
	}
	defer cli.Close()
	curentDir, _ := os.Getwd()
	networkName := "meli_network_" + api.GetCwdName(curentDir)
	networkID, _ := api.GetNetwork(ctx, networkName, cli)

	dc := &api.DockerContainer{
		ServiceName:       "myservicename",
		LogMedium:         ioutil.Discard,
		NetworkID:         networkID,
		DockerComposeFile: "/testdata/docker-compose.yml",
		ComposeService: api.ComposeService{
			Build: api.Buildstruct{Dockerfile: "Dockerfile"}},
	}
	api.GetAuth()

	for n := 0; n < b.N; n++ {
		wg.Add(1)
		startComposeServices(ctx, cli, &wg, dc)
	}
}
