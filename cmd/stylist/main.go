package main

import (
	"fmt"
	"os"

	"github.com/twelvelabs/stylist/internal/cmd"
	"github.com/twelvelabs/stylist/internal/stylist"
)

func main() {
	app, err := stylist.NewApp()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	command := cmd.NewRootCmd(app)
	if err := command.Execute(); err != nil {
		os.Exit(1)
	}
}
