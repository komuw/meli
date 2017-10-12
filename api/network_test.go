package api

import (
	"reflect"
	"testing"

	"golang.org/x/net/context"
)

func TestGetNetwork(t *testing.T) {
	tt := []struct {
		input       string
		expected    string
		expectedErr error
	}{
		{"myNetWorName", "string", nil},
	}
	for _, v := range tt {
		actual, err := GetNetwork(v.input)
		if err != nil {
			t.Errorf("\nran GetNetwork(%#+v) \ngot %s \nwanted %#+v", v.input, err, v.expected)
		}
		if reflect.TypeOf(actual).String() != "string" {
			t.Errorf("\nran GetNetwork(%#+v) \ngot %#+v \nwanted %#+v", v.input, reflect.TypeOf(actual).String(), v.expected)
		}
	}
}

func TestConnectNetwork(t *testing.T) {
	tt := []struct {
		input1      context.Context
		input2      string
		input3      string
		expectedErr error
	}{
		{context.Background(), "netID", "containerID", nil},
	}
	for _, v := range tt {
		err := ConnectNetwork(v.input1, v.input2, v.input3)
		if err != nil {
			t.Errorf("\nran ConnectNetwork(%#+v) \ngot %s \nwanted %#+v", v.input1, err, v.expectedErr)
		}

	}
}
