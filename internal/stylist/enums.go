package stylist

//go:generate go-enum -f=$GOFILE --marshal --names

// InputType represents how files are passed to a command.
// ENUM(arg, stdin, variadic).
type InputType string

// OutputType represents how to parse command output.
// ENUM(json, none, regexp, sarif).
type OutputType string

// ProcessorType represents the type of processor.
// ENUM(formatter, linter).
type ProcessorType string
