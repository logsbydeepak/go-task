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

type tab int
type taskScreenState struct {
	tasks     []task.Task
	taskInput textinput.Model

	tabs      []string
	activeTab tab

	outerWidth  int
	outerHeight int

	taskTable table.Model
	helpTable table.Model
}

const (
	pendingTaskTab = iota
	allTaskTab
	helpTaskTab

	maxWidthSize        = 60
	maxHeightSize       = 20
	borderLeftRightSize = 2
	taskInputPromptSize = 5
	taskInputHeightSize = 1
)

type helpKey struct {
	key         string
	description string
}

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

	taskInput := textinput.New()
	taskInput.Placeholder = "new task, press ? for help"
	taskInput.CharLimit = 156
	taskInput.PlaceholderStyle = lipgloss.NewStyle().Italic(true).Foreground(GrayColor)

	taskTable := table.New()
	taskColumns := []table.Column{
		{Title: "ID", Width: 2},
		{Title: "Description", Width: 28},
		{Title: "CreatedAt", Width: 12},
		{Title: "Status", Width: 6},
	}
	taskTable.SetColumns(taskColumns)

	helpTable := table.New()
	helpColumns := []table.Column{
		{Title: "Keys", Width: 18},
		{Title: "Description", Width: 20},
	}

	var helpRows []table.Row
	for _, each := range helpKeys {
		helpRows = append(helpRows, table.Row{each.key, each.description})
	}

	helpTable.SetColumns(helpColumns)
	helpTable.SetRows(helpRows)

	m.taskScreenState = taskScreenState{
		tabs:      []string{"pending", "all", "help"},
		activeTab: pendingTaskTab,
		taskInput: taskInput,

		outerWidth:  maxWidthSize,
		outerHeight: maxHeightSize,
		taskTable:   taskTable,
		helpTable:   helpTable,
	}
	m.updateContent()

	return m, nil
}

func (m model) TaskScreenView() string {
	tabStyle := lipgloss.NewStyle().Padding(0, 1)

	var tabs strings.Builder
	for i, currentTab := range m.tabs {
		if m.activeTab == tab(i) {
			tabs.WriteString(tabStyle.Bold(true).Render(currentTab))
		} else {
			tabs.WriteString(tabStyle.Foreground(GrayColor).Render(currentTab))
		}
	}

	innerWidth := m.taskScreenState.outerWidth - borderLeftRightSize
	innerHeight := m.taskScreenState.outerHeight - borderLeftRightSize - taskInputHeightSize

	var content string

	switch m.activeTab {
	case allTaskTab, pendingTaskTab:
		var taskView string
		if len(m.taskScreenState.tasks) == 0 {
			taskView = "No task to show"
		} else {
			taskView = m.taskScreenState.taskTable.View()
		}

		content = m.taskScreenState.taskInput.View() + "\n" + taskView

	case helpTaskTab:
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
			lastIndex := tab(len(m.tabs) - 1)
			if m.activeTab == lastIndex {
				m.taskScreenState.activeTab = 0
			} else if m.activeTab < lastIndex {
				m.taskScreenState.activeTab++
			} else {
				m.taskScreenState.activeTab--
			}
			m.updateContent()
		case tea.KeyShiftTab:
			lastIndex := tab(len(m.tabs) - 1)
			if m.activeTab == 0 {
				m.taskScreenState.activeTab = lastIndex
			} else if m.activeTab > 0 {
				m.taskScreenState.activeTab--
			} else {
				m.taskScreenState.activeTab++
			}
			m.updateContent()
		case tea.KeyEnter:
			if !m.taskScreenState.taskInput.Focused() {
				return m, nil
			}
			value := m.taskScreenState.taskInput.Value()
			if len(value) != 0 {
				db.Create(value)
			}
			m.ignoreQKey = false
			m.taskScreenState.taskInput.Reset()
			m.taskScreenState.taskInput.Blur()
			m.updateContent()
		case tea.KeyEscape:
			if !m.taskScreenState.taskInput.Focused() {
				return m, nil
			}
			m.ignoreQKey = false
			m.taskScreenState.taskInput.Reset()
			m.taskScreenState.taskInput.Blur()
		case tea.KeySpace:
			if !m.taskScreenState.taskTable.Focused() {
				return m, nil
			}
			selected := m.taskScreenState.taskTable.SelectedRow()
			if len(selected) == 0 {
				return m, nil
			}
			id, err := strconv.Atoi(selected[0])
			if err != nil {
				m.error = err
				return m, nil
			}
			db.MarkTaskCompleted(id)
			m.updateContent()
			return m, nil
		case tea.KeyRunes:
			currentKey := msg.String()
			if !m.taskScreenState.taskInput.Focused() && currentKey == "a" {
				m.ignoreQKey = true
				return m, m.taskScreenState.taskInput.Focus()
			} else if !m.taskScreenState.taskInput.Focused() && currentKey == "?" {
				m.taskScreenState.activeTab = 2
				return m, nil
			}
		}
	}

	var cmds []tea.Cmd

	var cmd tea.Cmd
	m.taskScreenState.taskInput, cmd = m.taskScreenState.taskInput.Update(msg)
	cmds = append(cmds, cmd)

	if m.taskScreenState.taskInput.Focused() {
		m.taskScreenState.taskTable.Blur()
	} else {
		m.taskScreenState.taskTable.Focus()
		m.taskScreenState.taskTable, cmd = m.taskTable.Update(msg)
		cmds = append(cmds, cmd)
	}

	if m.activeTab == 2 {
		m.taskScreenState.helpTable.Focus()
		m.taskScreenState.helpTable, cmd = m.helpTable.Update(msg)
		cmds = append(cmds, cmd)
	}

	return m, tea.Batch(cmds...)
}

func (m *model) updateContent() {
	switch m.activeTab {
	case allTaskTab:
		task, err := db.GetAllTask()
		if err != nil {
			m.error = err
			return
		}
		m.tasks = task
	case pendingTaskTab:
		task, err := db.GetAllPendingTask()
		if err != nil {
			m.error = err
			return
		}
		m.tasks = task
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
