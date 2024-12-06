package cmd

import (
	"example.com/pkg/db"
	"example.com/pkg/output"
	"github.com/spf13/cobra"
)

var addCmd = &cobra.Command{
	Use:   "add",
	Short: "Create new task",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			return
		}

		description := args[0]
		if len(description) == 0 {
			return
		}

		err := db.Create(description)
		if err != nil {
			output.Error("Failed to create task")
			return
		}
	},
}

func init() {
	rootCmd.AddCommand(addCmd)
}
