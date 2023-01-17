package stylist

import (
	"errors"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/bmatcuk/doublestar/v4"
	mapset "github.com/deckarep/golang-set/v2"
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
func (pi *PathIndexer) Index(pathSpecs ...string) error {
	for _, pathSpec := range pathSpecs {
		if err := pi.indexPathSpec(pathSpec); err != nil {
			return err
		}
	}
	return nil
}

// indexPathSpec dispatches to the appropriate index method for the given type.
func (pi *PathIndexer) indexPathSpec(pathSpec string) error {
	if strings.ContainsAny(pathSpec, patternChars) {
		return pi.indexPattern(pathSpec)
	}
	info, err := os.Stat(pathSpec)
	if err != nil {
		return err
	}
	if info.IsDir() {
		return pi.indexDir(pathSpec)
	}
	return pi.indexPath(pathSpec)
}

// indexPattern walks every path in pattern and calls indexWalkedPath().
func (pi *PathIndexer) indexPattern(pattern string) error {
	pattern = filepath.ToSlash(filepath.Clean(pattern))

	// Replicate the pattern validation logic from doublestar's Glob* functions.
	base, _ := doublestar.SplitPattern(pattern)
	if _, err := os.Lstat(base); err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return doublestar.ErrPatternNotExist
		}
		return err
	}

	// Note: We're intentionally using `filepath.WalkDir` (vs. using doublestar's
	// Glob or GlobWalk functions) because this allows us to skip excluded dirs
	// (see logic in indexWalkedPath).
	// For most projects this ends up being way faster, especially projects
	// containing large cache or node_modules directories.
	// Skipping those dirs can cut the time in half.
	return filepath.WalkDir(".", func(path string, d fs.DirEntry, err error) error {
		return pi.indexWalkedPath(path, d, err, pattern)
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
		ok, err := pi.exclude(path)
		if err != nil {
			return err
		}
		if ok {
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
	ok, err := pi.exclude(path)
	if err != nil {
		return err
	}
	if ok {
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
func (pi *PathIndexer) exclude(path string) (bool, error) {
	for exclude := range pi.Excludes.Iter() {
		ok, err := matchPattern(exclude, path)
		if err != nil {
			return false, err // doublestar parse error
		}
		if ok {
			return true, nil
		}
	}
	return false, nil
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
