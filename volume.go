package main

import (
	"context"
	"log"

	"github.com/docker/docker/api/types/volume"
	"github.com/docker/docker/client"
)

func CreateDockerVolume(ctx context.Context, name, driver string) {
	cli, err := client.NewEnvClient()
	if err != nil {
		log.Println(err, " :unable to intialize docker client")
	}
	defer cli.Close()

	volume, err := cli.VolumeCreate(
		ctx,
		volume.VolumesCreateBody{
			Driver: driver,
			Name:   name})
	if err != nil {
		log.Println(err, " :unable to create docker volume")
	}

	log.Printf("\ndocker volume: %s created succesfully.\n", volume.Name)
}
