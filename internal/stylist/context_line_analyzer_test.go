package stylist

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewContextLineAnalyzer(t *testing.T) {
	require.NotNil(t, NewContextLineAnalyzer())
}

func TestContextLineAnalyzer_DetectLanguage(t *testing.T) {
	tests := []struct {
		desc     string
		path     string
		lines    []string
		expected string
	}{
		{
			desc:     "should default to plaintext when all input values are empty",
			expected: "plaintext",
		},
		{
			desc:     "[path] should detect golang",
			path:     "example.go",
			expected: "go",
		},
		{
			desc:     "[path] should detect javascript",
			path:     "example.js",
			expected: "javascript",
		},
		{
			desc:     "[path] should detect python",
			path:     "example.py",
			expected: "python",
		},
		{
			desc:     "[path] should fallback to plaintext",
			path:     "README",
			expected: "plaintext",
		},
		{
			desc: "[lines] should detect golang",
			lines: []string{
				`package stylist`,
				``,
				`type Example struct {}`,
			},
			expected: "go",
		},
		{
			desc: "[lines] should fallback to plaintext",
			lines: []string{
				`no idea what this is`,
			},
			expected: "plaintext",
		},
	}
	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			analyzer := NewContextLineAnalyzer()
			actual := analyzer.DetectLanguage(tt.path, tt.lines)
			require.Equal(t, tt.expected, actual)
		})
	}
}
