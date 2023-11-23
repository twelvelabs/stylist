// Code generated by go-enum DO NOT EDIT.
// Version:
// Revision:
// Build Date:
// Built By:

package stylist

import (
	"fmt"
	"strings"
)

const (
	// CommandTypeCheck is a CommandType of type check.
	CommandTypeCheck CommandType = "check"
	// CommandTypeFix is a CommandType of type fix.
	CommandTypeFix CommandType = "fix"
)

var ErrInvalidCommandType = fmt.Errorf("not a valid CommandType, try [%s]", strings.Join(_CommandTypeNames, ", "))

var _CommandTypeNames = []string{
	string(CommandTypeCheck),
	string(CommandTypeFix),
}

// CommandTypeNames returns a list of possible string values of CommandType.
func CommandTypeNames() []string {
	tmp := make([]string, len(_CommandTypeNames))
	copy(tmp, _CommandTypeNames)
	return tmp
}

// String implements the Stringer interface.
func (x CommandType) String() string {
	return string(x)
}

// IsValid provides a quick way to determine if the typed value is
// part of the allowed enumerated values
func (x CommandType) IsValid() bool {
	_, err := ParseCommandType(string(x))
	return err == nil
}

var _CommandTypeValue = map[string]CommandType{
	"check": CommandTypeCheck,
	"fix":   CommandTypeFix,
}

// ParseCommandType attempts to convert a string to a CommandType.
func ParseCommandType(name string) (CommandType, error) {
	if x, ok := _CommandTypeValue[name]; ok {
		return x, nil
	}
	return CommandType(""), fmt.Errorf("%s is %w", name, ErrInvalidCommandType)
}

// MarshalText implements the text marshaller method.
func (x CommandType) MarshalText() ([]byte, error) {
	return []byte(string(x)), nil
}

// UnmarshalText implements the text unmarshaller method.
func (x *CommandType) UnmarshalText(text []byte) error {
	tmp, err := ParseCommandType(string(text))
	if err != nil {
		return err
	}
	*x = tmp
	return nil
}

// Set implements the Golang flag.Value interface func.
func (x *CommandType) Set(val string) error {
	v, err := ParseCommandType(val)
	*x = v
	return err
}

// Get implements the Golang flag.Getter interface func.
func (x *CommandType) Get() interface{} {
	return *x
}

// Type implements the github.com/spf13/pFlag Value interface.
func (x *CommandType) Type() string {
	return "CommandType"
}

const (
	// InputTypeArg is a InputType of type arg.
	InputTypeArg InputType = "arg"
	// InputTypeNone is a InputType of type none.
	InputTypeNone InputType = "none"
	// InputTypeStdin is a InputType of type stdin.
	InputTypeStdin InputType = "stdin"
	// InputTypeVariadic is a InputType of type variadic.
	InputTypeVariadic InputType = "variadic"
)

var ErrInvalidInputType = fmt.Errorf("not a valid InputType, try [%s]", strings.Join(_InputTypeNames, ", "))

var _InputTypeNames = []string{
	string(InputTypeArg),
	string(InputTypeNone),
	string(InputTypeStdin),
	string(InputTypeVariadic),
}

// InputTypeNames returns a list of possible string values of InputType.
func InputTypeNames() []string {
	tmp := make([]string, len(_InputTypeNames))
	copy(tmp, _InputTypeNames)
	return tmp
}

// String implements the Stringer interface.
func (x InputType) String() string {
	return string(x)
}

// IsValid provides a quick way to determine if the typed value is
// part of the allowed enumerated values
func (x InputType) IsValid() bool {
	_, err := ParseInputType(string(x))
	return err == nil
}

var _InputTypeValue = map[string]InputType{
	"arg":      InputTypeArg,
	"none":     InputTypeNone,
	"stdin":    InputTypeStdin,
	"variadic": InputTypeVariadic,
}

// ParseInputType attempts to convert a string to a InputType.
func ParseInputType(name string) (InputType, error) {
	if x, ok := _InputTypeValue[name]; ok {
		return x, nil
	}
	return InputType(""), fmt.Errorf("%s is %w", name, ErrInvalidInputType)
}

// MarshalText implements the text marshaller method.
func (x InputType) MarshalText() ([]byte, error) {
	return []byte(string(x)), nil
}

// UnmarshalText implements the text unmarshaller method.
func (x *InputType) UnmarshalText(text []byte) error {
	tmp, err := ParseInputType(string(text))
	if err != nil {
		return err
	}
	*x = tmp
	return nil
}

// Set implements the Golang flag.Value interface func.
func (x *InputType) Set(val string) error {
	v, err := ParseInputType(val)
	*x = v
	return err
}

// Get implements the Golang flag.Getter interface func.
func (x *InputType) Get() interface{} {
	return *x
}

// Type implements the github.com/spf13/pFlag Value interface.
func (x *InputType) Type() string {
	return "InputType"
}

const (
	// LogLevelError is a LogLevel of type error.
	LogLevelError LogLevel = "error"
	// LogLevelWarn is a LogLevel of type warn.
	LogLevelWarn LogLevel = "warn"
	// LogLevelInfo is a LogLevel of type info.
	LogLevelInfo LogLevel = "info"
	// LogLevelDebug is a LogLevel of type debug.
	LogLevelDebug LogLevel = "debug"
)

var ErrInvalidLogLevel = fmt.Errorf("not a valid LogLevel, try [%s]", strings.Join(_LogLevelNames, ", "))

var _LogLevelNames = []string{
	string(LogLevelError),
	string(LogLevelWarn),
	string(LogLevelInfo),
	string(LogLevelDebug),
}

// LogLevelNames returns a list of possible string values of LogLevel.
func LogLevelNames() []string {
	tmp := make([]string, len(_LogLevelNames))
	copy(tmp, _LogLevelNames)
	return tmp
}

// String implements the Stringer interface.
func (x LogLevel) String() string {
	return string(x)
}

// IsValid provides a quick way to determine if the typed value is
// part of the allowed enumerated values
func (x LogLevel) IsValid() bool {
	_, err := ParseLogLevel(string(x))
	return err == nil
}

var _LogLevelValue = map[string]LogLevel{
	"error": LogLevelError,
	"warn":  LogLevelWarn,
	"info":  LogLevelInfo,
	"debug": LogLevelDebug,
}

// ParseLogLevel attempts to convert a string to a LogLevel.
func ParseLogLevel(name string) (LogLevel, error) {
	if x, ok := _LogLevelValue[name]; ok {
		return x, nil
	}
	return LogLevel(""), fmt.Errorf("%s is %w", name, ErrInvalidLogLevel)
}

// MarshalText implements the text marshaller method.
func (x LogLevel) MarshalText() ([]byte, error) {
	return []byte(string(x)), nil
}

// UnmarshalText implements the text unmarshaller method.
func (x *LogLevel) UnmarshalText(text []byte) error {
	tmp, err := ParseLogLevel(string(text))
	if err != nil {
		return err
	}
	*x = tmp
	return nil
}

// Set implements the Golang flag.Value interface func.
func (x *LogLevel) Set(val string) error {
	v, err := ParseLogLevel(val)
	*x = v
	return err
}

// Get implements the Golang flag.Getter interface func.
func (x *LogLevel) Get() interface{} {
	return *x
}

// Type implements the github.com/spf13/pFlag Value interface.
func (x *LogLevel) Type() string {
	return "LogLevel"
}

const (
	// OutputFormatDiff is a OutputFormat of type diff.
	OutputFormatDiff OutputFormat = "diff"
	// OutputFormatJson is a OutputFormat of type json.
	OutputFormatJson OutputFormat = "json"
	// OutputFormatNone is a OutputFormat of type none.
	OutputFormatNone OutputFormat = "none"
	// OutputFormatRegexp is a OutputFormat of type regexp.
	OutputFormatRegexp OutputFormat = "regexp"
	// OutputFormatSarif is a OutputFormat of type sarif.
	OutputFormatSarif OutputFormat = "sarif"
)

var ErrInvalidOutputFormat = fmt.Errorf("not a valid OutputFormat, try [%s]", strings.Join(_OutputFormatNames, ", "))

var _OutputFormatNames = []string{
	string(OutputFormatDiff),
	string(OutputFormatJson),
	string(OutputFormatNone),
	string(OutputFormatRegexp),
	string(OutputFormatSarif),
}

// OutputFormatNames returns a list of possible string values of OutputFormat.
func OutputFormatNames() []string {
	tmp := make([]string, len(_OutputFormatNames))
	copy(tmp, _OutputFormatNames)
	return tmp
}

// String implements the Stringer interface.
func (x OutputFormat) String() string {
	return string(x)
}

// IsValid provides a quick way to determine if the typed value is
// part of the allowed enumerated values
func (x OutputFormat) IsValid() bool {
	_, err := ParseOutputFormat(string(x))
	return err == nil
}

var _OutputFormatValue = map[string]OutputFormat{
	"diff":   OutputFormatDiff,
	"json":   OutputFormatJson,
	"none":   OutputFormatNone,
	"regexp": OutputFormatRegexp,
	"sarif":  OutputFormatSarif,
}

// ParseOutputFormat attempts to convert a string to a OutputFormat.
func ParseOutputFormat(name string) (OutputFormat, error) {
	if x, ok := _OutputFormatValue[name]; ok {
		return x, nil
	}
	return OutputFormat(""), fmt.Errorf("%s is %w", name, ErrInvalidOutputFormat)
}

// MarshalText implements the text marshaller method.
func (x OutputFormat) MarshalText() ([]byte, error) {
	return []byte(string(x)), nil
}

// UnmarshalText implements the text unmarshaller method.
func (x *OutputFormat) UnmarshalText(text []byte) error {
	tmp, err := ParseOutputFormat(string(text))
	if err != nil {
		return err
	}
	*x = tmp
	return nil
}

// Set implements the Golang flag.Value interface func.
func (x *OutputFormat) Set(val string) error {
	v, err := ParseOutputFormat(val)
	*x = v
	return err
}

// Get implements the Golang flag.Getter interface func.
func (x *OutputFormat) Get() interface{} {
	return *x
}

// Type implements the github.com/spf13/pFlag Value interface.
func (x *OutputFormat) Type() string {
	return "OutputFormat"
}

const (
	// OutputTypeStdout is a OutputType of type stdout.
	OutputTypeStdout OutputType = "stdout"
	// OutputTypeStderr is a OutputType of type stderr.
	OutputTypeStderr OutputType = "stderr"
)

var ErrInvalidOutputType = fmt.Errorf("not a valid OutputType, try [%s]", strings.Join(_OutputTypeNames, ", "))

var _OutputTypeNames = []string{
	string(OutputTypeStdout),
	string(OutputTypeStderr),
}

// OutputTypeNames returns a list of possible string values of OutputType.
func OutputTypeNames() []string {
	tmp := make([]string, len(_OutputTypeNames))
	copy(tmp, _OutputTypeNames)
	return tmp
}

// String implements the Stringer interface.
func (x OutputType) String() string {
	return string(x)
}

// IsValid provides a quick way to determine if the typed value is
// part of the allowed enumerated values
func (x OutputType) IsValid() bool {
	_, err := ParseOutputType(string(x))
	return err == nil
}

var _OutputTypeValue = map[string]OutputType{
	"stdout": OutputTypeStdout,
	"stderr": OutputTypeStderr,
}

// ParseOutputType attempts to convert a string to a OutputType.
func ParseOutputType(name string) (OutputType, error) {
	if x, ok := _OutputTypeValue[name]; ok {
		return x, nil
	}
	return OutputType(""), fmt.Errorf("%s is %w", name, ErrInvalidOutputType)
}

// MarshalText implements the text marshaller method.
func (x OutputType) MarshalText() ([]byte, error) {
	return []byte(string(x)), nil
}

// UnmarshalText implements the text unmarshaller method.
func (x *OutputType) UnmarshalText(text []byte) error {
	tmp, err := ParseOutputType(string(text))
	if err != nil {
		return err
	}
	*x = tmp
	return nil
}

// Set implements the Golang flag.Value interface func.
func (x *OutputType) Set(val string) error {
	v, err := ParseOutputType(val)
	*x = v
	return err
}

// Get implements the Golang flag.Getter interface func.
func (x *OutputType) Get() interface{} {
	return *x
}

// Type implements the github.com/spf13/pFlag Value interface.
func (x *OutputType) Type() string {
	return "OutputType"
}

const (
	// ResultFormatCheckstyle is a ResultFormat of type checkstyle.
	ResultFormatCheckstyle ResultFormat = "checkstyle"
	// ResultFormatSarif is a ResultFormat of type sarif.
	ResultFormatSarif ResultFormat = "sarif"
	// ResultFormatTty is a ResultFormat of type tty.
	ResultFormatTty ResultFormat = "tty"
)

var ErrInvalidResultFormat = fmt.Errorf("not a valid ResultFormat, try [%s]", strings.Join(_ResultFormatNames, ", "))

var _ResultFormatNames = []string{
	string(ResultFormatCheckstyle),
	string(ResultFormatSarif),
	string(ResultFormatTty),
}

// ResultFormatNames returns a list of possible string values of ResultFormat.
func ResultFormatNames() []string {
	tmp := make([]string, len(_ResultFormatNames))
	copy(tmp, _ResultFormatNames)
	return tmp
}

// String implements the Stringer interface.
func (x ResultFormat) String() string {
	return string(x)
}

// IsValid provides a quick way to determine if the typed value is
// part of the allowed enumerated values
func (x ResultFormat) IsValid() bool {
	_, err := ParseResultFormat(string(x))
	return err == nil
}

var _ResultFormatValue = map[string]ResultFormat{
	"checkstyle": ResultFormatCheckstyle,
	"sarif":      ResultFormatSarif,
	"tty":        ResultFormatTty,
}

// ParseResultFormat attempts to convert a string to a ResultFormat.
func ParseResultFormat(name string) (ResultFormat, error) {
	if x, ok := _ResultFormatValue[name]; ok {
		return x, nil
	}
	return ResultFormat(""), fmt.Errorf("%s is %w", name, ErrInvalidResultFormat)
}

// MarshalText implements the text marshaller method.
func (x ResultFormat) MarshalText() ([]byte, error) {
	return []byte(string(x)), nil
}

// UnmarshalText implements the text unmarshaller method.
func (x *ResultFormat) UnmarshalText(text []byte) error {
	tmp, err := ParseResultFormat(string(text))
	if err != nil {
		return err
	}
	*x = tmp
	return nil
}

// Set implements the Golang flag.Value interface func.
func (x *ResultFormat) Set(val string) error {
	v, err := ParseResultFormat(val)
	*x = v
	return err
}

// Get implements the Golang flag.Getter interface func.
func (x *ResultFormat) Get() interface{} {
	return *x
}

// Type implements the github.com/spf13/pFlag Value interface.
func (x *ResultFormat) Type() string {
	return "ResultFormat"
}

const (
	// ResultLevelNone is a ResultLevel of type None.
	ResultLevelNone ResultLevel = iota
	// ResultLevelInfo is a ResultLevel of type Info.
	ResultLevelInfo
	// ResultLevelWarning is a ResultLevel of type Warning.
	ResultLevelWarning
	// ResultLevelError is a ResultLevel of type Error.
	ResultLevelError
)

var ErrInvalidResultLevel = fmt.Errorf("not a valid ResultLevel, try [%s]", strings.Join(_ResultLevelNames, ", "))

const _ResultLevelName = "noneinfowarningerror"

var _ResultLevelNames = []string{
	_ResultLevelName[0:4],
	_ResultLevelName[4:8],
	_ResultLevelName[8:15],
	_ResultLevelName[15:20],
}

// ResultLevelNames returns a list of possible string values of ResultLevel.
func ResultLevelNames() []string {
	tmp := make([]string, len(_ResultLevelNames))
	copy(tmp, _ResultLevelNames)
	return tmp
}

var _ResultLevelMap = map[ResultLevel]string{
	ResultLevelNone:    _ResultLevelName[0:4],
	ResultLevelInfo:    _ResultLevelName[4:8],
	ResultLevelWarning: _ResultLevelName[8:15],
	ResultLevelError:   _ResultLevelName[15:20],
}

// String implements the Stringer interface.
func (x ResultLevel) String() string {
	if str, ok := _ResultLevelMap[x]; ok {
		return str
	}
	return fmt.Sprintf("ResultLevel(%d)", x)
}

// IsValid provides a quick way to determine if the typed value is
// part of the allowed enumerated values
func (x ResultLevel) IsValid() bool {
	_, ok := _ResultLevelMap[x]
	return ok
}

var _ResultLevelValue = map[string]ResultLevel{
	_ResultLevelName[0:4]:   ResultLevelNone,
	_ResultLevelName[4:8]:   ResultLevelInfo,
	_ResultLevelName[8:15]:  ResultLevelWarning,
	_ResultLevelName[15:20]: ResultLevelError,
}

// ParseResultLevel attempts to convert a string to a ResultLevel.
func ParseResultLevel(name string) (ResultLevel, error) {
	if x, ok := _ResultLevelValue[name]; ok {
		return x, nil
	}
	return ResultLevel(0), fmt.Errorf("%s is %w", name, ErrInvalidResultLevel)
}

// MarshalText implements the text marshaller method.
func (x ResultLevel) MarshalText() ([]byte, error) {
	return []byte(x.String()), nil
}

// UnmarshalText implements the text unmarshaller method.
func (x *ResultLevel) UnmarshalText(text []byte) error {
	name := string(text)
	tmp, err := ParseResultLevel(name)
	if err != nil {
		return err
	}
	*x = tmp
	return nil
}

// Set implements the Golang flag.Value interface func.
func (x *ResultLevel) Set(val string) error {
	v, err := ParseResultLevel(val)
	*x = v
	return err
}

// Get implements the Golang flag.Getter interface func.
func (x *ResultLevel) Get() interface{} {
	return *x
}

// Type implements the github.com/spf13/pFlag Value interface.
func (x *ResultLevel) Type() string {
	return "ResultLevel"
}

const (
	// ResultSortLocation is a ResultSort of type location.
	ResultSortLocation ResultSort = "location"
	// ResultSortSeverity is a ResultSort of type severity.
	ResultSortSeverity ResultSort = "severity"
	// ResultSortSource is a ResultSort of type source.
	ResultSortSource ResultSort = "source"
)

var ErrInvalidResultSort = fmt.Errorf("not a valid ResultSort, try [%s]", strings.Join(_ResultSortNames, ", "))

var _ResultSortNames = []string{
	string(ResultSortLocation),
	string(ResultSortSeverity),
	string(ResultSortSource),
}

// ResultSortNames returns a list of possible string values of ResultSort.
func ResultSortNames() []string {
	tmp := make([]string, len(_ResultSortNames))
	copy(tmp, _ResultSortNames)
	return tmp
}

// String implements the Stringer interface.
func (x ResultSort) String() string {
	return string(x)
}

// IsValid provides a quick way to determine if the typed value is
// part of the allowed enumerated values
func (x ResultSort) IsValid() bool {
	_, err := ParseResultSort(string(x))
	return err == nil
}

var _ResultSortValue = map[string]ResultSort{
	"location": ResultSortLocation,
	"severity": ResultSortSeverity,
	"source":   ResultSortSource,
}

// ParseResultSort attempts to convert a string to a ResultSort.
func ParseResultSort(name string) (ResultSort, error) {
	if x, ok := _ResultSortValue[name]; ok {
		return x, nil
	}
	return ResultSort(""), fmt.Errorf("%s is %w", name, ErrInvalidResultSort)
}

// MarshalText implements the text marshaller method.
func (x ResultSort) MarshalText() ([]byte, error) {
	return []byte(string(x)), nil
}

// UnmarshalText implements the text unmarshaller method.
func (x *ResultSort) UnmarshalText(text []byte) error {
	tmp, err := ParseResultSort(string(text))
	if err != nil {
		return err
	}
	*x = tmp
	return nil
}

// Set implements the Golang flag.Value interface func.
func (x *ResultSort) Set(val string) error {
	v, err := ParseResultSort(val)
	*x = v
	return err
}

// Get implements the Golang flag.Getter interface func.
func (x *ResultSort) Get() interface{} {
	return *x
}

// Type implements the github.com/spf13/pFlag Value interface.
func (x *ResultSort) Type() string {
	return "ResultSort"
}
