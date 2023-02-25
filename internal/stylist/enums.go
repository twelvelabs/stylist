package stylist

//go:generate go-enum -f=$GOFILE --marshal --names --flag

// CommandType represents the type of command.
//
// ENUM(check, fix).
type CommandType string

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
// ENUM(diff, json, none, regexp, sarif).
type OutputFormat string

// ResultLevel represents the severity level of the result.
// These values were chosen to match those in the SARIF specification.
//
// ENUM(none, note, warning, error).
type ResultLevel string

// CoerceResultLevel returns the correct enum for the given value.
func CoerceResultLevel(value string) (ResultLevel, error) {
	switch value {
	case "", "<no value>":
		return ResultLevelNone, nil
	case "info":
		return ResultLevelNote, nil
	case "warn":
		return ResultLevelWarning, nil
	case "err":
		return ResultLevelError, nil
	default:
		return ParseResultLevel(value)
	}
}

// ResultFormat represents how to format the results.
//
// ENUM(sarif, tty).
type ResultFormat string

// LogLevel controls the log level.
//
// ENUM(error, warn, info, debug).
type LogLevel string
