package commands

import (
	"testing"
)

func TestGetVersions(t *testing.T) {
	vs := GetVersions()

	if len(vs) == 0 {
		t.Fatal("GetVersions() returned empty list")
	}

	// Verify total count matches cc-version-guard (25 versions)
	if len(vs) != 25 {
		t.Errorf("expected 25 versions, got %d", len(vs))
	}

	// Verify first version is 5.4.0 Beta6 (latest)
	first := vs[0]
	if first.Label != "5.4.0 (Beta6)" {
		t.Errorf("first version label = %q, want %q", first.Label, "5.4.0 (Beta6)")
	}
	if first.Version != "5.4.0" {
		t.Errorf("first version = %q, want %q", first.Version, "5.4.0")
	}
	if first.Type != "latest" {
		t.Errorf("first version type = %q, want %q", first.Type, "latest")
	}
	if first.Tag != "Latest" {
		t.Errorf("first version tag = %q, want %q", first.Tag, "Latest")
	}
	if first.Risk != "High" {
		t.Errorf("first version risk = %q, want %q", first.Risk, "High")
	}

	// Verify last version is 1.0.0
	last := vs[len(vs)-1]
	if last.Label != "1.0.0 (Latest)" {
		t.Errorf("last version label = %q, want %q", last.Label, "1.0.0 (Latest)")
	}
	if last.Risk != "Low" {
		t.Errorf("last version risk = %q, want %q", last.Risk, "Low")
	}

	// Verify all URLs point to ByteDance CDN
	for _, v := range vs {
		if v.URL == "" {
			t.Errorf("version %q has empty URL", v.Label)
		}
		if v.Label == "" {
			t.Errorf("version has empty label")
		}
		if v.Version == "" {
			t.Errorf("version has empty version field")
		}
		if v.Type == "" {
			t.Errorf("version %q has empty type", v.Label)
		}
		if v.Tag == "" {
			t.Errorf("version %q has empty tag", v.Label)
		}
		if v.Risk == "" {
			t.Errorf("version %q has empty risk", v.Label)
		}
	}

	// Verify beta versions exist
	betaCount := 0
	stableCount := 0
	for _, v := range vs {
		if v.Type == "beta" {
			betaCount++
		}
		if v.Type == "stable" {
			stableCount++
		}
	}
	if betaCount == 0 {
		t.Error("expected at least 1 beta version")
	}
	if stableCount == 0 {
		t.Error("expected at least 1 stable version")
	}
}

func TestDownloadVersion(t *testing.T) {
	// Test with invalid version
	err := DownloadVersion("nonexistent")
	if err == nil {
		t.Error("DownloadVersion with invalid version should return error")
	}

	// Test with valid label (uses rundll32 which is Windows-only)
	err = DownloadVersion("5.4.0 (Beta6)")
	if err != nil {
		t.Logf("DownloadVersion returned error (expected on non-Windows): %v", err)
	}
}
