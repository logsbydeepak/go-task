package main

import (
	"fmt"
	"os"

	"example.com/cmd"
	"example.com/db"
	"example.com/file"
)

func main() {
	err := db.Connect()
	if err != nil {
		fmt.Println(err)
		return
	}
	defer db.Close()

	err = db.Init()

	if err != nil {
		fmt.Println(err)
		return
	}

	f, err := file.LoadFile("./task.csv")
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return
	}

	cmd.Execute()
	defer file.CloseFile(f)
}
