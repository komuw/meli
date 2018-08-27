package meli

import (
	"context"
	"io/ioutil"
	"testing"
)

func TestPullDockerImage(t *testing.T) {
	tt := []struct {
		dc          *DockerContainer
		expectedErr error
	}{
		{&DockerContainer{ComposeService: ComposeService{Image: "busybox"}, LogMedium: ioutil.Discard}, nil},
	}
	var ctx = context.Background()
	cli := &mockDockerClient{}
	for _, v := range tt {
		err := PullDockerImage(ctx, cli, v.dc)
		if err != nil {
			t.Errorf("\nCalled PullDockerImage(%#+v) \ngot %s \nwanted %#+v", v.dc, err, v.expectedErr)
		}
	}
}

func TestBuildDockerImage(t *testing.T) {
	tt := []struct {
		dc          *DockerContainer
		expected    string
		expectedErr error
	}{
		{
			&DockerContainer{
				ServiceName:       "myservicename",
				DockerComposeFile: "docker-compose.yml",
				ComposeService: ComposeService{
					Build: Buildstruct{Dockerfile: "testdata/Dockerfile"}},
				LogMedium: ioutil.Discard},
			"meli_myservicename",
			nil},
		{
			&DockerContainer{
				ServiceName:       "myservicename",
				DockerComposeFile: "docker-compose.yml",
				ComposeService: ComposeService{
					Build: Buildstruct{Dockerfile: "testdata/Dockerfile"}},
				LogMedium: ioutil.Discard,
				Rebuild:   true,
			},
			"meli_myservicename",
			nil},
	}

	var ctx = context.Background()
	cli := &mockDockerClient{}
	LoadAuth()
	for _, v := range tt {
		actual, err := BuildDockerImage(ctx, cli, v.dc)
		if err != nil {
			t.Errorf("\nCalled BuildDockerImage(%#+v) \ngot %s \nwanted %#+v", v.dc, err, v.expectedErr)
		}
		if actual != v.expected {
			t.Errorf("\nCalled BuildDockerImage(%#+v) \ngot %s \nwanted %#+v", v.dc, actual, v.expected)
		}
	}
}

func BenchmarkPullDockerImage(b *testing.B) {
	var ctx = context.Background()
	cli := &mockDockerClient{}
	dc := &DockerContainer{ComposeService: ComposeService{Image: "busybox"}, LogMedium: ioutil.Discard}
	LoadAuth()
	for n := 0; n < b.N; n++ {
		_ = PullDockerImage(ctx, cli, dc)
	}
}

/*
If you run this benchmark:
  go test -memprofile mem.prof -bench=BenchmarkBuildDockerImage -run=XXX .
then
  pprof -sample_index=alloc_space meli.test mem.prof
and
  list walkFnClosure

 3.35GB     91:		readFile, err := ioutil.ReadAll(f)
         .          .     92:		if err != nil {
         .          .     93:			return err
         .          .     94:		}
         .     4.31GB     95:		_, err = tw.Write(readFile)
         .          .     96:		if err != nil {
         .          .     97:			return err
         .          .     98:		}
         .          .     99:		return nil
		 .          .    100:	}

We need to reduce the allocations in ioutil.ReadAll(f) and tw.Write(readFile)
*/
/*
with a newer implemenation that uses sync.Pool
.          .     93:
         .          .     94:		tr := io.TeeReader(f, tw)
         .   119.48MB     95:		_, err = poolReadFrom(tr)
         .          .     96:		if err != nil {
         .          .     97:			return err
         .          .     98:		}

*/
func BenchmarkBuildDockerImage(b *testing.B) {
	var ctx = context.Background()
	cli := &mockDockerClient{}
	dc := &DockerContainer{
		ServiceName:       "myservicename",
		DockerComposeFile: "docker-compose.yml",
		ComposeService: ComposeService{
			Build: Buildstruct{Dockerfile: "testdata/Dockerfile"}},
		LogMedium: ioutil.Discard,
		Rebuild:   true,
	}
	LoadAuth()
	for n := 0; n < b.N; n++ {
		_, _ = BuildDockerImage(ctx, cli, dc)
	}
}
