package stylist

type Processor struct {
	Name         string   `yaml:"name"`
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
