package tui

import "github.com/charmbracelet/lipgloss"

var (
	promptCardStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("#6366F1")).
			Padding(1, 3).
			Width(52)

	promptTitleStyle = lipgloss.NewStyle().
				Bold(true).
				Foreground(lipgloss.Color("#A5B4FC"))

	promptLabelStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("#D1D5DB"))

	spinnerLabelStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("#A5B4FC"))
	errorStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#EF4444"))
)
