package stylist

//go:generate go-enum -f=$GOFILE --marshal --names --flag

// InputType represents how files are passed to a command.
//
// ENUM(arg, none, stdin, variadic).
type InputType string

// OutputType represents where command output is sent.
//
// ENUM(stdout, stderr).
type OutputType string

// OutputFormat represents how to parse command output.
//
// ENUM(json, none, regexp, sarif).
type OutputFormat string

// ResultLevel represents the severity level of the result.
// These values were chosen to match those in the SARIF specification.
//
// ENUM(none, note, warning, error).
type ResultLevel string

// ResultFormat represents how to format the results.
//
// ENUM(sarif, tty).
type ResultFormat string

// LogLevel controls the log level.
//
// ENUM(error, warn, info, debug).
type LogLevel string
