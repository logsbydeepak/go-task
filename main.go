package main

import (
	"fmt"
	"os"

	"example.com/cmd"
	"example.com/file"
)

func main() {
	f, err := file.LoadFile("./task.csv")
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return
	}

	cmd.Execute()
	defer file.CloseFile(f)
}
