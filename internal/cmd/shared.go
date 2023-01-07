package cmd

import (
	"github.com/spf13/cobra"

	"github.com/twelvelabs/stylist/internal/stylist"
)

func addProcessorFilterFlags(cmd *cobra.Command, filter *stylist.ProcessorFilter) {
	cmd.Flags().StringSliceVarP(
		&filter.Names, "names", "n", filter.Names, "Comma separated list of processor names",
	)
	cmd.Flags().StringSliceVarP(
		&filter.Tags, "tags", "t", filter.Names, "Comma separated list of processor tags",
	)
}
