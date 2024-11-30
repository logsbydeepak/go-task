package cmd

import (
	"fmt"

	"example.com/db"
	"example.com/file"
	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all the task",
	Run: func(cmd *cobra.Command, args []string) {
		var err error
		showAll, err := cmd.Flags().GetBool("all")
		if err != nil {
			fmt.Println("Failed to parse flag value")
			return
		}

		var tasks []file.Task
		if showAll {
			tasks, err = db.GetAllTask()
		} else {
			tasks, err = db.GetAllPendingTask()
		}

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
