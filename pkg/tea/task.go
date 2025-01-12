package tea

import (
	"strings"

	"example.com/pkg/db"
	"example.com/pkg/task"
	bt "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type taskScreenState struct {
	tabs    []string
	content string
	active  int
}

func (m model) TaskScreenSwitch() (bt.Model, bt.Cmd) {
	m.screen = taskScreen
	m.taskScreenState = taskScreenState{
		tabs: []string{"pending", "all"},
	}
	return m, nil
}

func (m model) TaskScreenView() string {
	tabStyle := lipgloss.NewStyle().Padding(0, 1)

	var tabs strings.Builder
	for i, tab := range m.tabs {
		if i == m.active {
			tabs.WriteString(tabStyle.Bold(true).Render(tab + "_"))
		} else {
			tabs.WriteString(tabStyle.Render(tab + " "))
		}
	}

	width := 60
	height := 20

	innersqr := lipgloss.NewStyle().
		Width(width-2).
		Height(height-2).
		Border(lipgloss.RoundedBorder()).
		Padding(0, 1).
		Render(m.taskScreenState.content)

	outersqr := lipgloss.NewStyle().Width(width).
		Height(height).
		Border(lipgloss.HiddenBorder()).
		Render(tabs.String() + "\n" + innersqr)

	return lipgloss.Place(
		m.viewportWidth,
		m.viewportHeight,
		lipgloss.Center,
		lipgloss.Center,
		outersqr,
	)
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
		case bt.KeyTab.String():
			if m.active < len(m.tabs)-1 {
				m.taskScreenState.active++
				m.updateContent()
			} else {
				m.taskScreenState.active--
				m.updateContent()
			}
		case bt.KeyShiftTab.String():
			if m.active > 0 {
				m.taskScreenState.active--
				m.updateContent()
			} else {
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
		task, err := db.GetAllPendingTask()
		if err != nil {
			m.content = "ERROR"
		} else {
			m.content = showTask(task)
		}
	case 1:
		task, err := db.GetAllTask()
		if err != nil {
			m.content = "ERROR"
		} else {
			m.content = showTask(task)
		}
	}
}

func showTask(t []task.Task) string {
	var result strings.Builder

	for _, task := range t {
		result.WriteString(task.Description + "\n")
	}

	return result.String()
}
