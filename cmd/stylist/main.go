package main

import (
	"context"
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
	ctx := app.InitContext(context.Background())
	if err := command.ExecuteContext(ctx); err != nil {
		os.Exit(1)
	}
}
