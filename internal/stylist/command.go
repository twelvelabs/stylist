package stylist

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"runtime"

	mapset "github.com/deckarep/golang-set/v2"
	"github.com/google/shlex"
	"golang.org/x/sync/errgroup"
)

var (
	ErrCommandEmpty = errors.New("empty command")
)

// Command represents a check or fix command to be run by a Processor.
type Command struct {
	Template      string        `yaml:"template"`
	InputType     InputType     `yaml:"input"    default:"variadic"`
	OutputType    OutputType    `yaml:"output"   default:"stdout"`
	OutputFormat  OutputFormat  `yaml:"format"   default:"none"`
	ResultMapping ResultMapping `yaml:"mapping"`
	Parallelism   int           `yaml:"parallelism"`
	BatchSize     int           `yaml:"batch_size"`
}

// Execute executes paths concurrently in batches of 10.
func (c *Command) Execute(ctx context.Context, paths []string) ([]*Result, error) {
	results := []*Result{}

	group, ctx := errgroup.WithContext(ctx)
	group.SetLimit(c.parallelism())

	for _, batch := range c.partition(paths) {
		batch := batch
		group.Go(func() error {
			batchResults, err := c.executeBatch(ctx, batch)
			if err != nil {
				return err
			}
			// TODO: wrap access to this in a mutex
			results = append(results, batchResults...)
			return nil
		})
	}

	err := group.Wait()
	return results, err
}

// executes a single batch of paths.
func (c *Command) executeBatch(ctx context.Context, paths []string) ([]*Result, error) {
	if len(paths) == 0 {
		return nil, nil
	}

	logger := AppLogger(ctx)
	client := AppCmdClient(ctx)

	// TODO: render template string
	args, err := shlex.Split(c.Template)
	if err != nil {
		return nil, err
	}
	if len(args) == 0 {
		return nil, ErrCommandEmpty
	}

	if c.InputType == InputTypeArg {
		args = append(args, paths[0])
	}
	if c.InputType == InputTypeVariadic {
		args = append(args, paths...)
	}

	cmd := client.CommandContext(ctx, args[0], args[1:]...)

	// Setup the IO streams
	if c.InputType == InputTypeStdin {
		file, err := os.Open(paths[0])
		if err != nil {
			return nil, err
		}
		cmd.Stdin = file
	}
	stderr := &bytes.Buffer{}
	stdout := &bytes.Buffer{}
	cmd.Stderr = stderr
	cmd.Stdout = stdout

	logger.Debugln("Command:", cmd.String())

	err = cmd.Run()

	// Ignoring ExitError so we can parse the output.
	var exitErr *exec.ExitError
	if err != nil && !errors.As(err, &exitErr) {
		// non-ExitError (binary not found, permissions error, etc).
		return nil, err
	}

	// Build a CommandOutput struct...
	content := stdout
	if c.OutputType == OutputTypeStderr {
		content = stderr
	}
	output := CommandOutput{
		Content:  content,
		ExitCode: cmd.ProcessState.ExitCode(),
	}

	logger.Debugln("Output:", output.String())

	// Parse the output using the appropriate parser.
	parsed, err := NewOutputParser(c.OutputFormat).Parse(output, c.ResultMapping)
	if err != nil {
		return nil, err
	}

	// Do a little post processing on the results.
	pathSet := mapset.NewSet(paths...)
	transformed := []*Result{}
	for _, r := range parsed {
		// Add the command name to the results
		r.Source = args[0]
		// InputTypeNone doesn't pass `paths` to the command, so there may
		// be results for paths we don't care about. Filter those out.
		if pathSet.Contains(r.Location.Path) || r.Location.Path == "" {
			transformed = append(transformed, r)
		}
	}
	return transformed, nil
}

func (c *Command) parallelism() int {
	if c.Parallelism == 0 {
		return runtime.NumCPU()
	}
	return c.Parallelism
}

// Partitions paths into batches of 10.
func (c *Command) partition(paths []string) [][]string {
	if c.InputType == InputTypeNone {
		// InputTypeNone is used when a processor doesn't accept (or want)
		// individual file paths. Return a single batch containing everything.
		return [][]string{paths}
	}
	if c.InputType != InputTypeVariadic {
		// For non-variadic input we just return a slice of single-path batches.
		// This allows us to have a single code path, but at the expense
		// of a bunch of useless allocations.
		// Will need to profile to see whether that tradeoff is acceptable.
		batches := [][]string{}
		for _, path := range paths {
			batches = append(batches, []string{path})
		}
		return batches
	}

	batchSize := 10
	if c.BatchSize != 0 {
		batchSize = c.BatchSize
	}

	if len(paths) <= batchSize {
		return [][]string{paths}
	}

	batches := make([][]string, 0, (len(paths)+batchSize-1)/batchSize)
	for batchSize < len(paths) {
		paths, batches = paths[batchSize:], append(batches, paths[0:batchSize:batchSize])
	}
	if len(paths) > 0 {
		batches = append(batches, paths)
	}

	return batches
}

// CommandOutput contains the result of a single command invocation.
type CommandOutput struct {
	Content  io.Reader
	ExitCode int
}

func (o *CommandOutput) String() string {
	// io hijinks to reset the read offset
	buf := &bytes.Buffer{}
	reader := io.TeeReader(o.Content, buf)
	o.Content = buf

	content, _ := io.ReadAll(reader)
	return fmt.Sprintf(
		`<CommandOutput Content="%v" ExitCode="%v">`,
		string(content),
		o.ExitCode,
	)
}
