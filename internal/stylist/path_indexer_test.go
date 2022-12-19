package stylist

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewPathIndexer(t *testing.T) {
	indexer := NewPathIndexer(
		[]string{"go", "shell", "yaml"},
		[]string{"**/Dockerfile"},
		[]string{},
	)

	assert.Equal(t, 3, indexer.FileTypes.Cardinality())
	assert.Equal(t, 1, indexer.Includes.Cardinality())
	assert.Equal(t, 0, indexer.Excludes.Cardinality())
}

func TestPathIndexer_Index(t *testing.T) {
	tests := []struct {
		desc                    string
		indexer                 *PathIndexer
		pathSpec                string
		expectedPathsByFileType map[string][]string
		expectedPathsByInclude  map[string][]string
		err                     string
	}{
		{
			desc: "accepts a path as pathSpec",
			indexer: NewPathIndexer(
				[]string{"text"},
				[]string{"**/*.md", "**/*.txt"},
				[]string{},
			),
			pathSpec: "testdata/aaa.txt",
			expectedPathsByFileType: map[string][]string{
				"text": {}, // TODO
			},
			expectedPathsByInclude: map[string][]string{
				"**/*.md": {},
				"**/*.txt": {
					"testdata/aaa.txt",
				},
			},
			err: "",
		},

		{
			desc: "accepts a dir as pathSpec",
			indexer: NewPathIndexer(
				[]string{"text"},
				[]string{"**/*.md", "**/*.txt"},
				[]string{},
			),
			pathSpec: "testdata/001/011",
			expectedPathsByFileType: map[string][]string{
				"text": {}, // TODO
			},
			expectedPathsByInclude: map[string][]string{
				"**/*.md": {},
				"**/*.txt": {
					"testdata/001/011/111/aaa.txt",
					"testdata/001/011/111/bbb.txt",
					"testdata/001/011/111/ccc.txt",
					"testdata/001/011/aaa.txt",
					"testdata/001/011/bbb.txt",
					"testdata/001/011/ccc.txt",
				},
			},
			err: "",
		},

		{
			desc: "accepts a pattern as pathSpec",
			indexer: NewPathIndexer(
				[]string{"text"},
				[]string{"**/*.md", "**/*.txt"},
				[]string{},
			),
			pathSpec: "testdata/**/aaa.txt",
			expectedPathsByFileType: map[string][]string{
				"text": {}, // TODO
			},
			expectedPathsByInclude: map[string][]string{
				"**/*.md": {},
				"**/*.txt": {
					"testdata/001/011/111/aaa.txt",
					"testdata/001/011/aaa.txt",
					"testdata/001/aaa.txt",
					"testdata/002/022/222/aaa.txt",
					"testdata/002/022/aaa.txt",
					"testdata/002/aaa.txt",
					"testdata/003/033/333/aaa.txt",
					"testdata/003/033/aaa.txt",
					"testdata/003/aaa.txt",
					"testdata/aaa.txt",
				},
			},
			err: "",
		},
		{
			desc: "does not match excluded patterns",
			indexer: NewPathIndexer(
				[]string{"text"},
				[]string{"**/*.md", "**/*.txt"},
				[]string{"testdata/aaa.txt", "testdata/003/**"},
			),
			pathSpec: "testdata/**/aaa.txt",
			expectedPathsByFileType: map[string][]string{
				"text": {}, // TODO
			},
			expectedPathsByInclude: map[string][]string{
				"**/*.md": {},
				"**/*.txt": {
					"testdata/001/011/111/aaa.txt",
					"testdata/001/011/aaa.txt",
					"testdata/001/aaa.txt",
					"testdata/002/022/222/aaa.txt",
					"testdata/002/022/aaa.txt",
					"testdata/002/aaa.txt",
				},
			},
			err: "",
		},

		{
			desc: "does not match unless configured",
			indexer: NewPathIndexer(
				[]string{},
				[]string{},
				[]string{},
			),
			pathSpec:                "testdata/**/aaa.txt",
			expectedPathsByFileType: map[string][]string{},
			expectedPathsByInclude:  map[string][]string{},
			err:                     "",
		},
		{
			desc: "returns an error if the pathSpec pattern base does not exist",
			indexer: NewPathIndexer(
				[]string{"text"},
				[]string{"**/*.md", "**/*.txt"},
				[]string{},
			),
			pathSpec: "does/not/exist/**/aaa.txt",
			expectedPathsByFileType: map[string][]string{
				"text": {},
			},
			expectedPathsByInclude: map[string][]string{
				"**/*.md":  {},
				"**/*.txt": {},
			},
			err: "pattern does not exist",
		},
		{
			desc: "but no error if pathSpec pattern simply fails to match",
			indexer: NewPathIndexer(
				[]string{"text"},
				[]string{"**/*.md", "**/*.txt"},
				[]string{},
			),
			pathSpec: "testdata/**/nope/*.txt",
			expectedPathsByFileType: map[string][]string{
				"text": {},
			},
			expectedPathsByInclude: map[string][]string{
				"**/*.md":  {},
				"**/*.txt": {},
			},
			err: "",
		},
		{
			desc: "returns an error if any configured patterns are malformed",
			indexer: NewPathIndexer(
				[]string{"text"},
				[]string{"**/*.{txt,,"},
				[]string{},
			),
			pathSpec: "testdata/**/aaa.txt",
			expectedPathsByFileType: map[string][]string{
				"text": {},
			},
			expectedPathsByInclude: map[string][]string{
				"**/*.{txt,,": {},
			},
			err: "syntax error in pattern",
		},
	}
	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			err := tt.indexer.Index(tt.pathSpec)

			if tt.err == "" {
				assert.NoError(t, err)
			} else {
				assert.ErrorContains(t, err, tt.err)
			}

			for ft, expected := range tt.expectedPathsByFileType {
				assert.ElementsMatch(t, expected, tt.indexer.PathsByFileType[ft].ToSlice())
			}
			for p, expected := range tt.expectedPathsByInclude {
				assert.ElementsMatch(t, expected, tt.indexer.PathsByInclude[p].ToSlice())
			}
		})
	}
}

func TestPathIndexer_Match(t *testing.T) {
	tests := []struct {
		desc              string
		indexer           *PathIndexer
		path              string
		expectedFileTypes []string
		expectedPatterns  []string
		err               string
	}{
		{
			desc: "matches a single pattern",
			indexer: NewPathIndexer(
				[]string{},
				[]string{"**/*.md", "**/*.txt"},
				[]string{},
			),
			path:              "foo/bar/baz.txt",
			expectedFileTypes: []string(nil),
			expectedPatterns:  []string{"**/*.txt"},
			err:               "",
		},
		{
			desc: "matches multiple patterns",
			indexer: NewPathIndexer(
				[]string{},
				[]string{"**/*.txt", "**/baz.*"},
				[]string{},
			),
			path:              "foo/bar/baz.txt",
			expectedFileTypes: []string(nil),
			expectedPatterns:  []string{"**/*.txt", "**/baz.*"},
			err:               "",
		},
		// TODO: add filetype assertions
	}
	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			actualFileTypes, actualPatterns, err := tt.indexer.match(tt.path)

			if tt.err == "" {
				assert.NoError(t, err)
			} else {
				assert.ErrorContains(t, err, tt.err)
			}

			assert.ElementsMatch(t, tt.expectedFileTypes, actualFileTypes)
			assert.ElementsMatch(t, tt.expectedPatterns, actualPatterns)
		})
	}
}
