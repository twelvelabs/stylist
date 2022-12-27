package cmd

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/twelvelabs/stylist/internal/stylist"
)

func NewFixCmd(app *stylist.App) *cobra.Command {
	action := NewFixAction(app)

	cmd := &cobra.Command{
		Use:   "fix [flags] [PATH_OR_PATTERN...]",
		Short: "Run the fix command for each processor",
		Args:  cobra.ArbitraryArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := action.Validate(args); err != nil {
				return err
			}
			if err := action.Run(cmd.Context()); err != nil {
				return err
			}
			return nil
		},
		DisableFlagsInUseLine: true,
	}
	return cmd
}

func NewFixAction(app *stylist.App) *FixAction {
	return &FixAction{}
}

type FixAction struct {
	pathSpecs []string
}

func (a *FixAction) Validate(args []string) error {
	a.pathSpecs = args
	if len(a.pathSpecs) == 0 {
		a.pathSpecs = []string{"."}
	}
	return nil
}
func (a *FixAction) Run(ctx context.Context) error {
	configLoader := stylist.AppConfigLoader(ctx)
	logger := stylist.AppLogger(ctx)

	config, err := configLoader.Load()
	if err != nil {
		return err
	}

	pipeline := stylist.NewPipeline(config.Processors, config.Excludes)
	results, err := pipeline.Fix(ctx, a.pathSpecs)
	if err != nil {
		return err
	}

	for _, result := range results {
		logger.Debug(fmt.Sprintf("%#v", result))
	}

	return nil
}
