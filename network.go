package meli

import (
	"context"
	"fmt"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/network"

	"errors"
)

// GetNetwork gets or creates newtwork(if it doesn't exist yet.)
func GetNetwork(ctx context.Context, networkName string, cli APIclient) (string, error) {
	// return early if network exists
	netList, err := cli.NetworkList(ctx, types.NetworkListOptions{})
	if err != nil {
		return "", &popagateError{originalErr: err, newErr: errors.New(" :unable to intialize docker client")}

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
		return "", &popagateError{originalErr: err, newErr: errors.New(" :unable to create docker network")}
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
		return &popagateError{
			originalErr: err,
			newErr:      fmt.Errorf(" :unable to connect container %s to network %s", dc.ContainerID, dc.NetworkID)}

	}
	return nil
}
