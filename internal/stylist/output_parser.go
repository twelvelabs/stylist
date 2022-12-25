package stylist

import (
	"encoding/json"
	"fmt"
)

// OutputParser is the interface that wraps the Parse method.
//
// Parse parses command output into a slice of results.
type OutputParser interface {
	Parse(output CommandOutput, mapping OutputMapping) ([]*Result, error)
}

// NewOutputParser returns the appropriate parser for the given output type.
func NewOutputParser(format OutputFormat) OutputParser { //nolint:ireturn
	switch format {
	case OutputFormatJson:
		return &JSONOutputParser{}
	case OutputFormatNone:
		return &NoneOutputParser{}
	case OutputFormatRegexp:
		return &RegexpOutputParser{}
	case OutputFormatSarif:
		return &SarifOutputParser{}
	default:
		panic(fmt.Sprintf("unknown output format: %s", format))
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
	var items []outputData
	err := json.NewDecoder(output.Content).Decode(&items)
	if err != nil {
		return nil, err
	}
	return mapping.ToResultSlice(items)
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
