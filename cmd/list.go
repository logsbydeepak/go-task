package cmd

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
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

		w := tabwriter.NewWriter(os.Stdout, 0, 0, 4, ' ', tabwriter.StripEscape)
		for {
			header, err := csvReader.Read()
			if err == io.EOF {
				break
			} else if err != nil {
				fmt.Println("Error reading csv data")
				return
			}

			var text strings.Builder
			for _, each := range header {
				text.WriteString(each)
				text.WriteString("\t")
			}
			fmt.Fprintln(w, text.String())
			// fmt.Println(header)
		}
		w.Flush()
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
}
