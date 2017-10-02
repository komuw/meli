package main

import (
	"context"
	"strings"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/docker/client"

	"errors"
)

func getNetwork(networkName string) (string, error) {
	// create/get newtwork
	ctx := context.Background()
	cli, err := client.NewEnvClient()
	if err != nil {
		return "", &popagateError{originalErr: err, newErr: errors.New("unable to intialize docker client")}
	}

	// return early if network exists
	netList, err := cli.NetworkList(ctx, types.NetworkListOptions{})
	if err != nil {
		return "", &popagateError{originalErr: err, newErr: errors.New("unable to intialize docker client")}

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
		return "", &popagateError{originalErr: err, newErr: errors.New("unable to create docker network")}
	}
	return networkCreateResponse.ID, nil

}

func getCwdName(path string) string {
	//TODO: investigate if this will work cross platform
	// it might be unable to handle paths in windows OS
	f := func(c rune) bool {
		if c == 47 {
			// 47 is the '/' character
			return true
		}
		return false
	}
	pathSlice := strings.FieldsFunc(path, f)
	return pathSlice[len(pathSlice)-1]
}
