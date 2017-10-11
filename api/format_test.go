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

func BenchmarkFormatLabels(b *testing.B) {
	// run the FormatLabels function b.N times
	for n := 0; n < b.N; n++ {
		FormatLabels("traefik.backend=web")
	}
}

func BenchmarkFormatContainerName(b *testing.B) {
	for n := 0; n < b.N; n++ {
		FormatContainerName("build_with_no_specified_dockerfile")
	}

}
