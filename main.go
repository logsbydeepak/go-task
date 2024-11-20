package main

import (
	"fmt"

	"example.com/cmd"
	"example.com/file"
)

func main() {
	f, err := file.LoadFile("./task.csv")
	if err != nil {
		fmt.Println(err)
		return
	}

	cmd.Execute()

	defer file.CloseFile(f)
}
