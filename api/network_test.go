package api

import "testing"

func TestGetNetwork(t *testing.T) {
	tt := []struct {
		input       string
		expected    string
		expectedErr error
	}{
		{"myNetWorName", "ala", nil},
	}
	for _, v := range tt {
		actual, err := GetNetwork(v.input)
		if err != nil {
			t.Errorf("\nran GetNetwork(%#+v) \ngot %#+v \nwanted %#+v", v.input, err, v.expected)
		}
		if actual != v.expected {
			t.Errorf("\nran GetNetwork(%#+v) \ngot %#+v \nwanted %#+v", v.input, actual, v.expected)
		}
	}
}
