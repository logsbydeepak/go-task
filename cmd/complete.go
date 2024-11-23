package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var completeCmd = &cobra.Command{
	Use:   "complete",
	Short: "Mark task as complete",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("complete called")
	},
}

func init() {
	rootCmd.AddCommand(completeCmd)
}
