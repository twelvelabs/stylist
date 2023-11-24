package stylist

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/twelvelabs/termite/ui"
)

func TestNewResultPrinter(t *testing.T) {
	ios := ui.NewTestIOStreams()
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

func TestCheckstylePrinter_Print(t *testing.T) {
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
			},
		},
		{
			Source: "test-linter",
			Level:  ResultLevelWarning,
			Location: ResultLocation{
				Path:        "some/path/foo.go",
				StartLine:   2,
				StartColumn: 3,
			},
			Rule: ResultRule{
				ID:          "rule-id2",
				Name:        "rule-name2",
				Description: "valid start column",
			},
		},
		{
			Source: "test-linter",
			Level:  ResultLevelWarning,
			Location: ResultLocation{
				Path:        "some/path/bar.go",
				StartLine:   4,
				StartColumn: 5,
			},
			Rule: ResultRule{
				ID:          "rule-id2",
				Name:        "rule-name2",
				Description: "another valid start column",
			},
		},
	}

	tests := []struct {
		desc     string
		results  []*Result
		expected string
		err      string
	}{
		{
			desc:     "empty result set should print nothing",
			results:  []*Result{},
			expected: `<?xml version="1.0" encoding="UTF-8"?><checkstyle version="4.3"></checkstyle>`,
		},

		{
			desc:     "should print a minimal result set when config disabled",
			results:  results,
			expected: `<?xml version="1.0" encoding="UTF-8"?><checkstyle version="4.3"><file name="some/path/foo.go"><error line="1" column="0" message="no start column [rule-id1]" severity="error" source="test-linter"></error><error line="2" column="3" message="valid start column [rule-id2]" severity="warning" source="test-linter"></error></file><file name="some/path/bar.go"><error line="4" column="5" message="another valid start column [rule-id2]" severity="warning" source="test-linter"></error></file></checkstyle>`, //nolint: lll
		},
	}

	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			app := NewTestApp()
			printer := &CheckstylePrinter{
				ios:    app.IO,
				config: app.Config,
			}
			err := printer.Print(tt.results)

			if tt.err == "" {
				require.NoError(t, err)
			} else {
				require.ErrorContains(t, err, tt.err)
			}

			out := app.IO.Out.String()
			out = strings.ReplaceAll(out, "\n", "")
			assert.Equal(t, tt.expected, out)
		})
	}
}

func TestSarifPrinter_Print(t *testing.T) {
	app := NewTestApp()
	printer := &SarifPrinter{
		ios:    app.IO,
		config: app.Config,
	}
	results := []*Result{}

	err := printer.Print(results)
	assert.NoError(t, err)
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
			Level:  ResultLevelInfo,
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
				"some/path/foo.go:1:0: error: test-linter: no start column. [rule-id1]",
				"some/path/bar.go:2:10: warning: test-linter: valid start and end column. [rule-id2]",
				"some/path/baz.go:1:1: info: test-linter: single char indicator. [rule-id3]",
				"some/path/qux.go:1:99: none: test-linter: out of bounds indicator. [rule-id4]",
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
				"some/path/foo.go:1:0: error: test-linter: no start column. [rule-id1]",
				"context line one",
				"some/path/bar.go:2:10: warning: test-linter: valid start and end column. [rule-id2]",
				"\tcontext line two",
				"\t        ^^^^",
				"some/path/baz.go:1:1: info: test-linter: single char indicator. [rule-id3]",
				"context line three",
				"^",
				"some/path/qux.go:1:99: none: test-linter: out of bounds indicator. [rule-id4]",
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
				"some/path/foo.go:1:0: error: test-linter: no start column. " +
					"[rule-id1](https://test-linter.com/rule-id1)",
				"some/path/bar.go:2:10: warning: test-linter: valid start and end column. " +
					"[rule-id2](https://test-linter.com/rule-id2)",
				"some/path/baz.go:1:1: info: test-linter: single char indicator. " +
					"[rule-id3](https://test-linter.com/rule-id3)",
				"some/path/qux.go:1:99: none: test-linter: out of bounds indicator. " +
					"[rule-id4](https://test-linter.com/rule-id4)",
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

func TestTtyPrinter_Print_ColorEnabled(t *testing.T) {
	app := NewTestApp()

	// Enable color so we can exercise the syntax highlighting logic.
	app.IO.SetColorEnabled(true)
	app.Config.Output.ShowURL = false

	printer := &TtyPrinter{
		ios:    app.IO,
		config: app.Config,
	}
	err := printer.Print([]*Result{
		{
			Source: "test-linter",
			Level:  ResultLevelError,
			Location: ResultLocation{
				Path:        "some/path/foo.go",
				StartLine:   1,
				StartColumn: 1,
			},
			Rule: ResultRule{
				ID:          "rule-id1",
				Name:        "rule-name1",
				Description: "some issue",
				URI:         "https://test-linter.com/rule-id1",
			},
			ContextLines: []string{
				"package foo",
				"",
				"import \"os\"",
			},
		},
	})

	// Assert that the context lines have been highlighted.
	assert.NoError(t, err)
	assert.Equal(t, []string{
		// cspell: disable
		"\x1b[1m\x1b[38;5;129mpackage\x1b[0m foo",
		"",
		"\x1b[1m\x1b[38;5;129mimport\x1b[0m \x1b[38;5;131m\"os\"\x1b[0m",
		// cspell: enable
	}, app.IO.Out.Lines()[1:])
}
