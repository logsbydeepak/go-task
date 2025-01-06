package cmd

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/spf13/cobra"
)

var tuiCmd = &cobra.Command{
	Use:   "tui",
	Short: "open tui",
	Run: func(cmd *cobra.Command, args []string) {
		p := tea.NewProgram(model{}, tea.WithAltScreen())
		if _, err := p.Run(); err != nil {
			fmt.Println("Error running program:", err)
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(tuiCmd)
}

type model struct {
	viewportWidth  int
	viewportHeight int
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	fmt.Println("2")
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		}
	case tea.WindowSizeMsg:
		m.viewportHeight = msg.Height
		m.viewportWidth = msg.Width
	}
	return m, nil
}

func (m model) View() string {
	fmt.Println("3")
	return lipgloss.Place(m.viewportWidth, m.viewportHeight, lipgloss.Center, lipgloss.Center, "Splash")
}
