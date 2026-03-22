package commands

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"super-cc/app/config"
)

func GetDebugInfo() map[string]string {
	info := make(map[string]string)

	// System
	info["os"] = runtime.GOOS
	info["arch"] = runtime.GOARCH
	info["go_version"] = runtime.Version()
	info["num_cpu"] = fmt.Sprintf("%d", runtime.NumCPU())
	info["time"] = time.Now().Format("2006-01-02 15:04:05")

	// Paths
	home, _ := os.UserHomeDir()
	info["home"] = home
	info["config_dir"] = config.GetConfigDir()
	info["config_path"] = config.GetConfigPath()

	// Config state
	cfg, err := config.Load()
	if err == nil {
		info["config_loaded"] = "true"
		info["cfg_install_path"] = cfg.InstallPath
		info["cfg_locked"] = fmt.Sprintf("%v", cfg.Locked)
		info["cfg_locked_version"] = cfg.LockedVersion
		info["cfg_patch_applied"] = fmt.Sprintf("%v", cfg.PatchApplied)
		info["cfg_auto_lock"] = fmt.Sprintf("%v", cfg.AutoLock)
		info["cfg_auto_check"] = fmt.Sprintf("%v", cfg.AutoCheck)
	} else {
		info["config_loaded"] = fmt.Sprintf("error: %v", err)
	}

	// Environment vars
	envVars := []string{"LOCALAPPDATA", "APPDATA", "ProgramFiles", "ProgramFiles(x86)", "USERPROFILE", "COMSPEC"}
	for _, v := range envVars {
		info["env_"+v] = os.Getenv(v)
	}

	// Installation detection
	installed := DetectEditor()
	info["install_status"] = installed.Status
	info["install_version"] = installed.Version
	info["install_path"] = installed.Path
	info["install_size"] = installed.Size

	// Lock status
	info["is_locked"] = fmt.Sprintf("%v", IsLocked())
	info["locked_version"] = GetLockedVersion()

	// Crack status
	mode, vip, watermark := GetCrackStatus()
	info["crack_mode"] = mode
	info["crack_vip"] = fmt.Sprintf("%v", vip)
	info["crack_watermark"] = fmt.Sprintf("%v", watermark)

	// Hosts file
	hostsContent, err := os.ReadFile(HostsFilePath)
	if err == nil {
		lines := strings.Split(string(hostsContent), "\n")
		var blockedLines []string
		for _, line := range lines {
			trimmed := strings.TrimSpace(line)
			if strings.Contains(trimmed, "capcut") || strings.Contains(trimmed, "faceulv") {
				blockedLines = append(blockedLines, trimmed)
			}
		}
		info["hosts_blocked_entries"] = strings.Join(blockedLines, " | ")
		if len(blockedLines) == 0 {
			info["hosts_blocked_entries"] = "(none)"
		}
	} else {
		info["hosts_blocked_entries"] = fmt.Sprintf("error: %v", err)
	}

	// DLL paths
	if installed.Path != "" {
		var found []string

		// Check Apps/[version]/ first
		appsDir := filepath.Join(installed.Path, "Apps")
		if entries, err := os.ReadDir(appsDir); err == nil {
			for _, entry := range entries {
				if !entry.IsDir() {
					continue
				}
				dllPath := filepath.Join(appsDir, entry.Name(), "VECreator.dll")
				if _, err := os.Stat(dllPath); err == nil {
					found = append(found, dllPath)
				}
			}
		}

		// Fallback paths
		fallbackPaths := []string{
			filepath.Join(installed.Path, "VECreator.dll"),
			filepath.Join(installed.Path, "bin", "VECreator.dll"),
			filepath.Join(installed.Path, "resources", "app", "VECreator.dll"),
		}
		for _, p := range fallbackPaths {
			if _, err := os.Stat(p); err == nil {
				found = append(found, p)
			}
		}

		info["dll_found"] = strings.Join(found, " | ")
		if len(found) == 0 {
			info["dll_found"] = "(none)"
		}
	} else {
		info["dll_found"] = "(no install path)"
	}

	// Residual paths
	residuals := getResidualPaths()
	if len(residuals) > 0 {
		info["residual_paths"] = strings.Join(residuals, " | ")
	} else {
		info["residual_paths"] = "(none)"
	}

	// Download versions count
	info["download_versions"] = fmt.Sprintf("%d", len(GetVersions()))

	return info
}
