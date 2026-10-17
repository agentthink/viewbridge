package ota

import "testing"

func TestGetLastVersion(t *testing.T) {
	version, desc, err := GetLastVersion("gowvp/owl")
	if err != nil {
		t.Fatalf("GetLastVersion() error = %v", err)
	}
	t.Logf("version = %s", version)
	t.Logf("desc = %s", desc)
}
