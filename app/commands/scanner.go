package commands

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type InstallationInfo struct {
	Status  string `json:"status"`
	Version string `json:"version"`
	Path    string `json:"path"`
	Size    string `json:"size"`
	LastRun string `json:"lastRun"`
}

func GetAppVersion(exePath string) string {
	info, err := os.Stat(exePath)
	if err != nil {
		return ""
	}
	if info.IsDir() {
		exePath = filepath.Join(exePath, "CapCut.exe")
		info, err = os.Stat(exePath)
		if err != nil {
			return ""
		}
	}

	versionInfo, err := getVersionInfo(exePath)
	if err != nil {
		return extractVersionFromPath(exePath)
	}
	return versionInfo
}

func getVersionInfo(exePath string) (string, error) {
	return extractVersionFromPath(exePath), nil
}

func extractVersionFromPath(exePath string) string {
	base := filepath.Base(exePath)
	parts := strings.Split(base, ".")
	if len(parts) >= 2 {
		for i := 0; i < len(parts)-1; i++ {
			if strings.HasPrefix(parts[i], "v") || strings.HasPrefix(parts[i], "V") {
				version := strings.TrimPrefix(parts[i], "v")
				version = strings.TrimPrefix(version, "V")
				if isValidVersion(version) {
					return version
				}
			}
		}
	}
	return ""
}

func isValidVersion(s string) bool {
	for _, c := range s {
		if c < '0' || c > '9' {
			return false
		}
	}
	return len(s) > 0
}

func GetFolderSize(path string) (int64, error) {
	var totalSize int64

	err := filepath.Walk(path, func(filePath string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			totalSize += info.Size()
		}
		return nil
	})

	return totalSize, err
}

func GetLastRunTime(path string) string {
	exePath := path
	info, err := os.Stat(path)
	if err != nil {
		return ""
	}
	if info.IsDir() {
		exePath = filepath.Join(path, "CapCut.exe")
	}

	fileInfo, err := os.Stat(exePath)
	if err != nil {
		return ""
	}

	modTime := fileInfo.ModTime()
	return modTime.Format("2006-01-02 15:04:05")
}

func DetectEditor() InstallationInfo {
	paths := []string{
		"C:\\Program Files\\CapCut",
		"C:\\Program Files (x86)\\CapCut",
		filepath.Join(os.Getenv("LOCALAPPDATA"), "CapCut"),
		filepath.Join(os.Getenv("APPDATA"), "CapCut"),
	}

	for _, path := range paths {
		info, err := os.Stat(path)
		if err == nil && info.IsDir() {
			version := detectVersionFromApps(path)
			if version == "" {
				version = GetAppVersion(path)
			}
			if version == "" {
				version = extractVersionFromPath(path)
			}

			size, err := GetFolderSize(path)
			sizeStr := "Unknown"
			if err == nil {
				sizeStr = formatSize(size)
			}

			lastRun := GetLastRunTime(path)
			if lastRun == "" {
				lastRun = "Unknown"
			}

			return InstallationInfo{
				Status:  "Installed",
				Version: version,
				Path:    path,
				Size:    sizeStr,
				LastRun: lastRun,
			}
		}
	}

	return InstallationInfo{
		Status:  "Not Installed",
		Version: "",
		Path:    "",
		Size:    "",
		LastRun: "",
	}
}

// detectVersionFromApps checks [installPath]/Apps/ for version subdirectories
func detectVersionFromApps(installPath string) string {
	appsDir := filepath.Join(installPath, "Apps")
	info, err := os.Stat(appsDir)
	if err != nil || !info.IsDir() {
		return ""
	}

	entries, err := os.ReadDir(appsDir)
	if err != nil {
		return ""
	}

	// Look for version-numbered directories like "5.4.0", "5.3.0" etc.
	var versions []string
	for _, entry := range entries {
		if entry.IsDir() && looksLikeVersion(entry.Name()) {
			versions = append(versions, entry.Name())
		}
	}

	if len(versions) == 0 {
		return ""
	}

	// Return the latest version (last in sorted order)
	return latestVersion(versions)
}

// looksLikeVersion checks if a string looks like a version number (e.g., "5.4.0")
func looksLikeVersion(s string) bool {
	parts := strings.Split(s, ".")
	if len(parts) < 2 {
		return false
	}
	for _, part := range parts {
		if !isNumeric(part) {
			return false
		}
	}
	return true
}

func isNumeric(s string) bool {
	if len(s) == 0 {
		return false
	}
	for _, c := range s {
		if c < '0' || c > '9' {
			return false
		}
	}
	return true
}

// latestVersion returns the highest version from a list
func latestVersion(versions []string) string {
	if len(versions) == 0 {
		return ""
	}
	latest := versions[0]
	for _, v := range versions[1:] {
		if compareVersions(v, latest) > 0 {
			latest = v
		}
	}
	return latest
}

// compareVersions compares two version strings like "5.4.0" vs "5.3.0"
func compareVersions(a, b string) int {
	partsA := strings.Split(a, ".")
	partsB := strings.Split(b, ".")
	maxLen := len(partsA)
	if len(partsB) > maxLen {
		maxLen = len(partsB)
	}
	for i := 0; i < maxLen; i++ {
		var numA, numB int
		if i < len(partsA) {
			fmt.Sscanf(partsA[i], "%d", &numA)
		}
		if i < len(partsB) {
			fmt.Sscanf(partsB[i], "%d", &numB)
		}
		if numA > numB {
			return 1
		}
		if numA < numB {
			return -1
		}
	}
	return 0
}

func formatSize(bytes int64) string {
	const unit = 1024
	if bytes < unit {
		return fmt.Sprintf("%d B", bytes)
	}
	div, exp := int64(unit), 0
	for n := bytes / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	size := float64(bytes) / float64(div)
	return fmt.Sprintf("%.2f %cB", size, "KMGTPE"[exp])
}

type AppInfo struct {
	Installed bool   `json:"installed"`
	Path      string `json:"path"`
	Version   string `json:"version"`
}

func ScanApp() AppInfo {
	detected := DetectEditor()
	if detected.Status == "Installed" {
		return AppInfo{
			Installed: true,
			Path:      detected.Path,
			Version:   detected.Version,
		}
	}
	return AppInfo{
		Installed: false,
		Path:      "",
		Version:   "",
	}
}
