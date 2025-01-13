package main

import (
	"os"

	"example.com/cmd"
	"example.com/pkg/db"
	"example.com/pkg/output"
	"github.com/charmbracelet/log"
)

func main() {
	f, err := os.OpenFile("app.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	defer f.Close()

	log.SetOutput(f)

	log.Info("START")
	err = db.Connect()
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
