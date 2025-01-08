package tea

import (
	"github.com/charmbracelet/lipgloss"
)

func (m model) TaskScreenView() string {
	sqr := lipgloss.NewStyle().
		Width(60).
		Height(20).
		Border(lipgloss.RoundedBorder()).
		Render(lipgloss.NewStyle().Bold(true).Padding(0, 1).Render("pending_") + "|" + lipgloss.NewStyle().Padding(0, 1).Render("all"))

	return lipgloss.Place(m.viewportWidth, m.viewportHeight, lipgloss.Center, lipgloss.Center, sqr)
}
