package stylist

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewPathIndexer(t *testing.T) {
	indexer := NewPathIndexer(
		[]string{"**/Dockerfile"},
		[]string{},
	)

	assert.Equal(t, 1, indexer.Includes.Cardinality())
	assert.Equal(t, 0, indexer.Excludes.Cardinality())
}

func TestPathIndexer_Index(t *testing.T) {
	tests := []struct {
		desc                   string
		indexer                *PathIndexer
		pathSpec               string
		expectedPathsByInclude map[string][]string
		err                    string
	}{
		{
			desc: "accepts a path as pathSpec",
			indexer: NewPathIndexer(
				[]string{"**/*.md", "**/*.txt"},
				[]string{},
			),
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
			desc: "accepts a dir as pathSpec",
			indexer: NewPathIndexer(
				[]string{"**/*.md", "**/*.txt"},
				[]string{},
			),
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
			desc: "accepts a pattern as pathSpec",
			indexer: NewPathIndexer(
				[]string{"**/*.md", "**/*.txt"},
				[]string{},
			),
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
			desc: "does not match excluded patterns",
			indexer: NewPathIndexer(
				[]string{"**/*.md", "**/*.txt"},
				[]string{"testdata/txt/aaa.txt", "testdata/txt/003/**"},
			),
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
			desc: "does not match unless configured",
			indexer: NewPathIndexer(
				[]string{},
				[]string{},
			),
			pathSpec:               "testdata/txt/**/aaa.txt",
			expectedPathsByInclude: map[string][]string{},
			err:                    "",
		},
		{
			desc: "returns an error if the pathSpec pattern base does not exist",
			indexer: NewPathIndexer(
				[]string{"**/*.md", "**/*.txt"},
				[]string{},
			),
			pathSpec: "does/not/exist/**/aaa.txt",
			expectedPathsByInclude: map[string][]string{
				"**/*.md":  {},
				"**/*.txt": {},
			},
			err: "pattern does not exist",
		},
		{
			desc: "but no error if pathSpec pattern simply fails to match",
			indexer: NewPathIndexer(
				[]string{"**/*.md", "**/*.txt"},
				[]string{},
			),
			pathSpec: "testdata/txt/**/nope/*.txt",
			expectedPathsByInclude: map[string][]string{
				"**/*.md":  {},
				"**/*.txt": {},
			},
			err: "",
		},
		{
			desc: "returns an error if any configured patterns are malformed",
			indexer: NewPathIndexer(
				[]string{"**/*.{txt,,"},
				[]string{},
			),
			pathSpec: "testdata/txt/**/aaa.txt",
			expectedPathsByInclude: map[string][]string{
				"**/*.{txt,,": {},
			},
			err: "syntax error in pattern",
		},
	}
	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			app := NewTestApp()
			ctx := app.InitContext(context.Background())

			err := tt.indexer.Index(ctx, tt.pathSpec)

			if tt.err == "" {
				assert.NoError(t, err)
			} else {
				assert.ErrorContains(t, err, tt.err)
			}

			for p, expected := range tt.expectedPathsByInclude {
				assert.ElementsMatch(t, expected, tt.indexer.PathsByInclude[p].ToSlice())
			}
		})
	}
}

func TestPathIndexer_Match(t *testing.T) {
	tests := []struct {
		desc             string
		indexer          *PathIndexer
		path             string
		expectedPatterns []string
		err              string
	}{
		{
			desc: "matches a single pattern",
			indexer: NewPathIndexer(
				[]string{"**/*.md", "**/*.txt"},
				[]string{},
			),
			path:             "foo/bar/baz.txt",
			expectedPatterns: []string{"**/*.txt"},
			err:              "",
		},
		{
			desc: "matches multiple patterns",
			indexer: NewPathIndexer(
				[]string{"**/*.txt", "**/baz.*"},
				[]string{},
			),
			path:             "foo/bar/baz.txt",
			expectedPatterns: []string{"**/*.txt", "**/baz.*"},
			err:              "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			actualPatterns, err := tt.indexer.match(tt.path)

			if tt.err == "" {
				assert.NoError(t, err)
			} else {
				assert.ErrorContains(t, err, tt.err)
			}

			assert.ElementsMatch(t, tt.expectedPatterns, actualPatterns)
		})
	}
}
