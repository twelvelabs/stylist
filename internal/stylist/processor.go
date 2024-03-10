package stylist

import (
	"context"
	"fmt"
	"strings"

	mapset "github.com/deckarep/golang-set/v2"
	"github.com/imdario/mergo"
)

type Processor struct {
	Preset       string   `yaml:"preset,omitempty"`
	Name         string   `yaml:"name,omitempty"`
	Tags         []string `yaml:"tags,omitempty"`
	Includes     []string `yaml:"includes,omitempty"`
	Excludes     []string `yaml:"excludes,omitempty"`
	CheckCommand *Command `yaml:"check,omitempty"`
	FixCommand   *Command `yaml:"fix,omitempty"`
}

// Execute runs the given command for paths.
func (p *Processor) Execute(
	ctx context.Context, ct CommandType, paths []string,
) ([]*Result, error) {
	// Resolve the command to execute.
	var cmd *Command
	switch ct {
	case CommandTypeCheck:
		cmd = p.CheckCommand
	case CommandTypeFix:
		cmd = p.FixCommand
	}

	if cmd == nil {
		// Command not implemented - nothing to do.
		return nil, nil
	}

	// Delegate to the command.
	return cmd.Execute(ctx, p.Name, paths)
}

// Merge merges the receiver and arguments and returns a new processor
// Only exported fields are merged.
func (p *Processor) Merge(others ...*Processor) *Processor {
	dst := &Processor{}
	_ = mergo.Merge(dst, p)
	for _, other := range others {
		_ = mergo.Merge(dst, other, mergo.WithOverride)
	}
	return dst
}

// ProcessorFilter filters processors by name and/or tag.
type ProcessorFilter struct {
	Names []string
	Tags  []string
}

func (pf *ProcessorFilter) Cardinality() int {
	return len(pf.Names) + len(pf.Tags)
}

// Filter returns all processors matching the current name and tag filters,
// or an error if no processors were found.
func (pf *ProcessorFilter) Filter(processors []*Processor) ([]*Processor, error) {
	// Ensure a valid processor list
	err := pf.validate(processors)
	if err != nil {
		return nil, err
	}

	// If no filter criteria defined, then just return the unfiltered input.
	if pf.Cardinality() == 0 {
		return processors, nil
	}

	byName, byTag := pf.index(processors)
	found := mapset.NewSet[*Processor]()

	for _, name := range pf.Names {
		if p, ok := byName[name]; ok {
			found.Add(p)
		} else {
			return nil, fmt.Errorf("no processor named %s", name)
		}
	}
	for _, tag := range pf.Tags {
		if pSlice, ok := byTag[tag]; ok {
			for _, p := range pSlice {
				found.Add(p)
			}
		} else {
			return nil, fmt.Errorf("no processor tagged %s", tag)
		}
	}

	return found.ToSlice(), nil
}

func (pf *ProcessorFilter) index(processors []*Processor) (
	map[string]*Processor,
	map[string][]*Processor,
) {
	byName := map[string]*Processor{}
	byTag := map[string][]*Processor{}

	for _, p := range processors {
		byName[p.Name] = p
		for _, tag := range p.Tags {
			if _, ok := byTag[tag]; !ok {
				byTag[tag] = []*Processor{}
			}
			byTag[tag] = append(byTag[tag], p)
		}
	}

	return byName, byTag
}

func (pf *ProcessorFilter) validate(processors []*Processor) error {
	if len(processors) == 0 {
		return fmt.Errorf("no processors defined")
	}

	names := mapset.NewSet[string]()
	for idx, p := range processors {
		name := strings.TrimSpace(p.Name)
		if len(name) == 0 {
			return fmt.Errorf("processor at index %v is unnamed", idx)
		}
		if names.Contains(name) {
			return fmt.Errorf("processor at index %v has a duplicate name", idx)
		}
		names.Add(name)
	}

	return nil
}
