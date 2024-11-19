package cmd

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
	"text/tabwriter"

	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all the task",
	Run: func(cmd *cobra.Command, args []string) {
		_, err := os.Stat(FILE_PATH)
		isNewFile := os.IsNotExist(err)

		f, err := os.OpenFile(FILE_PATH, os.O_RDWR|os.O_APPEND|os.O_CREATE, 0644)
		if err != nil {
			fmt.Println("Failed to open task.csv file")
			return
		}

		defer f.Close()
		csvReader := csv.NewReader(f)

		if isNewFile {
			csvWriter := csv.NewWriter(f)
			err = csvWriter.Write([]string{"ID", "Description", "CreatedAt", "IsComplete"})
			csvWriter.Flush()
		}

		showAll, err := cmd.Flags().GetBool("all")
		if err != nil {
			fmt.Println("Failed to parse flag value")
			return
		}

		w := tabwriter.NewWriter(os.Stdout, 0, 0, 4, ' ', tabwriter.StripEscape)

		header, err := csvReader.Read()
		if err != nil {
			fmt.Println("Error reading csv data")
			return
		}

		var text strings.Builder
		for _, each := range header {
			text.WriteString(each)
			text.WriteString("\t")
		}
		fmt.Fprintln(w, text.String())

		for {
			header, err := csvReader.Read()
			if err == io.EOF {
				break
			} else if err != nil {
				fmt.Println("Error reading csv data")
				return
			}

			isCompleted := header[len(header)-1]
			value, err := strconv.ParseBool(isCompleted)

			if showAll {
				var text strings.Builder
				for _, each := range header {
					text.WriteString(each)
					text.WriteString("\t")
				}
				fmt.Fprintln(w, text.String())
			} else {
				var text strings.Builder
				for _, each := range header {
					if !value {
						text.WriteString(each)
						text.WriteString("\t")
					}
				}
				if !value {
					fmt.Fprintln(w, text.String())
				}
			}
		}

		w.Flush()
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
	listCmd.Flags().BoolP("all", "a", false, "List all complete and incomplete task")
}
