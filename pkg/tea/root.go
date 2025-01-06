package tea

import (
	"fmt"
	"os"

	bt "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

func Handler() {
	p := bt.NewProgram(model{}, bt.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}

type model struct {
	viewportWidth  int
	viewportHeight int
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
	return m, nil
}

func (m model) View() string {
	return lipgloss.Place(m.viewportWidth, m.viewportHeight, lipgloss.Center, lipgloss.Center, "Splash")
}