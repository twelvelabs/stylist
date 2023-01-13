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

	addFormatFlag(cmd, &action.Format)
	addProcessorFilterFlags(cmd, action.ProcessorFilter)

	return cmd
}

func NewFixAction(app *stylist.App) *FixAction {
	return &FixAction{
		App:             app,
		Format:          stylist.ResultFormatTty,
		ProcessorFilter: &stylist.ProcessorFilter{},
	}
}

type FixAction struct {
	*stylist.App

	Format          stylist.ResultFormat
	ProcessorFilter *stylist.ProcessorFilter

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
	excludes := a.Config.Excludes
	processors, err := a.ProcessorFilter.Filter(a.Config.Processors)
	if err != nil {
		return err
	}

	pipeline := stylist.NewPipeline(processors, excludes)
	results, err := pipeline.Fix(ctx, a.pathSpecs)
	if err != nil {
		return err
	}

	for _, result := range results {
		a.Logger.Debug(fmt.Sprintf("%#v", result))
	}

	err = stylist.NewResultPrinter(a.Format).Print(a.App.IO, results)
	if err != nil {
		return err
	}

	return nil
}
