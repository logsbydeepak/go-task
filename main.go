package main

import (
	"fmt"

	"example.com/cmd"
	"example.com/db"
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

	cmd.Execute()
}
