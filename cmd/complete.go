package cmd

import (
	"fmt"
	"strconv"

	"example.com/db"
	"github.com/spf13/cobra"
)

var completeCmd = &cobra.Command{
	Use:   "complete",
	Short: "Mark task as complete",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			return
		}

		id, err := strconv.Atoi(args[0])
		if err != nil {
			fmt.Println("Failed to update task")
			fmt.Println(err)
			return
		}

		err = db.MarkTaskCompleted(id)
		if err != nil {
			fmt.Println("Failed to update task")
			fmt.Println(err)
			return
		}
	},
}

func init() {
	rootCmd.AddCommand(completeCmd)
}
