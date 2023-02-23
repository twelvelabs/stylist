package fsutils

import (
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/twelvelabs/termite/testutil"
)

func TestNoPathExists(t *testing.T) {
	testutil.InTempDir(t, func(dir string) {
		assert.NoFileExists(t, "foo.txt")
		assert.Equal(t, true, NoPathExists("foo.txt"))

		testutil.WriteFile(t, "foo.txt", []byte(""), 0600)
		assert.Equal(t, false, NoPathExists("foo.txt"))
	})
}

func TestPathExists(t *testing.T) {
	testutil.InTempDir(t, func(dir string) {
		assert.NoFileExists(t, "foo.txt")
		assert.Equal(t, false, PathExists("foo.txt"))

		testutil.WriteFile(t, "foo.txt", []byte(""), 0600)
		assert.Equal(t, true, PathExists("foo.txt"))
	})
}

func TestRelativePath(t *testing.T) {
	tests := []struct {
		desc     string
		path     string
		expected string
		err      string
	}{
		{
			desc:     mustAbsPath("testdata/example.txt"),
			expected: "testdata/example.txt",
		},
		{
			desc:     "file://" + mustAbsPath("testdata/example.txt"),
			expected: "testdata/example.txt",
		},
		{
			desc:     "./../fsutils/testdata/example.txt",
			expected: "testdata/example.txt",
		},
		{
			desc:     "../fsutils",
			expected: ".",
		},
		{
			desc:     "testdata/example.txt",
			expected: "testdata/example.txt",
		},
		{
			desc:     "./testdata/example.txt",
			expected: "testdata/example.txt",
		},
		{
			desc:     "://broken",
			expected: "",
			err:      "missing protocol scheme",
		},
	}
	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			actual, err := RelativePath(tt.desc)

			if tt.err == "" {
				assert.NoError(t, err)
			} else {
				assert.ErrorContains(t, err, tt.err)
			}

			assert.Equal(t, tt.expected, actual)
		})
	}
}

func mustAbsPath(path string) string {
	absPath, err := filepath.Abs(path)
	if err != nil {
		panic(err)
	}
	return absPath
}
