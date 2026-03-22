package commands

import (
	"testing"
)

func TestGetCrackStatus(t *testing.T) {
	mode, vip, watermark := GetCrackStatus()

	// Mode should be one of known values or "unknown" when DLL not found
	switch mode {
	case "unknown", "free", "vip", "pro":
		// OK
	default:
		t.Errorf("unexpected crack mode: %q", mode)
	}

	t.Logf("GetCrackStatus() mode=%q vip=%v watermark=%v", mode, vip, watermark)
}

func TestFindVECreatorDLL(t *testing.T) {
	// Test with empty path (should scan system)
	dllPath := findVECreatorDLL("")
	t.Logf("findVECreatorDLL('') = %q", dllPath)
	// On Linux, this will be empty since the editor isn't installed

	// Test with non-existent path
	dllPath = findVECreatorDLL("/nonexistent/path")
	if dllPath != "" {
		t.Errorf("expected empty for non-existent path, got %q", dllPath)
	}
}

func TestKillEditor(t *testing.T) {
	err := KillEditor()
	// taskkill is Windows-specific, may fail on non-Windows
	if err != nil {
		t.Logf("KillEditor returned error (expected on non-Windows): %v", err)
	}
}
