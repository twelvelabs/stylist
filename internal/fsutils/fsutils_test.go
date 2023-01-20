package fsutils

import (
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
