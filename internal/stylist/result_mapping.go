package stylist

import (
	"strconv"
	"strings"

	"github.com/twelvelabs/termite/render"
)

const (
	strEmpty   string = ""
	strNoValue string = "<no value>"
)

// resultData is a map of key/value pairs parsed from CommandOutput
// and passed to a ResultMapping to be converted into a Result.
type resultData map[string]any

// ResultMapping is a set of rules for how to map command output to a result.
//
// Mappings are typically defined in stylist.yml when the output type
// has been set to "json" or "regexp".
type ResultMapping struct {
	Pattern         string           `yaml:"pattern,omitempty"`
	Level           *render.Template `yaml:"level,omitempty"`
	Path            *render.Template `yaml:"path,omitempty"`
	StartLine       *render.Template `yaml:"start_line,omitempty"`
	StartColumn     *render.Template `yaml:"start_column,omitempty"`
	EndLine         *render.Template `yaml:"end_line,omitempty"`
	EndColumn       *render.Template `yaml:"end_column,omitempty"`
	RuleID          *render.Template `yaml:"rule_id,omitempty"`
	RuleName        *render.Template `yaml:"rule_name,omitempty"`
	RuleDescription *render.Template `yaml:"rule_description,omitempty"`
	RuleURI         *render.Template `yaml:"rule_uri,omitempty"`
	Context         *render.Template `yaml:"context,omitempty"`
}

// ToResult converts a map of output data to a Result struct.
func (m ResultMapping) ToResult(item resultData) (*Result, error) {
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
	result.ContextLines, err = m.RenderStringSlice(m.Context, item)
	if err != nil {
		return nil, err
	}

	return result, nil
}

// ToResultSlice converts a slice of output data to a slice of results.
func (m ResultMapping) ToResultSlice(items []resultData) ([]*Result, error) {
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
func (m ResultMapping) RenderLevel(item resultData) (ResultLevel, error) {
	rendered, err := m.RenderString(m.Level, item)
	if err != nil {
		return ResultLevelNone, err
	}

	return CoerceResultLevel(rendered)
}

// RenderInt renders a template with the given output data and returns
// the rendered value as an int.
func (m ResultMapping) RenderInt(t *render.Template, item resultData) (int, error) {
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
func (m ResultMapping) RenderString(t *render.Template, item resultData) (string, error) {
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

// RenderStringSlice renders a template with the given output data and returns
// the rendered value as a slice of strings.
func (m ResultMapping) RenderStringSlice(t *render.Template, item resultData) ([]string, error) {
	if t == nil {
		return nil, nil
	}
	rendered, err := t.Render(item)
	if err != nil {
		return nil, err
	}
	if rendered == strEmpty || rendered == strNoValue {
		return nil, nil
	}
	return strings.Split(strings.TrimSuffix(rendered, "\n"), "\n"), nil
}
