package main

import (
	"super-cc/app/commands"
)

type App struct{}

func NewApp() *App {
	return &App{}
}

func (a *App) CheckApp() commands.AppInfo {
	return commands.ScanApp()
}

func (a *App) GetInstallationInfo() commands.InstallationInfo {
	return commands.DetectEditor()
}

func (a *App) GetVersions() []commands.Version {
	return commands.GetVersions()
}

func (a *App) DownloadVersion(version string) error {
	return commands.DownloadVersion(version)
}

func (a *App) GetCrackStatus() commands.CrackStatus {
	return commands.GetCrackStatusResult()
}

func (a *App) ApplyCrack(enable bool) error {
	return commands.ApplyCrack(enable)
}

func (a *App) GetCleanupStatus() commands.CleanupResult {
	return commands.GetCleanupStatus()
}

func (a *App) RunCleanup() (commands.CleanupResult, error) {
	return commands.RunCleanup()
}

func (a *App) IsLocked() bool {
	return commands.IsLocked()
}

func (a *App) GetLockedVersion() string {
	return commands.GetLockedVersion()
}

func (a *App) SetLock(locked bool, version string) error {
	return commands.SetLock(locked, version)
}

func (a *App) OpenFolder(path string) error {
	return commands.OpenFolder(path)
}

func (a *App) OpenURL(url string) error {
	return commands.OpenURL(url)
}

func (a *App) GetDebugInfo() map[string]string {
	return commands.GetDebugInfo()
}
