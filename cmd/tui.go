package cmd

import (
	"example.com/pkg/tea"

	"github.com/spf13/cobra"
)

var tuiCmd = &cobra.Command{
	Use:   "tui",
	Short: "open tui",
	Run: func(cmd *cobra.Command, args []string) {
		tea.Handler()
	},
}

func init() {
	rootCmd.AddCommand(tuiCmd)
}
