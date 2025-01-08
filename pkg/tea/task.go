package tea

import (
	"github.com/charmbracelet/lipgloss"
)

func (m model) TaskScreenView() string {
	tabs := lipgloss.NewStyle().Bold(true).Padding(0, 1).Render("pending_") + "|" + lipgloss.NewStyle().Padding(0, 1).Render("all")

	sqr := lipgloss.NewStyle().
		Width(60).
		Height(20).
		Border(lipgloss.RoundedBorder()).
		Render("Content")

	return lipgloss.Place(m.viewportWidth, m.viewportHeight, lipgloss.Center, lipgloss.Center, tabs+"\n"+sqr)
}
