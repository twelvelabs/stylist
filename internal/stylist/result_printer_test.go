package stylist

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/twelvelabs/termite/ioutil"
)

func TestNewResultPrinter(t *testing.T) {
	ios := ioutil.Test()
	config := NewConfig()
	// Ensure a printer exists for each enum value.
	for _, name := range ResultFormatNames() {
		config.Output.Format = ResultFormat(name)
		assert.NotPanics(t, func() {
			_ = NewResultPrinter(ios, config)
		})
	}
	assert.PanicsWithValue(t, "unknown result format: unknown", func() {
		config.Output.Format = ResultFormat("unknown")
		_ = NewResultPrinter(ios, config)
	})
}

func TestTtyPrinter_Print(t *testing.T) {
	results := []*Result{
		{
			Source: "test-linter",
			Level:  ResultLevelError,
			Location: ResultLocation{
				Path:        "some/path/foo.go",
				StartLine:   1,
				StartColumn: 0,
			},
			Rule: ResultRule{
				ID:          "rule-id1",
				Name:        "rule-name1",
				Description: "no start column",
				URI:         "https://test-linter.com/rule-id1",
			},
			ContextLines: []string{
				"context line one",
			},
		},
		{
			Source: "test-linter",
			Level:  ResultLevelWarning,
			Location: ResultLocation{
				Path:        "some/path/bar.go",
				StartLine:   2,
				StartColumn: 10,
				EndLine:     2,
				EndColumn:   14,
			},
			Rule: ResultRule{
				ID:          "rule-id2",
				Name:        "rule-name2",
				Description: "valid start and end column",
				URI:         "https://test-linter.com/rule-id2",
			},
			ContextLines: []string{
				"\tcontext line two",
			},
		},
		{
			Source: "test-linter",
			Level:  ResultLevelNote,
			Location: ResultLocation{
				Path:        "some/path/baz.go",
				StartLine:   1,
				StartColumn: 1,
			},
			Rule: ResultRule{
				ID:          "rule-id3",
				Name:        "rule-name3",
				Description: "single char indicator",
				URI:         "https://test-linter.com/rule-id3",
			},
			ContextLines: []string{
				"context line three",
			},
		},
		{
			Source: "test-linter",
			Level:  ResultLevelNone,
			Location: ResultLocation{
				Path:        "some/path/qux.go",
				StartLine:   1,
				StartColumn: 99,
				EndLine:     1,
				EndColumn:   999,
			},
			Rule: ResultRule{
				ID:          "rule-id4",
				Name:        "rule-name4",
				Description: "out of bounds indicator",
				URI:         "https://test-linter.com/rule-id4",
			},
			ContextLines: []string{
				"context line four",
			},
		},
	}

	tests := []struct {
		desc     string
		config   OutputConfig
		results  []*Result
		expected []string
		err      string
	}{
		{
			desc:     "empty result set should print nothing",
			results:  []*Result{},
			expected: []string{},
			err:      "",
		},

		{
			desc: "should print a minimal result set when config disabled",
			config: OutputConfig{
				ShowContext: false,
				ShowURL:     false,
			},
			results: results,
			expected: []string{
				"[test-linter] some/path/foo.go:1:0 no start column (rule-id1)",
				"[test-linter] some/path/bar.go:2:10 valid start and end column (rule-id2)",
				"[test-linter] some/path/baz.go:1:1 single char indicator (rule-id3)",
				"[test-linter] some/path/qux.go:1:99 out of bounds indicator (rule-id4)",
			},
			err: "",
		},
		{
			desc: "should print context when enabled",
			config: OutputConfig{
				ShowContext: true,
				ShowURL:     false,
			},
			results: results,
			expected: []string{
				"[test-linter] some/path/foo.go:1:0 no start column (rule-id1)",
				"context line one",
				"[test-linter] some/path/bar.go:2:10 valid start and end column (rule-id2)",
				"\tcontext line two",
				"\t        ^^^^",
				"[test-linter] some/path/baz.go:1:1 single char indicator (rule-id3)",
				"context line three",
				"^",
				"[test-linter] some/path/qux.go:1:99 out of bounds indicator (rule-id4)",
				"context line four",
				"                 ^",
			},
			err: "",
		},
		{
			desc: "should print urls when enabled",
			config: OutputConfig{
				ShowContext: false,
				ShowURL:     true,
			},
			results: results,
			expected: []string{
				"[test-linter] some/path/foo.go:1:0 no start column (rule-id1) " +
					"<https://test-linter.com/rule-id1>",
				"[test-linter] some/path/bar.go:2:10 valid start and end column (rule-id2) " +
					"<https://test-linter.com/rule-id2>",
				"[test-linter] some/path/baz.go:1:1 single char indicator (rule-id3) " +
					"<https://test-linter.com/rule-id3>",
				"[test-linter] some/path/qux.go:1:99 out of bounds indicator (rule-id4) " +
					"<https://test-linter.com/rule-id4>",
			},
			err: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			app := NewTestApp()
			app.Config.Output = tt.config

			printer := &TtyPrinter{
				ios:    app.IO,
				config: app.Config,
			}
			err := printer.Print(tt.results)

			if tt.err == "" {
				assert.NoError(t, err)
			} else {
				assert.ErrorContains(t, err, tt.err)
			}

			assert.Equal(t, tt.expected, app.IO.Out.Lines())
		})
	}
}
