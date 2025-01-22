package cmd

import (
	"example.com/pkg/tui"

	"github.com/spf13/cobra"
)

var tuiCmd = &cobra.Command{
	Use:   "tui",
	Short: "open tui",
	Run: func(cmd *cobra.Command, args []string) {
		tui.Handler()
	},
}

func init() {
	rootCmd.AddCommand(tuiCmd)
}
