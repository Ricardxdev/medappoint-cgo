package utils

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
)

func AlignW(s string, width int) string {
	aligned := lipgloss.PlaceHorizontal(
		width,
		lipgloss.Center,
		s,
	)
	return aligned
}

func Center(s string, width int, height int) string {
	alignedH := lipgloss.PlaceHorizontal(
		width,
		lipgloss.Center,
		s,
	)
	centered := lipgloss.PlaceVertical(
		height,
		lipgloss.Center,
		alignedH,
	)
	return centered
}

func BreadcrumbView(breadcrumb []string) string {
	colors := []lipgloss.Color{
		lipgloss.Color("#FFD700"), // gold
		lipgloss.Color("#00BFFF"), // deep sky blue
		lipgloss.Color("#32CD32"), // lime green
		lipgloss.Color("#FF69B4"), // hot pink
		lipgloss.Color("#FFA500"), // orange
	}
	var styled []string
	for i, item := range breadcrumb {
		style := lipgloss.NewStyle().Foreground(colors[i%len(colors)])
		styled = append(styled, style.Render(item))
	}
	return strings.Join(styled, " > ")
}
