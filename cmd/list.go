package cmd

import (
	"fmt"
	"os"
	"text/tabwriter"

	"example.com/pkg/db"
	"example.com/pkg/output"
	"example.com/pkg/task"
	"github.com/mergestat/timediff"
	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all the task",
	Run: func(cmd *cobra.Command, args []string) {

		var err error
		showAll, err := cmd.Flags().GetBool("all")
		if err != nil {
			output.Error("Failed to parse flag value")
			return
		}

		var tasks []task.Task
		if showAll {
			tasks, err = db.GetAllTask()
		} else {
			tasks, err = db.GetAllPendingTask()
		}

		if err != nil {
			output.Error("Failed get tasks")
			return
		}

		w := tabwriter.NewWriter(os.Stdout, 0, 0, 4, ' ', tabwriter.StripEscape)
		fmt.Fprintln(w, "ID\tDescription\tCreatedAt\tIsComplete\t")

		for _, task := range tasks {
			time := timediff.TimeDiff(task.CreatedAt)
			text := fmt.Sprintf("%v\t%s\t%s\t%v\t", task.ID, task.Description, time, task.IsComplete)
			fmt.Fprintln(w, text)
		}
		err = w.Flush()
		if err != nil {
			output.Error("Failed to print")
			return
		}
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
	listCmd.Flags().BoolP("all", "a", false, "List all complete and incomplete task")
}
