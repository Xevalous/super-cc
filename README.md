# Super CC

A desktop tool for managing video editor installations on Windows. Built with [Wails](https://wails.io/) (Go + WebView2).

## Features

- **Dashboard** — Scan and display installation status, version, size, and path
- **Download** — Browse and download 25+ versions with real CDN URLs, search/filter support. Version must be unlocked first (downloads are blocked while locked).
- **Lock** — Block update servers via hosts file to prevent auto-updates
- **Activate** — Patch DLL to unlock premium features (with automatic `.bak` backup)
- **Cleanup** — Remove residual files from uninstalled versions
- **Debug** — Full environment diagnostics and system log

## Important

> **This application must be run as Administrator.** It modifies system files (hosts file) and application DLLs which require elevated permissions.

> **When applying the patch (Activate), the application must be running.** The tool will close it before patching and restart it after. Make sure the application is open before clicking "Apply Patch".

> **To download a version, you must unlock it first.** Download URLs are blocked when the version is locked. Navigate to the Lock page and unlock before downloading.

## Requirements

- Windows 10/11 (64-bit)
- [Go](https://go.dev/dl/) 1.21+
- [Wails CLI](https://wails.io/docs/gettingstarted/installation) v2.x
- WebView2 Runtime (pre-installed on Windows 10/11)

## Build from Source

### Install Wails CLI

```bash
go install github.com/wailsapp/wails/v2/cmd/wails@latest
```

### Clone Repository

```bash
git clone https://github.com/Xevalous/super-cc.git
cd super-cc
```

### Development Mode

```bash
wails dev
```

This launches the app with hot-reload for frontend changes.

### Production Build

```bash
# Build for Windows (from any platform with GOOS=windows)
GOOS=windows GOARCH=amd64 wails build
```

The output binary is at `build/bin/Super CC.exe`.

### Build Options

```bash
# Default build (Windows)
wails build

# With specific output name
wails build -o "SuperCC.exe"

# Debug build with DevTools enabled
wails build -debug
```

## Project Structure

```
super-cc/
├── app.go                  # Wails bindings (Go -> JS bridge)
├── main.go                 # Wails app entry point
├── app/
│   ├── commands/
│   │   ├── scanner.go      # Installation detection
│   │   ├── downloader.go   # Version database & download
│   │   ├── protector.go    # Hosts file lock management
│   │   ├── crack.go        # DLL patching
│   │   ├── cleaner.go      # Residual file cleanup
│   │   └── debug.go        # Environment diagnostics
│   └── config/
│       └── config.go       # Config file management (~/.super-cc/config.json)
├── frontend/
│   ├── index.html          # App layout
│   ├── styles.css          # macOS 26 Tahoe design system
│   ├── app.js              # Frontend controller
│   ├── fonts/              # Phosphor Icons (bundled)
│   └── wailsjs/            # Auto-generated Wails bindings
├── build/
│   ├── appicon.png         # Application icon
│   ├── bin/                # Build output
│   └── windows/
│       ├── icon.ico        # Windows icon
│       ├── info.json       # Windows version info
│       └── wails.exe.manifest
└── wails.json              # Wails project config
```

## Configuration

App config is stored at `~/.super-cc/config.json`:

```json
{
  "install_path": "C:\\Program Files\\Editor",
  "auto_lock": true,
  "auto_check": true,
  "locked": false,
  "locked_version": "",
  "patch_applied": false
}
```

## Testing

```bash
go test ./... -v
```

All tests run on Linux/macOS CI (Windows-specific commands are tested with expected non-Windows errors).

## Tech Stack

- **Backend**: Go 1.21+
- **Frontend**: Vanilla JS, CSS
- **Desktop**: Wails v2 with WebView2

## Disclaimer

This software is provided **for educational and research purposes only**. The author does not endorse or encourage the violation of any software's terms of service, end-user license agreement (EULA), or applicable laws.

By using this software, you acknowledge that:

- You are solely responsible for how you use this tool.
- Modifying third-party software may violate its terms of service.
- The author assumes no liability for any consequences arising from the use of this software.
- This project is not affiliated with, endorsed by, or connected to any third-party software referenced herein.

**Use at your own risk.**

## Credit

Version download URLs and DLL pattern references sourced from [CapCut Project](https://iosgods.com/topic/190383-capcut-pro-crack-pc-open-source/).

## License

MIT License — see [LICENSE](LICENSE) for details.

Commercial redistribution of the compiled binary is not permitted without prior written consent.
