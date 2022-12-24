package stylist

import (
	"fmt"
)

// OutputParser is the interface that wraps the Parse method.
//
// Parse parses command output into a slice of results.
type OutputParser interface {
	Parse(output CommandOutput, mapping OutputMapping) ([]*Result, error)
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

// Parse parses command output into a slice of results.
func (p *JSONOutputParser) Parse(output CommandOutput, mapping OutputMapping) ([]*Result, error) {
	return nil, nil
}

/*
* NoneOutputParser
**/

// NoneOutputParser is a noop parser for commands that produce no output.
type NoneOutputParser struct {
}

// Parse parses command output into a slice of results.
func (p *NoneOutputParser) Parse(output CommandOutput, mapping OutputMapping) ([]*Result, error) {
	return nil, nil
}

/*
* RegexpOutputParser
**/

// RegexpOutputParser parses arbitrary text output using regular expressions.
type RegexpOutputParser struct {
}

// Parse parses command output into a slice of results.
func (p *RegexpOutputParser) Parse(output CommandOutput, mapping OutputMapping) ([]*Result, error) {
	return nil, nil
}

/*
* SarifOutputParser
**/

// SarifOutputParser parses SARIF formatted output.
type SarifOutputParser struct {
}

// Parse parses command output into a slice of results.
func (p *SarifOutputParser) Parse(output CommandOutput, mapping OutputMapping) ([]*Result, error) {
	return nil, nil
}
