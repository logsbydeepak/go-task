package cmd

import (
	"fmt"
	"io"
	"os"

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
				fmt.Fprintln(os.Stderr, "Failed to write to file")
				return
			}
		} else {
			csvReader := file.NewReader()
			header, err := csvReader.Read()

			if err != nil {
				fmt.Fprintln(os.Stderr, "Failed to read file")
				return
			}

			err = file.ParseHeader(header)
			if err != nil {
				fmt.Fprintln(os.Stderr, "Failed to parse header")
				return
			}

			lineNumber := 1

			for {
				line, err := csvReader.Read()
				if err == io.EOF {
					break
				} else if err != nil {
					fmt.Fprintln(os.Stderr, "Failed to read file")
					return
				}

				data, err := file.ParseLine(line)
				if err != nil {
					fmt.Fprintln(os.Stderr, "Failed to parse line")
					return
				}

				id = data.ID
				lineNumber++
			}
			if lineNumber != 1 {
				id++
			}

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
				fmt.Fprintln(os.Stderr, "Failed to write line into file")
				return
			}
		}

		csvWriter.Flush()
		if err := csvWriter.Error(); err != nil {
			fmt.Fprintln(os.Stderr, "Failed to write into file")
			return
		}
	},
}

func init() {
	rootCmd.AddCommand(addCmd)
}
