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
	formatHelp := fmt.Sprintf(
		"Output format [`FORMAT`: %s]",
		strings.Join(formatNames, ", "),
	)
	formatCompFunc := func(cmd *cobra.Command, args []string, toComplete string) (
		[]string, cobra.ShellCompDirective,
	) {
		return formatNames, cobra.ShellCompDirectiveNoFileComp
	}

	// Since go-enum generates `flag.Value` methods we can use it directly,
	// and the generated `.Set()` method will take care of validation and type casting.
	cmd.Flags().VarP(&oc.Format, "format", "f", formatHelp)
	if err := cmd.RegisterFlagCompletionFunc("format", formatCompFunc); err != nil {
		panic(err)
	}

	sortNames := stylist.ResultSortNames()
	sortHelp := fmt.Sprintf(
		"Sort issues by [`SORT`: %s]",
		strings.Join(sortNames, ", "),
	)
	sortCompFunc := func(cmd *cobra.Command, args []string, toComplete string) (
		[]string, cobra.ShellCompDirective,
	) {
		return sortNames, cobra.ShellCompDirectiveNoFileComp
	}

	cmd.Flags().VarP(&oc.Sort, "sort", "s", sortHelp)
	if err := cmd.RegisterFlagCompletionFunc("sort", sortCompFunc); err != nil {
		panic(err)
	}

	severityNames := stylist.ResultLevelNames()
	severityHelp := "Comma separated list of severities to display"
	severityCompFunc := func(cmd *cobra.Command, args []string, toComplete string) (
		[]string, cobra.ShellCompDirective,
	) {
		return severityNames, cobra.ShellCompDirectiveNoFileComp
	}

	cmd.Flags().StringSliceVar(&oc.Severity, "severity", oc.Severity, severityHelp)
	if err := cmd.RegisterFlagCompletionFunc("severity", severityCompFunc); err != nil {
		panic(err)
	}

	boolCompFunc := func(cmd *cobra.Command, args []string, toComplete string) (
		[]string, cobra.ShellCompDirective,
	) {
		return []string{"false", "true"}, cobra.ShellCompDirectiveNoFileComp
	}

	cmd.Flags().BoolVar(
		&oc.ShowContext, "show-context", oc.ShowContext, "Show the lines of code affected",
	)
	if err := cmd.RegisterFlagCompletionFunc("show-context", boolCompFunc); err != nil {
		panic(err)
	}

	cmd.Flags().BoolVar(
		&oc.ShowURL, "show-url", oc.ShowURL, "Show issue URLs when available",
	)
	if err := cmd.RegisterFlagCompletionFunc("show-url", boolCompFunc); err != nil {
		panic(err)
	}

	cmd.Flags().BoolVar(
		&oc.SyntaxHighlight, "highlight", oc.SyntaxHighlight, "Syntax highlight context lines",
	)
	if err := cmd.RegisterFlagCompletionFunc("highlight", boolCompFunc); err != nil {
		panic(err)
	}
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
