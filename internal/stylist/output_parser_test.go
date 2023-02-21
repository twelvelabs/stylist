package stylist

import (
	"bytes"
	"io"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/twelvelabs/stylist/internal/render"
)

func TestNewOutputParser(t *testing.T) {
	// Ensure a parser exists for each enum value.
	for _, name := range OutputFormatNames() {
		assert.NotPanics(t, func() {
			_ = NewOutputParser(OutputFormat(name))
		})
	}
	assert.PanicsWithValue(t, "unknown output format: unknown", func() {
		_ = NewOutputParser(OutputFormat("unknown"))
	})
}

func TestJSONOutputParser_Parse(t *testing.T) {
	file, err := os.Open("testdata/output/shellcheck.json")
	assert.NoError(t, err)

	parser := &JSONOutputParser{}
	results, err := parser.Parse(
		CommandOutput{
			Content: file,
		},
		ResultMapping{
			Level:           render.MustCompile(`{{ .level }}`),
			Path:            render.MustCompile(`{{ .file }}`),
			StartLine:       render.MustCompile(`{{ .line }}`),
			StartColumn:     render.MustCompile(`{{ .column }}`),
			EndLine:         render.MustCompile(`{{ .endLine }}`),
			EndColumn:       render.MustCompile(`{{ .endColumn }}`),
			RuleID:          render.MustCompile(`SC{{ .code }}`),
			RuleName:        render.MustCompile(`SC{{ .code }}`),
			RuleDescription: render.MustCompile(`{{ .message }}`),
			RuleURI:         render.MustCompile(`https://www.shellcheck.net/wiki/SC{{ .code }}`),
		},
	)

	assert.NoError(t, err)
	require.Equal(t, 1, len(results))
	assert.Equal(t, &Result{
		Level: ResultLevelNote,
		Location: ResultLocation{
			Path:        "entrypoint.sh",
			StartLine:   15,
			StartColumn: 6,
			EndLine:     15,
			EndColumn:   19,
		},
		Rule: ResultRule{
			ID:          "SC2086",
			Name:        "SC2086",
			Description: "Double quote to prevent globbing and word splitting.",
			URI:         "https://www.shellcheck.net/wiki/SC2086",
		},
	}, results[0])
}

func TestRegexpOutputParser_Parse(t *testing.T) {
	fixtureFile := func(path string) io.Reader {
		f, err := os.Open(path)
		if err != nil {
			panic(err)
		}
		return f
	}

	mapping := ResultMapping{
		Pattern: strings.Join([]string{
			`(?m)`,
			`Secret:\s+(?P<secret>.*)`,
			`RuleID:\s+(?P<rule_id>.*)`,
			`Entropy:\s+(?P<entropy>.*)`,
			`File:\s+(?P<file>.*)`,
			`Line:\s+(?P<line>.*)`,
		}, `\n`),
		Level:           render.MustCompile(`error`),
		Path:            render.MustCompile(`{{ .file }}`),
		StartLine:       render.MustCompile(`{{ .line }}`),
		RuleID:          render.MustCompile(`{{ .rule_id }}`),
		RuleName:        render.MustCompile(`{{ .rule_id }}`),
		RuleDescription: render.MustCompile(`Secret detected: {{ .secret }}`),
	}

	tests := []struct {
		desc     string
		content  io.Reader
		mapping  ResultMapping
		expected []*Result
		err      string
	}{
		{
			desc:    "returns an error if pattern is missing",
			content: fixtureFile("testdata/output/gitleaks.txt"),
			mapping: ResultMapping{
				Pattern: "",
			},
			expected: nil,
			err:      "pattern is required",
		},
		{
			desc:    "returns an error if pattern is malformed",
			content: fixtureFile("testdata/output/gitleaks.txt"),
			mapping: ResultMapping{
				Pattern: "(?lol)",
			},
			expected: nil,
			err:      "error parsing regexp",
		},
		{
			desc:     "returns an empty slice when no content",
			content:  bytes.NewBufferString(""),
			mapping:  mapping,
			expected: nil,
			err:      "",
		},
		{
			desc:     "returns an empty slice when nothing matches",
			content:  bytes.NewBufferString("welp"),
			mapping:  mapping,
			expected: nil,
			err:      "",
		},
		{
			desc:    "parses text using a regexp pattern",
			content: fixtureFile("testdata/output/gitleaks.txt"),
			mapping: mapping,
			expected: []*Result{
				{
					Level: ResultLevelError,
					Location: ResultLocation{
						Path:      ".env",
						StartLine: 8,
					},
					Rule: ResultRule{
						ID:          "github-oauth",
						Name:        "github-oauth",
						Description: "Secret detected: REDACTED",
					},
				},
				{
					Level: ResultLevelError,
					Location: ResultLocation{
						Path:      "bin/run-in-container.sh",
						StartLine: 22,
					},
					Rule: ResultRule{
						ID:          "github-oauth",
						Name:        "github-oauth",
						Description: "Secret detected: REDACTED",
					},
				},
				{
					Level: ResultLevelError,
					Location: ResultLocation{
						Path:      "etc/githooks/commit-msg",
						StartLine: 10,
					},
					Rule: ResultRule{
						ID:          "github-oauth",
						Name:        "github-oauth",
						Description: "Secret detected: REDACTED",
					},
				},
			},
			err: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			actual, err := (&RegexpOutputParser{}).Parse(
				CommandOutput{
					Content: tt.content,
				},
				tt.mapping,
			)

			if tt.err == "" {
				assert.NoError(t, err)
			} else {
				assert.ErrorContains(t, err, tt.err)
			}

			assert.Equal(t, tt.expected, actual)
		})
	}
}
