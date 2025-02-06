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
	helpTable         table.Model
}

const (
	maxWidthSize        = 60
	maxHeightSize       = 20
	borderLeftRightSize = 2
	taskInputPromptSize = 5
	taskInputHeightSize = 1

	GrayColor   = lipgloss.Color("8")
	WhiteColor  = lipgloss.Color("7")
	AccentColor = lipgloss.Color("212")
)

var shortHelpViewKeys = []key.Binding{
	key.NewBinding(
		key.WithKeys("a"),
		key.WithHelp("a", "new task"),
	),
	key.NewBinding(
		key.WithKeys("tab"),
		key.WithHelp("<tab>", "move to next tab"),
	),
	key.NewBinding(
		key.WithKeys("tab"),
		key.WithHelp("<tab>", "move to previous tab"),
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
	ti := textinput.New()
	ti.Placeholder = "new task, press ? for help"
	ti.CharLimit = 156
	ti.PlaceholderStyle = lipgloss.NewStyle().Italic(true).Foreground(GrayColor)

	t := table.New()
	h := help.New()
	h.Styles.ShortKey = lipgloss.NewStyle().Bold(true).Foreground(WhiteColor)
	h.Styles.ShortDesc = lipgloss.NewStyle().Foreground(WhiteColor)
	m.screen = taskScreen

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

	columns := []table.Column{
		{Title: "Keys", Width: 4},
		{Title: "Description", Width: 10},
	}

	m.taskScreenState.table.SetColumns(columns)

	m.taskScreenState = taskScreenState{
		tabs:      []string{"pending", "all", "help"},
		taskInput: ti,

		outerWidth:        maxWidthSize,
		outerHeight:       maxHeightSize,
		table:             t,
		help:              h,
		helpTable:         helpTable,
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
			taskView = m.taskScreenState.table.View()
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
	m.taskScreenState.table.SetHeight(innerHeight - 1) // 1 for column
	m.taskScreenState.table.SetWidth(innerWidth)
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
			m.taskScreenState.table.Blur()
		} else {
			m.taskScreenState.table.Focus()
			m.taskScreenState.table, cmd = m.table.Update(msg)
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

	columns := []table.Column{
		{Title: "ID", Width: 2},
		{Title: "Description", Width: 28},
		{Title: "CreatedAt", Width: 12},
		{Title: "Status", Width: 6},
	}

	m.taskScreenState.table.SetColumns(columns)
	m.taskScreenState.table.SetRows(rows)
}

type helpKey struct {
	key         string
	description string
}
