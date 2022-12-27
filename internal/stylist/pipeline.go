package stylist

import (
	"context"
	"sort"

	"github.com/bmatcuk/doublestar/v4"
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
	// Aggregate each processor's file types and include patterns
	fileTypes := []string{}
	includes := []string{}
	for _, processor := range p.processors {
		fileTypes = append(fileTypes, processor.Types...)
		includes = append(includes, processor.Includes...)
	}

	// TODO: support passing Context to the indexer
	// Create an index of paths (resolved from the path specs),
	// matching any of the types and include patterns used by our processors.
	// Doing this once is _much_ faster than once per-processor,
	// especially when dealing w/ very large projects and many processors or patterns.
	indexer := NewPathIndexer(fileTypes, includes, p.excludes)
	if err := indexer.Index(pathSpecs...); err != nil {
		return err
	}

	// For each processor...
	for _, processor := range p.processors {
		// Gather all paths matching the file types and include patterns
		// configured for this processor.
		pathSet := NewPathSet()
		for _, ft := range processor.Types {
			if ftPaths, ok := indexer.PathsByFileType[ft]; ok {
				pathSet = pathSet.Union(ftPaths)
			}
		}
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

	return nil
}

// Check executes the check command for each processor in the pipeline.
func (p *Pipeline) Check(ctx context.Context, pathSpecs []string) ([]*Result, error) {
	if err := p.Index(ctx, pathSpecs); err != nil {
		return nil, err
	}
	results := []*Result{}
	for _, processor := range p.processors {
		if processor.CheckCommand == nil {
			continue
		}
		pr, err := processor.CheckCommand.Execute(ctx, processor.Paths())
		if err != nil {
			return nil, err
		}
		results = append(results, pr...)
	}
	return results, nil
}

func (p *Pipeline) Fix(ctx context.Context, pathSpecs []string) ([]*Result, error) {
	return nil, nil
}
