package stylist

// Result describes a single result detected by a processor.
type Result struct {
	Level    ResultLevel
	Location ResultLocation
	Rule     ResultRule
}

// ResultLocation describes the physical location where the result occurred.
type ResultLocation struct {
	Path        string
	StartLine   int
	StartColumn int
	EndLine     int
	EndColumn   int
}

// ResultRule describes the rule that was evaluated to produce the result.
type ResultRule struct {
	ID          string
	Name        string
	Description string
	URI         string
}
