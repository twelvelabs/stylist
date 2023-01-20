package fsutils

import (
	"errors"
	"os"
)

func NoPathExists(path string) bool {
	_, err := os.Stat(path)
	return errors.Is(err, os.ErrNotExist)
}

func PathExists(path string) bool {
	return !NoPathExists(path)
}
