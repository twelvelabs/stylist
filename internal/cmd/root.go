package cmd

import (
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	"github.com/twelvelabs/stylist/internal/stylist"
)

func NewRootCmd(app *stylist.App) *cobra.Command {
	// Panic: 0
	// Fatal: 1
	// Error: 2  (default)
	// Warn:  3  -v
	// Info:  4  -vv
	// Debug: 5  -vvv
	// Trace: 6  -vvvv
	verbosity := 0

	cmd := &cobra.Command{
		Use:   "stylist",
		Short: "Lint and format with style",
		Args:  cobra.ArbitraryArgs,
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			level := logrus.Level(verbosity + 2)
			if level >= logrus.TraceLevel {
				level = logrus.TraceLevel
			}
			app.Logger.SetLevel(level)
			app.Logger.Debug("Set log level to " + level.String())

			app.Logger.Debug("Loading config from " + app.ConfigLoader.Path)
			_, err := app.ConfigLoader.Load()
			if err != nil {
				return err
			}
			return nil
		},
		Version:      "X.X.X",
		SilenceUsage: true,
	}

	flags := cmd.PersistentFlags()
	flags.StringVarP(
		&app.ConfigLoader.Path, "config", "c", app.ConfigLoader.Path, "Config path",
	)
	flags.CountVarP(&verbosity, "verbose", "v", "Set the log level")

	cmd.AddCommand(NewCheckCmd(app))
	cmd.AddCommand(NewFixCmd(app))
	cmd.AddCommand(NewFilesCmd(app))

	return cmd
}
