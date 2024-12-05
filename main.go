package main

import (
	"fmt"
	"os"

	"example.com/cmd"
	"example.com/db"
)

func main() {
	err := db.Connect()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Failed to connect to DB")
		return
	}
	defer db.Close()

	err = db.Init()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Failed to initialize DB")
		return
	}

	cmd.Execute()
}
