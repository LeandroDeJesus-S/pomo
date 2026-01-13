# ⏳ Pomo

A beautiful and efficient Pomodoro timer for the terminal, built with [Bubbletea](https://github.com/charmbracelet/bubbletea) and [Lipgloss](https://github.com/charmbracelet/lipgloss). Features a modern TUI interface with color-coded sessions, ASCII art timer, progress tracking, and system notifications.

## ✨ Features

- **Session Types**: Focus Time (Study), Short Break, Long Break with distinct color labels
- **Visual Timer**: Large ASCII art countdown with progress bar
- **Statistics**: Track completed sessions, total study time, and long break cycles
- **In-App Configuration**: Edit session durations without restarting
- **Help System**: Interactive keyboard shortcuts and config editor
- **Notifications**: Desktop notifications at session transitions
- **Cross-Platform**: Works on Linux, macOS, and Windows (with terminal support)

## 📦 Installation

### Option 1: Go Install (Recommended)
```bash
go install github.com/LeandroDeJesus-S/pomo@v0.1.0
```

### Option 2: Build from Source
```bash
git clone https://github.com/LeandroDeJesus-S/pomo.git
cd pomo
go build
```

**Prerequisites**: Go 1.25.4 or later.

## 🚀 Usage

### Basic Timer
```bash
pomo
# Uses default durations: 25 min study, 5 min break, 15 min long break
```

### Custom Durations
```bash
pomo -study 30 -break 10 -lbreak 20
```

### Keyboard Shortcuts
- `SPACE` / `P`: Pause/resume timer
- `N`: Skip to next session
- `R`: Reset current session
- `+` / `-`: Adjust time by 1 minute (when not paused)
- `?`: Toggle help menu
- `ESC`: Close help/config (if open)
- `Q`: Quit

### In-App Configuration
1. Press `?` to open help
2. Press `C` to enter config mode
3. Use `S`, `B`, `L` to select session type
4. Use `+` / `-` to adjust duration
5. Press `C` or `ESC` to exit config

### Debug Mode
Set the environment variable for detailed logging:
```bash
POMO_DEBUG=1 pomo
```
Logs are saved to `debug.log`.

## 🛠 Development

**Dependencies**:
- [Bubbletea](https://github.com/charmbracelet/bubbletea) - TUI framework
- [Lipgloss](https://github.com/charmbracelet/lipgloss) - Styling
- [beeep](https://github.com/gen2brain/beeep) - Notifications
- [bubbletea-overlay](https://github.com/rmhubbert/bubbletea-overlay) - Overlay components

**Building**:
```bash
go mod tidy
go build
```

**Contributing**: PRs welcome! Please test on multiple terminals.

## 📄 License

MIT License - see [LICENSE](LICENSE) for details.
