package stylist

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/twelvelabs/termite/run"
)

func TestNewPipeline(t *testing.T) {
	pipeline := NewPipeline([]*Processor{}, []string{})
	assert.IsType(t, &Pipeline{}, pipeline)
}

func TestPipeline_Check(t *testing.T) {
	tests := []struct {
		desc       string
		setupFunc  func(c *run.Client)
		processors []*Processor
		excludes   []string
		pathSpecs  []string
		expected   []*Result
		err        string
	}{
		{
			desc: "should ignore processors missing check commands",
			processors: []*Processor{
				{
					CheckCommand: nil,
				},
			},
			pathSpecs: []string{
				"testdata/txt",
			},
			expected: []*Result{},
		},
		{
			desc: "should run the check commands for processors that have them",
			setupFunc: func(c *run.Client) {
				c.RegisterStub(
					run.MatchString("pretend-linter testdata/txt/aaa.txt"),
					run.StdoutResponse([]byte(""), 0),
				)
				c.RegisterStub(
					run.MatchString("pretend-linter testdata/txt/bbb.txt"),
					run.StdoutResponse([]byte(""), 0),
				)
				c.RegisterStub(
					run.MatchString("pretend-linter testdata/txt/ccc.txt"),
					run.StdoutResponse([]byte("lint failure"), 1),
				)
			},
			processors: []*Processor{
				{
					Includes: []string{"testdata/txt/*.txt"},
					CheckCommand: &Command{
						Template:     "pretend-linter",
						InputType:    InputTypeArg,
						OutputType:   OutputTypeStdout,
						OutputFormat: OutputFormatNone,
					},
				},
			},
			pathSpecs: []string{
				"testdata/txt",
			},
			expected: []*Result{
				{
					Level: ResultLevelError,
					Location: ResultLocation{
						Path: "testdata/txt/ccc.txt",
					},
					Rule: ResultRule{
						Description: "Unknown issue",
					},
					ContextLang: "plaintext",
					ContextLines: []string{
						"lint failure",
					},
				},
			},
		},
		{
			desc: "should return errors for any commands that fail unexpectedly",
			setupFunc: func(c *run.Client) {
				c.RegisterStub(
					run.MatchString("pretend-linter testdata/txt/aaa.txt"),
					run.StdoutResponse([]byte(""), 0),
				)
				c.RegisterStub(
					run.MatchString("pretend-linter testdata/txt/bbb.txt"),
					run.StdoutResponse([]byte(""), 0),
				)
				c.RegisterStub(
					run.MatchString("pretend-linter testdata/txt/ccc.txt"),
					run.ErrorResponse(errors.New("boom")),
				)
			},
			processors: []*Processor{
				{
					Includes: []string{"testdata/txt/*.txt"},
					CheckCommand: &Command{
						Template:     "pretend-linter",
						InputType:    InputTypeArg,
						OutputType:   OutputTypeStdout,
						OutputFormat: OutputFormatNone,
					},
				},
			},
			pathSpecs: []string{
				"testdata/txt",
			},
			err: "boom",
		},
	}

	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			app := NewTestApp()
			defer app.CmdClient.VerifyStubs(t)

			if tt.setupFunc != nil {
				tt.setupFunc(app.CmdClient)
			}

			ctx := app.InitContext(context.Background())

			pipeline := NewPipeline(tt.processors, tt.excludes)
			actual, err := pipeline.Check(ctx, tt.pathSpecs)

			if tt.err == "" {
				assert.NoError(t, err)
			} else {
				assert.ErrorContains(t, err, tt.err)
			}

			assert.Equal(t, tt.expected, actual)
		})
	}
}

func TestPipeline_Index(t *testing.T) {
	tests := []struct {
		desc       string
		processors []*Processor
		excludes   []string
		pathSpecs  []string
		expectFunc func(t *testing.T, pipeline *Pipeline)
		err        string
	}{
		{
			desc:       "should handle empty values",
			processors: []*Processor{},
			excludes:   []string{},
			pathSpecs: []string{
				"testdata/txt",
			},
		},
		{
			desc: "should filter files for each processor",
			processors: []*Processor{
				{
					Includes: []string{"testdata/txt/**/*.txt"},
					Excludes: []string{"testdata/txt/002/**"},
				},
				{
					Includes: []string{"testdata/txt/**/*.txt"},
					Excludes: []string{"testdata/txt/003/**"},
				},
			},
			excludes: []string{
				"testdata/txt/**/bbb.txt",
				"testdata/txt/**/ccc.txt",
			},
			pathSpecs: []string{
				"testdata/txt",
			},
			expectFunc: func(t *testing.T, pipeline *Pipeline) {
				t.Helper()

				p1 := pipeline.processors[0]
				assert.ElementsMatch(t, []string{
					"testdata/txt/001/011/111/aaa.txt",
					"testdata/txt/001/011/aaa.txt",
					"testdata/txt/001/aaa.txt",
					"testdata/txt/003/033/333/aaa.txt",
					"testdata/txt/003/033/aaa.txt",
					"testdata/txt/003/aaa.txt",
					"testdata/txt/aaa.txt",
				}, p1.Paths())

				p2 := pipeline.processors[1]
				assert.ElementsMatch(t, []string{
					"testdata/txt/001/011/111/aaa.txt",
					"testdata/txt/001/011/aaa.txt",
					"testdata/txt/001/aaa.txt",
					"testdata/txt/002/022/222/aaa.txt",
					"testdata/txt/002/022/aaa.txt",
					"testdata/txt/002/aaa.txt",
					"testdata/txt/aaa.txt",
				}, p2.Paths())
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			app := NewTestApp()
			ctx := app.InitContext(context.Background())

			pipeline := NewPipeline(tt.processors, tt.excludes)
			err := pipeline.Index(ctx, tt.pathSpecs)

			if tt.err == "" {
				assert.NoError(t, err)
			} else {
				assert.ErrorContains(t, err, tt.err)
			}

			if tt.expectFunc != nil {
				tt.expectFunc(t, pipeline)
			}
		})
	}
}

func TestEnsureContextLines(t *testing.T) {
	tests := []struct {
		desc     string
		config   *OutputConfig
		results  []*Result
		expected []*Result
		err      string
	}{
		{
			desc: "ensures context by default",
			results: []*Result{
				{
					Location: ResultLocation{
						Path:      "testdata/txt/aaa.txt",
						StartLine: 1,
						EndLine:   1,
					},
				},
				{
					Location: ResultLocation{
						Path:      "testdata/txt/bbb.txt",
						StartLine: 1,
						EndLine:   1,
					},
				},
				{
					Location: ResultLocation{
						Path:      "testdata/txt/ccc.txt",
						StartLine: 1,
						EndLine:   1,
					},
				},
			},
			expected: []*Result{
				{
					Location: ResultLocation{
						Path:      "testdata/txt/aaa.txt",
						StartLine: 1,
						EndLine:   1,
					},
					ContextLang: "plaintext",
					ContextLines: []string{
						"aaa content",
					},
				},
				{
					Location: ResultLocation{
						Path:      "testdata/txt/bbb.txt",
						StartLine: 1,
						EndLine:   1,
					},
					ContextLang: "plaintext",
					ContextLines: []string{
						"bbb content",
					},
				},
				{
					Location: ResultLocation{
						Path:      "testdata/txt/ccc.txt",
						StartLine: 1,
						EndLine:   1,
					},
					ContextLang: "plaintext",
					ContextLines: []string{
						"ccc content",
					},
				},
			},
			err: "",
		},

		{
			desc: "strips context when disabled via config",
			config: &OutputConfig{
				ShowContext: false,
			},
			results: []*Result{
				{
					Location: ResultLocation{
						Path:      "testdata/txt/aaa.txt",
						StartLine: 1,
						EndLine:   1,
					},
					ContextLang: "plaintext",
					ContextLines: []string{
						"aaa content",
					},
				},
				{
					Location: ResultLocation{
						Path:      "testdata/txt/bbb.txt",
						StartLine: 1,
						EndLine:   1,
					},
					ContextLang: "plaintext",
					ContextLines: []string{
						"bbb content",
					},
				},
				{
					Location: ResultLocation{
						Path:      "testdata/txt/ccc.txt",
						StartLine: 1,
						EndLine:   1,
					},
					ContextLang: "plaintext",
					ContextLines: []string{
						"ccc content",
					},
				},
			},
			expected: []*Result{
				{
					Location: ResultLocation{
						Path:      "testdata/txt/aaa.txt",
						StartLine: 1,
						EndLine:   1,
					},
					ContextLang:  "",
					ContextLines: nil,
				},
				{
					Location: ResultLocation{
						Path:      "testdata/txt/bbb.txt",
						StartLine: 1,
						EndLine:   1,
					},
					ContextLang:  "",
					ContextLines: nil,
				},
				{
					Location: ResultLocation{
						Path:      "testdata/txt/ccc.txt",
						StartLine: 1,
						EndLine:   1,
					},
					ContextLang:  "",
					ContextLines: nil,
				},
			},
			err: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			app := NewTestApp()
			if tt.config != nil {
				app.Config.Output = *tt.config
			}
			ctx := app.InitContext(context.Background())

			actual, err := EnsureContextLines(ctx, tt.results)

			if tt.err == "" {
				require.NoError(t, err)
			} else {
				require.ErrorContains(t, err, tt.err)
			}

			require.Equal(t, tt.expected, actual)
		})
	}
}

func TestSortResults(t *testing.T) {
	tests := []struct {
		desc     string
		config   *OutputConfig
		results  []*Result
		expected []*Result
		err      string
	}{
		{
			desc: "sorts by location by default",
			results: []*Result{
				{
					Location: ResultLocation{
						Path: "testdata/txt/ccc.txt",
					},
				},
				{
					Location: ResultLocation{
						Path: "testdata/txt/aaa.txt",
					},
				},
				{
					Location: ResultLocation{
						Path: "testdata/txt/bbb.txt",
					},
				},
			},
			expected: []*Result{
				{
					Location: ResultLocation{
						Path: "testdata/txt/aaa.txt",
					},
				},
				{
					Location: ResultLocation{
						Path: "testdata/txt/bbb.txt",
					},
				},
				{
					Location: ResultLocation{
						Path: "testdata/txt/ccc.txt",
					},
				},
			},
			err: "",
		},

		{
			desc: "sorts by severity",
			config: &OutputConfig{
				Sort: ResultSortSeverity,
			},
			results: []*Result{
				{
					Level: ResultLevelWarning,
				},
				{
					Level: ResultLevelError,
				},
				{
					Level: ResultLevelInfo,
				},
			},
			expected: []*Result{
				{
					Level: ResultLevelError,
				},
				{
					Level: ResultLevelWarning,
				},
				{
					Level: ResultLevelInfo,
				},
			},
			err: "",
		},

		{
			desc: "sorts by source",
			config: &OutputConfig{
				Sort: ResultSortSource,
			},
			results: []*Result{
				{
					Source: "ccc",
				},
				{
					Source: "aaa",
				},
				{
					Source: "bbb",
				},
			},
			expected: []*Result{
				{
					Source: "aaa",
				},
				{
					Source: "bbb",
				},
				{
					Source: "ccc",
				},
			},
			err: "",
		},

		{
			desc: "returns an error when unknown source",
			config: &OutputConfig{
				Sort: ResultSort("nope"),
			},
			results:  []*Result{},
			expected: nil,
			err:      "unknown sort type: nope",
		},
	}
	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			app := NewTestApp()
			if tt.config != nil {
				app.Config.Output = *tt.config
			}
			ctx := app.InitContext(context.Background())

			actual, err := SortResults(ctx, tt.results)

			if tt.err == "" {
				require.NoError(t, err)
			} else {
				require.ErrorContains(t, err, tt.err)
			}

			require.Equal(t, tt.expected, actual)
		})
	}
}
