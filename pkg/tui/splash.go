package tui

import (
	"time"

	"github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type splashScreenState struct {
	cursorOn bool
}

func (m model) SplashScreenView() string {
	text := "task"

	if m.splashScreenState.cursorOn == true {
		text = text + "_"
	} else {
		text = text + " "
	}

	return lipgloss.Place(
		m.viewportWidth,
		m.viewportHeight,
		lipgloss.Center,
		lipgloss.Center,
		lipgloss.NewStyle().Bold(true).Render(text),
	)
}

type SplashScrrenTimeOverMsg struct{}
type SplashScrrenCursorToggleMsg struct{}

func (m model) SplashScrrenUpdate(msg tea.Msg) (tea.Model, tea.Cmd) {
	screenTimeoutCmd := tea.Tick(time.Second*4, func(t time.Time) tea.Msg {
		return SplashScrrenTimeOverMsg{}
	})
	toggleCursorCmd := tea.Tick(time.Millisecond*500, func(t time.Time) tea.Msg {
		return SplashScrrenCursorToggleMsg{}
	})

	switch msg.(type) {
	case SplashScrrenTimeOverMsg:
		return m.TaskScreenSwitch()
	case SplashScrrenCursorToggleMsg:
		m.splashScreenState.cursorOn = !m.splashScreenState.cursorOn
		return m, toggleCursorCmd
	}

	return m, tea.Batch(toggleCursorCmd, screenTimeoutCmd)
}
