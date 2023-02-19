package stylist

import (
	"context"
	"runtime"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/twelvelabs/termite/run"
)

func TestCommand_Execute(t *testing.T) {
	tests := []struct {
		desc     string
		command  *Command
		paths    []string
		setup    func(c *run.Client)
		expected []*Result
		err      string
	}{
		{
			desc: "empty path set is a noop",
			command: &Command{
				Template:     "test-linter --verbose",
				InputType:    InputTypeArg,
				OutputFormat: OutputFormatNone,
			},
			paths:    []string{},
			expected: []*Result{},
			err:      "",
		},

		{
			desc: "empty command returns an error",
			command: &Command{
				Template: "",
			},
			paths: []string{
				"testdata/txt/aaa.txt",
				"testdata/txt/bbb.txt",
			},
			expected: []*Result{},
			err:      ErrCommandEmpty.Error(),
		},
		{
			desc: "malformed command returns an error",
			command: &Command{
				Template: `"`,
			},
			paths: []string{
				"testdata/txt/aaa.txt",
				"testdata/txt/bbb.txt",
			},
			expected: []*Result{},
			err:      "EOF found when expecting closing quote",
		},

		{
			desc: "[arg] runs command once per path",
			command: &Command{
				Template:     "test-linter --verbose",
				InputType:    InputTypeArg,
				OutputFormat: OutputFormatNone,
			},
			paths: []string{
				"testdata/txt/aaa.txt",
				"testdata/txt/bbb.txt",
			},
			setup: func(c *run.Client) {
				c.RegisterStub(
					run.MatchString("test-linter --verbose testdata/txt/aaa.txt"),
					run.StringResponse(""),
				)
				c.RegisterStub(
					run.MatchString("test-linter --verbose testdata/txt/bbb.txt"),
					run.StringResponse(""),
				)
			},
			expected: []*Result{},
			err:      "",
		},

		{
			desc: "[none] runs command once without paths",
			command: &Command{
				Template:     "test-linter --verbose",
				InputType:    InputTypeNone,
				OutputFormat: OutputFormatNone,
			},
			paths: []string{
				"testdata/txt/aaa.txt",
				"testdata/txt/bbb.txt",
			},
			setup: func(c *run.Client) {
				c.RegisterStub(
					run.MatchString("test-linter --verbose"),
					run.StringResponse(""),
				)
			},
			expected: []*Result{},
			err:      "",
		},

		{
			desc: "[stdin] runs command once per path with content passed to stdin",
			command: &Command{
				Template:     "test-linter --verbose",
				InputType:    InputTypeStdin,
				OutputFormat: OutputFormatNone,
			},
			paths: []string{
				"testdata/txt/aaa.txt",
				"testdata/txt/bbb.txt",
			},
			setup: func(c *run.Client) {
				c.RegisterStub(
					run.MatchAll(
						run.MatchString("test-linter --verbose"),
						run.MatchStdin("aaa content\n"),
					),
					run.StringResponse(""),
				)
				c.RegisterStub(
					run.MatchAll(
						run.MatchString("test-linter --verbose"),
						run.MatchStdin("bbb content\n"),
					),
					run.StringResponse(""),
				)
			},
			expected: []*Result{},
			err:      "",
		},

		{
			desc: "[variadic] runs command once per batch of paths",
			command: &Command{
				Template:     "test-linter --verbose",
				InputType:    InputTypeVariadic,
				OutputFormat: OutputFormatNone,
			},
			paths: []string{
				"testdata/txt/aaa.txt",
				"testdata/txt/bbb.txt",
			},
			setup: func(c *run.Client) {
				c.RegisterStub(
					run.MatchString("test-linter --verbose testdata/txt/aaa.txt testdata/txt/bbb.txt"),
					run.StringResponse(""),
				)
			},
			expected: []*Result{},
			err:      "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			app := NewTestApp()
			defer app.CmdClient.VerifyStubs(t)

			if tt.setup != nil {
				tt.setup(app.CmdClient)
			}

			ctx := app.InitContext(context.Background())
			actual, err := tt.command.Execute(ctx, "test-linter", tt.paths)

			if tt.err == "" {
				assert.NoError(t, err)
			} else {
				assert.ErrorContains(t, err, tt.err)
			}
			assert.Equal(t, tt.expected, actual)
		})
	}
}

func TestCommand_executeBatch_ErrorCases(t *testing.T) {
	command := &Command{}

	results, err := command.executeBatch(context.Background(), "", []string{})
	assert.Equal(t, []*Result(nil), results)
	assert.NoError(t, err)
}

func TestCommand_parallelism(t *testing.T) {
	// Defaults to total CPU cores
	command := &Command{}
	assert.Equal(t, runtime.NumCPU(), command.parallelism())

	// But can be explicitly set
	command = &Command{
		Parallelism: 8675309,
	}
	assert.Equal(t, 8675309, command.parallelism())
}

func TestCommand_partition(t *testing.T) {
	tests := []struct {
		desc     string
		command  *Command
		paths    []string
		expected [][]string
	}{
		{
			desc: "[arg] partitions into single-item batches",
			command: &Command{
				InputType: InputTypeArg,
			},
			paths: []string{
				"aaa.txt",
				"bbb.txt",
				"ccc.txt",
			},
			expected: [][]string{
				{"aaa.txt"},
				{"bbb.txt"},
				{"ccc.txt"},
			},
		},

		{
			desc: "[stdin] partitions into single-item batches",
			command: &Command{
				InputType: InputTypeStdin,
			},
			paths: []string{
				"aaa.txt",
				"bbb.txt",
				"ccc.txt",
			},
			expected: [][]string{
				{"aaa.txt"},
				{"bbb.txt"},
				{"ccc.txt"},
			},
		},

		{
			desc: "[variadic] partitions into batches",
			command: &Command{
				InputType: InputTypeVariadic,
				BatchSize: 3,
			},
			paths: []string{
				"aaa.txt",
				"bbb.txt",
				"ccc.txt",

				"ddd.txt",
				"eee.txt",
				"fff.txt",

				"ggg.txt",
				"hhh.txt",
				"iii.txt",

				"jjj.txt",
			},
			expected: [][]string{
				{"aaa.txt", "bbb.txt", "ccc.txt"},
				{"ddd.txt", "eee.txt", "fff.txt"},
				{"ggg.txt", "hhh.txt", "iii.txt"},
				{"jjj.txt"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			actual := tt.command.partition(tt.paths)
			assert.Equal(t, tt.expected, actual)
		})
	}
}
