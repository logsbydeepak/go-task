package cmd

import (
	"fmt"

	"example.com/db"
	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all the task",
	Run: func(cmd *cobra.Command, args []string) {
		tasks, err := db.ListAll()
		if err != nil {
			fmt.Println("Failed get tasks")
			return
		}

		for _, task := range tasks {
			fmt.Println(task)
		}
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
	listCmd.Flags().BoolP("all", "a", false, "List all complete and incomplete task")
}
