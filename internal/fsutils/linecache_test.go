package fsutils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewLineCache(t *testing.T) {
	assert.IsType(t, &LineCache{}, NewLineCache(NewFileCache()))
}

func TestLineCache_GetLine(t *testing.T) {
	tests := []struct {
		desc     string
		path     string
		index    int
		expected string
		err      string
	}{
		{
			desc:     "returns the line at the given index",
			path:     "testdata/example.txt",
			index:    2,
			expected: "two",
			err:      "",
		},
		{
			desc:     "returns the first line if the index is 0",
			path:     "testdata/example.txt",
			index:    0,
			expected: "one",
			err:      "",
		},
		{
			desc:     "returns an error if index is out of bounds",
			path:     "testdata/example.txt",
			index:    99,
			expected: "",
			err:      "index out of bounds",
		},
		{
			desc:     "returns an error if path does not exist",
			path:     "testdata/does-not-exist.txt",
			index:    1,
			expected: "",
			err:      "no such file or directory",
		},
	}
	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			cache := NewLineCache(NewFileCache())
			actual, err := cache.GetLine(tt.path, tt.index)

			if tt.err == "" {
				assert.NoError(t, err)
			} else {
				assert.ErrorContains(t, err, tt.err)
			}

			assert.Equal(t, tt.expected, actual)
		})
	}
}

func TestLineCache_getRawLinesIsCached(t *testing.T) {
	cache := NewLineCache(NewFileCache())
	expected := [][]byte{
		[]byte("one"),
		[]byte("two"),
		[]byte("three"),
		[]byte(""),
	}

	lines, ok, err := cache.getRawLines("testdata/example.txt")
	assert.Equal(t, expected, lines)
	assert.Equal(t, false, ok) // cache miss
	assert.NoError(t, err)

	lines, ok, err = cache.getRawLines("testdata/example.txt")
	assert.Equal(t, expected, lines)
	assert.Equal(t, true, ok) // cache hit
	assert.NoError(t, err)
}
