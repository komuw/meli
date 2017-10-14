package api

import "testing"

// func TestGetRegistryAuth(t *testing.T) {
// 	tt := []struct {
// 		input       string
// 		expected    string
// 		expectedErr error
// 	}{
// 		{"ImageName", "RegistryAuth", nil},
// 	}
// 	for _, v := range tt {
// 		actual, err := GetRegistryAuth(v.input)
// 		if err != nil {
// 			t.Errorf("\nran GetRegistryAuth(%#+v) \ngot %s \nwanted %#+v", v.input, err, v.expected)
// 		}
// 		if actual != v.expected {
// 			t.Errorf("\nran GetRegistryAuth(%#+v) \ngot %#+v \nwanted %#+v", v.input, actual, v.expected)
// 		}
// 	}
// }

func TestGetAuth(t *testing.T) {
	tt := []struct {
		input       string
		expected    string
		expectedErr error
	}{
		{"myImageName", "https://index.docker.io/v1/", nil},
		{"quay.io/quayImage", "quay.io", nil},
	}
	for _, v := range tt {
		registryURL, _, _, err := GetAuth(v.input)
		if err != nil {
			t.Errorf("\nran GetAuth(%#+v) \ngot %s \nwanted %#+v", v.input, err, v.expectedErr)
		}
		if registryURL != v.expected {
			t.Errorf("\nran GetAuth(%#+v) \ngot %#+v \nwanted %#+v", v.input, registryURL, v.expected)
		}
	}
}

func BenchmarkGetAuth(b *testing.B) {
	for n := 0; n < b.N; n++ {
		_, _, _, _ = GetAuth("myImageName")
	}
}
