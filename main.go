package main

import (
	"example.com/cmd"
	"example.com/pkg/db"
	"example.com/pkg/output"
)

func main() {
	err := db.Connect()
	if err != nil {
		output.Error("Failed to connect to DB")
		return
	}

	defer db.Close()

	err = db.Init()
	if err != nil {
		output.Error("Failed to initialize DB")
		return
	}

	cmd.Execute()
}
