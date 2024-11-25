package cmd

import (
	"fmt"
	"io"
	"os"

	"example.com/file"
	"github.com/spf13/cobra"
)

var completeCmd = &cobra.Command{
	Use:   "complete",
	Short: "Mark task as complete",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			return
		}

		isNewFile := file.IsNewFile()
		var data [][]string

		csvReader := file.NewReader()
		header, err := csvReader.Read()
		if err != nil {
			fmt.Fprintln(os.Stderr, "Failed to read file")
			return
		}
		if isNewFile {
			fmt.Fprintln(os.Stderr, "Not found")
			return
		} else {

			err = file.ParseHeader(header)
			if err != nil {
				fmt.Fprintln(os.Stderr, "Failed to parse header")
				return
			}
			data = append(data, header)
		}

		for {
			line, err := csvReader.Read()
			if err == io.EOF {
				break
			} else if err != nil {
				fmt.Fprintln(os.Stderr, "Failed to read file")
				return
			}

			_, err = file.ParseLine(line)
			if err != nil {
				fmt.Fprintln(os.Stderr, "Failed to parse line")
				return
			}

			if line[0] == args[0] {
				line[3] = "true"
			}

			data = append(data, line)
		}

		err = file.Truncate()
		if err != nil {
			fmt.Fprintln(os.Stderr, "Failed to truncate file")
			return
		}

		csvWriter := file.NewWriter()
		err = csvWriter.WriteAll(data)
		if err != nil {
			fmt.Fprintln(os.Stderr, "Failed to write")
			return
		}
		csvWriter.Flush()
	},
}

func init() {
	rootCmd.AddCommand(completeCmd)
}
