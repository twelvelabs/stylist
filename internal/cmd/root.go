package cmd

import (
	"context"
	"fmt"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/twelvelabs/termite/conf"
	"github.com/twelvelabs/termite/ioutil"
	"github.com/twelvelabs/termite/run"
	"github.com/twelvelabs/termite/ui"

	"github.com/twelvelabs/stylist/internal/stylist"
)

func NewRootCmd(app *stylist.App) *cobra.Command {
	action := NewRootAction(app)

	cmd := &cobra.Command{
		Use:   "stylist",
		Short: "Lint and format with style",
		Args:  cobra.ArbitraryArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := action.Setup(cmd, args); err != nil {
				return err
			}
			if err := action.Validate(); err != nil {
				return err
			}
			if err := action.Run(cmd.Context()); err != nil {
				return err
			}
			return nil
		},
		Version:      "X.X.X",
		SilenceUsage: true,
	}

	cmd.Flags().StringVarP(
		&action.ConfigLoader.Path, "config", "c", action.ConfigLoader.Path, "Config path",
	)
	cmd.Flags().CountVarP(&action.Verbosity, "verbose", "v", "Set the log level")

	cmd.Context()

	return cmd
}

func NewRootAction(app *stylist.App) *RootAction {
	return &RootAction{
		App:          app,
		IO:           app.IO,
		ConfigLoader: app.ConfigLoader,
		Messenger:    app.Messenger,
		CmdClient:    app.CmdClient,
	}
}

type RootAction struct {
	App          *stylist.App
	IO           *ioutil.IOStreams
	ConfigLoader *conf.Loader[*stylist.Config]
	Messenger    *ui.Messenger
	CmdClient    *run.Client

	Verbosity int

	pathSpecs []string
}

func (a *RootAction) Setup(cmd *cobra.Command, args []string) error {
	// Panic: 0
	// Fatal: 1
	// Error: 2  (default)
	// Warn:  3  -v
	// Info:  4  -vv
	// Debug: 5  -vvv
	// Trace: 6  -vvvv
	a.App.SetLogLevel(logrus.Level(a.Verbosity + 2))

	a.pathSpecs = args
	if len(a.pathSpecs) == 0 {
		a.pathSpecs = []string{"."}
	}
	return nil
}

func (a *RootAction) Validate() error {
	// TODO: validate pathSpecs using doublestar
	return nil
}

func (a *RootAction) Run(ctx context.Context) error {
	ctx = a.App.InitContext(ctx)

	config, err := a.ConfigLoader.Load()
	if err != nil {
		return err
	}
	processors := config.Processors

	a.Messenger.Info("Indexing...\n")

	err = processors.Index(a.pathSpecs, config.Excludes)
	if err != nil {
		return err
	}

	fmt.Println("")
	for _, processor := range processors {
		a.Messenger.Success("%s:\n", processor.Name)
		fmt.Println("")
		for _, path := range processor.Paths() {
			a.Messenger.Info("%s\n", path)
		}
		fmt.Println("")
		results, err := processor.Check(ctx)
		if err != nil {
			return err
		}
		for _, result := range results {
			a.App.Logger.Debug(fmt.Sprintf("%#v", result))
		}
		// _, _ = processor.Fix(ctx)
		fmt.Println("")
	}
	fmt.Println("")

	return nil
}
