package stylist

import (
	"github.com/twelvelabs/stylist/internal/fsutils"
)

func NewContextLineLoader() *ContextLineLoader {
	return &ContextLineLoader{
		lineCache: fsutils.NewLineCache(fsutils.NewFileCache()),
	}
}

type ContextLineLoader struct {
	lineCache *fsutils.LineCache
}

func (l *ContextLineLoader) Load(loc ResultLocation) ([]string, error) {
	if loc.Path == "" || loc.StartLine == 0 {
		return nil, nil
	}

	start, end := loc.LineRange()
	lines := []string{}

	for i := start; i <= end; i++ {
		line, err := l.lineCache.GetLine(loc.Path, i)
		if err != nil {
			return nil, err
		}
		lines = append(lines, line)
	}

	return lines, nil
}
