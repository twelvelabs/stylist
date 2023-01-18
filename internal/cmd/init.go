package cmd

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/twelvelabs/stylist/internal/stylist"
)

func NewInitCmd(app *stylist.App) *cobra.Command {
	action := NewInitAction(app)

	cmd := &cobra.Command{
		Use:   "init",
		Short: "Initialize a new Stylist config",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := action.Validate(args); err != nil {
				return err
			}
			if err := action.Run(cmd.Context()); err != nil {
				return err
			}
			return nil
		},
	}

	return cmd
}

func NewInitAction(app *stylist.App) *InitAction {
	return &InitAction{
		App: app,
	}
}

type InitAction struct {
	*stylist.App
}

func (a *InitAction) Validate(args []string) error {
	return nil
}

func (a *InitAction) Run(ctx context.Context) error {
	config := stylist.NewConfig()
	excludes := config.Excludes

	store, err := stylist.NewPresetStore()
	if err != nil {
		return err
	}

	pipeline := stylist.NewPipeline(store.All(), excludes)
	processors, err := pipeline.Match(ctx, []string{"."})
	if err != nil {
		return err
	}

	fmt.Fprintf(a.IO.Out, "Found %d matching presets:\n", len(processors))
	for _, p := range processors {
		fmt.Fprintf(a.IO.Out, "- %s\n", p.Name)
	}
	fmt.Fprintf(a.IO.Out, "\n")

	return nil
}
