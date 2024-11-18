package cmd

import (
	"encoding/csv"
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

const FILE_PATH = "./task.csv"

var addCmd = &cobra.Command{
	Use:   "add",
	Short: "Create new task",
	Run: func(cmd *cobra.Command, args []string) {
		_, err := os.Stat(FILE_PATH)
		isNewFile := os.IsNotExist(err)

		f, err := os.OpenFile(FILE_PATH, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0644)

		if err != nil {
			if os.IsExist(err) {
				fmt.Println("hi")
			}

			fmt.Println("Failed to open task.csv file")
			return
		}
		defer f.Close()

		csvWriter := csv.NewWriter(f)

		if isNewFile {
			err = csvWriter.Write([]string{"ID", "Description", "CreatedAt", "IsComplete"})

			if err != nil {
				fmt.Println("Failed to write to task.csv file")
				return
			}
		}

		for _, arg := range args {
			if len(arg) == 0 {
				return
			}

			err = csvWriter.Write([]string{"1", arg, "time", "false"})
			if err != nil {
				fmt.Println("Failed to write to task.csv file")
				return
			}
		}

		csvWriter.Flush()
		if err := csvWriter.Error(); err != nil {
			fmt.Println("Failed to flush csv writter")
			return
		}
	},
}

func init() {
	rootCmd.AddCommand(addCmd)
}
