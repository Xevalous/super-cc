package config

import (
	"encoding/json"
	"os"
	"path/filepath"
)

type AppConfig struct {
	InstallPath   string `json:"install_path"`
	AutoLock      bool   `json:"auto_lock"`
	AutoCheck     bool   `json:"auto_check"`
	LastCheckPath string `json:"last_check_path"`
	Locked        bool   `json:"locked"`
	LockedVersion string `json:"locked_version"`
	PatchApplied  bool   `json:"patch_applied"`
}

var DefaultConfig = AppConfig{
	AutoLock:  true,
	AutoCheck: true,
}

func GetConfigDir() string {
	configDir, err := os.UserConfigDir()
	if err != nil || configDir == "" {
		home, err := os.UserHomeDir()
		if err != nil || home == "" {
			home = os.TempDir()
		}
		// Fallback to ~/.config if UserConfigDir fails
		configDir = filepath.Join(home, ".config")
	}
	return filepath.Join(configDir, "super-cc")
}

func GetConfigPath() string {
	return filepath.Join(GetConfigDir(), "config.json")
}

func Load() (*AppConfig, error) {
	configPath := GetConfigPath()

	data, err := os.ReadFile(configPath)
	if err != nil {
		if os.IsNotExist(err) {
			return &DefaultConfig, nil
		}
		return nil, err
	}

	var cfg AppConfig
	if err := json.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}

func Save(cfg *AppConfig) error {
	configDir := GetConfigDir()
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return err
	}

	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(GetConfigPath(), data, 0644)
}
