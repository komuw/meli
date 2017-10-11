package api

import (
	"reflect"
	"testing"
)

func TestFomatLabels(t *testing.T) {
	fomatLabelsTests := []struct {
		input    string
		expected []string
	}{
		{"traefik.backend=web", []string{"traefik.backend", "web"}},
		{"env:prod", []string{"env", "prod"}},
	}
	for _, tt := range fomatLabelsTests {
		actual := fomatLabels(tt.input)
		if !reflect.DeepEqual(actual, tt.expected) {
			t.Errorf("\nran fomatLabels(%#+v) \ngot %#+v \nwanted %#+v", tt.input, actual, tt.expected)
		}
	}
}

func BenchmarkFomatLabels(b *testing.B) {
	// run the Fib function b.N times
	for n := 0; n < b.N; n++ {
		fomatLabels("traefik.backend=web")
	}
}

func BenchmarkFormatContainerName(b *testing.B) {
	// run the Fib function b.N times
	for n := 0; n < b.N; n++ {
		FormatContainerName("build_with_no_specified_dockerfile")
	}

}
