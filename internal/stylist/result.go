package stylist

import (
	"fmt"
)

// Result describes a single result detected by a processor.
type Result struct {
	Source       string
	Level        ResultLevel
	Location     ResultLocation
	Rule         ResultRule
	ContextLines []string
}

// ResultLocation describes the physical location where the result occurred.
type ResultLocation struct {
	Path        string
	StartLine   int
	StartColumn int
	EndLine     int
	EndColumn   int
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
	ID          string
	Name        string
	Description string
	URI         string
}
