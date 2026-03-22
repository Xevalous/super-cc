package commands

import (
	"os"
	"path/filepath"
	"testing"
)

func TestDetectEditor(t *testing.T) {
	info := DetectEditor()

	// On Linux, the editor is not installed, so status should be "Not Installed"
	if info.Status != "Installed" && info.Status != "Not Installed" {
		t.Errorf("unexpected status: %q", info.Status)
	}

	if info.Status == "Not Installed" {
		if info.Path != "" {
			t.Errorf("path should be empty when not installed, got %q", info.Path)
		}
	}
}

func TestGetAppVersion(t *testing.T) {
	// Test with non-existent path
	version := GetAppVersion("/nonexistent/path")
	if version != "" {
		t.Errorf("expected empty version for non-existent path, got %q", version)
	}
}

func TestExtractVersionFromPath(t *testing.T) {
	tests := []struct {
		path string
		want string
	}{
		// Function splits filename by "." and checks each part for "v" or "V" prefix
		{"/opt/Editor/v5.4.0.exe", "5"},    // ["v5", "4", "0", "exe"] → "v5" found
		{"/opt/Editor/V450.exe", "450"},    // ["V450", "exe"] → "V450" found
		{"/opt/Editor/App_v5.4.0.exe", ""}, // ["App_v5", "4", "0", "exe"] → no part starts with v
		{"/opt/Editor/App.exe", ""},        // no version in filename
		{"/opt/Editor/v5.4.0/App.exe", ""}, // version in dir, not filename
		{"", ""},
	}

	for _, tt := range tests {
		got := extractVersionFromPath(tt.path)
		if got != tt.want {
			t.Errorf("extractVersionFromPath(%q) = %q, want %q", tt.path, got, tt.want)
		}
	}
}

func TestFormatSize(t *testing.T) {
	tests := []struct {
		bytes int64
		want  string
	}{
		{0, "0 B"},
		{1024, "1.00 KB"},
		{1048576, "1.00 MB"},
		{1073741824, "1.00 GB"},
		{500, "500 B"},
	}

	for _, tt := range tests {
		got := formatSize(tt.bytes)
		if got != tt.want {
			t.Errorf("formatSize(%d) = %q, want %q", tt.bytes, got, tt.want)
		}
	}
}

func TestGetFolderSize(t *testing.T) {
	// Test with a real temp directory
	tmpDir, err := os.MkdirTemp("", "supercc-test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	// Create a file with known size
	testFile := filepath.Join(tmpDir, "test.txt")
	content := "hello world"
	if err := os.WriteFile(testFile, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}

	size, err := GetFolderSize(tmpDir)
	if err != nil {
		t.Fatalf("GetFolderSize returned error: %v", err)
	}
	if size != int64(len(content)) {
		t.Errorf("GetFolderSize = %d, want %d", size, len(content))
	}
}

func TestScanApp(t *testing.T) {
	info := ScanApp()

	// Should return a valid struct regardless of installation status
	if info.Installed && info.Path == "" {
		t.Error("if installed, path should not be empty")
	}
}

func TestIsValidVersion(t *testing.T) {
	tests := []struct {
		input string
		want  bool
	}{
		{"540", true},
		{"5.4.0", false}, // contains dots
		{"", false},
		{"abc", false},
	}

	for _, tt := range tests {
		got := isValidVersion(tt.input)
		if got != tt.want {
			t.Errorf("isValidVersion(%q) = %v, want %v", tt.input, got, tt.want)
		}
	}
}

func TestLooksLikeVersion(t *testing.T) {
	tests := []struct {
		input string
		want  bool
	}{
		{"5.4.0", true},
		{"5.3.0", true},
		{"1.0.0", true},
		{"10.2.3", true},
		{"Editor", false},
		{"beta", false},
		{"5", false}, // needs at least 2 parts
		{"", false},
		{"5.4.0rc1", false}, // contains letters
	}

	for _, tt := range tests {
		got := looksLikeVersion(tt.input)
		if got != tt.want {
			t.Errorf("looksLikeVersion(%q) = %v, want %v", tt.input, got, tt.want)
		}
	}
}

func TestCompareVersions(t *testing.T) {
	tests := []struct {
		a, b string
		want int
	}{
		{"5.4.0", "5.3.0", 1},
		{"5.3.0", "5.4.0", -1},
		{"5.4.0", "5.4.0", 0},
		{"10.0.0", "9.0.0", 1},
		{"1.0.0", "1.0.1", -1},
	}

	for _, tt := range tests {
		got := compareVersions(tt.a, tt.b)
		if got != tt.want {
			t.Errorf("compareVersions(%q, %q) = %d, want %d", tt.a, tt.b, got, tt.want)
		}
	}
}

func TestLatestVersion(t *testing.T) {
	tests := []struct {
		versions []string
		want     string
	}{
		{[]string{"5.3.0", "5.4.0", "5.2.0"}, "5.4.0"},
		{[]string{"1.0.0"}, "1.0.0"},
		{[]string{}, ""},
		{[]string{"9.0.0", "10.0.0", "2.0.0"}, "10.0.0"},
	}

	for _, tt := range tests {
		got := latestVersion(tt.versions)
		if got != tt.want {
			t.Errorf("latestVersion(%v) = %q, want %q", tt.versions, got, tt.want)
		}
	}
}

func TestDetectVersionFromApps(t *testing.T) {
	// Create temp dir structure: /tmp/test/Apps/5.3.0, /tmp/test/Apps/5.4.0
	tmpDir, err := os.MkdirTemp("", "supercc-apps-test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	appsDir := filepath.Join(tmpDir, "Apps")
	os.MkdirAll(filepath.Join(appsDir, "5.3.0"), 0755)
	os.MkdirAll(filepath.Join(appsDir, "5.4.0"), 0755)
	os.MkdirAll(filepath.Join(appsDir, "notaversion"), 0755)

	version := detectVersionFromApps(tmpDir)
	if version != "5.4.0" {
		t.Errorf("detectVersionFromApps = %q, want %q", version, "5.4.0")
	}

	// Test with no Apps dir
	version = detectVersionFromApps("/nonexistent")
	if version != "" {
		t.Errorf("detectVersionFromApps with no Apps dir = %q, want empty", version)
	}
}
