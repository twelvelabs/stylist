package stylist

import (
	"fmt"

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
	return nil
}
