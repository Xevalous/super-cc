package commands

import (
	"os"
	"path/filepath"
	"testing"
)

func TestGetCleanupStatus(t *testing.T) {
	result := GetCleanupStatus()

	// Status should be one of the known values
	switch result.Status {
	case "idle", "residual_found", "clean":
		// OK
	default:
		t.Errorf("unexpected status: %q", result.Status)
	}

	// LastCheck should be set
	if result.LastCheck == "" {
		t.Error("LastCheck should not be empty")
	}

	// FilesFound should match residual paths count
	if result.FilesFound != len(result.ResidualPaths) {
		t.Errorf("FilesFound (%d) != len(ResidualPaths) (%d)", result.FilesFound, len(result.ResidualPaths))
	}

	// FilesCleaned should be 0 before cleanup runs
	if result.FilesCleaned != 0 {
		t.Errorf("FilesCleaned should be 0 before cleanup, got %d", result.FilesCleaned)
	}
}

func TestGetResidualPaths(t *testing.T) {
	paths := getResidualPaths()

	// On Linux, this should return empty since it checks Windows paths
	t.Logf("getResidualPaths() returned %d paths", len(paths))
	for _, p := range paths {
		t.Logf("  residual path: %s", p)
	}
}

func TestRunCleanup(t *testing.T) {
	// Create a temp dir to simulate install
	tmpDir, err := os.MkdirTemp("", "supercc-cleanup-test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	testFile := filepath.Join(tmpDir, "dummy.dll")
	if err := os.WriteFile(testFile, []byte("test"), 0644); err != nil {
		t.Fatal(err)
	}

	result, err := RunCleanup()
	if err != nil {
		t.Fatalf("RunCleanup returned error: %v", err)
	}

	// Status should be updated
	if result.Status != "completed" && result.Status != "already_clean" {
		t.Errorf("unexpected status after cleanup: %q", result.Status)
	}
}

func TestOpenFolder(t *testing.T) {
	// OpenFolder uses "explorer" which is Windows-only
	// On non-Windows, this will return an error
	tmpDir, err := os.MkdirTemp("", "supercc-open-test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	err = OpenFolder(tmpDir)
	if err != nil {
		t.Logf("OpenFolder returned error (expected on non-Windows): %v", err)
	}
}
