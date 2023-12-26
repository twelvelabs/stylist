package stylist

import (
	"errors"
	"testing"

	"github.com/prashantv/gostub"
	"github.com/stretchr/testify/require"
)

func TestNewPathIgnorer(t *testing.T) {
	require := require.New(t)

	ignorer, err := NewPathIgnorer("", nil)
	require.NoError(err)
	require.Empty(ignorer.gitIgnore)
	require.Empty(ignorer.gitIgnorePath)
	require.Empty(ignorer.patterns)

	ignorer, err = NewPathIgnorer("testdata/.gitignore", []string{
		`\Documents\**\*.bat`,
	})
	require.NoError(err)
	require.NotNil(ignorer.gitIgnore)
	require.Equal("testdata/.gitignore", ignorer.gitIgnorePath)
	require.Equal([]string{
		"/Documents/**/*.bat",
	}, ignorer.patterns)
}

func TestNewPathIgnorer_WhenGitIgnoreParseError(t *testing.T) {
	stubs := gostub.StubFunc(&newGitIgnoreParser, nil, errors.New("boom"))
	defer stubs.Reset()

	require := require.New(t)

	ignorer, err := NewPathIgnorer("testdata/.gitignore", nil)
	require.ErrorContains(err, "boom")
	require.Nil(ignorer)
}

func TestNewPathIgnorer_WhenPatternValidateError(t *testing.T) {
	require := require.New(t)

	ignorer, err := NewPathIgnorer("", []string{
		`*{{*}`,
	})
	require.ErrorContains(err, "error in pattern")
	require.Nil(ignorer)
}

func TestPathIgnorer_ShouldIgnore(t *testing.T) {
	tests := []struct {
		path     string
		dir      bool
		expected bool
	}{
		{
			path:     "example.json",
			expected: false,
		},
		{
			path:     "example.yaml",
			expected: true,
		},
		{
			path:     "ignored_file.txt",
			expected: true,
		},
		{
			path:     "ignored_dir",
			dir:      true,
			expected: true,
		},
		{
			path:     "ignored_dir",
			dir:      false,
			expected: false,
		},
		{
			path:     "ignored_dir/allowed.txt",
			expected: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.path, func(t *testing.T) {
			require := require.New(t)

			ignorer, err := NewPathIgnorer("testdata/.gitignore", []string{
				"**/*.{yaml,yml}",
			})
			require.NoError(err)

			actual := ignorer.ShouldIgnore(tt.path, tt.dir)
			require.Equal(tt.expected, actual)
		})
	}
}
