package stylist

import (
	"context"
	"runtime"
	"sort"
	"time"

	"github.com/bmatcuk/doublestar/v4"
	"golang.org/x/sync/errgroup"
)

func NewPipeline(processors []*Processor, excludes []string) *Pipeline {
	// Always ignore git dirs.
	excludes = append(excludes, ".git/**")
	return &Pipeline{
		processors: processors,
		excludes:   excludes,
	}
}

type Pipeline struct {
	processors []*Processor
	excludes   []string
}

// Index populates the paths for each processor in the pipeline.
//
// The source paths are resolved from each path spec; matched against
// any global exclude patterns; then matched against each processor's
// individual type, include, and exclude patterns.
func (p *Pipeline) Index(ctx context.Context, pathSpecs []string) error {
	logger := AppLogger(ctx)

	// Aggregate each processor's include patterns
	includes := []string{}
	for _, processor := range p.processors {
		includes = append(includes, processor.Includes...)
	}

	startedAt := time.Now()
	logger.Debugf(
		"Indexing: includes=%v excludes=%v",
		includes,
		p.excludes,
	)

	// TODO: support passing Context to the indexer
	// Create an index of paths (resolved from the path specs),
	// matching any of the include patterns used by our processors.
	// Doing this once is _much_ faster than once per-processor,
	// especially when dealing w/ very large projects and many processors or patterns.
	indexer := NewPathIndexer(includes, p.excludes)
	if err := indexer.Index(pathSpecs...); err != nil {
		return err
	}

	// For each processor...
	for _, processor := range p.processors {
		// Gather all paths matching the include patterns
		// configured for this processor.
		pathSet := NewPathSet()
		for _, inc := range processor.Includes {
			if incPaths, ok := indexer.PathsByInclude[inc]; ok {
				pathSet = pathSet.Union(incPaths)
			}
		}

		// Now, filter out anything _this processor_ is configured to ignore.
		// (vs. the global excludes we passed in to the indexer).
		paths := []string{}
		for path := range pathSet.Iter() {
			excluded := false
			for _, pattern := range processor.Excludes {
				ok, err := doublestar.PathMatch(pattern, path)
				if err != nil {
					return err
				}
				if ok {
					excluded = true
				}
			}
			if !excluded {
				paths = append(paths, path)
			}
		}

		sort.Strings(paths)
		processor.paths = paths
	}

	logger.Debugf("Indexed in %s", time.Since(startedAt))
	return nil
}

// Match returns all processors that match the given path specs.
func (p *Pipeline) Match(ctx context.Context, pathSpecs []string) ([]*Processor, error) {
	if err := p.Index(ctx, pathSpecs); err != nil {
		return nil, err
	}
	matches := []*Processor{}
	for _, processor := range p.processors {
		if len(processor.Paths()) > 0 {
			matches = append(matches, processor)
		}
	}
	return matches, nil
}

// Check executes the check command for each processor in the pipeline.
func (p *Pipeline) Check(ctx context.Context, pathSpecs []string) ([]*Result, error) {
	return p.execute(ctx, pathSpecs, CommandTypeCheck)
}

// Check executes the fix command for each processor in the pipeline.
func (p *Pipeline) Fix(ctx context.Context, pathSpecs []string) ([]*Result, error) {
	return p.execute(ctx, pathSpecs, CommandTypeFix)
}

func (p *Pipeline) execute(
	ctx context.Context, pathSpecs []string, ct CommandType,
) ([]*Result, error) {
	// Get all the processors that match the pathSpecs.
	processors, err := p.Match(ctx, pathSpecs)
	if err != nil {
		return nil, err
	}

	// Setup an errgroup w/ the correct level of parallelism.
	// Fix commands mutate files, so each processor needs to run serially.
	group, ctx := errgroup.WithContext(ctx)
	if ct == CommandTypeFix {
		group.SetLimit(1)
	} else {
		// TODO: once we have a good test case (lots of processors and files),
		// check to see whether this can be safely removed.
		// Might run better un-throttled.
		group.SetLimit(runtime.NumCPU())
	}

	// Execute the processors in goroutines and aggregate their results.
	results := []*Result{}
	for _, processor := range processors {
		processor := processor
		group.Go(func() error {
			pr, err := processor.Execute(ctx, ct)
			if err != nil {
				return err
			}
			// TODO: add mutex
			results = append(results, pr...)
			return nil
		})
	}

	err = group.Wait()
	if err != nil {
		return nil, err
	}

	// Return the transformed results.
	return p.transform(results)
}

func (p *Pipeline) transform(results []*Result) ([]*Result, error) {
	// Maybe this should be controlled by a flag (and done by the caller)?
	sort.Slice(results, func(i, j int) bool {
		if results[i].Source != results[j].Source {
			return results[i].Source < results[j].Source
		}
		if results[i].Location.Path != results[j].Location.Path {
			return results[i].Location.Path < results[j].Location.Path
		}
		if results[i].Location.StartLine != results[j].Location.StartLine {
			return results[i].Location.StartLine < results[j].Location.StartLine
		}
		return results[i].Location.StartColumn < results[j].Location.StartColumn
	})

	loader := NewContextLineLoader()
	for _, result := range results {
		if result.ContextLines == nil {
			lines, err := loader.Load(result.Location)
			if err != nil {
				return nil, err
			}
			result.ContextLines = lines
		}
	}

	return results, nil
}
