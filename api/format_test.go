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
			t.Errorf("\nran FormatLabels(%#+v) \ngot %#+v \nwanted %#+v", v.input, actual, v.expected)
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

func BenchmarkFormatContainerName(b *testing.B) {
	for n := 0; n < b.N; n++ {
		_ = FormatContainerName("build_with_no_specified_dockerfile")
	}

}
