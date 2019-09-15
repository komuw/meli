package meli

import (
	"context"
	"fmt"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/network"
)

// GetNetwork gets or creates newtwork(if it doesn't exist yet.)
func GetNetwork(ctx context.Context, networkName string, cli APIclient) (string, error) {
	// return early if network exists
	netList, err := cli.NetworkList(ctx, types.NetworkListOptions{})
	if err != nil {
		return "", fmt.Errorf("unable to list docker networks: %w", err)

	}
	for _, v := range netList {
		if v.Name == networkName {
			return v.ID, nil
		}
	}

	var typeNetworkCreate = types.NetworkCreate{
		CheckDuplicate: true,
		Driver:         "bridge",
		EnableIPv6:     false,
		IPAM:           &network.IPAM{Driver: "default"},
		Internal:       false,
		Attachable:     true,
	}
	networkCreateResponse, err := cli.NetworkCreate(
		ctx,
		networkName,
		typeNetworkCreate)
	if err != nil {
		return "", fmt.Errorf("unable to create docker network %v: %w", networkName, err)

	}
	return networkCreateResponse.ID, nil

}

// ConnectNetwork connects a container to an existent docker network.
func ConnectNetwork(ctx context.Context, cli APIclient, dc *DockerContainer) error {
	err := cli.NetworkConnect(
		ctx,
		dc.NetworkID,
		dc.ContainerID,
		&network.EndpointSettings{})
	if err != nil {
		return fmt.Errorf("unable to connect container %s of service %v to network %s: %w", dc.ContainerID, dc.ServiceName, dc.NetworkID, err)

	}
	return nil
}
