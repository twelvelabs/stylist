package fsutils

import (
	"bytes"
	"fmt"
	"sync"
)

// NewLineCache returns a new line cache.
func NewLineCache(fc *FileCache) *LineCache {
	return &LineCache{
		fileCache: fc,
	}
}

// LineCache is a utility for caching individual lines from a file.
type LineCache struct {
	rawLines  sync.Map
	fileCache *FileCache
}

// GetLine returns the index1-th (1-based index) line from the file at path.
func (lc *LineCache) GetLine(path string, index1 int) (string, error) {
	if index1 == 0 {
		index1 = 1
	}

	const index1To0Offset = -1
	line, err := lc.getRawLine(path, index1+index1To0Offset)
	if err != nil {
		return "", err
	}

	return string(bytes.Trim(line, "\r")), nil
}

func (lc *LineCache) getRawLine(path string, index0 int) ([]byte, error) {
	lines, _, err := lc.getRawLines(path)
	if err != nil {
		return nil, fmt.Errorf("line cache: %w", err)
	}

	if index0 < 0 || index0 >= len(lines) {
		return nil, fmt.Errorf(
			"line cache: index out of bounds: index0=%d, len(lines)=%d",
			index0,
			len(lines),
		)
	}

	return lines[index0], nil
}

func (lc *LineCache) getRawLines(path string) ([][]byte, bool, error) {
	if loaded, ok := lc.rawLines.Load(path); ok {
		return loaded.([][]byte), true, nil
	}

	fileBytes, _, err := lc.fileCache.GetFileBytes(path)
	if err != nil {
		return nil, false, err
	}

	lines := bytes.Split(fileBytes, []byte("\n"))
	lc.rawLines.Store(path, lines)
	return lines, false, nil
}
