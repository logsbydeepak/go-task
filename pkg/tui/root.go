package tui

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/log"
)

func Handler() {
	m := model{}
	p := tea.NewProgram(m, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}

type screen int

const (
	splashScreen = iota
	taskScreen

	GrayColor   = lipgloss.Color("8")
	WhiteColor  = lipgloss.Color("7")
	AccentColor = lipgloss.Color("212")
)

type model struct {
	viewportWidth  int
	viewportHeight int

	error error

	ignoreQKey bool
	screen
	splashScreenState
	taskScreenState
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC:
			return m, tea.Quit
		case tea.KeyEsc:
			if m.ignoreQKey == false {
				return m, tea.Quit
			}
		case tea.KeyRunes:
			if m.ignoreQKey == false && msg.String() == "q" {
				return m, tea.Quit
			}
		}
	case tea.WindowSizeMsg:
		log.Infof("Height: %v Width: %v", msg.Height, msg.Width)
		m.viewportHeight = msg.Height
		m.viewportWidth = msg.Width
	}

	if m.error != nil {
		return m, nil
	}

	switch m.screen {
	case splashScreen:
		return m.SplashScrrenUpdate(msg)
	case taskScreen:
		return m.TaskScrrenUpdate(msg)
	}

	return m, nil
}

func (m model) View() string {
	if m.error != nil {
		return lipgloss.Place(
			m.viewportWidth,
			m.viewportHeight,
			lipgloss.Center,
			lipgloss.Center,
			"ERROR",
		)
	}

	switch m.screen {
	case splashScreen:
		return m.SplashScreenView()
	case taskScreen:
		return m.TaskScreenView()
	}

	return "ERROR"
}
