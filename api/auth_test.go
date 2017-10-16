package api

import "testing"

func TestGetAuth(t *testing.T) {
	if len(AuthInfo) != 0 {
		t.Errorf("AuthInfo should be empty, it is %d", len(AuthInfo))
	}
	GetAuth()
	if len(AuthInfo) != 2 {
		t.Errorf("AuthInfo should be len 2, it is %d", len(AuthInfo))
	}
}
