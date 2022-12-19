package stylist

import (
	"fmt"
	"io"
)

// OutputParser is the interface that wraps the Parse method.
//
// Parse parses output into a slice of diagnostic pointers.
type OutputParser interface {
	Parse(cmd *Command, r io.Reader) ([]*Diagnostic, error)
}

// NewOutputParser returns the appropriate parser for the given output type.
func NewOutputParser(ot OutputType) OutputParser { //nolint:ireturn
	switch ot {
	case OutputTypeJson:
		return &JSONOutputParser{}
	case OutputTypeNone:
		return &NoneOutputParser{}
	case OutputTypeRegexp:
		return &RegexpOutputParser{}
	case OutputTypeSarif:
		return &SarifOutputParser{}
	default:
		panic(fmt.Sprintf("unknown output type: %s", ot))
	}
}

/*
* JSONOutputParser
**/

// JSONOutputParser parses JSON formatted output.
type JSONOutputParser struct {
}

// Parse parses output into a slice of diagnostic pointers.
func (p *JSONOutputParser) Parse(cmd *Command, r io.Reader) ([]*Diagnostic, error) {
	return nil, nil
}

/*
* NoneOutputParser
**/

// NoneOutputParser is a noop parser for commands that produce no output.
type NoneOutputParser struct {
}

// Parse parses output into a slice of diagnostic pointers.
func (p *NoneOutputParser) Parse(cmd *Command, r io.Reader) ([]*Diagnostic, error) {
	return nil, nil
}

/*
* RegexpOutputParser
**/

// RegexpOutputParser parses arbitrary text output using regular expressions.
type RegexpOutputParser struct {
}

// Parse parses output into a slice of diagnostic pointers.
func (p *RegexpOutputParser) Parse(cmd *Command, r io.Reader) ([]*Diagnostic, error) {
	return nil, nil
}

/*
* SarifOutputParser
**/

// SarifOutputParser parses SARIF formatted output.
type SarifOutputParser struct {
}

// Parse parses output into a slice of diagnostic pointers.
func (p *SarifOutputParser) Parse(cmd *Command, r io.Reader) ([]*Diagnostic, error) {
	return nil, nil
}
