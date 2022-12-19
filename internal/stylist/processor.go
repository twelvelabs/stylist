package stylist

import (
	"context"
	"fmt"
	"sort"

	"github.com/bmatcuk/doublestar/v4"
)

type Processor struct {
	Name         string   `yaml:"name"`
	Types        []string `yaml:"types"`
	Includes     []string `yaml:"includes"`
	Excludes     []string `yaml:"excludes"`
	CheckCommand *Command `yaml:"check"`
	FixCommand   *Command `yaml:"fix"`

	paths []string
}

func (p *Processor) Check(ctx context.Context) ([]*Diagnostic, error) {
	if p.CheckCommand == nil {
		fmt.Println("p.CheckCommand == nil")
		return nil, nil
	}
	return p.CheckCommand.Execute(ctx, p.Paths())
}

func (p *Processor) Fix(ctx context.Context) ([]*Diagnostic, error) {
	if p.FixCommand == nil {
		fmt.Println("p.FixCommand == nil")
		return nil, nil
	}
	return p.FixCommand.Execute(ctx, p.Paths())
}

// Paths returns all paths matched by the processor.
func (p *Processor) Paths() []string {
	return p.paths
}

type ProcessorList []*Processor

// Index searches for the paths each processor should handle,
// resolving each path spec and ignoring anything matching the global exclude patterns.
func (pl ProcessorList) Index(pathSpecs []string, globalExcludes []string) error {
	// Always ignore git dirs.
	// Probably should do node_modules too, but waiting until I hear feedback from others.
	globalExcludes = append(globalExcludes, ".git/**")

	// Aggregate each processor's file types and include patterns
	fileTypes := []string{}
	includes := []string{}
	for _, p := range pl {
		fileTypes = append(fileTypes, p.Types...)
		includes = append(includes, p.Includes...)
	}

	// Create an index of paths (resolved from the path specs),
	// matching any of the types and include patterns used by our processors.
	// Doing this once is _much_ faster than once per-processor,
	// especially when dealing w/ very large projects and many processors or patterns.
	indexer := NewPathIndexer(fileTypes, includes, globalExcludes)
	if err := indexer.Index(pathSpecs...); err != nil {
		return err
	}

	// For each processor...
	for _, p := range pl {
		// Gather all paths matching the file types and include patterns
		// configured for this processor.
		pathSet := NewPathSet()
		for _, ft := range p.Types {
			if ftPaths, ok := indexer.PathsByFileType[ft]; ok {
				pathSet = pathSet.Union(ftPaths)
			}
		}
		for _, inc := range p.Includes {
			if incPaths, ok := indexer.PathsByInclude[inc]; ok {
				pathSet = pathSet.Union(incPaths)
			}
		}

		// Now, filter out anything _this processor_ is configured to ignore.
		// (vs. the global excludes we passed in to the indexer).
		paths := []string{}
		for path := range pathSet.Iter() {
			excluded := false
			for _, pattern := range p.Excludes {
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
		p.paths = paths
	}

	return nil
}
