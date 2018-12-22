package meli

import (
	"context"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/network"
	"github.com/pkg/errors"
)

// GetNetwork gets or creates newtwork(if it doesn't exist yet.)
func GetNetwork(ctx context.Context, networkName string, cli APIclient) (string, error) {
	// return early if network exists
	netList, err := cli.NetworkList(ctx, types.NetworkListOptions{})
	if err != nil {
		return "", errors.Wrap(err, "unable to intialize docker client")

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
		return "", errors.Wrapf(err, "unable to create docker network %v", networkName)
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
		return errors.Wrapf(err, "unable to connect container %s of service %v to network %s", dc.ContainerID, dc.ServiceName, dc.NetworkID)
	}
	return nil
}
