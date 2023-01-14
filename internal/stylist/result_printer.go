package stylist

import (
	"fmt"

	"github.com/twelvelabs/termite/ioutil"
)

// ResultPrinter is the interface that wraps the Print method.
type ResultPrinter interface {
	// Print writes the results to Stdout.
	Print(results []*Result) error
}

// NewResultPrinter returns the appropriate printer for the given format.
func NewResultPrinter(ios *ioutil.IOStreams, config *Config) ResultPrinter { //nolint:ireturn
	format := config.Output.Format
	switch format {
	case ResultFormatSarif:
		return &SarifPrinter{ios: ios, config: config}
	case ResultFormatTty:
		return &TtyPrinter{ios: ios, config: config}
	default:
		panic(fmt.Sprintf("unknown result format: %s", format))
	}
}

/*
* SarifPrinter
**/

// SarifPrinter generates SARIF formatted output.
type SarifPrinter struct {
	ios    *ioutil.IOStreams
	config *Config
}

// Print writes the SARIF formatted results to Stdout.
func (p *SarifPrinter) Print(results []*Result) error {
	return nil
}

/*
* TtyPrinter
**/

// TtyPrinter generates TTY formatted output.
// The output will contain ANSI color codes if the terminal allows them.
type TtyPrinter struct {
	ios    *ioutil.IOStreams
	config *Config
}

// Print writes the TTY formatted results to Stdout.
func (p *TtyPrinter) Print(results []*Result) error {
	formatter := p.ios.Formatter()
	for _, result := range results {
		p.printLocation(result, formatter)
		p.printContext(result)
		p.printUnderLinePointer(result, formatter)
	}
	return nil
}

func (p *TtyPrinter) printLocation(result *Result, formatter *ioutil.Formatter) {
	msg := ""
	if result.Rule.Description != "" {
		msg = formatter.Red(result.Rule.Description)
	}

	rule := ""
	if result.Rule.ID != "" {
		rule = fmt.Sprintf("(%s)", result.Rule.ID)
	}

	fmt.Fprintf(
		p.ios.Out,
		"[%s] %s %s %s\n",
		formatter.Bold(result.Source),
		formatter.Bold(result.Location.String()),
		msg,
		rule,
	)
}

func (p *TtyPrinter) printContext(result *Result) {
	for _, line := range result.ContextLines {
		fmt.Fprintln(p.ios.Out, line)
	}
}

// Copied from golangci-lint
func (p *TtyPrinter) printUnderLinePointer(result *Result, formatter *ioutil.Formatter) {
	// StartColumn == 0 means "unknown".
	if len(result.ContextLines) != 1 || result.Location.StartColumn == 0 {
		return
	}

	col0 := result.Location.StartColumn - 1
	line := result.ContextLines[0]
	prefixRunes := make([]rune, 0, len(line))
	for j := 0; j < len(line) && j < col0; j++ {
		if line[j] == '\t' {
			prefixRunes = append(prefixRunes, '\t')
		} else {
			prefixRunes = append(prefixRunes, ' ')
		}
	}

	fmt.Fprintf(p.ios.Out, "%s%s\n", string(prefixRunes), formatter.Yellow("^"))
}
