package api

import (
	"os"
	"reflect"
	"testing"
)

func TestFormatContainerName(t *testing.T) {
	tt := []struct {
		input    string
		expected string
	}{
		{"redis", "meli_redis."},
		{"nats:", "meli_nats."},
		{"yolo:ala", "meli_yolo."},
	}
	for _, v := range tt {
		actual := FormatContainerName(v.input, ".")
		if actual != v.expected {
			t.Errorf("\nCalled FormatContainerName(%#+v) \ngot %#+v \nwanted %#+v", v.input, actual, v.expected)
		}
	}
}

func TestFormatLabels(t *testing.T) {
	tt := []struct {
		input    string
		expected []string
	}{
		{"traefik.backend=web", []string{"traefik.backend", "web"}},
		{"env:prod", []string{"env", "prod"}},
	}
	for _, v := range tt {
		actual := FormatLabels(v.input)
		if !reflect.DeepEqual(actual, v.expected) {
			t.Errorf("\nCalled FormatLabels(%#+v) \ngot %#+v \nwanted %#+v", v.input, actual, v.expected)
		}
	}
}

func TestFormatPorts(t *testing.T) {
	tt := []struct {
		input    string
		expected []string
	}{
		{"6300:6379", []string{"6300", "6379"}},
	}
	for _, v := range tt {
		actual := FormatPorts(v.input)
		if !reflect.DeepEqual(actual, v.expected) {
			t.Errorf("\nCalled TestFormatPorts(%#+v) \ngot %#+v \nwanted %#+v", v.input, actual, v.expected)
		}
	}
}

func TestFormatServiceVolumes(t *testing.T) {
	currentDir, _ := os.Getwd()
	tt := []struct {
		volume            string
		dockerComposeFile string
		expected          []string
	}{
		{"data-volume:/home", "composefile", []string{"data-volume", "/home"}},
		{"./:/mydir", "composefile", []string{currentDir, "/mydir"}},
		{"/var/run/docker.sock:/var/run/docker.sock", "composefile", []string{"/var/run/docker.sock", "/var/run/docker.sock"}},
		{".startWithDot:/home/.startWithDot", "composefile", []string{currentDir + "/.startWithDot", "/home/.startWithDot"}},
	}
	for _, v := range tt {
		actual := FormatServiceVolumes(v.volume, v.dockerComposeFile)
		if !reflect.DeepEqual(actual, v.expected) {
			t.Errorf("\nCalled FormatServiceVolumes(%#+v) \ngot %#+v \nwanted %#+v", v.volume, actual, v.expected)
		}
	}
}

func TestFormatRegistryAuth(t *testing.T) {
	tt := []struct {
		input    string
		expected []string
	}{
		{"myUsername:myPassword001", []string{"myUsername", "myPassword001"}},
	}
	for _, v := range tt {
		actual := FormatRegistryAuth(v.input)
		if !reflect.DeepEqual(actual, v.expected) {
			t.Errorf("\nCalled FormatRegistryAuth(%#+v) \ngot %#+v \nwanted %#+v", v.input, actual, v.expected)
		}
	}
}

func TestFormatComposePath(t *testing.T) {
	tt := []struct {
		input    string
		expected []string
	}{
		{"testdata/dockerFile", []string{"testdata", "dockerFile"}},
	}
	for _, v := range tt {
		actual := FormatComposePath(v.input)
		if !reflect.DeepEqual(actual, v.expected) {
			t.Errorf("\nCalled FormatComposePath(%#+v) \ngot %#+v \nwanted %#+v", v.input, actual, v.expected)
		}
	}
}

func BenchmarkFormatLabels(b *testing.B) {
	// run the FormatLabels function b.N times
	for n := 0; n < b.N; n++ {
		_ = FormatLabels("traefik.backend=web")
	}
}

func BenchmarkFormatPorts(b *testing.B) {
	for n := 0; n < b.N; n++ {
		_ = FormatPorts("6300:6379")
	}
}

func BenchmarkFormatServiceVolumes(b *testing.B) {
	for n := 0; n < b.N; n++ {
		_ = FormatServiceVolumes("data-volume:/home", "composeFile")
	}
}

func BenchmarkFormatContainerName(b *testing.B) {
	for n := 0; n < b.N; n++ {
		_ = FormatContainerName("build_with_no_specified_dockerfile", ".")
	}

}
