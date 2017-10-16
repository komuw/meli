package api

import (
	"context"
	"fmt"
	"strings"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/network"

	"errors"
)

// GetNetwork gets or creates newtwork(if it doesn't exist yet.)
func GetNetwork(ctx context.Context, networkName string, cli MeliAPiClient) (string, error) {
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

func ConnectNetwork(ctx context.Context, networkID, containerID string, cli MeliAPiClient) error {
	err := cli.NetworkConnect(
		ctx,
		networkID,
		containerID,
		&network.EndpointSettings{})
	if err != nil {
		return &popagateError{
			originalErr: err,
			newErr:      fmt.Errorf(" :unable to connect container %s to network %s", containerID, networkID)}

	}
	return nil
}

func GetCwdName(path string) string {
	//TODO: investigate if this will work cross platform
	// it might be  :unable to handle paths in windows OS
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
