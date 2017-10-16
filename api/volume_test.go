package api

import (
	"context"
	"testing"
)

func TestCreateDockerVolume(t *testing.T) {
	tt := []struct {
		name        string
		driver      string
		expected    string
		expectedErr error
	}{
		{"MyVolumeName", "local", "MyVolume007", nil},
	}
	var ctx = context.Background()
	cli := &MockDockerClient{}
	for _, v := range tt {
		actual, err := CreateDockerVolume(ctx, v.name, v.driver, cli)
		if err != nil {
			t.Errorf("\nran CreateDockerVolume(%#+v) \ngot %s \nwanted %#+v", v.name, err, v.expectedErr)
		}
		if actual != v.expected {
			t.Errorf("\nran CreateDockerVolume(%#+v) \ngot %#+v \nwanted %#+v", v.name, actual, v.expected)
		}
	}
}

func BenchmarkCreateDockerVolume(b *testing.B) {
	var ctx = context.Background()
	cli := &MockDockerClient{}
	for n := 0; n < b.N; n++ {
		_, _ = CreateDockerVolume(ctx, "name", "local", cli)
	}
}
