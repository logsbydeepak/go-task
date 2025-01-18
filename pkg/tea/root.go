package tea

import (
	"fmt"
	"os"

	bt "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/log"
)

func Handler() {
	m := model{}
	p := bt.NewProgram(m, bt.WithAltScreen())
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
	taskScreenState
}

func (m model) Init() bt.Cmd {
	return nil
}

func (m model) Update(msg bt.Msg) (bt.Model, bt.Cmd) {
	switch msg := msg.(type) {
	case bt.KeyMsg:
		switch msg.Type {
		case bt.KeyCtrlC:
			return m, bt.Quit
		case bt.KeyRunes:
			if msg.String() == "q" {
				return m, bt.Quit
			}
		}
	case bt.WindowSizeMsg:
		log.Infof("Height: %v Width: %v", msg.Height, msg.Width)
		m.viewportHeight = msg.Height
		m.viewportWidth = msg.Width
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
	switch m.screen {
	case splashScreen:
		return m.SplashScreenView()
	case taskScreen:
		return m.TaskScreenView()
	}

	return "ERROR"
}
