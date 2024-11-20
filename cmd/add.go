package cmd

import (
	"fmt"
	"io"
	"os"
	"strconv"

	"example.com/file"
	"github.com/spf13/cobra"
)

const FILE_PATH = "./task.csv"

var addCmd = &cobra.Command{
	Use:   "add",
	Short: "Create new task",
	Run: func(cmd *cobra.Command, args []string) {
		_, err := os.Stat(FILE_PATH)
		isNewFile := os.IsNotExist(err)
		id := 1

		csvWriter := file.NewWriter()

		if isNewFile {
			id = 1
			err = csvWriter.Write([]string{"ID", "Description", "CreatedAt", "IsComplete"})

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
				id, err = strconv.Atoi(line[0])
				if err != nil {
					fmt.Println("Failed to parse ID")
					return
				}
			}
			id++
		}

		for _, arg := range args {
			if len(arg) == 0 {
				return
			}

			err = csvWriter.Write([]string{strconv.Itoa(id), arg, "time", "false"})
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
