package tea

import (
	"time"

	bt "github.com/charmbracelet/bubbletea"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type splashScreenState struct {
	cursorOn bool
}

func (m model) SplashScreenView() string {
	text := "task"

	if m.splashScreenState.cursorOn == true {
		text = text + "_"
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

func (m model) SplashScrrenUpdate(msg bt.Msg) (bt.Model, bt.Cmd) {
	screenTimeoutCmd := bt.Tick(time.Second*4, func(t time.Time) bt.Msg {
		return SplashScrrenTimeOverMsg{}
	})
	toggleCursorCmd := bt.Tick(time.Millisecond*500, func(t time.Time) bt.Msg {
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
