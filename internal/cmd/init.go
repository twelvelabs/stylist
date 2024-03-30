package cmd

import (
	"context"
	"os"

	"github.com/spf13/cobra"

	"github.com/twelvelabs/stylist/internal/fsutils"
	"github.com/twelvelabs/stylist/internal/stylist"
)

func NewInitCmd(app *stylist.App) *cobra.Command {
	action := NewInitAction(app)

	cmd := &cobra.Command{
		Use:   "init",
		Short: "Initialize a new Stylist config",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := action.Validate(args); err != nil {
				return err
			}
			return action.Run(cmd.Context())
		},
	}

	return cmd
}

func NewInitAction(app *stylist.App) *InitAction {
	return &InitAction{
		App: app,
	}
}

type InitAction struct {
	*stylist.App
}

func (a *InitAction) Validate(_ []string) error {
	return nil
}

func (a *InitAction) Run(ctx context.Context) error {
	config := stylist.NewConfig()
	configPath := config.ConfigPath
	verb := "Created"

	// Handle existing config file.
	if fsutils.PathExists(configPath) {
		verb = "Replaced"
		a.UI.Out(a.UI.WarningIcon()+" %s already exists\n", configPath)
		ok, err := a.Prompter.Confirm("Overwrite?", false)
		if err != nil {
			return err
		}
		if !ok {
			return nil // user said "no", so bail
		}
		err = os.Remove(configPath)
		if err != nil {
			return err
		}
	}

	// Find all presets that match files in the current working dir.
	store, err := stylist.NewPresetStore()
	if err != nil {
		return err
	}
	presets := store.All()
	excludes := config.Excludes
	pipeline := stylist.NewPipeline(presets, excludes)

	cwd, _ := os.Getwd()
	matches, err := pipeline.Match(ctx, cwd, []string{"."})
	if err != nil {
		return err
	}

	processors := []*stylist.Processor{}
	for _, match := range matches {
		processors = append(processors, match.Processor)
	}

	// Generate a new config file containing all matching presets.
	// Commenting out everything but the `preset: foo` line so that
	// users can see what the preset is doing and how to override.
	config = &stylist.Config{
		Processors: processors,
	}
	if err := stylist.WriteConfig(config, configPath); err != nil {
		return err
	}
	if err := stylist.CommentOutConfigPresets(configPath); err != nil {
		return err
	}

	a.UI.Out(a.UI.SuccessIcon()+" %s %s\n", verb, configPath)
	return nil
}
