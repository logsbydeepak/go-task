package tea

import (
	"strings"

	bt "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type taskScreenState struct {
	tabs    []string
	content string
	active  int
}

var (
	tabStyle       = lipgloss.NewStyle().Padding(0, 1)
	activeTabStyle = lipgloss.NewStyle().Bold(true).Padding(0, 1)
)

func (m model) TaskScreenView() string {
	var tabs strings.Builder

	for i, tab := range m.tabs {
		if i == m.active {
			tabs.WriteString(activeTabStyle.Render(tab + "_"))
		} else {
			tabs.WriteString(tabStyle.Render(tab + " "))
		}
	}

	width := 60
	height := 20

	innersqr := lipgloss.NewStyle().
		Width(width - 2).
		Height(height - 2).
		Border(lipgloss.RoundedBorder()).
		Render(m.taskScreenState.content)

	outersqr := lipgloss.NewStyle().Width(width).
		Height(height).
		Border(lipgloss.HiddenBorder()).
		Render(tabs.String() + "\n" + innersqr)

	return lipgloss.Place(m.viewportWidth, m.viewportHeight, lipgloss.Center, lipgloss.Center, outersqr)
}

func (m model) TaskScrrenUpdate(msg bt.Msg) (bt.Model, bt.Cmd) {
	switch msg := msg.(type) {
	case bt.KeyMsg:
		switch msg.String() {
		case bt.KeyLeft.String():
			if m.active > 0 {
				m.taskScreenState.active--
				m.updateContent()
			}
		case bt.KeyRight.String():
			if m.active < len(m.tabs)-1 {
				m.taskScreenState.active++
				m.updateContent()
			}
		}
	}

	return m, nil
}

func (m *model) updateContent() {
	switch m.active {
	case 0:
		m.content = "pending"
	case 1:
		m.content = "all"
	}
}
