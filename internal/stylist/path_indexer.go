package stylist

import (
	"context"
	"errors"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
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

func NormalizePath(basePath, path string) string {
	if !filepath.IsAbs(path) {
		path = filepath.Join(basePath, path)
	}
	return path
}

func NewNormalizedPathSet(base string, paths ...string) *NormalizedPathSet {
	base = strings.TrimSuffix(base, string(filepath.Separator))

	ps := &NormalizedPathSet{
		paths:    NewPathSet(),
		basePath: base,
	}
	for _, path := range paths {
		ps.Add(path)
	}

	return ps
}

type NormalizedPathSet struct {
	paths    PathSet
	basePath string
}

// Add adds the path to the set, first converting it to an
// absolute path (relative to the configured base) if needed.
// Returns whether the path was added.
func (ps *NormalizedPathSet) Add(path string) bool {
	return ps.paths.Add(ps.normalize(path))
}

// Contains returns whether the given paths are all
// in the set, first converting to absolute paths
// (relative to the configured base) if needed.
func (ps *NormalizedPathSet) Contains(paths ...string) bool {
	normalized := []string{}
	for _, path := range paths {
		normalized = append(normalized, ps.normalize(path))
	}
	return ps.paths.Contains(normalized...)
}

func (ps *NormalizedPathSet) normalize(path string) string {
	return NormalizePath(ps.basePath, path)
}

// AbsolutePaths returns a slice of all paths in the set
// as absolute paths.
func (ps *NormalizedPathSet) AbsolutePaths() []string {
	paths := ps.paths.ToSlice()
	sort.Strings(paths)
	return paths
}

// RelativePaths returns a slice of all paths in the set
// as relative paths.
func (ps *NormalizedPathSet) RelativePaths() []string {
	paths := []string{}
	for _, absPath := range ps.paths.ToSlice() {
		relPath, _ := filepath.Rel(ps.basePath, absPath)
		paths = append(paths, relPath)
	}
	sort.Strings(paths)
	return paths
}

// NewPathIndex returns a new, empty path index.
func NewPathIndex(basePath string) *PathIndex {
	return &PathIndex{
		basePath: basePath,
		pathSets: map[string]*NormalizedPathSet{},
	}
}

type PathIndex struct {
	basePath string
	pathSets map[string]*NormalizedPathSet
}

// Add adds the given pattern and path tuple to the index.
// Returns false if the tuple already exists in the index.
func (pi *PathIndex) Add(pattern, path string) bool {
	pattern = pi.normalize(pattern)
	path = pi.normalize(path)

	if _, ok := pi.pathSets[pattern]; !ok {
		pi.pathSets[pattern] = NewNormalizedPathSet(pi.basePath)
	}

	return pi.pathSets[pattern].Add(path)
}

func (pi *PathIndex) PathsFor(pattern string) *NormalizedPathSet {
	pattern = pi.normalize(pattern)
	if _, ok := pi.pathSets[pattern]; !ok {
		pi.pathSets[pattern] = NewNormalizedPathSet(pi.basePath)
	}
	return pi.pathSets[pattern]
}

func (pi *PathIndex) normalize(path string) string {
	return NormalizePath(pi.basePath, path)
}

// NewPathIndexer returns a new path index.
func NewPathIndexer(basePath string, includes, excludes []string) *PathIndexer {
	// Normalize include/exclude patterns to abs paths.
	includes = NewNormalizedPathSet(basePath, includes...).AbsolutePaths()
	excludes = NewNormalizedPathSet(basePath, excludes...).AbsolutePaths()

	indexer := &PathIndexer{
		basePath: basePath,
		includes: mapset.NewSet(includes...),
		excludes: mapset.NewSet(excludes...),
	}
	return indexer
}

// PathIndexer is a utility for indexing paths and grouping them by wildcard pattern.
type PathIndexer struct {
	// Set of patterns to include in the index.
	includes PathSet

	// Set of patterns to exclude from the index (even if a path would normally match).
	excludes PathSet

	basePath string
	ignorer  *PathIgnorer
	index    *PathIndex
	logger   *logrus.Logger
}

// Index resolves each pathSpec (a path or a wildcard pattern)
// to a list of paths and attempts to add them to the index.
// Paths will only be added to the index if they match
// the types and/or patterns registered with the indexer.
func (pi *PathIndexer) Index(ctx context.Context, pathSpecs ...string) (*PathIndex, error) {
	ignorer, err := NewPathIgnorer(".gitignore", pi.excludes.ToSlice())
	if err != nil {
		return nil, err
	}
	pi.ignorer = ignorer
	pi.index = NewPathIndex(pi.basePath)
	pi.logger = AppLogger(ctx)

	pi.logger.Debugf("[index] Includes=%v", pi.includes)
	pi.logger.Debugf("[index] Excludes=%v", pi.excludes)

	files, dirs, patterns, err := pi.partitionPathSpecs(pathSpecs)
	if err != nil {
		return nil, err
	}

	pi.logger.Debugf(
		"[index] Partition: files=%v dirs=%v patterns=%s",
		files,
		dirs,
		patterns,
	)

	for _, f := range files {
		if err := pi.indexPath(f); err != nil {
			return nil, err
		}
	}
	for _, d := range dirs {
		if err := pi.indexDir(d); err != nil {
			return nil, err
		}
	}
	if len(patterns) > 0 {
		if err := pi.indexPatterns(patterns); err != nil {
			return nil, err
		}
	}

	return pi.index, nil
}

func (pi *PathIndexer) partitionPathSpecs(pathSpecs []string) (
	[]string, []string, []string, error,
) {
	var fileSpecs []string
	var dirSpecs []string
	var patternSpecs []string

	for _, pathSpec := range pathSpecs {
		// Ensure abs path
		pathSpec = NormalizePath(pi.basePath, pathSpec)

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
	return filepath.WalkDir(pi.basePath, func(path string, d fs.DirEntry, err error) error {
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
	for _, pattern := range matchedPatterns {
		pi.index.Add(pattern, path)
	}

	return nil
}

// exclude returns true if path should be excluded.
func (pi *PathIndexer) exclude(path string, isDir bool) bool {
	if pi.ignorer.ShouldIgnore(path, isDir) {
		pi.logger.Debugf("[index] Ignore path=%s", path)
		return true
	}
	return false
}

// match returns the patterns that match path.
func (pi *PathIndexer) match(path string) ([]string, error) {
	var matchedPatterns []string

	for pattern := range pi.includes.Iter() {
		ok, err := matchPattern(pattern, path)
		if err != nil {
			return matchedPatterns, err
		}
		if ok {
			pi.logger.Debugf(
				"[index] Match: path=%v pattern=%v",
				path,
				pattern,
			)
			matchedPatterns = append(matchedPatterns, pattern)
		}
	}

	return matchedPatterns, nil
}
