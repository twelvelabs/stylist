package stylist

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/twelvelabs/stylist/internal/render"
)

func newOutputDataFixture() outputData {
	return outputData{
		"level":            "warning",
		"path":             "<path>",
		"start_line":       1,
		"start_column":     11,
		"end_line":         2,
		"end_column":       22,
		"rule_id":          "<id>",
		"rule_name":        "<name>",
		"rule_description": "<description>",
		"rule_uri":         "<uri>",
	}
}
func newOutputMappingFixture() OutputMapping {
	return OutputMapping{
		Level:           render.MustCompile(`{{ .level }}`),
		Path:            render.MustCompile(`{{ .path }}`),
		StartLine:       render.MustCompile(`{{ .start_line }}`),
		StartColumn:     render.MustCompile(`{{ .start_column }}`),
		EndLine:         render.MustCompile(`{{ .end_line }}`),
		EndColumn:       render.MustCompile(`{{ .end_column }}`),
		RuleID:          render.MustCompile(`{{ .rule_id }}`),
		RuleName:        render.MustCompile(`{{ .rule_name }}`),
		RuleDescription: render.MustCompile(`{{ .rule_description }}`),
		RuleURI:         render.MustCompile(`{{ .rule_uri }}`),
	}
}

func newResultFixture() *Result {
	return &Result{
		Level: ResultLevelWarning,
		Location: ResultLocation{
			Path:        "<path>",
			StartLine:   1,
			StartColumn: 11,
			EndLine:     2,
			EndColumn:   22,
		},
		Rule: ResultRule{
			ID:          "<id>",
			Name:        "<name>",
			Description: "<description>",
			URI:         "<uri>",
		},
	}
}

func TestOutputMapping_ToResult(t *testing.T) {
	tests := []struct {
		desc     string
		data     outputData
		mapping  OutputMapping
		expected *Result
		setup    func(d *outputData, m *OutputMapping, r *Result)
		err      string
	}{
		{
			desc:     "should convert data into a result struct",
			data:     newOutputDataFixture(),
			mapping:  newOutputMappingFixture(),
			expected: newResultFixture(),
		},
		{
			desc:    "should return an empty result when data is empty",
			data:    nil,
			mapping: newOutputMappingFixture(),
			expected: &Result{
				Level:    ResultLevelNone,
				Location: ResultLocation{},
				Rule:     ResultRule{},
			},
		},
		{
			desc:    "should return an empty result when mapping is empty",
			data:    newOutputDataFixture(),
			mapping: OutputMapping{},
			expected: &Result{
				Level:    ResultLevelNone,
				Location: ResultLocation{},
				Rule:     ResultRule{},
			},
		},
		{
			desc:     "should handle error when rendering Level",
			data:     newOutputDataFixture(),
			mapping:  newOutputMappingFixture(),
			expected: nil,
			setup: func(d *outputData, m *OutputMapping, r *Result) {
				m.Level = render.MustCompile(`{{ fail "boom" }}`)
			},
			err: "fail: boom",
		},
		{
			desc:     "should handle error when rendering Path",
			data:     newOutputDataFixture(),
			mapping:  newOutputMappingFixture(),
			expected: nil,
			setup: func(d *outputData, m *OutputMapping, r *Result) {
				m.Path = render.MustCompile(`{{ fail "boom" }}`)
			},
			err: "fail: boom",
		},
		{
			desc:     "should handle error when rendering StartLine",
			data:     newOutputDataFixture(),
			mapping:  newOutputMappingFixture(),
			expected: nil,
			setup: func(d *outputData, m *OutputMapping, r *Result) {
				m.StartLine = render.MustCompile(`{{ fail "boom" }}`)
			},
			err: "fail: boom",
		},
		{
			desc:     "should handle error when rendering StartColumn",
			data:     newOutputDataFixture(),
			mapping:  newOutputMappingFixture(),
			expected: nil,
			setup: func(d *outputData, m *OutputMapping, r *Result) {
				m.StartColumn = render.MustCompile(`{{ fail "boom" }}`)
			},
			err: "fail: boom",
		},
		{
			desc:     "should handle error when rendering EndLine",
			data:     newOutputDataFixture(),
			mapping:  newOutputMappingFixture(),
			expected: nil,
			setup: func(d *outputData, m *OutputMapping, r *Result) {
				m.EndLine = render.MustCompile(`{{ fail "boom" }}`)
			},
			err: "fail: boom",
		},
		{
			desc:     "should handle error when rendering EndColumn",
			data:     newOutputDataFixture(),
			mapping:  newOutputMappingFixture(),
			expected: nil,
			setup: func(d *outputData, m *OutputMapping, r *Result) {
				m.EndColumn = render.MustCompile(`{{ fail "boom" }}`)
			},
			err: "fail: boom",
		},

		{
			desc:     "should handle error when rendering RuleID",
			data:     newOutputDataFixture(),
			mapping:  newOutputMappingFixture(),
			expected: nil,
			setup: func(d *outputData, m *OutputMapping, r *Result) {
				m.RuleID = render.MustCompile(`{{ fail "boom" }}`)
			},
			err: "fail: boom",
		},
		{
			desc:     "should handle error when rendering RuleName",
			data:     newOutputDataFixture(),
			mapping:  newOutputMappingFixture(),
			expected: nil,
			setup: func(d *outputData, m *OutputMapping, r *Result) {
				m.RuleName = render.MustCompile(`{{ fail "boom" }}`)
			},
			err: "fail: boom",
		},
		{
			desc:     "should handle error when rendering RuleDescription",
			data:     newOutputDataFixture(),
			mapping:  newOutputMappingFixture(),
			expected: nil,
			setup: func(d *outputData, m *OutputMapping, r *Result) {
				m.RuleDescription = render.MustCompile(`{{ fail "boom" }}`)
			},
			err: "fail: boom",
		},
		{
			desc:     "should handle error when rendering RuleURI",
			data:     newOutputDataFixture(),
			mapping:  newOutputMappingFixture(),
			expected: nil,
			setup: func(d *outputData, m *OutputMapping, r *Result) {
				m.RuleURI = render.MustCompile(`{{ fail "boom" }}`)
			},
			err: "fail: boom",
		},
	}

	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			if tt.setup != nil {
				tt.setup(&tt.data, &tt.mapping, tt.expected)
			}

			actual, err := tt.mapping.ToResult(tt.data)

			if tt.err == "" {
				assert.NoError(t, err)
			} else {
				assert.ErrorContains(t, err, tt.err)
			}

			assert.Equal(t, tt.expected, actual)
		})
	}
}

func TestOutputMapping_ToResultSlice(t *testing.T) {
	// happy path
	items := []outputData{
		newOutputDataFixture(),
		newOutputDataFixture(),
	}
	mapping := newOutputMappingFixture()
	results, err := mapping.ToResultSlice(items)
	assert.NoError(t, err)
	assert.Equal(t, 2, len(results))

	// error path
	mapping = OutputMapping{
		Level: render.MustCompile(`{{ fail "boom" }}`),
	}
	results, err = mapping.ToResultSlice(items)
	assert.ErrorContains(t, err, "fail: boom")
	assert.Equal(t, 0, len(results))
}

func TestOutputMapping_RenderLevel(t *testing.T) {
	tests := []struct {
		desc     string
		template *render.Template
		data     outputData
		expected ResultLevel
		err      string
	}{
		{
			desc:     "enum values should be accepted",
			template: render.MustCompile(`{{ .level }}`),
			data: outputData{
				"level": "error",
			},
			expected: ResultLevelError,
		},

		{
			desc:     "missing template should be normalized to none",
			template: nil,
			data:     outputData{},
			expected: ResultLevelNone,
		},
		{
			desc:     "missing key should be normalized to none",
			template: render.MustCompile(`{{ .level }}`),
			data:     outputData{},
			expected: ResultLevelNone,
		},
		{
			desc:     "empty string should be normalized to none",
			template: render.MustCompile(`{{ .level }}`),
			data: outputData{
				"level": "",
			},
			expected: ResultLevelNone,
		},
		{
			desc:     "info should be normalized to note",
			template: render.MustCompile(`{{ .level }}`),
			data: outputData{
				"level": "info",
			},
			expected: ResultLevelNote,
		},
		{
			desc:     "warn should be normalized to warning",
			template: render.MustCompile(`{{ .level }}`),
			data: outputData{
				"level": "warn",
			},
			expected: ResultLevelWarning,
		},
		{
			desc:     "err should be normalized to error",
			template: render.MustCompile(`{{ .level }}`),
			data: outputData{
				"level": "err",
			},
			expected: ResultLevelError,
		},

		{
			desc:     "should return an error when not an enum or normalized value",
			template: render.MustCompile(`{{ .level }}`),
			data: outputData{
				"level": "unknown",
			},
			expected: ResultLevel(""),
			err:      "unknown is not a valid ResultLevel",
		},
		{
			desc:     "should return an error when template fails to render",
			template: render.MustCompile(`{{ fail "boom" }}`),
			data: outputData{
				"level": "unknown",
			},
			expected: ResultLevel(""),
			err:      "fail: boom",
		},
	}
	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			mapping := OutputMapping{
				Level: tt.template,
			}
			actual, err := mapping.RenderLevel(tt.data)

			if tt.err == "" {
				assert.NoError(t, err)
			} else {
				assert.ErrorContains(t, err, tt.err)
			}

			assert.Equal(t, tt.expected, actual)
		})
	}
}

func TestOutputMapping_RenderString(t *testing.T) {
	tests := []struct {
		desc     string
		template *render.Template
		data     outputData
		expected string
		err      string
	}{
		{
			desc:     "normal values should be rendered",
			template: render.MustCompile(`{{ .something }}`),
			data: outputData{
				"something": "foo",
			},
			expected: "foo",
		},
		{
			desc:     "missing template should return an empty value",
			template: nil,
			data:     nil,
			expected: "",
		},
		{
			desc:     "missing key should return an empty value",
			template: render.MustCompile(`{{ .something }}`),
			data:     nil,
			expected: "",
		},
		{
			desc:     "should return an error when template fails to render",
			template: render.MustCompile(`{{ fail "boom" }}`),
			data:     nil,
			expected: "",
			err:      "fail: boom",
		},
	}

	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			mapping := OutputMapping{}
			actual, err := mapping.RenderString(tt.template, tt.data)

			if tt.err == "" {
				assert.NoError(t, err)
			} else {
				assert.ErrorContains(t, err, tt.err)
			}

			assert.Equal(t, tt.expected, actual)
		})
	}
}

func TestOutputMapping_RenderInt(t *testing.T) {
	tests := []struct {
		desc     string
		template *render.Template
		data     outputData
		expected int
		err      string
	}{
		{
			desc:     "normal values should be rendered",
			template: render.MustCompile(`{{ .something }}`),
			data: outputData{
				"something": 12,
			},
			expected: 12,
		},
		{
			desc:     "missing template should return an empty value",
			template: nil,
			data:     nil,
			expected: 0,
		},
		{
			desc:     "missing key should return an empty value",
			template: render.MustCompile(`{{ .something }}`),
			data:     nil,
			expected: 0,
		},
		{
			desc:     "should return an error when template fails to render",
			template: render.MustCompile(`{{ fail "boom" }}`),
			data:     nil,
			expected: 0,
			err:      "fail: boom",
		},
		{
			desc:     "should return an error when value fails to cast",
			template: render.MustCompile(`{{ .something }}`),
			data: outputData{
				"something": "not-an-int",
			},
			expected: 0,
			err:      "invalid syntax",
		},
	}

	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			mapping := OutputMapping{}
			actual, err := mapping.RenderInt(tt.template, tt.data)

			if tt.err == "" {
				assert.NoError(t, err)
			} else {
				assert.ErrorContains(t, err, tt.err)
			}

			assert.Equal(t, tt.expected, actual)
		})
	}
}