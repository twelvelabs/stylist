package fsutils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewFileCache(t *testing.T) {
	assert.IsType(t, &FileCache{}, NewFileCache())
}

func TestFileCache_GetFileBytes(t *testing.T) {
	cache := NewFileCache()

	_, _, err := cache.GetFileBytes("does/not/exist.txt")
	assert.Error(t, err)

	data, ok, err := cache.GetFileBytes("testdata/example.txt")
	assert.Equal(t, "one\ntwo\nthree\n", string(data))
	assert.Equal(t, false, ok)
	assert.NoError(t, err)

	data, ok, err = cache.GetFileBytes("testdata/example.txt")
	assert.Equal(t, "one\ntwo\nthree\n", string(data))
	assert.Equal(t, true, ok)
	assert.NoError(t, err)
}
