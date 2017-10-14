package api

import (
	"reflect"
	"testing"
)

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
			t.Errorf("\nran FormatLabels(%#+v) \ngot %#+v \nwanted %#+v", v.input, actual, v.expected)
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
			t.Errorf("\nran TestFormatPorts(%#+v) \ngot %#+v \nwanted %#+v", v.input, actual, v.expected)
		}
	}
}

func TestFormatServiceVolumes(t *testing.T) {
	tt := []struct {
		input    string
		expected []string
	}{
		{"data-volume:/home", []string{"data-volume", "/home"}},
	}
	for _, v := range tt {
		actual := FormatServiceVolumes(v.input)
		if !reflect.DeepEqual(actual, v.expected) {
			t.Errorf("\nran FormatServiceVolumes(%#+v) \ngot %#+v \nwanted %#+v", v.input, actual, v.expected)
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
			t.Errorf("\nran FormatRegistryAuth(%#+v) \ngot %#+v \nwanted %#+v", v.input, actual, v.expected)
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
			t.Errorf("\nran FormatComposePath(%#+v) \ngot %#+v \nwanted %#+v", v.input, actual, v.expected)
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
		_ = FormatServiceVolumes("data-volume:/home")
	}
}

func BenchmarkFormatContainerName(b *testing.B) {
	for n := 0; n < b.N; n++ {
		_ = FormatContainerName("build_with_no_specified_dockerfile")
	}

}
