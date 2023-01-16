package stylist

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestContextLineLoader(t *testing.T) {
	tests := []struct {
		desc     string
		location ResultLocation
		expected []string
		err      string
	}{
		{
			desc: "returns nil when path is empty",
			location: ResultLocation{
				Path: "",
			},
			expected: nil,
		},
		{
			desc: "returns nil when start line is empty",
			location: ResultLocation{
				Path:      "testdata/example.txt",
				StartLine: 0,
			},
			expected: nil,
		},
		{
			desc: "returns a single line when end line is empty",
			location: ResultLocation{
				Path:      "testdata/example.txt",
				StartLine: 1,
				EndLine:   0,
			},
			expected: []string{
				"one",
			},
		},
		{
			desc: "returns the lines for the given range",
			location: ResultLocation{
				Path:      "testdata/example.txt",
				StartLine: 1,
				EndLine:   2,
			},
			expected: []string{
				"one",
				"two",
			},
		},
		{
			desc: "returns an error if path does not exist",
			location: ResultLocation{
				Path:      "testdata/does-not-exist.txt",
				StartLine: 1,
			},
			expected: nil,
			err:      "no such file or directory",
		},
	}
	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			actual, err := NewContextLineLoader().Load(tt.location)

			if tt.err == "" {
				assert.NoError(t, err)
			} else {
				assert.ErrorContains(t, err, tt.err)
			}

			assert.Equal(t, tt.expected, actual)
		})
	}
}
