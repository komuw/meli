package api

import (
	"bytes"
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
	dst := bytes.NewBuffer(make([]byte, 0, 0))
	for _, v := range tt {
		actual, err := CreateDockerVolume(ctx, cli, v.name, v.driver, dst)
		if err != nil {
			t.Errorf("\nCalled CreateDockerVolume(%#+v) \ngot %s \nwanted %#+v", v.name, err, v.expectedErr)
		}
		if actual != v.expected {
			t.Errorf("\nCalled CreateDockerVolume(%#+v) \ngot %#+v \nwanted %#+v", v.name, actual, v.expected)
		}
	}
}

func BenchmarkCreateDockerVolume(b *testing.B) {
	var ctx = context.Background()
	cli := &MockDockerClient{}
	dst := bytes.NewBuffer(make([]byte, 0, 0))
	for n := 0; n < b.N; n++ {
		_, _ = CreateDockerVolume(ctx, cli, "name", "local", dst)
	}
}
