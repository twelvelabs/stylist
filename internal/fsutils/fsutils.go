package fsutils

import (
	"errors"
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"strings"
)

func NoPathExists(path string) bool {
	_, err := os.Stat(path)
	return errors.Is(err, os.ErrNotExist)
}

func PathExists(path string) bool {
	return !NoPathExists(path)
}

// RelativePath ensures path is relative to the current working dir.
// Also trims URI scheme or leading ./ segments.
func RelativePath(path string) (string, error) {
	parsed, err := url.Parse(path)
	if err != nil {
		return "", fmt.Errorf("relative path: %w", err)
	}

	absPath, _ := filepath.Abs(parsed.Path)
	cwd, _ := os.Getwd()
	relPath, _ := filepath.Rel(cwd, absPath)

	return strings.TrimPrefix(relPath, "./"), nil
}
