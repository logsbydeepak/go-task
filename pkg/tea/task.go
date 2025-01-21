package tea

import (
	"strconv"
	"strings"

	"example.com/pkg/db"
	"example.com/pkg/task"
	"github.com/charmbracelet/bubbles/textinput"
	bt "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type taskScreenState struct {
	tabs      []string
	tasks     []task.Task
	active    int
	taskInput textinput.Model

	outerWidth  int
	outerHeight int
}

const (
	maxWidthSize        = 60
	maxHeightSize       = 20
	borderLeftRightSize = 2
	taskInputPromptSize = 5
	taskInputHeightSize = 1

	GrayColor = lipgloss.Color("8")
)

func (m model) TaskScreenSwitch() (bt.Model, bt.Cmd) {
	ti := textinput.New()
	ti.Placeholder = "new task"
	ti.CharLimit = 156
	ti.PlaceholderStyle = lipgloss.NewStyle().Italic(true).Foreground(GrayColor)

	m.screen = taskScreen
	m.taskScreenState = taskScreenState{
		tabs:      []string{"pending", "all"},
		taskInput: ti,

		outerWidth:  maxWidthSize,
		outerHeight: maxHeightSize,
	}
	m.updateContent()

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

	innerWidth := m.taskScreenState.outerWidth - borderLeftRightSize
	innerHeight := m.taskScreenState.outerHeight - borderLeftRightSize - taskInputHeightSize
	tasks := m.taskScreenState.tasks

	if len(tasks) > innerHeight-taskInputHeightSize {
		tasks = tasks[:innerHeight]
	}

	zeroCount := len(strconv.Itoa(int(tasks[0].ID)))
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
		Render(m.taskScreenState.taskInput.View() + "\n" + strings.Join(content, "\n"))

	outersqr := lipgloss.NewStyle().Width(m.taskScreenState.outerWidth).
		Height(m.taskScreenState.outerHeight).
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
	if m.viewportWidth <= maxWidthSize {
		m.taskScreenState.outerWidth = m.viewportWidth - borderLeftRightSize
	} else {
		m.taskScreenState.outerWidth = maxWidthSize
	}

	if m.viewportHeight <= maxHeightSize {
		m.taskScreenState.outerHeight = m.viewportHeight - borderLeftRightSize
	} else {
		m.taskScreenState.outerHeight = maxHeightSize
	}

	m.taskScreenState.taskInput.Width = m.taskScreenState.outerWidth - borderLeftRightSize - taskInputPromptSize

	switch msg := msg.(type) {
	case bt.KeyMsg:
		switch msg.Type {
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
		case bt.KeyEscape:
			m.ignoreQKey = false
			m.taskScreenState.taskInput.Reset()
			m.taskScreenState.taskInput.Blur()
			return m, nil
		case bt.KeyRunes:
			if msg.String() == "a" {
				m.ignoreQKey = true
				return m, m.taskScreenState.taskInput.Focus()
			}
		}
	}

	var cmd bt.Cmd
	m.taskScreenState.taskInput, cmd = m.taskScreenState.taskInput.Update(msg)
	return m, cmd
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
