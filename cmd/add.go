package cmd

import (
	"errors"
	"fmt"
	"io"
	"os"
	"time"

	"example.com/file"
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
			var err error
			id, err = getID()
			if err != nil {
				fmt.Fprintln(os.Stderr, err)
				return
			}
		}

		task := file.Task{
			ID:          id,
			Description: description,
			CreatedAt:   time.Now(),
			IsComplete:  false,
		}

		err := file.WriteTask(csvWriter, task)
		if err != nil {
			fmt.Fprintln(os.Stderr, "Failed to write line into file")
			return
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

func getID() (int64, error) {
	var id int64 = 1
	csvReader := file.NewReader()
	header, err := csvReader.Read()

	if err != nil {
		return id, errors.New("Failed to read file")
	}

	err = file.ParseHeader(header)
	if err != nil {
		return id, errors.New("Failed to parse header")
	}

	lineNumber := 1
	for {
		line, err := csvReader.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			return id, errors.New("Failed to read file")
		}

		data, err := file.ParseLine(line)
		if err != nil {
			return id, errors.New("Failed to parse line")
		}

		id = data.ID
		lineNumber++
	}
	if lineNumber != 1 {
		id++
	}

	return id, nil
}
