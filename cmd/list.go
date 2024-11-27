package cmd

import (
	"fmt"
	"io"
	"os"
	"strings"
	"text/tabwriter"

	"example.com/file"
	"github.com/mergestat/timediff"
	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all the task",
	Run: func(cmd *cobra.Command, args []string) {
		isNewFile := file.IsNewFile()

		if isNewFile {
			csvWriter := file.NewWriter()
			err := file.WriteHeader(csvWriter)
			if err != nil {
				fmt.Println("Failed to write header to file")
				return
			}

			csvWriter.Flush()
		}

		showAll, err := cmd.Flags().GetBool("all")
		if err != nil {
			fmt.Println("Failed to parse flag value")
			return
		}

		w := tabwriter.NewWriter(os.Stdout, 0, 0, 4, ' ', tabwriter.StripEscape)

		csvReader := file.NewReader()
		header, err := csvReader.Read()
		if err != nil {
			fmt.Println("Error reading csv data")
			return
		}

		err = file.ParseHeader(header)
		if err != nil {
			fmt.Fprintln(os.Stderr, "Failed to parse header")
			return
		}

		var text strings.Builder
		for _, each := range header {
			text.WriteString(each)
			text.WriteString("\t")
		}
		fmt.Fprintln(w, text.String())

	lineLoop:
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
				fmt.Println(err)
				return
			}

			time := timediff.TimeDiff(data.CreatedAt)

			var text strings.Builder
			for i, each := range line {
				if showAll == false && data.IsComplete == true {
					continue lineLoop
				}

				if i == file.CREATED_AT {
					text.WriteString(time)
				} else {
					text.WriteString(each)
				}

				text.WriteString("\t")
			}
			fmt.Fprintln(w, text.String())
		}
		w.Flush()
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
	listCmd.Flags().BoolP("all", "a", false, "List all complete and incomplete task")
}
