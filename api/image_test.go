package api

import (
	"bytes"
	"context"
	"io"
	"io/ioutil"
	"testing"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
)

func BenchmarkPullDockerImage(b *testing.B) {
	var ctx = context.Background()
	var cli, _ = client.NewEnvClient()
	defer cli.Close()
	for n := 0; n < b.N; n++ {
		err := PullDockerImage(ctx, "busybox", cli)
		b.Log(err)
	}
}

// https://medium.com/@zach_4342/dependency-injection-in-golang-e587c69478a8
type MockDockerClient struct{}

func (m *MockDockerClient) ImagePull(ctx context.Context, refStr string, options types.ImagePullOptions) (io.ReadCloser, error) {
	return ioutil.NopCloser(bytes.NewBuffer([]byte("Pulling from library/testImage"))), nil
}
func TestGetPullDockerImage(t *testing.T) {
	tt := []struct {
		input       string
		expectedErr error
	}{
		{"busybox", nil},
	}
	var ctx = context.Background()
	cli := &MockDockerClient{}
	for _, v := range tt {
		err := PullDockerImage(ctx, v.input, cli)
		if err != nil {
			t.Errorf("\nran PullDockerImage(%#+v) \ngot %s \nwanted %#+v", v.input, err, v.expectedErr)
		}
	}
}
