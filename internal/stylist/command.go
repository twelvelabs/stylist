package stylist

import (
	"context"
	"errors"
	"os"
	"runtime"

	"github.com/google/shlex"
	"github.com/twelvelabs/termite/run"
	"golang.org/x/sync/errgroup"
)

var (
	ErrCommandEmpty = errors.New("empty command")
)

// NewCommand returns a new Command.
func NewCommand(client *run.Client) *Command {
	return &Command{
		client: client,
	}
}

// Command represents a check or fix command to be run by a Processor.
type Command struct {
	Template    string     `yaml:"template"`
	Input       InputType  `yaml:"input"`
	Output      OutputType `yaml:"output"`
	Parallelism int        `yaml:"parallelism"`
	BatchSize   int        `yaml:"batch_size"`

	client *run.Client
}

// Execute executes paths concurrently in batches of 10.
func (c *Command) Execute(ctx context.Context, paths []string) ([]*Diagnostic, error) {
	diagnostics := []*Diagnostic{}

	group, ctx := errgroup.WithContext(ctx)
	group.SetLimit(c.parallelism())

	for _, batch := range c.partition(paths) {
		batch := batch
		group.Go(func() error {
			results, err := c.executeBatch(ctx, batch)
			if err != nil {
				return err
			}
			// TODO: wrap access to this in a mutex
			diagnostics = append(diagnostics, results...)
			return nil
		})
	}

	err := group.Wait()
	return diagnostics, err
}

// executes a single batch of paths.
func (c *Command) executeBatch(ctx context.Context, paths []string) ([]*Diagnostic, error) {
	if len(paths) == 0 {
		return nil, nil
	}

	// TODO: render template string
	args, err := shlex.Split(c.Template)
	if err != nil {
		return nil, err
	}
	if len(args) == 0 {
		return nil, ErrCommandEmpty
	}

	if c.Input == InputTypeArg {
		args = append(args, paths[0])
	}
	if c.Input == InputTypeVariadic {
		args = append(args, paths...)
	}

	cmd := c.client.CommandContext(ctx, args[0], args[1:]...)

	if c.Input == InputTypeStdin {
		file, err := os.Open(paths[0])
		if err != nil {
			return nil, err
		}
		cmd.Stdin = file
	}

	err = cmd.Run()
	if err != nil {
		return nil, err
	}

	// TODO: parse output...
	return []*Diagnostic{}, nil
}

func (c *Command) parallelism() int {
	if c.Parallelism == 0 {
		return runtime.NumCPU()
	}
	return c.Parallelism
}

// Partitions paths into batches of 10.
func (c *Command) partition(paths []string) [][]string {
	if c.Input != InputTypeVariadic {
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
