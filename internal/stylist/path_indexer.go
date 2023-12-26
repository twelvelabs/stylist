package stylist

import (
	"context"
	"errors"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/bmatcuk/doublestar/v4"
	mapset "github.com/deckarep/golang-set/v2"
	"github.com/sirupsen/logrus"
)

var (
	// Characters used in doublestar patterns.
	patternChars = "*?[{"
	// For stubbing.
	matchPattern = matchPatternFunc
)

// Wrapper around doublestar.Match that ensures both args are using "/".
// This allows the patterns in `stylist.yml` to be written in either
// Posix or Windows style, but still be usable cross-platform.
func matchPatternFunc(pattern string, path string) (bool, error) {
	ok, err := doublestar.Match(filepath.ToSlash(pattern), filepath.ToSlash(path))
	if err != nil {
		return false, err
	}
	return ok, nil
}

func NewPathSet(paths ...string) PathSet { //nolint:ireturn
	return mapset.NewSet(paths...)
}

// PathSet is a unique set of filesystem paths.
type PathSet mapset.Set[string]

// NewPathIndexer returns a new path index.
func NewPathIndexer(includes, excludes []string) *PathIndexer {
	indexer := &PathIndexer{
		Includes:       mapset.NewSet(includes...),
		Excludes:       mapset.NewSet(excludes...),
		PathsByInclude: map[string]PathSet{},
	}
	for p := range indexer.Includes.Iter() {
		indexer.PathsByInclude[p] = NewPathSet()
	}
	return indexer
}

// PathIndexer is a utility for indexing paths and grouping them by wildcard pattern.
type PathIndexer struct {
	// Set of patterns to include in the index.
	Includes mapset.Set[string]

	// Set of patterns to exclude from the index (even if a path would normally match).
	Excludes mapset.Set[string]

	// Paths grouped by wildcard pattern.
	PathsByInclude map[string]PathSet

	ignorer *PathIgnorer
	logger  *logrus.Logger
}

// Cardinality returns the total number of patterns
// that the indexer is configured to match.
func (pi *PathIndexer) Cardinality() int {
	return pi.Includes.Cardinality()
}

// Index resolves each pathSpec (a path or a wildcard pattern)
// to a list of paths and attempts to add them to the index.
// Paths will only be added to the index if they match
// the types and/or patterns registered with the indexer.
func (pi *PathIndexer) Index(ctx context.Context, pathSpecs ...string) error {
	ignorer, err := NewPathIgnorer(".gitignore", pi.Excludes.ToSlice())
	if err != nil {
		return err
	}
	pi.ignorer = ignorer
	pi.logger = AppLogger(ctx)

	files, dirs, patterns, err := pi.partitionPathSpecs(pathSpecs)
	if err != nil {
		return err
	}

	for _, f := range files {
		if err := pi.indexPath(f); err != nil {
			return err
		}
	}
	for _, d := range dirs {
		if err := pi.indexDir(d); err != nil {
			return err
		}
	}
	if len(patterns) > 0 {
		if err := pi.indexPatterns(patterns); err != nil {
			return err
		}
	}
	return nil
}

func (pi *PathIndexer) partitionPathSpecs(pathSpecs []string) (
	[]string, []string, []string, error,
) {
	var fileSpecs []string
	var dirSpecs []string
	var patternSpecs []string

	for _, pathSpec := range pathSpecs {
		if strings.ContainsAny(pathSpec, patternChars) {
			pattern := filepath.ToSlash(filepath.Clean(pathSpec))

			base, _ := doublestar.SplitPattern(pattern)
			if _, err := os.Lstat(base); err != nil {
				if errors.Is(err, os.ErrNotExist) {
					return nil, nil, nil, doublestar.ErrPatternNotExist
				}
				return nil, nil, nil, err
			}

			patternSpecs = append(patternSpecs, pattern)
			continue
		}

		info, err := os.Lstat(pathSpec)
		if err != nil {
			return nil, nil, nil, err
		}
		if info.IsDir() {
			dirSpecs = append(dirSpecs, pathSpec)
		} else {
			fileSpecs = append(fileSpecs, pathSpec)
		}
	}

	return fileSpecs, dirSpecs, patternSpecs, nil
}

func (pi *PathIndexer) indexPatterns(patterns []string) error {
	return filepath.WalkDir(".", func(path string, d fs.DirEntry, err error) error {
		return pi.indexWalkedPath(path, d, err, patterns...)
	})
}

// indexDir walks every path in dir and calls indexWalkedPath().
func (pi *PathIndexer) indexDir(dir string) error {
	return filepath.WalkDir(dir, func(path string, d fs.DirEntry, err error) error {
		return pi.indexWalkedPath(path, d, err)
	})
}

// indexWalkedPath is called by indexPattern and indexDir as they walk
// the filesystem.
func (pi *PathIndexer) indexWalkedPath(
	path string, d fs.DirEntry, err error, pathSpecPatterns ...string,
) error {
	if err != nil {
		return err // walk error, likely FS permissions
	}

	// For directories, check to see if the path matches any exclude patterns,
	// and if so, return `fs.SkipDir`. This is a sentinel that `filepath.WalkDir`
	// uses to determine whether to recurse into subdirectories.
	if d.IsDir() {
		if pi.exclude(path, true) {
			return fs.SkipDir
		}
		return nil
	}

	// Otherwise, figure out whether we should index the path.
	var shouldIndex bool
	if len(pathSpecPatterns) == 0 {
		shouldIndex = true
	} else {
		for _, pattern := range pathSpecPatterns {
			ok, err := matchPattern(pattern, path)
			if err != nil {
				return err // doublestar parse error
			}
			if ok {
				shouldIndex = true
			}
		}
	}

	if shouldIndex {
		return pi.indexPath(path)
	}
	return nil
}

// indexPath adds path to the index if it matches.
func (pi *PathIndexer) indexPath(path string) error {
	// Check if it's excluded
	if pi.exclude(path, false) {
		return nil
	}

	// Check if it matches
	matchedPatterns, err := pi.match(path)
	if err != nil {
		return err
	}
	for _, p := range matchedPatterns {
		pi.PathsByInclude[p].Add(path)
	}

	return nil
}

// exclude returns true if path should be excluded.
func (pi *PathIndexer) exclude(path string, isDir bool) bool {
	if pi.ignorer.ShouldIgnore(path, isDir) {
		pi.logger.Debugf("Index ignore dir=%v path=%s", isDir, path)
		return true
	}
	return false
}

// match returns the patterns that match path.
func (pi *PathIndexer) match(path string) ([]string, error) {
	var matchedPatterns []string

	// Optimization: nothing to match, so don't bother checking.
	if pi.Cardinality() == 0 {
		return matchedPatterns, nil
	}

	for pattern := range pi.Includes.Iter() {
		ok, err := matchPattern(pattern, path)
		if err != nil {
			return matchedPatterns, err
		}
		if ok {
			matchedPatterns = append(matchedPatterns, pattern)
		}
	}

	return matchedPatterns, nil
}
