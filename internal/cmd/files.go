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

func NewFilesAction(app *stylist.App) *FilesAction {
	return &FilesAction{
		ProcessorFilter: &stylist.ProcessorFilter{},
	}
}

type FilesAction struct {
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
	configLoader := stylist.AppConfigLoader(ctx)
	// logger := stylist.AppLogger(ctx)

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
	err = pipeline.Index(ctx, a.pathSpecs)
	if err != nil {
		return err
	}

	for _, processor := range processors {
		fmt.Printf("Processor: %s\n", processor.Name)
		paths := processor.Paths()
		if len(paths) == 0 {
			fmt.Printf(" [no matching files]\n")
		} else {
			for _, path := range processor.Paths() {
				fmt.Printf(" - %s\n", path)
			}
		}
		fmt.Printf("\n")
	}

	return nil
}
