package cmd

import (
	"context"
	"fmt"
	"time"

	"github.com/spf13/cobra"

	"github.com/twelvelabs/stylist/internal/stylist"
)

func NewVersionCmd(app *stylist.App) *cobra.Command {
	action := NewVersionAction(app)

	cmd := &cobra.Command{
		Use:   "version",
		Short: "Show full version info",
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

func NewVersionAction(app *stylist.App) *VersionAction {
	return &VersionAction{
		App: app,
	}
}

type VersionAction struct {
	*stylist.App
}

func (a *VersionAction) Validate(args []string) error {
	return nil
}

func (a *VersionAction) Run(ctx context.Context) error {
	fmt.Fprintln(a.IO.Out, "Version:", a.Meta.Version)
	fmt.Fprintln(a.IO.Out, "GOOS:", a.Meta.GOOS)
	fmt.Fprintln(a.IO.Out, "GOARCH:", a.Meta.GOARCH)
	fmt.Fprintln(a.IO.Out, "")
	fmt.Fprintln(a.IO.Out, "Build Time:", a.Meta.BuildTime.Format(time.RFC3339))
	fmt.Fprintln(a.IO.Out, "Build Commit:", a.Meta.BuildCommit)
	fmt.Fprintln(a.IO.Out, "Build Version:", a.Meta.BuildVersion)
	fmt.Fprintln(a.IO.Out, "Build Checksum:", a.Meta.BuildChecksum)
	fmt.Fprintln(a.IO.Out, "Build Go Version:", a.Meta.BuildGoVersion)
	return nil
}
