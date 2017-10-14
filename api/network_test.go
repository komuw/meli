package api

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
	cli := &MockDockerClient{}
	for _, v := range tt {
		actual, err := GetNetwork(ctx, v.input, cli)
		if err != nil {
			t.Errorf("\nran GetNetwork(%#+v) \ngot %s \nwanted %#+v", v.input, err, v.expectedErr)
		}
		if actual != v.expected {
			t.Errorf("\nran GetNetwork(%#+v) \ngot %#+v \nwanted %#+v", v.input, actual, v.expected)
		}
	}
}

func TestConnectNetwork(t *testing.T) {
	tt := []struct {
		netWorkID   string
		containerID string
		expectedErr error
	}{
		{"myNetWorkID", "myContainerID003", nil},
	}
	var ctx = context.Background()
	cli := &MockDockerClient{}
	for _, v := range tt {
		err := ConnectNetwork(ctx, v.netWorkID, v.containerID, cli)
		if err != nil {
			t.Errorf("\nran ConnectNetwork(%#+v) \ngot %s \nwanted %#+v", v.netWorkID, err, v.expectedErr)
		}
	}
}

func BenchmarkGetNetwork(b *testing.B) {
	var ctx = context.Background()
	cli := &MockDockerClient{}
	for n := 0; n < b.N; n++ {
		_, _ = GetNetwork(ctx, "myNetWorkName", cli)
	}
}
func BenchmarkConnectNetwork(b *testing.B) {
	var ctx = context.Background()
	cli := &MockDockerClient{}
	for n := 0; n < b.N; n++ {
		_ = ConnectNetwork(ctx, "netWorkID", "containerID", cli)
	}
}
