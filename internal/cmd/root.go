package cmd

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	"github.com/twelvelabs/stylist/internal/stylist"
)

func NewRootCmd(app *stylist.App) *cobra.Command {
	cmd := &cobra.Command{
		Use:          "stylist",
		Short:        "Lint and format with style",
		Version:      app.Meta.Version,
		SilenceUsage: true,
	}

	cfg := app.Config
	flags := cmd.PersistentFlags()
	flags.StringVarP(&cfg.ConfigPath, "config", "c", cfg.ConfigPath, "Config path")

	levelNames := stylist.LogLevelNames()
	levelHelp := fmt.Sprintf("Log level [`LEVEL`: %s]", strings.Join(levelNames, ", "))
	flags.Var(&cfg.LogLevel, "log-level", levelHelp)

	levelCompFunc := func(cmd *cobra.Command, args []string, toComplete string) (
		[]string, cobra.ShellCompDirective,
	) {
		return levelNames, cobra.ShellCompDirectiveNoFileComp
	}
	if err := cmd.RegisterFlagCompletionFunc("log-level", levelCompFunc); err != nil {
		panic(err)
	}

	cmd.AddCommand(NewCheckCmd(app))
	cmd.AddCommand(NewFixCmd(app))
	cmd.AddCommand(NewFilesCmd(app))
	cmd.AddCommand(NewInitCmd(app))
	cmd.AddCommand(NewServerCmd(app))
	cmd.AddCommand(NewVersionCmd(app))

	return cmd
}
