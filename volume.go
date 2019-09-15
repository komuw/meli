package meli

import (
	"context"
	"fmt"
	"io"

	"github.com/docker/docker/api/types/volume"
)

// CreateDockerVolume creates a docker volume
func CreateDockerVolume(ctx context.Context, cli APIclient, name, driver string, dst io.Writer) (string, error) {
	volume, err := cli.VolumeCreate(
		ctx,
		volume.VolumeCreateBody{
			Driver: driver,
			Name:   name})
	if err != nil {
		return "", fmt.Errorf("unable to create docker volume %v: %w", name, err)

	}
	fmt.Fprintf(dst, "\ndocker volume: %s created successfully.\n", volume.Name)

	return volume.Name, nil
}
