package stylist

import (
	"strconv"

	"github.com/twelvelabs/stylist/internal/render"
)

const (
	strEmpty   string = ""
	strNoValue string = "<no value>"
)

// outputData is a map of key/value pairs parsed from CommandOutput
// and passed to an OutputMapping to be converted into a Result.
type outputData map[string]any

// OutputMapping is a set of rules for how to map command output to a result.
//
// Mappings are typically defined in stylist.yml when the output type
// has been set to "json" or "regexp".
type OutputMapping struct {
	Level           *render.Template `yaml:"level"`
	Path            *render.Template `yaml:"path"`
	StartLine       *render.Template `yaml:"start_line"`
	StartColumn     *render.Template `yaml:"start_column"`
	EndLine         *render.Template `yaml:"end_line"`
	EndColumn       *render.Template `yaml:"end_column"`
	RuleID          *render.Template `yaml:"rule_id"`
	RuleName        *render.Template `yaml:"rule_name"`
	RuleDescription *render.Template `yaml:"rule_description"`
	RuleURI         *render.Template `yaml:"rule_url"`
}

// ToResult converts a map of output data to a Result struct.
func (m OutputMapping) ToResult(item outputData) (*Result, error) {
	var err error

	result := &Result{
		Location: ResultLocation{},
		Rule:     ResultRule{},
	}

	result.Level, err = m.RenderLevel(item)
	if err != nil {
		return nil, err
	}

	result.Location.Path, err = m.RenderString(m.Path, item)
	if err != nil {
		return nil, err
	}
	result.Location.StartLine, err = m.RenderInt(m.StartLine, item)
	if err != nil {
		return nil, err
	}
	result.Location.StartColumn, err = m.RenderInt(m.StartColumn, item)
	if err != nil {
		return nil, err
	}
	result.Location.EndLine, err = m.RenderInt(m.EndLine, item)
	if err != nil {
		return nil, err
	}
	result.Location.EndColumn, err = m.RenderInt(m.EndColumn, item)
	if err != nil {
		return nil, err
	}

	result.Rule.ID, err = m.RenderString(m.RuleID, item)
	if err != nil {
		return nil, err
	}
	result.Rule.Name, err = m.RenderString(m.RuleName, item)
	if err != nil {
		return nil, err
	}
	result.Rule.Description, err = m.RenderString(m.RuleDescription, item)
	if err != nil {
		return nil, err
	}
	result.Rule.URI, err = m.RenderString(m.RuleURI, item)
	if err != nil {
		return nil, err
	}

	return result, nil
}

// ToResultSlice converts a slice of output data to a slice of results.
func (m OutputMapping) ToResultSlice(items []outputData) ([]*Result, error) {
	results := []*Result{}

	for _, item := range items {
		result, err := m.ToResult(item)
		if err != nil {
			return nil, err
		}
		results = append(results, result)
	}

	return results, nil
}

// RenderLevel renders the Level template using item.
// The rendered value will be normalized to one of the valid ResultLevel
// enum values.
func (m OutputMapping) RenderLevel(item outputData) (ResultLevel, error) {
	rendered, err := m.RenderString(m.Level, item)
	if err != nil {
		return ResultLevel(rendered), err
	}

	switch rendered {
	case "", strNoValue:
		return ResultLevelNone, nil
	case "info":
		return ResultLevelNote, nil
	case "warn":
		return ResultLevelWarning, nil
	case "err":
		return ResultLevelError, nil
	default:
		return ParseResultLevel(rendered)
	}
}

// RenderInt renders a template with the given output data and returns
// the rendered value as an int.
func (m OutputMapping) RenderInt(t *render.Template, item outputData) (int, error) {
	if t == nil {
		return 0, nil
	}
	rendered, err := t.Render(item)
	if err != nil {
		return 0, err
	}
	if rendered == strEmpty || rendered == strNoValue {
		return 0, nil
	}
	casted, err := strconv.Atoi(rendered)
	if err != nil {
		return 0, err
	}
	return casted, nil
}

// RenderString renders a template with the given output data and returns
// the rendered value as a string.
func (m OutputMapping) RenderString(t *render.Template, item outputData) (string, error) {
	if t == nil {
		return strEmpty, nil
	}
	rendered, err := t.Render(item)
	if err != nil {
		return strEmpty, err
	}
	if rendered == strEmpty || rendered == strNoValue {
		return strEmpty, nil
	}
	return rendered, nil
}
