package stylist

import (
	"fmt"
	"strings"

	mapset "github.com/deckarep/golang-set/v2"
	"github.com/imdario/mergo"
)

type Processor struct {
	Name         string   `yaml:"name,omitempty"`
	Preset       string   `yaml:"preset,omitempty"`
	Tags         []string `yaml:"tags,omitempty"`
	Includes     []string `yaml:"includes,omitempty"`
	Excludes     []string `yaml:"excludes,omitempty"`
	CheckCommand *Command `yaml:"check,omitempty"`
	FixCommand   *Command `yaml:"fix,omitempty"`

	paths []string
}

// Merge merges the receiver and arguments and returns a new processor.
// Only exported fields are merged.
func (p *Processor) Merge(others ...*Processor) *Processor {
	dst := &Processor{}
	_ = mergo.Merge(dst, p)
	for _, other := range others {
		_ = mergo.Merge(dst, other, mergo.WithOverride)
	}
	return dst
}

// Paths returns all paths matched by the processor.
func (p *Processor) Paths() []string {
	return p.paths
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
