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

	addFormatFlag(cmd, &app.Config.Output.Format)
	addProcessorFilterFlags(cmd, action.ProcessorFilter)

	return cmd
}

func NewCheckAction(app *stylist.App) *CheckAction {
	return &CheckAction{
		App:             app,
		ProcessorFilter: &stylist.ProcessorFilter{},
	}
}

type CheckAction struct {
	*stylist.App

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
	excludes := a.Config.Excludes
	processors, err := a.ProcessorFilter.Filter(a.Config.Processors)
	if err != nil {
		return err
	}

	pipeline := stylist.NewPipeline(processors, excludes)
	results, err := pipeline.Check(ctx, a.pathSpecs)
	if err != nil {
		return err
	}

	for _, result := range results {
		a.Logger.Debug(fmt.Sprintf("%#v", result))
	}

	err = stylist.NewResultPrinter(a.IO, a.Config).Print(results)
	if err != nil {
		return err
	}

	return nil
}
