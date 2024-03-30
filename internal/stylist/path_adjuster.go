package stylist

import (
	"fmt"
	"path/filepath"
)

// NewPathAdjuster returns a new path adjuster for the given args.
func NewPathAdjuster(basePath string, pathType ResultPath) *PathAdjuster {
	return &PathAdjuster{
		basePath: basePath,
		pathType: pathType,
	}
}

// PathAdjuster converts paths to the configured path type.
type PathAdjuster struct {
	basePath string
	pathType ResultPath
}

func (pa *PathAdjuster) Convert(path string) (string, error) {
	var err error

	switch pa.pathType {
	case ResultPathAbsolute:
		if !filepath.IsAbs(path) {
			path, err = filepath.Abs(filepath.Join(pa.basePath, path))
		} else {
			path = filepath.Clean(path)
		}
	case ResultPathRelative:
		if filepath.IsAbs(path) {
			path, err = filepath.Rel(pa.basePath, path)
		} else {
			path = filepath.Clean(path)
		}
	}

	if err != nil {
		err = fmt.Errorf("unable to convert to %s path: %w", pa.pathType.String(), err)
	}

	return path, err
}
