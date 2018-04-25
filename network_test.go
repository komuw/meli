package meli

import (
	"context"
	"testing"
)

func TestGetNetwork(t *testing.T) {
	tt := []struct {
		input       string
		expected    string
		expectedErr error
	}{
		{"myNetWorkName", "myNetworkId002", nil},
	}
	var ctx = context.Background()
	cli := &mockDockerClient{}
	for _, v := range tt {
		actual, err := GetNetwork(ctx, v.input, cli)
		if err != nil {
			t.Errorf("\nCalled GetNetwork(%#+v) \ngot %s \nwanted %#+v", v.input, err, v.expectedErr)
		}
		if actual != v.expected {
			t.Errorf("\nCalled GetNetwork(%#+v) \ngot %#+v \nwanted %#+v", v.input, actual, v.expected)
		}
	}
}

func TestConnectNetwork(t *testing.T) {
	tt := []struct {
		dc          *DockerContainer
		expectedErr error
	}{
		{
			&DockerContainer{
				NetworkID:   "myNetWorkID",
				ContainerID: "myContainerID003"},
			nil},
	}
	var ctx = context.Background()
	cli := &mockDockerClient{}

	for _, v := range tt {
		err := ConnectNetwork(ctx, cli, v.dc)
		if err != nil {
			t.Errorf("\nCalled ConnectNetwork(%#+v) \ngot %s \nwanted %#+v", v.dc, err, v.expectedErr)
		}
	}
}

func BenchmarkGetNetwork(b *testing.B) {
	var ctx = context.Background()
	cli := &mockDockerClient{}
	for n := 0; n < b.N; n++ {
		_, _ = GetNetwork(ctx, "myNetWorkName", cli)
	}
}

func BenchmarkConnectNetwork(b *testing.B) {
	var ctx = context.Background()
	cli := &mockDockerClient{}
	dc := &DockerContainer{NetworkID: "myNetWorkID", ContainerID: "myContainerID003"}
	for n := 0; n < b.N; n++ {
		_ = ConnectNetwork(ctx, cli, dc)
	}
}
