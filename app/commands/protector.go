package commands

import (
	"os"
	"strings"

	"super-cc/app/config"
)

var HostsFilePath = "C:\\Windows\\System32\\drivers\\etc\\hosts"

var blockedHosts = []string{
	"lf16-capcut.faceulv.com",
	"api.capcut.com",
	"update.capcut.com",
}

type LockStatus struct {
	Locked  bool   `json:"locked"`
	Version string `json:"version"`
}

func IsLocked() bool {
	// Check config first (most reliable)
	cfg, err := config.Load()
	if err == nil && cfg.Locked {
		return true
	}

	// Then check hosts file
	content, err := os.ReadFile(HostsFilePath)
	if err != nil {
		return false
	}

	lines := strings.Split(string(content), "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "#") {
			continue
		}
		if strings.HasPrefix(line, "127.0.0.1") {
			for _, host := range blockedHosts {
				if strings.Contains(line, host) {
					return true
				}
			}
		}
	}
	return false
}

func GetLockedVersion() string {
	cfg, err := config.Load()
	if err != nil {
		return ""
	}
	return cfg.LockedVersion
}

func SetLock(locked bool, version string) error {
	cfg, err := config.Load()
	if err != nil {
		cfg = &config.AppConfig{
			AutoLock:  true,
			AutoCheck: true,
		}
	}

	if locked {
		cfg.Locked = true
		cfg.LockedVersion = version
	} else {
		cfg.Locked = false
		cfg.LockedVersion = ""
	}

	if err := config.Save(cfg); err != nil {
		return err
	}

	if locked {
		return enableLock()
	}
	return disableLock()
}

func enableLock() error {
	hostsFilePath := HostsFilePath
	content, err := os.ReadFile(hostsFilePath)
	if err != nil {
		return err
	}

	existingContent := string(content)

	// Check if block already exists
	if strings.Contains(existingContent, "# Update Block") {
		return nil
	}

	newContent := existingContent + "\n\n# Update Block\n"
	for _, host := range blockedHosts {
		newContent += "127.0.0.1 " + host + "\n"
	}

	return os.WriteFile(HostsFilePath, []byte(newContent), 0644)
}

func disableLock() error {
	hostsFilePath := HostsFilePath
	content, err := os.ReadFile(hostsFilePath)
	if err != nil {
		return err
	}

	lines := strings.Split(string(content), "\n")
	var newLines []string
	skipBlock := false

	for _, line := range lines {
		trimmed := strings.TrimSpace(line)

		// Start of our block
		if strings.Contains(line, "# Update Block") {
			skipBlock = true
			continue
		}

		if skipBlock {
			// Empty line ends the block
			if trimmed == "" {
				skipBlock = false
				newLines = append(newLines, line)
				continue
			}
			// Skip any 127.0.0.1 line that matches our blocked hosts
			if strings.HasPrefix(trimmed, "127.0.0.1") {
				isBlocked := false
				for _, host := range blockedHosts {
					if strings.Contains(trimmed, host) {
						isBlocked = true
						break
					}
				}
				if isBlocked {
					continue
				}
			}
			// Non-matching line ends the block
			skipBlock = false
		}

		newLines = append(newLines, line)
	}

	return os.WriteFile(hostsFilePath, []byte(strings.Join(newLines, "\n")), 0644)
}

type ProtectorError struct {
	Message string
}

func (e *ProtectorError) Error() string {
	return e.Message
}
