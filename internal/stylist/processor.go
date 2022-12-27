package stylist

import (
	mapset "github.com/deckarep/golang-set/v2"
)

type Processor struct {
	Name         string   `yaml:"name"`
	Tags         []string `yaml:"tags"`
	Types        []string `yaml:"types"`
	Includes     []string `yaml:"includes"`
	Excludes     []string `yaml:"excludes"`
	CheckCommand *Command `yaml:"check"`
	FixCommand   *Command `yaml:"fix"`

	paths []string
}

// Paths returns all paths matched by the processor.
func (p *Processor) Paths() []string {
	return p.paths
}

type ProcessorList []*Processor

func (pl ProcessorList) All() []*Processor {
	return pl
}

func (pl ProcessorList) Named(names ...string) []*Processor {
	found := []*Processor{}
	for _, p := range pl {
		for _, name := range names {
			if p.Name == name {
				found = append(found, p)
			}
		}
	}
	return found
}

func (pl ProcessorList) Tagged(tags ...string) []*Processor {
	found := []*Processor{}
	for _, p := range pl {
		if len(p.Tags) == 0 {
			continue
		}
		if mapset.NewSet(p.Tags...).Contains(tags...) {
			found = append(found, p)
		}
	}
	return found
}
