package stylist

import (
	"bytes"
	"fmt"
	"io"
	"regexp"

	"github.com/tidwall/gjson"
)

var (
	ansiRegexpStr = "[\u001B\u009B][[\\]()#;?]*(?:" +
		"(?:(?:[a-zA-Z\\d]*(?:;[a-zA-Z\\d]*)*)?\u0007)|" +
		"(?:(?:\\d{1,4}(?:;\\d{0,4})*)?[\\dA-PRZcf-ntqry=><~]))"
	ansiRegexp = regexp.MustCompile(ansiRegexpStr)
)

// OutputParser is the interface that wraps the Parse method.
//
// Parse parses command output into a slice of results.
type OutputParser interface {
	Parse(output CommandOutput, mapping ResultMapping) ([]*Result, error)
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
func (p *JSONOutputParser) Parse(output CommandOutput, mapping ResultMapping) ([]*Result, error) {
	buf := &bytes.Buffer{}
	_, err := buf.ReadFrom(output.Content)
	if err != nil {
		return nil, err
	}

	json := buf.String()
	if json == "" {
		// No results
		return nil, nil
	}

	// Ensure valid JSON
	if !gjson.Valid(json) {
		return nil, fmt.Errorf("invalid json: %s", json)
	}

	// Parse the JSON and return the data @ pattern.
	// Note: `@this` is how GJSON addresses the root element.
	pattern := "@this"
	if mapping.Pattern != "" {
		pattern = mapping.Pattern
	}
	result := gjson.Get(json, pattern)

	// Transform the GJSON results into resultData
	items := []resultData{}
	if !result.IsArray() {
		return nil, fmt.Errorf(
			"invalid output: pattern=%v is not an array, json=%v",
			pattern, json,
		)
	}
	for idx, r := range result.Array() {
		if !r.IsObject() {
			return nil, fmt.Errorf(
				"invalid output: pattern=%v.%v is not an object, json=%v",
				pattern, idx, json,
			)
		}
		item := r.Value().(map[string]any)
		items = append(items, resultData(item))
	}

	// Transform the resultData into `Result` structs.
	return mapping.ToResultSlice(items)
}

/*
* NoneOutputParser
**/

// NoneOutputParser is a noop parser for commands that produce no output.
type NoneOutputParser struct {
}

// Parse parses command output into a slice of results.
func (p *NoneOutputParser) Parse(output CommandOutput, mapping ResultMapping) ([]*Result, error) {
	return nil, nil
}

/*
* RegexpOutputParser
**/

// RegexpOutputParser parses arbitrary text output using regular expressions.
type RegexpOutputParser struct {
}

// Parse parses command output into a slice of results.
func (p *RegexpOutputParser) Parse(output CommandOutput, mapping ResultMapping) ([]*Result, error) {
	// Validate the regexp pattern.
	if mapping.Pattern == "" {
		return nil, fmt.Errorf("mapping pattern is required when output format is regexp")
	}
	r, err := regexp.Compile(mapping.Pattern)
	if err != nil {
		return nil, fmt.Errorf("mapping pattern: %w", err)
	}

	// Read the content.
	buf, err := io.ReadAll(output.Content)
	if err != nil {
		return nil, err
	}
	content := ansiRegexp.ReplaceAllString(string(buf), "")
	if content == "" {
		return nil, nil // nothing to parse
	}

	// Run the regexp
	matches := r.FindAllStringSubmatch(content, -1)
	if len(matches) == 0 {
		return nil, nil // nothing found
	}
	keys := r.SubexpNames()

	// Convert the regexp matches into a slice of resultData maps.
	items := []resultData{}
	for _, match := range matches {
		item := resultData{}
		for i := 1; i < len(keys); i++ {
			item[keys[i]] = match[i]
		}
		items = append(items, item)
	}

	// Transform the resultData into `Result` structs.
	return mapping.ToResultSlice(items)
}

/*
* SarifOutputParser
**/

// SarifOutputParser parses SARIF formatted output.
type SarifOutputParser struct {
}

// Parse parses command output into a slice of results.
func (p *SarifOutputParser) Parse(output CommandOutput, mapping ResultMapping) ([]*Result, error) {
	return nil, nil
}
