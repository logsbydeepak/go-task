package tea

import (
	"fmt"
	"os"

	bt "github.com/charmbracelet/bubbletea"
)

func Handler() {
	p := bt.NewProgram(model{}, bt.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}

type screen int

const (
	splashScreen = iota
	taskScreen
)

type model struct {
	viewportWidth  int
	viewportHeight int
	screen
	splashScreenState
}

func (m model) Init() bt.Cmd {
	return nil
}

func (m model) Update(msg bt.Msg) (bt.Model, bt.Cmd) {
	switch msg := msg.(type) {
	case bt.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, bt.Quit
		}
	case bt.WindowSizeMsg:
		m.viewportHeight = msg.Height
		m.viewportWidth = msg.Width
	}

	switch m.screen {
	case splashScreen:
		return m.SplashScrrenUpdate(msg)
	case taskScreen:
		return m, nil
	}

	return m, nil
}

func (m model) View() string {
	switch m.screen {
	case splashScreen:
		return m.SplashScreenView()
	case taskScreen:
		return m.TaskScreenView()
	}

	return "ERROR"
}
