package main

import (
	"context"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/docker/client"
	"github.com/pkg/errors"
)

func getNetwork() (string, error) {
	// create/get newtwork
	// TODO: move this into it's own file
	networkName := "melinetworkname"
	ctx := context.Background()
	cli, err := client.NewEnvClient()
	if err != nil {
		return "", errors.Wrap(err, "unable to intialize docker client")
	}

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
		return "", errors.Wrap(err, "unable to create docker network")
	}
	return networkCreateResponse.ID, nil

}
