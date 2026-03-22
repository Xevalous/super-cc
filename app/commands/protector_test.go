package commands

import (
	"testing"
)

func TestIsLocked(t *testing.T) {
	locked := IsLocked()
	// On Linux, this should return false since hosts file check is Windows-only
	t.Logf("IsLocked() = %v", locked)
}

func TestGetLockedVersion(t *testing.T) {
	version := GetLockedVersion()
	t.Logf("GetLockedVersion() = %q", version)
	// Version can be empty if no config exists yet, that's valid
}

func TestSetLock(t *testing.T) {
	// SetLock uses hosts file which is Windows-specific path
	err := SetLock(false, "")
	if err != nil {
		t.Logf("SetLock returned error (expected on non-Windows): %v", err)
	}
}
