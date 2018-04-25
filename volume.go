package meli

import (
	"context"
	"errors"
	"fmt"
	"io"

	"github.com/docker/docker/api/types/volume"
)

// CreateDockerVolume creates a docker volume
func CreateDockerVolume(ctx context.Context, cli MeliAPiClient, name, driver string, dst io.Writer) (string, error) {
	volume, err := cli.VolumeCreate(
		ctx,
		volume.VolumesCreateBody{
			Driver: driver,
			Name:   name})
	if err != nil {
		return "", &popagateError{originalErr: err, newErr: errors.New(" :unable to create docker volume")}
	}
	fmt.Fprintf(dst, "\ndocker volume: %s created succesfully.\n", volume.Name)

	return volume.Name, nil
}
