package stylist

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/bmatcuk/doublestar/v4"
	"github.com/denormal/go-gitignore"

	"github.com/twelvelabs/stylist/internal/fsutils"
)

var (
	newGitIgnoreParser = gitignore.NewFromFile
)

// NewPathIgnorer returns a new path ignorer.
func NewPathIgnorer(gitIgnorePath string, patterns []string) (*PathIgnorer, error) {
	cleanPatterns := []string{}
	for _, pattern := range patterns {
		p := strings.ReplaceAll(pattern, "\\", "/")
		if !doublestar.ValidatePattern(p) {
			return nil, doublestar.ErrBadPattern
		}
		cleanPatterns = append(cleanPatterns, p)
	}

	pi := &PathIgnorer{
		patterns: cleanPatterns,
	}
	if gitIgnorePath != "" && fsutils.PathExists(gitIgnorePath) {
		gitIgnore, err := newGitIgnoreParser(gitIgnorePath)
		if err != nil {
			return nil, fmt.Errorf("gitignore parse: %w", err)
		}
		pi.gitIgnore = gitIgnore
		pi.gitIgnorePath = gitIgnorePath
	}
	return pi, nil
}

// PathIgnorer is responsible for ignoring paths during the index process.
type PathIgnorer struct {
	gitIgnore     gitignore.GitIgnore
	gitIgnorePath string
	patterns      []string
}

// ShouldIgnore returns true if path should be ignored.
func (pi *PathIgnorer) ShouldIgnore(path string, isDir bool) bool {
	if pi.gitIgnore != nil {
		// Note: Specifically using `.Relative` for speed.
		//       The higher level GitIgnore methods (Ignore, Match)
		//       stat the path to figure out isDir.
		if match := pi.gitIgnore.Relative(path, isDir); match != nil {
			if match.Ignore() {
				return true
			}
		}
	}
	for _, pattern := range pi.patterns {
		if ok, _ := doublestar.Match(pattern, filepath.ToSlash(path)); ok {
			return true
		}
	}
	return false
}
