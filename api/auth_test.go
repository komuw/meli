package api

import "testing"

func TestGetAuth(t *testing.T) {
	_, okDocker := AuthInfo.Load("dockerhub")
	_, okQuay := AuthInfo.Load("quay")

	if okDocker != false {
		t.Errorf("AuthInfo should not be loaded, we got %t", okDocker)
	}
	if okQuay != false {
		t.Errorf("AuthInfo should not be loaded, we got %t", okQuay)
	}

	GetAuth()
	_, okDocker = AuthInfo.Load("dockerhub")
	_, okQuay = AuthInfo.Load("quay")
	if okDocker != true {
		t.Errorf("AuthInfo should be loaded, we got %t", okDocker)
	}
	if okQuay != true {
		t.Errorf("AuthInfo should be loaded, we got %t", okQuay)
	}
}

func BenchmarkGetAuth(b *testing.B) {
	for n := 0; n < b.N; n++ {
		GetAuth()
	}
}
