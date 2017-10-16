package api

import (
	"context"
	"errors"
	"log"

	"github.com/docker/docker/api/types/volume"
)

func CreateDockerVolume(ctx context.Context, name, driver string, cli MeliAPiClient) (string, error) {
	volume, err := cli.VolumeCreate(
		ctx,
		volume.VolumesCreateBody{
			Driver: driver,
			Name:   name})
	if err != nil {
		return "", &popagateError{originalErr: err, newErr: errors.New(" :unable to create docker volume")}
	}

	log.Printf("\ndocker volume: %s created succesfully.\n", volume.Name)
	return volume.Name, nil
}
