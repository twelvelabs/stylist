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

func mustOpenFile(path string) *os.File {
	f, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	return f
}

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

func TestDiffOutputParser_Parse(t *testing.T) {
	tests := []struct {
		desc     string
		content  io.Reader
		expected []*Result
		err      string
	}{
		{
			desc:     "returns an empty slice when no content",
			content:  bytes.NewBufferString(""),
			expected: nil,
			err:      "",
		},
		{
			desc:     "returns an empty slice when not a valid diff",
			content:  bytes.NewBufferString("not a diff"),
			expected: []*Result{},
			err:      "",
		},
		{
			desc:    "parses diffs",
			content: mustOpenFile("testdata/output/shfmt.diff"),
			expected: []*Result{
				{
					Level: ResultLevelError,
					Location: ResultLocation{
						Path:      "bin/command.sh",
						StartLine: 4,
					},
					Rule: ResultRule{
						ID:          "diff",
						Name:        "diff",
						Description: "Formatting error",
					},
					ContextLines: []string{
						"@@ -1,10 +1,10 @@",
						" #!/usr/bin/env bash",
						" set -o errexit -o errtrace -o nounset -o pipefail",
						"",
						"-",
						"-if [",
						"+if",
						"+    [",
						"     $foo == \"bar\"",
						"-]",
						"+    ]",
						" then",
						"     echo \"lol\"",
						" fi",
					},
					ContextLang: "diff",
				},
				{
					Level: ResultLevelError,
					Location: ResultLocation{
						Path:      "bin/entrypoint.sh",
						StartLine: 19,
					},
					Rule: ResultRule{
						ID:          "diff",
						Name:        "diff",
						Description: "Formatting error",
					},
					ContextLines: []string{
						"@@ -16,8 +16,8 @@",
						"     # fix permissions",
						"     sudo chown -R app:app \\",
						"         /app \\",
						"-            /home/app \\",
						"-                /run/host-services/ssh-auth.sock",
						"+        /home/app \\",
						"+        /run/host-services/ssh-auth.sock",
						" fi",
						"",
					},
					ContextLang: "diff",
				},
			},
			err: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			actual, err := (&DiffOutputParser{}).Parse(
				CommandOutput{
					Content: tt.content,
				},
				ResultMapping{},
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
			content: mustOpenFile("testdata/output/gitleaks.txt"),
			mapping: ResultMapping{
				Pattern: "",
			},
			expected: nil,
			err:      "pattern is required",
		},
		{
			desc:    "returns an error if pattern is malformed",
			content: mustOpenFile("testdata/output/gitleaks.txt"),
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
			content:  bytes.NewBufferString("something"),
			mapping:  mapping,
			expected: nil,
			err:      "",
		},
		{
			desc:    "parses text using a regexp pattern",
			content: mustOpenFile("testdata/output/gitleaks.txt"),
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

func TestSarifOutputParser_Parse(t *testing.T) {
	tests := []struct {
		desc     string
		content  io.Reader
		expected []*Result
		err      string
	}{
		{
			desc:     "returns an empty slice when no content",
			content:  bytes.NewBufferString(""),
			expected: nil,
			err:      "",
		},
		{
			desc:     "returns an error when not sarif content",
			content:  bytes.NewBufferString("blah"),
			expected: nil,
			err:      "invalid sarif",
		},
		{
			desc:    "parses sarif output",
			content: mustOpenFile("testdata/output/hadolint.sarif"),
			expected: []*Result{
				{
					Level: ResultLevelWarning,
					Location: ResultLocation{
						Path:        "Dockerfile",
						StartLine:   44,
						StartColumn: 1,
						EndLine:     44,
						EndColumn:   1,
					},
					Rule: ResultRule{
						ID:   "DL4006",
						Name: "DL4006",
						Description: "Set the SHELL option -o pipefail before RUN with a pipe in it. " +
							"If you are using /bin/sh in an alpine image or if your shell is " +
							"symlinked to busybox then consider explicitly setting your SHELL " +
							"to /bin/ash, or disable this check",
					},
				},
				{
					Level: ResultLevelWarning,
					Location: ResultLocation{
						Path:        "Dockerfile",
						StartLine:   44,
						StartColumn: 1,
						EndLine:     44,
						EndColumn:   1,
					},
					Rule: ResultRule{
						ID:          "SC3044",
						Name:        "SC3044",
						Description: "In POSIX sh, 'popd' is undefined.",
					},
				},
			},
			err: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			actual, err := (&SarifOutputParser{}).Parse(
				CommandOutput{
					Content: tt.content,
				},
				ResultMapping{},
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
