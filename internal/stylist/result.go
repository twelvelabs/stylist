package stylist

import (
	"fmt"
)

// Result describes a single result detected by a processor.
type Result struct {
	Source       string         `json:"source"`
	Level        ResultLevel    `json:"level"`
	Location     ResultLocation `json:"location"`
	Rule         ResultRule     `json:"rule"`
	ContextLines []string       `json:"context_lines"`
	ContextLang  string         `json:"context_lang"`
}

// ResultLocation describes the physical location where the result occurred.
type ResultLocation struct {
	Path        string `json:"path"`
	StartLine   int    `json:"start_line"`
	StartColumn int    `json:"start_column"`
	EndLine     int    `json:"end_line"`
	EndColumn   int    `json:"end_column"`
}

// Returns the start and end lines.
// If the end line is 0, then it defaults to the start line.
func (r ResultLocation) LineRange() (int, int) {
	if r.EndLine != 0 {
		return r.StartLine, r.EndLine
	}
	return r.StartLine, r.StartLine
}

// Returns a display string in the form of "path:line:col".
func (r ResultLocation) String() string {
	path := r.Path
	if path == "" {
		path = "<none>"
	}
	return fmt.Sprintf("%s:%d:%d", path, r.StartLine, r.StartColumn)
}

// ResultRule describes the rule that was evaluated to produce the result.
type ResultRule struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	URI         string `json:"uri"`
}

// Results is a sortable collection of results.
type Results []*Result

// Len implements sort.Interface.
func (r Results) Len() int { return len(r) }

// Swap implements sort.Interface.
func (r Results) Swap(i, j int) { r[i], r[j] = r[j], r[i] }

type ResultsByLocation struct{ Results }

func (r ResultsByLocation) Less(i, j int) bool {
	return resultsByLocation(r.Results, i, j)
}

type ResultsBySeverity struct{ Results }

func (r ResultsBySeverity) Less(i, j int) bool {
	if r.Results[i].Level != r.Results[j].Level {
		return r.Results[i].Level > r.Results[j].Level
	}
	return resultsByLocation(r.Results, i, j)
}

type ResultsBySource struct{ Results }

func (r ResultsBySource) Less(i, j int) bool {
	if r.Results[i].Source != r.Results[j].Source {
		return r.Results[i].Source < r.Results[j].Source
	}
	return resultsByLocation(r.Results, i, j)
}

func resultsByLocation(results []*Result, i, j int) bool {
	if results[i].Location.Path != results[j].Location.Path {
		return results[i].Location.Path < results[j].Location.Path
	}
	if results[i].Location.StartLine != results[j].Location.StartLine {
		return results[i].Location.StartLine < results[j].Location.StartLine
	}
	if results[i].Location.StartColumn != results[j].Location.StartColumn {
		return results[i].Location.StartColumn < results[j].Location.StartColumn
	}
	// Additional tie-breakers for deterministic results.
	if results[i].Source != results[j].Source {
		return results[i].Source < results[j].Source
	}
	return results[i].Level > results[j].Level
}

// NewResultsError returns a new error when the results slice is non-empty.
func NewResultsError(results []*Result) error {
	if len(results) == 0 {
		return nil
	}
	return &ResultsError{
		results: results,
	}
}

// ResultsError is a sentinel type returned by actions when there are results.
type ResultsError struct {
	results []*Result
}

// Error implements the error interface.
func (re *ResultsError) Error() string {
	return fmt.Sprintf("%d issue(s)", len(re.results))
}
