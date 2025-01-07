package tea

import (
	"time"

	bt "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

func (m model) SplashScreenView() string {
	return lipgloss.Place(m.viewportWidth, m.viewportHeight, lipgloss.Center, lipgloss.Center, "Splash")
}

type TickMsg time.Time

type SplashScrrenTimeOver struct{}

func (m model) SplashScrrenUpdate(msg bt.Msg) (bt.Model, bt.Cmd) {
	switch msg.(type) {
	case SplashScrrenTimeOver:
		m.screen = taskScreen
		return m, nil
	}

	cmd := bt.Tick(time.Second*2, func(t time.Time) bt.Msg {
		return SplashScrrenTimeOver{}
	})

	return m, cmd
}
