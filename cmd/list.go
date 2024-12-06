package cmd

import (
	"fmt"
	"os"
	"strings"
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
			var text strings.Builder
			text.WriteString(fmt.Sprintf("%v", task.ID))
			text.WriteString("\t")
			text.WriteString(task.Description)
			text.WriteString("\t")
			time := timediff.TimeDiff(task.CreatedAt)
			text.WriteString(time)
			text.WriteString("\t")
			text.WriteString(fmt.Sprintf("%v", task.IsComplete))
			text.WriteString("\t")
			fmt.Fprintln(w, text.String())
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
