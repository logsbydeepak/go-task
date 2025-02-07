package tui

import (
	"fmt"
	"strconv"
	"strings"

	"example.com/pkg/db"
	"example.com/pkg/task"
	"github.com/charmbracelet/bubbles/table"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/mergestat/timediff"
)

type taskScreenState struct {
	tabs   []string
	tasks  []task.Task
	active int

	taskInput textinput.Model

	outerWidth  int
	outerHeight int

	taskTable table.Model
	helpTable table.Model
}

const (
	maxWidthSize        = 60
	maxHeightSize       = 20
	borderLeftRightSize = 2
	taskInputPromptSize = 5
	taskInputHeightSize = 1
)

var helpKeys = []helpKey{
	helpKey{
		key:         "a",
		description: "add new task",
	},
	helpKey{
		key:         "<tab>",
		description: "move to next tab",
	},
	helpKey{
		key:         "<shift+tab>",
		description: "move to previous tab",
	},
	helpKey{
		key:         "<space>",
		description: "mark task as done",
	},
	helpKey{
		key:         "↑/k",
		description: "move up",
	},
	helpKey{
		key:         "↓/j",
		description: "move down",
	},
	helpKey{
		key:         "↓/j",
		description: "move down",
	},
	helpKey{
		key:         "?",
		description: "switch to help",
	},
	helpKey{
		key:         "q, <esc>, <ctrl+c>",
		description: "quit",
	},
}

func (m model) TaskScreenSwitch() (tea.Model, tea.Cmd) {
	m.screen = taskScreen

	ti := textinput.New()
	ti.Placeholder = "new task, press ? for help"
	ti.CharLimit = 156
	ti.PlaceholderStyle = lipgloss.NewStyle().Italic(true).Foreground(GrayColor)

	columns := []table.Column{
		{Title: "ID", Width: 2},
		{Title: "Description", Width: 28},
		{Title: "CreatedAt", Width: 12},
		{Title: "Status", Width: 6},
	}
	t := table.New()
	t.SetColumns(columns)

	helpColumns := []table.Column{
		{Title: "Keys", Width: 18},
		{Title: "Description", Width: 20},
	}

	var helpRows []table.Row

	for _, each := range helpKeys {
		helpRows = append(helpRows, table.Row{each.key, each.description})
	}
	helpTable := table.New()
	helpTable.SetColumns(helpColumns)
	helpTable.SetRows(helpRows)

	helpTableColumns := []table.Column{
		{Title: "Keys", Width: 4},
		{Title: "Description", Width: 10},
	}

	m.taskScreenState.taskTable.SetColumns(helpTableColumns)

	m.taskScreenState = taskScreenState{
		tabs:      []string{"pending", "all", "help"},
		taskInput: ti,

		outerWidth:  maxWidthSize,
		outerHeight: maxHeightSize,
		taskTable:   t,
		helpTable:   helpTable,
	}
	m.updateContent()

	return m, nil
}

func (m model) TaskScreenView() string {
	tabStyle := lipgloss.NewStyle().Padding(0, 1)

	var tabs strings.Builder
	for i, tab := range m.tabs {
		if i == m.active {
			tabs.WriteString(tabStyle.Bold(true).Render(tab))
		} else {
			tabs.WriteString(tabStyle.Foreground(GrayColor).Render(tab))
		}
	}

	innerWidth := m.taskScreenState.outerWidth - borderLeftRightSize
	innerHeight := m.taskScreenState.outerHeight - borderLeftRightSize - taskInputHeightSize

	var content string

	if m.active == 0 || m.active == 1 {
		tasks := m.taskScreenState.tasks
		tasksLen := len(tasks)

		if tasksLen > innerHeight-taskInputHeightSize {
			tasks = tasks[:innerHeight]
		}

		var taskView string

		if tasksLen == 0 {
			taskView = "No task to show"
		} else {
			taskView = m.taskScreenState.taskTable.View()
		}

		content = m.taskScreenState.taskInput.View() + "\n" + taskView
	}

	if m.active == 2 {
		content = m.taskScreenState.helpTable.View()
	}

	innersqr := lipgloss.NewStyle().
		Width(innerWidth).
		Height(innerHeight).
		Border(lipgloss.RoundedBorder()).
		Padding(0, 1).
		Render(content)

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
	m.taskScreenState.taskTable.SetHeight(innerHeight - 1) // 1 for column
	m.taskScreenState.taskTable.SetWidth(innerWidth)
	m.taskScreenState.helpTable.SetHeight(innerHeight - 1) // 1 for column
	m.taskScreenState.helpTable.SetWidth(innerWidth)

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyTab:
			lastIndex := len(m.tabs) - 1
			if m.active == lastIndex {
				m.taskScreenState.active = 0
			} else if m.active < lastIndex {
				m.taskScreenState.active++
			} else {
				m.taskScreenState.active--
			}
			m.updateContent()
		case tea.KeyShiftTab:
			lastIndex := len(m.tabs) - 1
			if m.active == 0 {
				m.taskScreenState.active = lastIndex
			} else if m.active > 0 {
				m.taskScreenState.active--
			} else {
				m.taskScreenState.active++
			}
			m.updateContent()
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
			if m.taskScreenState.taskTable.Focused() {
				selected := m.taskScreenState.taskTable.SelectedRow()
				if len(selected) != 0 {
					id, err := strconv.Atoi(m.taskScreenState.taskTable.SelectedRow()[0])
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
			if !m.taskScreenState.taskInput.Focused() && currentKey == "?" {
				m.taskScreenState.active = 2
				return m, nil
			}
		}
	}

	var cmds []tea.Cmd

	var cmd tea.Cmd
	{
		m.taskScreenState.taskInput, cmd = m.taskScreenState.taskInput.Update(msg)
		cmds = append(cmds, cmd)
	}

	{
		if m.taskScreenState.taskInput.Focused() {
			m.taskScreenState.taskTable.Blur()
		} else {
			m.taskScreenState.taskTable.Focus()
			m.taskScreenState.taskTable, cmd = m.taskTable.Update(msg)
			cmds = append(cmds, cmd)
		}

		if m.active == 2 {
			m.taskScreenState.helpTable.Focus()
			m.taskScreenState.helpTable, cmd = m.helpTable.Update(msg)
			cmds = append(cmds, cmd)
		}
	}

	return m, tea.Batch(cmds...)
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

	m.taskScreenState.taskTable.SetRows(rows)
}

type helpKey struct {
	key         string
	description string
}
