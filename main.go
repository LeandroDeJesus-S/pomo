package main

import (
	"flag"
	"fmt"
	"os"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/gen2brain/beeep"
	overlay "github.com/rmhubbert/bubbletea-overlay"
)

type SessionType int

const (
	Study SessionType = iota
	Break
	LongBreak
	UnknownSession

	DefaultStudyMins  = 25
	DefaultBreakMins  = 5
	DefaultLBreakMins = 15
)

func (s SessionType) String() string {
	switch s {
	case Study:
		return "FOCUS TIME"
	case Break:
		return "SHORT BREAK"
	case LongBreak:
		return "LONG BREAK"
	default:
		return "IDLE"
	}
}

// Color palette
var (
	purpleColor = lipgloss.Color("#A855F7")
	mutedColor  = lipgloss.Color("#6B7280")
	accentColor = lipgloss.Color("#EC4899")
	greenColor  = lipgloss.Color("#10B981")
)

var (
	// Giant purple label
	labelStyle = lipgloss.NewStyle().
			Foreground(purpleColor).
			Bold(true).
			Align(lipgloss.Center).
			MarginBottom(2)

	// Timer style
	timerStyle = lipgloss.NewStyle().
			Foreground(purpleColor).
			Align(lipgloss.Center).
			MarginBottom(2)

	// Progress bar container
	progressStyle = lipgloss.NewStyle().
			Align(lipgloss.Center).
			MarginBottom(3)

	progressFilledStyle = lipgloss.NewStyle().
				Foreground(purpleColor)

	progressEmptyStyle = lipgloss.NewStyle().
				Foreground(mutedColor)

	// Minimalist stats
	statsStyle = lipgloss.NewStyle().
			Align(lipgloss.Center).
			Foreground(mutedColor).
			MarginBottom(2)

	statValueStyle = lipgloss.NewStyle().
			Foreground(purpleColor).
			Bold(true)

	// Pause indicator
	pauseStyle = lipgloss.NewStyle().
			Foreground(accentColor).
			Bold(true).
			Align(lipgloss.Center).
			MarginTop(1).
			MarginBottom(1)

	// Help section
	helpStyle = lipgloss.NewStyle().
			Align(lipgloss.Center).
			Foreground(mutedColor).
			MarginTop(2)

	helpKeyStyle = lipgloss.NewStyle().
			Foreground(purpleColor).
			Bold(true)

	// Floating help window
	floatingHelpStyle = lipgloss.NewStyle().
				Border(lipgloss.RoundedBorder()).
				BorderForeground(purpleColor).
				Padding(2, 4).
		// Align(lipgloss.Center).
		Foreground(mutedColor)

	helpTitleStyle = lipgloss.NewStyle().
			Foreground(purpleColor).
			Bold(true).
			Align(lipgloss.Center).
			MarginBottom(1)

	// Config editing
	configStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(purpleColor).
			Padding(2, 4).
		// Align(lipgloss.Center).
		Foreground(mutedColor)
)

type Config struct {
	StudyDuration     time.Duration
	BreakDuration     time.Duration
	LongBreakDuration time.Duration
}

type Stats struct {
	NumStudySessions          int
	NumStudySessionsCompleted int
	TotalStudyTime            time.Duration
	SessionsUntilLongBreak    int
}

type Model struct {
	config             Config
	currentSession     SessionType
	initialDuration    time.Duration
	timeLeft           time.Duration
	paused             bool
	quitting           bool
	stats              Stats
	helpVisible        bool
	editingConfig      bool
	editingConfigField SessionType
	width              int
	height             int
}

func initialModel(cfg Config) Model {
	return Model{
		config:          cfg,
		currentSession:  Study,
		initialDuration: cfg.StudyDuration,
		timeLeft:        cfg.StudyDuration,
		stats: Stats{
			SessionsUntilLongBreak: 4,
		},
	}
}

func (m Model) Init() tea.Cmd {
	return tickCmd()
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

	case tea.KeyMsg:
		key := msg.String()

		if key == "q" || key == "ctrl+c" {
			m.quitting = true
			return m, tea.Quit
		}

		if key == "?" {
			m.helpVisible = !m.helpVisible
			if !m.helpVisible {
				m.editingConfig = false
				m.editingConfigField = UnknownSession
			}
			return m, nil
		}

		if key == "esc" && !m.editingConfig {
			if m.helpVisible {
				m.helpVisible = false
			}
			return m, nil
		}

		if m.helpVisible && m.editingConfig {
			switch key {
			case "c", "esc":
				m.editingConfig = false
				m.editingConfigField = UnknownSession
			case "s":
				m.editingConfigField = Study
			case "b":
				m.editingConfigField = Break
			case "l":
				m.editingConfigField = LongBreak
			case "+":
				if m.editingConfigField != UnknownSession {
					m.updateConfigDuration(m.editingConfigField, time.Minute)
				}
			case "-":
				if m.editingConfigField != UnknownSession {
					m.updateConfigDuration(m.editingConfigField, -time.Minute)
				}
			}
			return m, nil
		}

		if m.helpVisible && key == "c" {
			m.editingConfig = true
			return m, nil
		}

		switch key {
		case "p", " ":
			m.paused = !m.paused
			if !m.paused {
				return m, tickCmd()
			}
		case "+":
			if !m.paused {
				m.timeLeft += time.Minute
			}
		case "-":
			if !m.paused && m.timeLeft > time.Minute {
				m.timeLeft -= time.Minute
			}
		case "n":
			return m, m.nextSessionCmd()
		case "r":
			m.timeLeft = m.initialDuration
			return m, tickCmd()
		}

	case tickMsg:
		if m.paused || m.quitting {
			return m, nil
		}

		m.timeLeft -= time.Second
		if m.timeLeft <= 0 {
			return m, m.nextSessionCmd()
		}
		return m, tickCmd()
	}

	return m, nil
}

func (m *Model) updateConfigDuration(sessionType SessionType, change time.Duration) {
	switch sessionType {
	case Study:
		if m.config.StudyDuration+change >= time.Minute {
			m.config.StudyDuration += change
			if m.currentSession == Study {
				m.initialDuration = m.config.StudyDuration
			}
		}
	case Break:
		if m.config.BreakDuration+change >= time.Minute {
			m.config.BreakDuration += change
			if m.currentSession == Break {
				m.initialDuration = m.config.BreakDuration
			}
		}
	case LongBreak:
		if m.config.LongBreakDuration+change >= time.Minute {
			m.config.LongBreakDuration += change
			if m.currentSession == LongBreak {
				m.initialDuration = m.config.LongBreakDuration
			}
		}
	}
}

func (m Model) View() string {
	if m.quitting {
		goodbye := lipgloss.NewStyle().
			Foreground(purpleColor).
			Bold(true).
			Align(lipgloss.Center).
			Padding(3).
			Render("Thanks for staying focused!\n\nSee you next time 👋")

		return lipgloss.Place(m.width, m.height,
			lipgloss.Center, lipgloss.Center,
			goodbye)
	}

	// Render main content
	mainContent := m.renderMainContent()

	mainView := lipgloss.Place(m.width, m.height,
		lipgloss.Center, lipgloss.Center,
		mainContent)

	// If help is visible, render overlay on top
	if m.helpVisible {
		var helpContent string
		if m.editingConfig {
			helpContent = m.renderConfigHelp()
		} else {
			helpContent = m.renderHelp()
		}

		return overlay.Composite(helpContent, mainView, overlay.Center, overlay.Center, 0, 0)
	}

	return mainView
}

func (m Model) renderMainContent() string {
	var content strings.Builder

	// 1. Big purple label centered
	content.WriteString(m.renderBigLabel())
	content.WriteString("\n\n")

	// 2. Giant timer centered
	content.WriteString(m.renderGiantTimer())
	content.WriteString("\n")

	// Pause indicator if paused
	if m.paused {
		content.WriteString(pauseStyle.Render("⏸  PAUSED"))
		content.WriteString("\n")
	}

	// 3. Progress bar below timer
	content.WriteString(m.renderProgressBar())
	content.WriteString("\n")

	// 4. Minimalist stats below bar
	content.WriteString(m.renderMinimalistStats())
	content.WriteString("\n")

	// 5. Help hint at the end
	if !m.helpVisible {
		content.WriteString(helpStyle.Render("Press ? for help"))
	}

	return content.String()
}

func (m Model) renderBigLabel() string {
	label := m.currentSession.String()

	// Create ASCII art style label
	art := []string{
		"╔" + strings.Repeat("═", len(label)+2) + "╗",
		"║ " + label + " ║",
		"╚" + strings.Repeat("═", len(label)+2) + "╝",
	}

	return labelStyle.Render(strings.Join(art, "\n"))
}

func (m Model) renderGiantTimer() string {
	minutes := int(m.timeLeft.Minutes())
	seconds := int(m.timeLeft.Seconds()) % 60
	timeStr := fmt.Sprintf("%02d:%02d", minutes, seconds)

	return timerStyle.Render(m.renderASCIITime(timeStr))
}

func (m Model) renderASCIITime(timeStr string) string {
	digits := make([]rune, 0)
	for _, ch := range timeStr {
		digits = append(digits, ch)
	}

	digitMap := map[rune][]string{
		'0': {
			" ██████╗ ",
			"██╔═══██╗",
			"██║   ██║",
			"██║   ██║",
			"╚██████╔╝",
			" ╚═════╝ ",
		},
		'1': {
			" ██╗",
			"███║",
			"╚██║",
			" ██║",
			" ██║",
			" ╚═╝",
		},
		'2': {
			"██████╗ ",
			"╚════██╗",
			" █████╔╝",
			"██╔═══╝ ",
			"███████╗",
			"╚══════╝",
		},
		'3': {
			"██████╗ ",
			"╚════██╗",
			" █████╔╝",
			" ╚═══██╗",
			"██████╔╝",
			"╚═════╝ ",
		},
		'4': {
			"██╗  ██╗",
			"██║  ██║",
			"███████║",
			"╚════██║",
			"     ██║",
			"     ╚═╝",
		},
		'5': {
			"███████╗",
			"██╔════╝",
			"███████╗",
			"╚════██║",
			"███████║",
			"╚══════╝",
		},
		'6': {
			" ██████╗ ",
			"██╔════╝ ",
			"███████╗ ",
			"██╔═══██╗",
			"╚██████╔╝",
			" ╚═════╝ ",
		},
		'7': {
			"███████╗",
			"╚════██║",
			"    ██╔╝",
			"   ██╔╝ ",
			"   ██║  ",
			"   ╚═╝  ",
		},
		'8': {
			" ██████╗ ",
			"██╔═══██╗",
			"╚█████╔╝ ",
			"██╔═══██╗",
			"╚██████╔╝",
			" ╚═════╝ ",
		},
		'9': {
			" ██████╗ ",
			"██╔═══██╗",
			"╚██████╔╝",
			" ╚═══██║ ",
			" █████╔╝ ",
			" ╚════╝  ",
		},
		':': {
			"   ",
			"██╗",
			"╚═╝",
			"██╗",
			"╚═╝",
			"   ",
		},
	}

	var result []string
	for i := range 6 {
		var line strings.Builder
		for _, digit := range digits {
			if art, ok := digitMap[digit]; ok {
				line.WriteString(art[i])
				line.WriteString(" ")
			}
		}
		result = append(result, line.String())
	}

	return strings.Join(result, "\n")
}

func (m Model) renderProgressBar() string {
	width := 50
	elapsed := m.initialDuration - m.timeLeft
	progress := float64(elapsed) / float64(m.initialDuration)
	if progress > 1 {
		progress = 1
	}
	if progress < 0 {
		progress = 0
	}

	filled := int(float64(width) * progress)
	empty := width - filled

	bar := progressFilledStyle.Render(strings.Repeat("━", filled)) +
		progressEmptyStyle.Render(strings.Repeat("━", empty))

	percentage := int(progress * 100)

	return progressStyle.Render(fmt.Sprintf("%s  %d%%", bar, percentage))
}

func (m Model) renderMinimalistStats() string {
	var stats strings.Builder

	// Single line, minimal design
	stats.WriteString(statsStyle.Render(fmt.Sprintf(
		"🍅 %s  •  ⏱️  %s  •  📊 %s until long break",
		statValueStyle.Render(fmt.Sprintf("%d", m.stats.NumStudySessions)),
		statValueStyle.Render(formatDuration(m.stats.TotalStudyTime)),
		statValueStyle.Render(fmt.Sprintf("%d", m.stats.SessionsUntilLongBreak)),
	)))

	return stats.String()
}

func (m Model) renderHelp() string {
	var help strings.Builder

	help.WriteString(helpTitleStyle.Render("KEYBOARD SHORTCUTS"))
	help.WriteString("\n\n")

	shortcuts := []string{
		"SPACE / P  pause/resume",
		"N  next session",
		"R  reset timer",
		"+/-  adjust time",
		"C  configure",
		"?  close help",
		"Q  quit",
	}

	for _, s := range shortcuts {
		parts := strings.SplitN(s, "  ", 2)
		help.WriteString(
			helpKeyStyle.Render(fmt.Sprintf("%-15s", parts[0])) +
				"  " + parts[1] + "\n",
		)
	}

	return floatingHelpStyle.Render(help.String())
}

func (m Model) renderConfigHelp() string {
	var cfg strings.Builder

	cfg.WriteString(helpTitleStyle.Render("CONFIGURATION"))
	cfg.WriteString("\n\n")

	sessions := []struct {
		key      string
		sType    SessionType
		duration time.Duration
	}{
		{"S", Study, m.config.StudyDuration},
		{"B", Break, m.config.BreakDuration},
		{"L", LongBreak, m.config.LongBreakDuration},
	}

	for _, s := range sessions {
		selected := " "
		if m.editingConfigField == s.sType {
			selected = "▶"
		}
		cfg.WriteString(fmt.Sprintf(
			"%s %s  %-15s  %s\n",
			selected,
			helpKeyStyle.Render(s.key),
			s.sType.String(),
			statValueStyle.Render(fmt.Sprintf("%d min", int(s.duration.Minutes()))),
		))
	}

	cfg.WriteString("\n")
	cfg.WriteString(lipgloss.NewStyle().Foreground(mutedColor).Italic(true).
		Render("Use +/- to adjust  •  C to exit"))

	return configStyle.Render(cfg.String())
}

func formatDuration(d time.Duration) string {
	hours := int(d.Hours())
	minutes := int(d.Minutes()) % 60
	if hours > 0 {
		return fmt.Sprintf("%dh%dm", hours, minutes)
	}
	return fmt.Sprintf("%dm", minutes)
}

type tickMsg time.Time

func tickCmd() tea.Cmd {
	return tea.Tick(time.Second, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}

func sendNotificationCmd(title, body string) tea.Cmd {
	return func() tea.Msg {
		_ = beeep.Notify(title, body, "")
		return nil
	}
}

func (m *Model) nextSessionCmd() tea.Cmd {
	var notificationTitle, notificationBody string

	if m.currentSession == Study {
		m.stats.NumStudySessions++
		m.stats.NumStudySessionsCompleted++
		m.stats.TotalStudyTime += m.config.StudyDuration
		m.stats.SessionsUntilLongBreak--
	}

	switch m.currentSession {
	case Study:
		if m.stats.SessionsUntilLongBreak <= 0 {
			m.currentSession = LongBreak
			m.initialDuration = m.config.LongBreakDuration
			m.timeLeft = m.config.LongBreakDuration
			m.stats.SessionsUntilLongBreak = 4
			notificationTitle = "🎯 Focus Session Complete"
			notificationBody = "Time for a long break!"
		} else {
			m.currentSession = Break
			m.initialDuration = m.config.BreakDuration
			m.timeLeft = m.config.BreakDuration
			notificationTitle = "🎯 Focus Session Complete"
			notificationBody = "Time for a short break."
		}
	case Break:
		m.currentSession = Study
		m.initialDuration = m.config.StudyDuration
		m.timeLeft = m.config.StudyDuration
		notificationTitle = "☕ Break Over"
		notificationBody = "Ready to focus again?"
	case LongBreak:
		m.currentSession = Study
		m.initialDuration = m.config.StudyDuration
		m.timeLeft = m.config.StudyDuration
		notificationTitle = "🌟 Long Break Over"
		notificationBody = "Let's get back to work!"
	}

	return tea.Batch(tickCmd(), sendNotificationCmd(notificationTitle, notificationBody))
}

func main() {
	studyMins := flag.Int("study", DefaultStudyMins, "Duration of the study session in minutes")
	breakMins := flag.Int("break", DefaultBreakMins, "Duration of the short break in minutes")
	longBreakMins := flag.Int("lbreak", DefaultLBreakMins, "Duration of the long break in minutes")

	flag.Parse()

	cfg := Config{
		StudyDuration:     time.Duration(*studyMins) * time.Minute,
		BreakDuration:     time.Duration(*breakMins) * time.Minute,
		LongBreakDuration: time.Duration(*longBreakMins) * time.Minute,
	}

	p := tea.NewProgram(initialModel(cfg), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
}
