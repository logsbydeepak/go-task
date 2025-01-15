package tea

import (
	"strconv"
	"strings"

	"example.com/pkg/db"
	"example.com/pkg/task"
	bt "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type taskScreenState struct {
	tabs   []string
	tasks  []task.Task
	active int
}

func (m model) TaskScreenSwitch() (bt.Model, bt.Cmd) {
	m.screen = taskScreen
	m.taskScreenState = taskScreenState{
		tabs: []string{"pending", "all"},
	}
	m.updateContent()
	return m, nil
}

const (
	GrayColor = lipgloss.Color("8")
)

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

	outerWidth := 60
	outerHeight := 20

	if m.viewportWidth <= outerWidth {
		outerWidth = m.viewportWidth - 2
	}

	if m.viewportHeight <= outerHeight {
		outerHeight = m.viewportHeight - 2
	}

	innerWidth := outerWidth - 2
	innerHeight := outerHeight - 3
	tasks := m.taskScreenState.tasks

	if len(tasks) > innerHeight {
		tasks = tasks[:innerHeight]
	}

	zeroCount := len(strconv.Itoa(int(tasks[len(tasks)-1].ID)))
	content := make([]string, len(tasks))

	for i, each := range tasks {
		id := strconv.Itoa(int(each.ID))
		idLen := len(id)
		if zeroCount > idLen {
			id = strings.Repeat(" ", zeroCount-idLen) + id
		}
		content[i] = lipgloss.NewStyle().Foreground(GrayColor).Render(id) + " " + each.Description
	}

	innersqr := lipgloss.NewStyle().
		Width(innerWidth).
		Height(innerHeight).
		Border(lipgloss.RoundedBorder()).
		Padding(0, 1).
		Render(strings.Join(content, "\n"))

	outersqr := lipgloss.NewStyle().Width(outerWidth).
		Height(outerHeight).
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
		switch msg.Type {
		case bt.KeyLeft:
			if m.active > 0 {
				m.taskScreenState.active--
				m.updateContent()
			}
		case bt.KeyRight:
			if m.active < len(m.tabs)-1 {
				m.taskScreenState.active++
				m.updateContent()
			}
		case bt.KeyTab:
			if m.active < len(m.tabs)-1 {
				m.taskScreenState.active++
				m.updateContent()
			} else {
				m.taskScreenState.active--
				m.updateContent()
			}
		case bt.KeyShiftTab:
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
		if err == nil {
			m.tasks = task
		}
	case 1:
		task, err := db.GetAllTask()
		if err == nil {
			m.tasks = task
		}
	}
}
