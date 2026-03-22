package commands

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"super-cc/app/config"
)

type CleanupResult struct {
	LastCheck     string   `json:"lastCheck"`
	Status        string   `json:"status"`
	FilesFound    int      `json:"filesFound"`
	FilesCleaned  int      `json:"filesCleaned"`
	ResidualPaths []string `json:"residualPaths"`
	Errors        []string `json:"errors"`
}

func GetCleanupStatus() CleanupResult {
	residualPaths := getResidualPaths()

	result := CleanupResult{
		LastCheck:     time.Now().Format("2006-01-02 15:04:05"),
		Status:        "idle",
		FilesFound:    len(residualPaths),
		FilesCleaned:  0,
		ResidualPaths: residualPaths,
	}

	if len(residualPaths) > 0 {
		result.Status = "residual_found"
	} else {
		result.Status = "clean"
	}

	return result
}

func getResidualPaths() []string {
	var paths []string

	envVars := map[string]string{
		"APPDATA":           os.Getenv("APPDATA"),
		"LOCALAPPDATA":      os.Getenv("LOCALAPPDATA"),
		"ProgramFiles":      os.Getenv("ProgramFiles"),
		"ProgramFiles(x86)": os.Getenv("ProgramFiles(x86)"),
	}

	locations := []struct {
		envVar string
		subDir string
	}{
		{"APPDATA", "CapCut"},
		{"APPDATA", "capcut"},
		{"LOCALAPPDATA", "CapCut"},
		{"LOCALAPPDATA", "capcut"},
		{"LOCALAPPDATA", "Programs\\CapCut"},
		{"ProgramFiles", "CapCut"},
		{"ProgramFiles(x86)", "CapCut"},
	}

	for _, loc := range locations {
		basePath, ok := envVars[loc.envVar]
		if !ok || basePath == "" {
			continue
		}

		fullPath := filepath.Join(basePath, loc.subDir)
		if _, err := os.Stat(fullPath); err == nil {
			paths = append(paths, fullPath)
		}
	}

	return paths
}

func RunCleanup() (CleanupResult, error) {
	result := GetCleanupStatus()

	cleaned := 0
	var cleanupErrors []string
	for _, path := range result.ResidualPaths {
		if err := os.RemoveAll(path); err != nil {
			cleanupErrors = append(cleanupErrors, fmt.Sprintf("failed to remove %s: %v", path, err))
			continue
		}
		cleaned++
	}

	cfg, err := config.Load()
	if err == nil {
		cfg.InstallPath = ""
		if err := config.Save(cfg); err != nil {
			cleanupErrors = append(cleanupErrors, fmt.Sprintf("failed to save config: %v", err))
		}
	}

	result.FilesCleaned = cleaned
	result.Errors = cleanupErrors
	result.Status = "completed"
	result.LastCheck = time.Now().Format("2006-01-02 15:04:05")

	if len(result.ResidualPaths) == 0 {
		result.Status = "already_clean"
	}

	return result, nil
}

func OpenFolder(path string) error {
	return exec.Command("explorer", path).Start()
}

func OpenURL(url string) error {
	if !strings.HasPrefix(url, "https://") && !strings.HasPrefix(url, "http://") {
		return fmt.Errorf("invalid URL scheme: only http:// and https:// are allowed")
	}
	return exec.Command("rundll32", "url.dll,FileProtocolHandler", url).Start()
}
