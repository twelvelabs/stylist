package stylist

import (
	"bytes"
	"errors"
	"io"
	"os"
	"strings"
	"testing"
	"testing/iotest"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/twelvelabs/termite/render"
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

func TestCheckstyleOutputParser_Parse(t *testing.T) {
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
			desc:     "returns an error when unable to read content",
			content:  iotest.ErrReader(errors.New("boom")),
			expected: nil,
			err:      "boom",
		},
		{
			desc:     "returns an error when not checkstyle content",
			content:  bytes.NewBufferString("not a checkstyle document"),
			expected: nil,
			err:      "invalid checkstyle XML",
		},
		{
			desc:    "parses checkstyle",
			content: mustOpenFile("testdata/output/golangci.xml"),
			expected: []*Result{
				{
					Level: ResultLevelError,
					Location: ResultLocation{
						Path:        "internal/stylist/output_parser.go",
						StartLine:   33,
						StartColumn: 76,
					},
					Rule: ResultRule{
						ID:          "godot",
						Name:        "godot",
						Description: "Comment should end in a period",
					},
				},
				{
					Level: ResultLevelError,
					Location: ResultLocation{
						Path:        "internal/stylist/output_parser.go",
						StartLine:   57,
						StartColumn: 48,
					},
					Rule: ResultRule{
						ID:          "godot",
						Name:        "godot",
						Description: "Comment should end in a period",
					},
				},
				{
					Level: ResultLevelError,
					Location: ResultLocation{
						Path:        "internal/stylist/output_parser_test.go",
						StartLine:   35,
						StartColumn: 7,
					},
					Rule: ResultRule{
						ID:          "godot",
						Name:        "godot",
						Description: "Comment should end in a period",
					},
				},
			},
			err: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			actual, err := (&CheckstyleOutputParser{}).Parse(
				CommandOutput{
					Content: tt.content,
				},
				ResultMapping{},
			)

			if tt.err == "" {
				require.NoError(t, err)
			} else {
				require.ErrorContains(t, err, tt.err)
			}

			require.Equal(t, tt.expected, actual)
		})
	}
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
			desc:     "returns an empty slice when not a diff",
			content:  bytes.NewBufferString("not a diff"),
			expected: []*Result{},
			err:      "",
		},
		{
			desc:     "returns an error when unable to read content",
			content:  iotest.ErrReader(errors.New("boom")),
			expected: nil,
			err:      "boom",
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
			Content: bytes.NewBufferString(""),
		},
		ResultMapping{},
	)
	require.NoError(t, err)
	require.Nil(t, results)

	results, err = parser.Parse(
		CommandOutput{
			Content: bytes.NewBufferString("{{}"),
		},
		ResultMapping{},
	)
	require.ErrorContains(t, err, "invalid json")
	require.Nil(t, results)

	results, err = parser.Parse(
		CommandOutput{
			Content: iotest.ErrReader(errors.New("boom")),
		},
		ResultMapping{},
	)
	require.ErrorContains(t, err, "boom")
	require.Nil(t, results)

	results, err = parser.Parse(
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
		Level: ResultLevelInfo,
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
			desc:     "returns an error when unable to read content",
			content:  iotest.ErrReader(errors.New("boom")),
			mapping:  mapping,
			expected: nil,
			err:      "boom",
		},
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
			desc:     "returns an error when unable to read content",
			content:  iotest.ErrReader(errors.New("boom")),
			expected: nil,
			err:      "boom",
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
