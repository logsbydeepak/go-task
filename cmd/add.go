package cmd

import (
	"fmt"
	"io"
	"strconv"

	"example.com/file"
	"github.com/spf13/cobra"
)

var addCmd = &cobra.Command{
	Use:   "add",
	Short: "Create new task",
	Run: func(cmd *cobra.Command, args []string) {
		isNewFile := file.IsNewFile()
		var id int64
		id = 1

		csvWriter := file.NewWriter()
		if isNewFile {
			id = 1
			err := file.WriteHeader(csvWriter)

			if err != nil {
				fmt.Println("Failed to write to task.csv file")
				return
			}
		} else {
			csvReader := file.NewReader()
			header, err := csvReader.Read()
			if err != nil {
				fmt.Println("Failed to read file")
				return
			}
			fmt.Println(header)

			for {
				line, err := csvReader.Read()
				if err == io.EOF {
					break
				} else if err != nil {
					fmt.Println("Error reading csv data")
					return
				}

				data, err := file.ParseLine(line)
				if err != nil {
					fmt.Println("Failed to parse ID")
					return
				}

				id = data.ID
			}
			id++
		}

		for _, arg := range args {
			if len(arg) == 0 {
				return
			}

			task := file.Task{
				ID:          id,
				Description: arg,
				CreatedAt:   "time",
				IsComplete:  false,
			}

			err := file.WriteTask(csvWriter, task)
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
