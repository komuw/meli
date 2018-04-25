package api

import (
	"os"
	"reflect"
	"testing"
)

func TestformatContainerName(t *testing.T) {
	tt := []struct {
		input    string
		expected string
	}{
		{"redis", "meli_redis."},
		{"nats:", "meli_nats."},
		{"yolo:ala", "meli_yolo."},
	}
	for _, v := range tt {
		actual := formatContainerName(v.input, ".")
		if actual != v.expected {
			t.Errorf("\nCalled formatContainerName(%#+v) \ngot %#+v \nwanted %#+v", v.input, actual, v.expected)
		}
	}
}

func TestformatLabels(t *testing.T) {
	tt := []struct {
		input    string
		expected []string
	}{
		{"traefik.backend=web", []string{"traefik.backend", "web"}},
		{"env:prod", []string{"env", "prod"}},
	}
	for _, v := range tt {
		actual := formatLabels(v.input)
		if !reflect.DeepEqual(actual, v.expected) {
			t.Errorf("\nCalled formatLabels(%#+v) \ngot %#+v \nwanted %#+v", v.input, actual, v.expected)
		}
	}
}

func TestformatPorts(t *testing.T) {
	tt := []struct {
		input    string
		expected []string
	}{
		{"6300:6379", []string{"6300", "6379"}},
	}
	for _, v := range tt {
		actual := formatPorts(v.input)
		if !reflect.DeepEqual(actual, v.expected) {
			t.Errorf("\nCalled TestformatPorts(%#+v) \ngot %#+v \nwanted %#+v", v.input, actual, v.expected)
		}
	}
}

func TestformatServiceVolumes(t *testing.T) {
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
		{"$HOME/.aws:/root/.aws", "composefile", []string{os.ExpandEnv("$HOME/.aws"), "/root/.aws"}},
	}
	for _, v := range tt {
		actual := formatServiceVolumes(v.volume, v.dockerComposeFile)
		if !reflect.DeepEqual(actual, v.expected) {
			t.Errorf("\nCalled formatServiceVolumes(%#+v) \ngot %#+v \nwanted %#+v", v.volume, actual, v.expected)
		}
	}
}

func TestformatRegistryAuth(t *testing.T) {
	tt := []struct {
		input    string
		expected []string
	}{
		{"myUsername:myPassword001", []string{"myUsername", "myPassword001"}},
	}
	for _, v := range tt {
		actual := formatRegistryAuth(v.input)
		if !reflect.DeepEqual(actual, v.expected) {
			t.Errorf("\nCalled formatRegistryAuth(%#+v) \ngot %#+v \nwanted %#+v", v.input, actual, v.expected)
		}
	}
}

func TestformatComposePath(t *testing.T) {
	tt := []struct {
		input    string
		expected []string
	}{
		{"testdata/dockerFile", []string{"testdata", "dockerFile"}},
	}
	for _, v := range tt {
		actual := formatComposePath(v.input)
		if !reflect.DeepEqual(actual, v.expected) {
			t.Errorf("\nCalled formatComposePath(%#+v) \ngot %#+v \nwanted %#+v", v.input, actual, v.expected)
		}
	}
}

func BenchmarkformatLabels(b *testing.B) {
	// run the formatLabels function b.N times
	for n := 0; n < b.N; n++ {
		_ = formatLabels("traefik.backend=web")
	}
}

func BenchmarkformatPorts(b *testing.B) {
	for n := 0; n < b.N; n++ {
		_ = formatPorts("6300:6379")
	}
}

func BenchmarkformatServiceVolumes(b *testing.B) {
	for n := 0; n < b.N; n++ {
		_ = formatServiceVolumes("data-volume:/home", "composeFile")
	}
}

func BenchmarkformatContainerName(b *testing.B) {
	for n := 0; n < b.N; n++ {
		_ = formatContainerName("build_with_no_specified_dockerfile", ".")
	}

}
