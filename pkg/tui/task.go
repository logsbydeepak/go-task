package tui

import (
	"strconv"
	"strings"

	"example.com/pkg/db"
	"example.com/pkg/task"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
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

func (m model) TaskScreenSwitch() (tea.Model, tea.Cmd) {
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

	validDescriptionWidth := innerWidth - zeroCount - 3 // 1+2=3, 1 for space and 2 for left and right padding
	for i, task := range tasks {
		id := strconv.Itoa(int(task.ID))
		idLen := len(id)
		if zeroCount > idLen {
			id = strings.Repeat(" ", zeroCount-idLen) + id
		}

		description := task.Description
		if task.IsComplete {
			description += "✓"
		} else {
			description += "⨯"
		}
		if len(description) > validDescriptionWidth {
			description = description[:validDescriptionWidth-3] + "..."
		}

		content[i] = lipgloss.NewStyle().Foreground(GrayColor).Render(id) + " " + description
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

func (m model) TaskScrrenUpdate(msg tea.Msg) (tea.Model, tea.Cmd) {
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
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyTab:
			if m.active < len(m.tabs)-1 {
				m.taskScreenState.active++
				m.updateContent()
			} else {
				m.taskScreenState.active--
				m.updateContent()
			}
		case tea.KeyShiftTab:
			if m.active > 0 {
				m.taskScreenState.active--
				m.updateContent()
			} else {
				m.taskScreenState.active++
				m.updateContent()
			}
		case tea.KeyEnter:
			if m.taskScreenState.taskInput.Focused() {
				value := m.taskScreenState.taskInput.Value()
				if len(value) != 0 {
					db.Create(value)
				}
				m.ignoreQKey = false
				m.taskScreenState.taskInput.Reset()
				m.taskScreenState.taskInput.Blur()
				m.updateContent()
			}
		case tea.KeyEscape:
			m.ignoreQKey = false
			m.taskScreenState.taskInput.Reset()
			m.taskScreenState.taskInput.Blur()
			return m, nil
		case tea.KeyRunes:
			if !m.taskScreenState.taskInput.Focused() && msg.String() == "a" {
				m.ignoreQKey = true
				return m, m.taskScreenState.taskInput.Focus()
			}
		}
	}

	var cmd tea.Cmd
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
