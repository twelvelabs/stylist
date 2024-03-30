package stylist

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNormalizedPathSet(t *testing.T) {
	pathSet := NewNormalizedPathSet(
		// The base path
		"/aaa/bbb",
		// Some abs paths...
		"/aaa/bbb/333.txt",
		"/aaa/111.txt",
		// ... and a few relative ones.
		"ccc/111.txt",
		"222.txt",
		"111.txt",
	)

	assert.Equal(t, []string{
		"/aaa/111.txt",
		"/aaa/bbb/111.txt",
		"/aaa/bbb/222.txt",
		"/aaa/bbb/333.txt",
		"/aaa/bbb/ccc/111.txt",
	}, pathSet.AbsolutePaths())
	assert.Equal(t, true, pathSet.Contains("/aaa/111.txt"))
	assert.Equal(t, true, pathSet.Contains("/aaa/bbb/111.txt"))
	assert.Equal(t, true, pathSet.Contains("/aaa/bbb/ccc/111.txt"))
	assert.Equal(t, false, pathSet.Contains("/000.txt"))

	assert.Equal(t, []string{
		"../111.txt",
		"111.txt",
		"222.txt",
		"333.txt",
		"ccc/111.txt",
	}, pathSet.RelativePaths())
	assert.Equal(t, true, pathSet.Contains("../111.txt"))
	assert.Equal(t, true, pathSet.Contains("111.txt"))
	assert.Equal(t, true, pathSet.Contains("ccc/111.txt"))
	assert.Equal(t, false, pathSet.Contains("000.txt"))
}

func TestPathIndex(t *testing.T) {
	index := NewPathIndex("")

	// Populate the index...
	assert.Equal(t, true, index.Add("**/*", "/aaa.txt"))
	assert.Equal(t, true, index.Add("**/*.txt", "/aaa.txt"))
	assert.Equal(t, true, index.Add("**/*", "/bbb.go"))
	assert.Equal(t, true, index.Add("**/*.go", "/bbb.go"))
	// Already in the index.
	assert.Equal(t, false, index.Add("**/*.go", "/bbb.go"))

	assert.ElementsMatch(t, []string{
		"/aaa.txt",
		"/bbb.go",
	}, index.PathsFor("**/*").AbsolutePaths())

	assert.ElementsMatch(t, []string{
		"/aaa.txt",
	}, index.PathsFor("**/*.txt").AbsolutePaths())

	assert.ElementsMatch(t, []string{
		"/bbb.go",
	}, index.PathsFor("**/*.go").AbsolutePaths())

	assert.ElementsMatch(t, []string{}, index.PathsFor("unknown/pattern").AbsolutePaths())
}

func TestPathIndexer_Index(t *testing.T) {
	tests := []struct {
		desc                   string
		basePath               string
		includes               []string
		excludes               []string
		pathSpec               string
		expectedPathsByInclude map[string][]string
		err                    string
	}{
		{
			desc:     "accepts a path as pathSpec",
			basePath: ".",
			includes: []string{"**/*.md", "**/*.txt"},
			pathSpec: "testdata/txt/aaa.txt",
			expectedPathsByInclude: map[string][]string{
				"**/*.md": {},
				"**/*.txt": {
					"testdata/txt/aaa.txt",
				},
			},
			err: "",
		},

		{
			desc:     "accepts a dir as pathSpec",
			basePath: ".",
			includes: []string{"**/*.md", "**/*.txt"},
			pathSpec: "testdata/txt/001/011",
			expectedPathsByInclude: map[string][]string{
				"**/*.md": {},
				"**/*.txt": {
					"testdata/txt/001/011/111/aaa.txt",
					"testdata/txt/001/011/111/bbb.txt",
					"testdata/txt/001/011/111/ccc.txt",
					"testdata/txt/001/011/aaa.txt",
					"testdata/txt/001/011/bbb.txt",
					"testdata/txt/001/011/ccc.txt",
				},
			},
			err: "",
		},

		{
			desc:     "accepts a pattern as pathSpec",
			basePath: ".",
			includes: []string{"**/*.md", "**/*.txt"},
			pathSpec: "testdata/txt/**/aaa.txt",
			expectedPathsByInclude: map[string][]string{
				"**/*.md": {},
				"**/*.txt": {
					"testdata/txt/001/011/111/aaa.txt",
					"testdata/txt/001/011/aaa.txt",
					"testdata/txt/001/aaa.txt",
					"testdata/txt/002/022/222/aaa.txt",
					"testdata/txt/002/022/aaa.txt",
					"testdata/txt/002/aaa.txt",
					"testdata/txt/003/033/333/aaa.txt",
					"testdata/txt/003/033/aaa.txt",
					"testdata/txt/003/aaa.txt",
					"testdata/txt/aaa.txt",
				},
			},
			err: "",
		},
		{
			desc:     "does not match excluded patterns",
			basePath: ".",
			includes: []string{"**/*.md", "**/*.txt"},
			excludes: []string{"testdata/txt/aaa.txt", "testdata/txt/003/**"},
			pathSpec: "testdata/txt/**/aaa.txt",
			expectedPathsByInclude: map[string][]string{
				"**/*.md": {},
				"**/*.txt": {
					"testdata/txt/001/011/111/aaa.txt",
					"testdata/txt/001/011/aaa.txt",
					"testdata/txt/001/aaa.txt",
					"testdata/txt/002/022/222/aaa.txt",
					"testdata/txt/002/022/aaa.txt",
					"testdata/txt/002/aaa.txt",
				},
			},
			err: "",
		},

		{
			desc:                   "does not match unless configured",
			basePath:               ".",
			pathSpec:               "testdata/txt/**/aaa.txt",
			expectedPathsByInclude: map[string][]string{},
			err:                    "",
		},
		{
			desc:                   "returns an error if the pathSpec does not exist",
			basePath:               ".",
			includes:               []string{"**/*.md", "**/*.txt"},
			pathSpec:               "does/not/exist/aaa.txt",
			expectedPathsByInclude: map[string][]string{},
			err:                    "no such file or directory",
		},
		{
			desc:                   "returns an error if the pathSpec pattern base does not exist",
			basePath:               ".",
			includes:               []string{"**/*.md", "**/*.txt"},
			pathSpec:               "does/not/exist/**/aaa.txt",
			expectedPathsByInclude: map[string][]string{},
			err:                    "pattern does not exist",
		},
		{
			desc:     "but no error if pathSpec pattern simply fails to match",
			basePath: ".",
			includes: []string{"**/*.md", "**/*.txt"},
			pathSpec: "testdata/txt/**/nope/*.txt",
			expectedPathsByInclude: map[string][]string{
				"**/*.md":  {},
				"**/*.txt": {},
			},
			err: "",
		},
		{
			desc:                   "returns an error if pathSpec patterns are malformed",
			basePath:               ".",
			includes:               []string{"**/*"},
			pathSpec:               "testdata/txt/**/aaa.{txt,,",
			expectedPathsByInclude: map[string][]string{},
			err:                    "syntax error in pattern",
		},
		{
			desc:                   "returns an error if any configured patterns are malformed",
			basePath:               ".",
			includes:               []string{"**/*.{txt,,"},
			pathSpec:               "testdata/txt/**/aaa.txt",
			expectedPathsByInclude: map[string][]string{},
			err:                    "syntax error in pattern",
		},
	}
	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			app := NewTestApp()
			ctx := app.InitContext(context.Background())

			indexer := NewPathIndexer(tt.basePath, tt.includes, tt.excludes)
			index, err := indexer.Index(ctx, tt.pathSpec)

			if tt.err != "" {
				assert.ErrorContains(t, err, tt.err)
			} else {
				assert.NoError(t, err)
				for pattern, expected := range tt.expectedPathsByInclude {
					assert.ElementsMatch(t, expected, index.PathsFor(pattern).RelativePaths())
				}
			}
		})
	}
}

func TestPathIndexer_Match(t *testing.T) {
	tests := []struct {
		desc             string
		basePath         string
		includes         []string
		excludes         []string
		path             string
		expectedPatterns []string
		err              string
	}{
		{
			desc:             "matches a single pattern",
			includes:         []string{"**/*.md", "**/*.txt"},
			path:             "foo/bar/baz.txt",
			expectedPatterns: []string{"**/*.txt"},
			err:              "",
		},
		{
			desc:             "matches multiple patterns",
			includes:         []string{"**/*.txt", "**/baz.*"},
			path:             "foo/bar/baz.txt",
			expectedPatterns: []string{"**/*.txt", "**/baz.*"},
			err:              "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			app := NewTestApp()

			indexer := NewPathIndexer(tt.basePath, tt.includes, tt.excludes)
			indexer.logger = app.Logger

			actualPatterns, err := indexer.match(tt.path)

			if tt.err == "" {
				assert.NoError(t, err)
			} else {
				assert.ErrorContains(t, err, tt.err)
			}

			assert.ElementsMatch(t, tt.expectedPatterns, actualPatterns)
		})
	}
}
