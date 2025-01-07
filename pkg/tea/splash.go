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

var (
	baseText   = lipgloss.NewStyle().Bold(true).Render("task ")
	cursorText = lipgloss.NewStyle().Bold(true).Render("task_")
)

func (m model) SplashScreenView() string {
	if m.splashScreenState.cursorOn == true {
		return lipgloss.Place(m.viewportWidth, m.viewportHeight, lipgloss.Center, lipgloss.Center, cursorText)
	} else {
		return lipgloss.Place(m.viewportWidth, m.viewportHeight, lipgloss.Center, lipgloss.Center, baseText)
	}

}

type SplashScrrenTimeOver struct{}
type SplashScrrenCursorToggle struct{}

func (m model) SplashScrrenUpdate(msg bt.Msg) (bt.Model, bt.Cmd) {
	screenTimeoutCmd := bt.Tick(time.Second*4, func(t time.Time) bt.Msg {
		return SplashScrrenTimeOver{}
	})
	toggleCursorCmd := bt.Tick(time.Millisecond*500, func(t time.Time) bt.Msg {
		return SplashScrrenCursorToggle{}
	})

	switch msg.(type) {
	case SplashScrrenTimeOver:
		m.screen = taskScreen
		return m, nil
	case SplashScrrenCursorToggle:
		m.splashScreenState.cursorOn = !m.splashScreenState.cursorOn
		return m, toggleCursorCmd
	}

	return m, tea.Batch(toggleCursorCmd, screenTimeoutCmd)
}
