package stylist

import (
	"fmt"
	"sort"

	"github.com/twelvelabs/termite/ioutil"
)

// ResultPrinter is the interface that wraps the Print method.
type ResultPrinter interface {
	// Print writes the results to Stdout.
	Print(ios *ioutil.IOStreams, results []*Result) error
}

// NewResultPrinter returns the appropriate printer for the given format.
func NewResultPrinter(format ResultFormat) ResultPrinter { //nolint:ireturn
	switch format {
	case ResultFormatSarif:
		return &SarifPrinter{}
	case ResultFormatTty:
		return &TtyPrinter{}
	default:
		panic(fmt.Sprintf("unknown result format: %s", format))
	}
}

/*
* SarifPrinter
**/

// SarifPrinter generates SARIF formatted output.
type SarifPrinter struct {
}

// Print writes the SARIF formatted results to Stdout.
func (f *SarifPrinter) Print(ios *ioutil.IOStreams, results []*Result) error {
	return nil
}

/*
* TtyPrinter
**/

// TtyPrinter generates TTY formatted output.
// The output will contain ANSI color codes if the terminal allows them.
type TtyPrinter struct {
}

// Print writes the TTY formatted results to Stdout.
func (f *TtyPrinter) Print(ios *ioutil.IOStreams, results []*Result) error {
	// Maybe this should be controlled by a flag (and done by the caller)?
	sort.Slice(results, func(i, j int) bool {
		if results[i].Source != results[j].Source {
			return results[i].Source < results[j].Source
		}
		if results[i].Location.Path != results[j].Location.Path {
			return results[i].Location.Path < results[j].Location.Path
		}
		if results[i].Location.StartLine != results[j].Location.StartLine {
			return results[i].Location.StartLine < results[j].Location.StartLine
		}
		return results[i].Location.StartColumn < results[j].Location.StartColumn
	})

	formatter := ios.Formatter()
	for _, result := range results {
		fmt.Fprintf(
			ios.Out,
			"[%s] %s:%d:%d %s (%s)\n",
			formatter.Bold(result.Source),
			formatter.Bold(result.Location.Path),
			result.Location.StartLine,
			result.Location.StartColumn,
			formatter.Red(result.Rule.Description),
			result.Rule.ID,
		)
	}

	return nil
}
