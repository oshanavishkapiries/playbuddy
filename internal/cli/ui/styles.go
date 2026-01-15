package ui

import "github.com/charmbracelet/lipgloss"

var (
	// Colors
	Green = lipgloss.Color("#00FF41")
	Gray  = lipgloss.Color("#888888")

	// Styles
	TitleStyle = lipgloss.NewStyle().
			Foreground(Green).
			Bold(true).
			Padding(0, 1).
			MarginBottom(1)

	BorderStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(Green).
			Padding(1, 2).
			MarginLeft(2)

	DescStyle = lipgloss.NewStyle().
			Foreground(Gray).
			Italic(true)

	KeyStyle = lipgloss.NewStyle().
			Foreground(Gray).
			Width(10).
			Align(lipgloss.Right)

	InputBoxStyle = lipgloss.NewStyle().
			Border(lipgloss.NormalBorder(), true).
			BorderForeground(Green).
			Padding(0, 1)
)
