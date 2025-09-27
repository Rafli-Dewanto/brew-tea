# Brew-Tea Package Manager

A user-friendly Terminal User Interface (TUI) for managing Homebrew packages on macOS, built with [Bubble Tea](https://github.com/charmbracelet/bubbletea) and [Lip Gloss](https://github.com/charmbracelet/lipgloss).

```
   ( (
    ) )
  ........
  |      |]
  \      /
   '----'
🍵 Brew-Tea Package Manager
```

## Features

- **Formula & Cask Support**: Manage both command-line tools (formulae) and GUI applications (casks)
- **Search & Install**: Search for packages and install them with ease
- **View Installed**: Browse and uninstall your installed packages
- **Real-time Feedback**: Loading indicators and status messages keep you informed

## Requirements

- macOS with [Homebrew](https://brew.sh/) installed
- Go 1.19 or later

## Installation

### Quick Install (Recommended)

Install Brew-Tea using Go:

```bash
go install github.com/Rafli-Dewanto/brew-tea@latest
```

Make sure `$GOPATH/bin` (or `$GOBIN`) is in your PATH.

### Manual Installation

1. Clone or download this repository
2. Build the application:
   ```bash
   go build -o brew-tea main.go
   ```
3. Run the application:
   ```bash
   ./brew-tea
   ```

## Usage

### Main Menu

- `1` - View installed packages
- `2` - Search and install packages
- `q` - Quit

### Installed Packages View

- `u` - Uninstall selected package
- `b` - Back to main menu
- `q` - Quit

### Search View

- Type your search query
- `Enter` - Perform search
- `Ctrl+i` - Install selected package
- `b` - Back to main menu
- `q` - Quit

## Key Bindings

| Action    | Key             |
| --------- | --------------- |
| Quit      | `q` or `Ctrl+c` |
| Back      | `b`             |
| Search    | `Enter`         |
| Install   | `Ctrl+i`        |
| Uninstall | `u`             |

## Dependencies

This project uses the following Go modules:

- `github.com/charmbracelet/bubbletea`
- `github.com/charmbracelet/bubbles`
- `github.com/charmbracelet/lipgloss`

## Contributing

Feel free to open issues or submit pull requests to improve Brew-Tea!

## License

This project is open source. Feel free to use and modify as needed.
