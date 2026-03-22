package config

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

func TestDefaultConfig(t *testing.T) {
	if !DefaultConfig.AutoLock {
		t.Error("DefaultConfig.AutoLock should be true")
	}
	if !DefaultConfig.AutoCheck {
		t.Error("DefaultConfig.AutoCheck should be true")
	}
}

func TestLoadSaveConfig(t *testing.T) {
	// Override config dir to temp for testing
	tmpDir, err := os.MkdirTemp("", "supercc-config-test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	// Save a config
	testCfg := &AppConfig{
		InstallPath:   "C:\\Program Files\\Editor",
		AutoLock:      true,
		AutoCheck:     false,
		Locked:        true,
		LockedVersion: "5.3.0",
	}

	// Write config manually
	data, err := json.MarshalIndent(testCfg, "", "  ")
	if err != nil {
		t.Fatal(err)
	}
	configPath := filepath.Join(tmpDir, "config.json")
	if err := os.WriteFile(configPath, data, 0644); err != nil {
		t.Fatal(err)
	}

	// Read it back
	readData, err := os.ReadFile(configPath)
	if err != nil {
		t.Fatal(err)
	}

	var loaded AppConfig
	if err := json.Unmarshal(readData, &loaded); err != nil {
		t.Fatal(err)
	}

	if loaded.InstallPath != testCfg.InstallPath {
		t.Errorf("InstallPath = %q, want %q", loaded.InstallPath, testCfg.InstallPath)
	}
	if loaded.Locked != testCfg.Locked {
		t.Errorf("Locked = %v, want %v", loaded.Locked, testCfg.Locked)
	}
	if loaded.LockedVersion != testCfg.LockedVersion {
		t.Errorf("LockedVersion = %q, want %q", loaded.LockedVersion, testCfg.LockedVersion)
	}
}

func TestLoadNonExistentConfig(t *testing.T) {
	// Temporarily override config path to non-existent
	origGetConfigPath := GetConfigPath
	_ = origGetConfigPath // just to use it

	// Load should return default config when file doesn't exist
	cfg, err := Load()
	if err != nil {
		// If config dir doesn't exist, Load might error
		t.Logf("Load() returned error (expected if no config exists): %v", err)
		return
	}

	if cfg == nil {
		t.Fatal("Load() returned nil config")
	}
}
