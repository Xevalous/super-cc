package commands

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"super-cc/app/config"
)

var (
	vipPattern = []byte("\x00vip_entrance\x00")
	proPattern = []byte("\x00pro_fortnite\x00")
)

type CrackStatus struct {
	Mode             string `json:"mode"`
	VIPEnabled       bool   `json:"vipEnabled"`
	WatermarkRemoved bool   `json:"watermarkRemoved"`
	DLLFound         bool   `json:"dllFound"`
}

func GetCrackStatusResult() CrackStatus {
	mode, vip, watermark := GetCrackStatus()
	return CrackStatus{
		Mode:             mode,
		VIPEnabled:       vip,
		WatermarkRemoved: watermark,
	}
}

func GetCrackStatus() (string, bool, bool) {
	cfg, _ := config.Load()

	dllPath := findVECreatorDLL(cfg.InstallPath)
	if dllPath == "" {
		return "unknown", false, false
	}

	data, err := os.ReadFile(dllPath)
	if err != nil {
		return "unknown", false, false
	}

	hasVipPattern := bytes.Contains(data, vipPattern)
	hasProPattern := bytes.Contains(data, proPattern)

	// Pattern FOUND = DLL is unpatched (restrictions active)
	// Pattern NOT FOUND = DLL is patched (patterns nulled out)
	enabled := !hasVipPattern && !hasProPattern

	var mode string
	if enabled {
		mode = "patched"
	} else {
		mode = "free"
	}

	watermarkRemoved := !bytes.Contains(data, []byte("watermark"))

	return mode, enabled, watermarkRemoved
}

func ApplyCrack(enable bool) error {
	cfg, err := config.Load()
	if err != nil {
		return err
	}

	dllPath := findVECreatorDLL(cfg.InstallPath)
	if dllPath == "" {
		return &CrackError{Message: "VECreator.dll not found. Please ensure the application is installed."}
	}

	if err := KillEditor(); err != nil {
		return err
	}

	time.Sleep(500 * time.Millisecond)

	data, err := os.ReadFile(dllPath)
	if err != nil {
		return err
	}

	backupPath := dllPath + ".bak"
	if _, err := os.Stat(backupPath); os.IsNotExist(err) {
		if err := os.WriteFile(backupPath, data, 0644); err != nil {
			return err
		}
	}

	modified := make([]byte, len(data))
	copy(modified, data)

	if enable {
		modified = bytes.Replace(modified, vipPattern, bytes.Repeat([]byte{0}, len(vipPattern)), -1)
		modified = bytes.Replace(modified, proPattern, bytes.Repeat([]byte{0}, len(proPattern)), -1)
	} else {
		backupData, err := os.ReadFile(backupPath)
		if err != nil {
			return err
		}
		modified = backupData
	}

	if err := os.WriteFile(dllPath, modified, 0644); err != nil {
		return err
	}

	// Save patch state to config
	cfg, err = config.Load()
	if err != nil {
		cfg = &config.AppConfig{
			AutoLock:  true,
			AutoCheck: true,
		}
	}
	cfg.PatchApplied = enable
	if err := config.Save(cfg); err != nil {
		return err
	}

	return nil
}

func KillEditor() error {
	processes := []string{"CapCut.exe", "CapCut Video Editor.exe"}
	var errs []string
	for _, proc := range processes {
		cmd := exec.Command("taskkill", "/F", "/IM", proc)
		if err := cmd.Run(); err != nil {
			errs = append(errs, fmt.Sprintf("failed to kill %s: %v", proc, err))
		}
	}
	if len(errs) > 0 {
		return fmt.Errorf("failed to terminate editor processes: %s", strings.Join(errs, "; "))
	}
	return nil
}

func findVECreatorDLL(installPath string) string {
	if installPath == "" {
		detected := DetectEditor()
		installPath = detected.Path
	}

	if installPath == "" {
		return ""
	}

	// Search in Apps/[version]/ subdirectories first (primary location)
	appsDir := filepath.Join(installPath, "Apps")
	if entries, err := os.ReadDir(appsDir); err == nil {
		// Check version folders in reverse order (latest first)
		for i := len(entries) - 1; i >= 0; i-- {
			if !entries[i].IsDir() {
				continue
			}
			dllPath := filepath.Join(appsDir, entries[i].Name(), "VECreator.dll")
			if _, err := os.Stat(dllPath); err == nil {
				return dllPath
			}
		}
	}

	// Fallback: check root and common subdirectories
	searchPaths := []string{
		filepath.Join(installPath, "VECreator.dll"),
		filepath.Join(installPath, "bin", "VECreator.dll"),
		filepath.Join(installPath, "resources", "app", "VECreator.dll"),
	}

	for _, path := range searchPaths {
		if _, err := os.Stat(path); err == nil {
			return path
		}
	}

	return ""
}

type CrackError struct {
	Message string
}

func (e *CrackError) Error() string {
	return e.Message
}
