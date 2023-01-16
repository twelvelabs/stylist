package fsutils

import (
	"fmt"
	"os"
	"sync"
)

// NewFileCache returns a new file cache.
func NewFileCache() *FileCache {
	return &FileCache{}
}

// FileCache is a utility for caching file contents.
type FileCache struct {
	files sync.Map
}

// GetFileBytes returns the cached bytes for path.
func (fc *FileCache) GetFileBytes(path string) ([]byte, bool, error) {
	cachedBytes, ok := fc.files.Load(path)
	if ok {
		return cachedBytes.([]byte), true, nil
	}

	fileBytes, err := os.ReadFile(path)
	if err != nil {
		return nil, false, fmt.Errorf("file cache: %w", err)
	}

	fc.files.Store(path, fileBytes)
	return fileBytes, false, nil
}
