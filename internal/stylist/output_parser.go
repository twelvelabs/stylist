package stylist

import (
	"bytes"
	"fmt"
	"io"
	"regexp"
	"strings"

	"github.com/owenrumney/go-sarif/sarif"
	"github.com/sourcegraph/go-diff/diff"
	"github.com/tidwall/gjson"

	"github.com/twelvelabs/stylist/internal/fsutils"
)

var (
	ansiRegexpStr = "[\u001B\u009B][[\\]()#;?]*(?:" +
		"(?:(?:[a-zA-Z\\d]*(?:;[a-zA-Z\\d]*)*)?\u0007)|" +
		"(?:(?:\\d{1,4}(?:;\\d{0,4})*)?[\\dA-PRZcf-ntqry=><~]))" // cspell: disable-line
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
	case OutputFormatDiff:
		return &DiffOutputParser{}
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
* DiffOutputParser
**/

// DiffOutputParser parses unified diffs.
type DiffOutputParser struct {
}

// Parse parses command output into a slice of results.
func (p *DiffOutputParser) Parse(output CommandOutput, mapping ResultMapping) ([]*Result, error) {
	// Read the content.
	buf, err := io.ReadAll(output.Content)
	if err != nil {
		return nil, err
	}
	content := ansiRegexp.ReplaceAll(buf, []byte(""))
	if len(content) == 0 {
		return nil, nil // nothing to parse
	}

	// Parse.
	diffs, err := diff.ParseMultiFileDiff(content)
	if err != nil {
		return nil, fmt.Errorf("invalid diff: %w", err)
	}

	// Map the diffs to a slice of `Result` structs.
	results := []*Result{}
	for _, d := range diffs {
		var startLine int
		var contextLines []string

		if len(d.Hunks) > 0 {
			// Hunks often start w/ a few preceding context lines.
			// Calculate the actual start of the changeset.
			changeStart := 0
			bodyLines := bytes.Split(d.Hunks[0].Body, []byte{'\n'})
			for _, line := range bodyLines {
				if len(line) > 0 && (line[0] == '+' || line[0] == '-') {
					break
				}
				changeStart++
			}
			startLine = int(d.Hunks[0].OrigStartLine) + changeStart

			// Printing just the hunks (vs full diff) so we don't have
			// redundant file names at the top of the context.
			hunks, _ := diff.PrintHunks(d.Hunks)
			contextLines = strings.Split(strings.TrimSuffix(string(hunks), "\n"), "\n")
		}

		result := &Result{
			Level: ResultLevelError,
			Location: ResultLocation{
				Path:      d.NewName,
				StartLine: startLine,
			},
			Rule: ResultRule{
				ID:          "diff",
				Name:        "diff",
				Description: "Formatting error",
			},
			ContextLines: contextLines,
			ContextLang:  "diff",
		}
		results = append(results, result)
	}

	return results, nil
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
	// Read the content.
	buf, err := io.ReadAll(output.Content)
	if err != nil {
		return nil, err
	}
	content := ansiRegexp.ReplaceAllString(string(buf), "")
	if content == "" {
		return nil, nil // nothing to parse
	}

	// Parse.
	report, err := sarif.FromString(content)
	if err != nil {
		return nil, fmt.Errorf("invalid sarif: %w", err)
	}

	// Map the report to a slice of `Result` structs.
	issues := []*Result{}
	for _, run := range report.Runs {
		for _, result := range run.Results {
			level, err := CoerceResultLevel(*result.Level)
			if err != nil {
				return nil, fmt.Errorf("invalid sarif level: %w", err)
			}

			rule := ResultRule{}
			rule.ID = *result.RuleID
			rule.Name = *result.RuleID
			rule.Description = *result.Message.Text

			loc := ResultLocation{}
			if len(result.Locations) > 0 {
				loc, err = resultLocationFromSarif(result.Locations[0])
				if err != nil {
					return nil, fmt.Errorf("invalid sarif location: %w", err)
				}
			}

			issue := &Result{
				Level:    level,
				Location: loc,
				Rule:     rule,
			}

			issues = append(issues, issue)
		}
	}

	return issues, nil
}

func resultLocationFromSarif(loc *sarif.Location) (ResultLocation, error) {
	var err error

	resultLocation := ResultLocation{}

	pl := loc.PhysicalLocation
	if pl != nil {
		al := pl.ArtifactLocation
		if al != nil {
			resultLocation.Path, err = fsutils.RelativePath(*al.URI)
			if err != nil {
				return resultLocation, err
			}
		}
	}
	reg := pl.Region
	if reg != nil {
		resultLocation.StartLine = *reg.StartLine
		resultLocation.StartColumn = *reg.StartColumn
		resultLocation.EndLine = *reg.EndLine
		resultLocation.EndColumn = *reg.EndColumn
	}

	return resultLocation, nil
}
