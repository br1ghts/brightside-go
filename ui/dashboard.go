package ui

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// Terminal UI Model
type model struct {
	title    string
	selected int
	options  []string
}

// Initialize the UI
func (m model) Init() tea.Cmd {
	return nil
}

// Handle Keyboard Input
func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q": // Quit
			fmt.Println("\n👋 Jack: See ya later, boss.")
			return m, tea.Quit
		case "down": // Move down
			if m.selected < len(m.options)-1 {
				m.selected++
			}
		case "up": // Move up
			if m.selected > 0 {
				m.selected--
			}
		case "enter": // Select option
			fmt.Println("\n✅ Jack: Executing", m.options[m.selected])
		}
	}
	return m, nil
}

// Render UI
func (m model) View() string {
	title := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#ff66b2")).
		Render("🚀 Brightside Jack - AI Terminal Dashboard")

	optionsList := ""
	for i, option := range m.options {
		cursor := "  "
		if i == m.selected {
			cursor = "👉"
		}
		optionsList += fmt.Sprintf("%s %s\n", cursor, option)
	}

	return fmt.Sprintf("%s\n\n%s\n\n%s", title, "Use ↑ ↓ to navigate, Enter to select, Q to quit.", optionsList)
}

// Start the UI
func StartDashboard() {
	model := model{
		title:    "Brightside Jack",
		selected: 0,
		options:  []string{"📡 Live Twitch Chat", "🖥 System Stats", "📰 News Feeds", "🤖 Jack AI", "❌ Exit"},
	}

	if _, err := tea.NewProgram(model).Run(); err != nil {
		fmt.Fprintf(os.Stderr, "❌ Error launching Brightside Jack: %v\n", err)
		os.Exit(1)
	}
}
