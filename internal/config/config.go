package config

import "github.com/charmbracelet/lipgloss"

var (
	ErrorStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("9"))
	InfoStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("12"))
	SuccessStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#04B575"))
)
