package cmd

import (
	"fmt"
	"strings"

	mapset "github.com/deckarep/golang-set/v2"
	"github.com/spf13/cobra"

	"github.com/twelvelabs/stylist/internal/stylist"
)

func addOutputFlags(cmd *cobra.Command, oc *stylist.OutputConfig) {
	formatNames := stylist.ResultFormatNames()
	formatHelp := fmt.Sprintf("Output format [`FORMAT`: %s]", strings.Join(formatNames, ", "))
	// Since go-enum generates `flag.Value` methods we can use it directly,
	// and the generated `.Set()` method will take care of validation and type casting.
	cmd.Flags().VarP(&oc.Format, "format", "f", formatHelp)

	compFunc := func(cmd *cobra.Command, args []string, toComplete string) (
		[]string, cobra.ShellCompDirective,
	) {
		return formatNames, cobra.ShellCompDirectiveNoFileComp
	}
	if err := cmd.RegisterFlagCompletionFunc("format", compFunc); err != nil {
		panic(err)
	}

	cmd.Flags().BoolVar(
		&oc.ShowContext, "show-context", oc.ShowContext, "Show the lines of code affected",
	)
}

func addProcessorFilterFlags(cmd *cobra.Command, filter *stylist.ProcessorFilter) {
	cmd.Flags().StringSliceVarP(
		&filter.Names, "names", "n", filter.Names, "Comma separated list of processor names",
	)
	namesCompFunc := func(cmd *cobra.Command, args []string, toComplete string) (
		[]string, cobra.ShellCompDirective,
	) {
		names, _ := processorFilterFlagValues(cmd)
		return names, cobra.ShellCompDirectiveNoFileComp
	}
	if err := cmd.RegisterFlagCompletionFunc("names", namesCompFunc); err != nil {
		panic(err)
	}

	cmd.Flags().StringSliceVarP(
		&filter.Tags, "tags", "t", filter.Names, "Comma separated list of processor tags",
	)
	tagsCompFunc := func(cmd *cobra.Command, args []string, toComplete string) (
		[]string, cobra.ShellCompDirective,
	) {
		_, tags := processorFilterFlagValues(cmd)
		return tags, cobra.ShellCompDirectiveNoFileComp
	}
	if err := cmd.RegisterFlagCompletionFunc("tags", tagsCompFunc); err != nil {
		panic(err)
	}
}

// Returns all the processor names and tags defined in the config file.
func processorFilterFlagValues(cmd *cobra.Command) ([]string, []string) {
	names := mapset.NewSet[string]()
	tags := mapset.NewSet[string]()

	config := stylist.AppConfig(cmd.Context())
	for _, p := range config.Processors {
		names.Add(p.Name)
		for _, tag := range p.Tags {
			tags.Add(tag)
		}
	}

	return names.ToSlice(), tags.ToSlice()
}
