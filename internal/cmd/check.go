package cmd

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/twelvelabs/stylist/internal/stylist"
)

func NewCheckCmd(app *stylist.App) *cobra.Command {
	action := NewCheckAction(app)

	cmd := &cobra.Command{
		Use:   "check [flags] [PATH_OR_PATTERN...]",
		Short: "Run the check command for each processor",
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

	addProcessorFilterFlags(cmd, action.ProcessorFilter)

	return cmd
}

func NewCheckAction(app *stylist.App) *CheckAction {
	return &CheckAction{
		ProcessorFilter: &stylist.ProcessorFilter{},
	}
}

type CheckAction struct {
	ProcessorFilter *stylist.ProcessorFilter

	pathSpecs []string
}

func (a *CheckAction) Validate(args []string) error {
	a.pathSpecs = args
	if len(a.pathSpecs) == 0 {
		a.pathSpecs = []string{"."}
	}
	return nil
}
func (a *CheckAction) Run(ctx context.Context) error {
	configLoader := stylist.AppConfigLoader(ctx)
	logger := stylist.AppLogger(ctx)

	config, err := configLoader.Load()
	if err != nil {
		return err
	}

	excludes := config.Excludes
	processors, err := a.ProcessorFilter.Filter(config.Processors)
	if err != nil {
		return err
	}

	pipeline := stylist.NewPipeline(processors, excludes)
	results, err := pipeline.Check(ctx, a.pathSpecs)
	if err != nil {
		return err
	}

	for _, result := range results {
		logger.Debug(fmt.Sprintf("%#v", result))
	}

	return nil
}
