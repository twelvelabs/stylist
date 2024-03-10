package cmd

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/twelvelabs/stylist/internal/stylist"
)

func NewFilesCmd(app *stylist.App) *cobra.Command {
	action := NewFilesAction(app)

	cmd := &cobra.Command{
		Use:   "files [flags] [PATH_OR_PATTERN...]",
		Short: "List the files that would be processed by a check of fix command",
		Args:  cobra.ArbitraryArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := action.Validate(args); err != nil {
				return err
			}
			return action.Run(cmd.Context())
		},
		DisableFlagsInUseLine: true,
	}

	addProcessorFilterFlags(cmd, action.ProcessorFilter)

	return cmd
}

func NewFilesAction(app *stylist.App) *FilesAction {
	return &FilesAction{
		App:             app,
		ProcessorFilter: &stylist.ProcessorFilter{},
	}
}

type FilesAction struct {
	*stylist.App

	ProcessorFilter *stylist.ProcessorFilter

	pathSpecs []string
}

func (a *FilesAction) Validate(args []string) error {
	a.pathSpecs = args
	if len(a.pathSpecs) == 0 {
		a.pathSpecs = []string{"."}
	}
	return nil
}
func (a *FilesAction) Run(ctx context.Context) error {
	excludes := a.Config.Excludes
	processors, err := a.ProcessorFilter.Filter(a.Config.Processors)
	if err != nil {
		return err
	}

	pipeline := stylist.NewPipeline(processors, excludes)
	matches, err := pipeline.Match(ctx, a.pathSpecs)
	if err != nil {
		return err
	}

	for _, match := range matches {
		fmt.Printf("Processor: %s\n", match.Processor.Name)
		for _, path := range match.Paths {
			fmt.Printf(" - %s\n", path)
		}
		fmt.Printf("\n")
	}

	return nil
}
