package stylist

import (
	"os"
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
	file, err := os.Open("testdata/json/shellcheck.json")
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
