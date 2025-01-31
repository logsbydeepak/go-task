package tui

import (
	"fmt"
	"strconv"
	"strings"

	"example.com/pkg/db"
	"example.com/pkg/task"
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/table"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/mergestat/timediff"
)

type taskScreenState struct {
	tabs      []string
	tasks     []task.Task
	active    int
	taskInput textinput.Model

	outerWidth  int
	outerHeight int
	table       table.Model

	help              help.Model
	shortHelpViewKeys []key.Binding
}

const (
	maxWidthSize        = 60
	maxHeightSize       = 20
	borderLeftRightSize = 2
	taskInputPromptSize = 5
	taskInputHeightSize = 1

	GrayColor = lipgloss.Color("8")
)

var shortHelpViewKeys = []key.Binding{
	key.NewBinding(
		key.WithKeys("a"),
		key.WithHelp("a", "new task"),
	),
	key.NewBinding(
		key.WithKeys("space"),
		key.WithHelp("<space>", "mark task as done"),
	),
	key.NewBinding(
		key.WithKeys("up", "k"),
		key.WithHelp("↑/k", "move up"),
	),
	key.NewBinding(
		key.WithKeys("down", "j"),
		key.WithHelp("↓/j", "move down"),
	),
	key.NewBinding(
		key.WithKeys("left", "h"),
		key.WithHelp("←/h", "move left"),
	),
	key.NewBinding(
		key.WithKeys("right", "l"),
		key.WithHelp("→/l", "move right"),
	),
	key.NewBinding(
		key.WithKeys("?"),
		key.WithHelp("?", "toggle help"),
	),
	key.NewBinding(
		key.WithKeys("q", "esc", "ctrl+c"),
		key.WithHelp("q", "quit"),
	),
}

func (m model) TaskScreenSwitch() (tea.Model, tea.Cmd) {
	ti := textinput.New()
	ti.Placeholder = "new task, press ? for help"
	ti.CharLimit = 156
	ti.PlaceholderStyle = lipgloss.NewStyle().Italic(true).Foreground(GrayColor)

	t := table.New()
	h := help.New()
	m.screen = taskScreen
	m.taskScreenState = taskScreenState{
		tabs:      []string{"pending", "all"},
		taskInput: ti,

		outerWidth:        maxWidthSize,
		outerHeight:       maxHeightSize,
		table:             t,
		help:              h,
		shortHelpViewKeys: shortHelpViewKeys,
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
	tasksLen := len(tasks)

	if tasksLen > innerHeight-taskInputHeightSize {
		tasks = tasks[:innerHeight]
	}

	var taskView string

	if tasksLen == 0 {
		taskView = "No task to show"
	} else {
		taskView = m.taskScreenState.table.View()
	}

	innersqr := lipgloss.NewStyle().
		Width(innerWidth).
		Height(innerHeight).
		Border(lipgloss.RoundedBorder()).
		Padding(0, 1).
		Render(m.taskScreenState.taskInput.View() + "\n" + taskView)

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

	innerHeight := m.taskScreenState.outerHeight - borderLeftRightSize - taskInputHeightSize
	innerWidth := m.taskScreenState.outerWidth - borderLeftRightSize
	m.taskScreenState.table.SetHeight(innerHeight - 1) // 1 for column
	m.taskScreenState.table.SetWidth(innerWidth)

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
		case tea.KeySpace:
			if m.taskScreenState.table.Focused() {
				selected := m.taskScreenState.table.SelectedRow()
				if len(selected) != 0 {
					id, err := strconv.Atoi(m.taskScreenState.table.SelectedRow()[0])
					if err == nil {
						db.MarkTaskCompleted(id)
						m.updateContent()
						return m, nil
					}
				}
			}
		case tea.KeyRunes:
			currentKey := msg.String()
			if !m.taskScreenState.taskInput.Focused() && currentKey == "a" {
				m.ignoreQKey = true
				return m, m.taskScreenState.taskInput.Focus()
			}
		}
	}

	var taskCmd tea.Cmd
	var tableCmd tea.Cmd
	m.taskScreenState.taskInput, taskCmd = m.taskScreenState.taskInput.Update(msg)
	m.taskScreenState.table, tableCmd = m.table.Update(msg)

	if !m.taskScreenState.taskInput.Focused() {
		m.taskScreenState.table.Focus()
	} else {
		m.taskScreenState.table.Blur()
	}

	return m, tea.Batch(taskCmd, tableCmd)
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

	var rows []table.Row

	for _, each := range m.tasks {
		time := timediff.TimeDiff(each.CreatedAt)
		var isComplete string
		if each.IsComplete {
			isComplete += "✓"
		} else {
			isComplete += "⨯"
		}

		rows = append(rows, table.Row{fmt.Sprintf("%v", each.ID), each.Description, time, isComplete})
	}

	columns := []table.Column{
		{Title: "ID", Width: 2},
		{Title: "Description", Width: 28},
		{Title: "CreatedAt", Width: 12},
		{Title: "Status", Width: 6},
	}
	m.taskScreenState.table.SetColumns(columns)
	m.taskScreenState.table.SetRows(rows)
}
