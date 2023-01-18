package stylist

import (
	"embed"
	"fmt"
	"io"
	"sort"

	"gopkg.in/yaml.v3"
)

//go:embed presets.yml
var presetsFS embed.FS
var presetsPath = "presets.yml"

// ResolvePresets searches the given slice for any processor extending a preset
// and reverse merges that preset into processor.
func ResolvePresets(processors []*Processor) ([]*Processor, error) {
	store, err := NewPresetStore()
	if err != nil {
		return nil, err
	}
	for idx, p := range processors {
		if p.Preset != "" {
			preset, err := store.Get(p.Preset)
			if err != nil {
				return nil, err
			}
			processors[idx] = preset.Merge(p)
		}
	}
	return processors, nil
}

// PresetStore manages preset processor configurations.
type PresetStore struct {
	presets map[string]*Processor
}

// NewPresetStore returns a new preset store.
func NewPresetStore() (*PresetStore, error) {
	f, err := presetsFS.Open(presetsPath)
	if err != nil {
		return nil, err
	}
	return NewPresetStoreFromReader(f)
}

// NewPresetStore returns a new preset store using data provided by the reader.
func NewPresetStoreFromReader(reader io.Reader) (*PresetStore, error) {
	var presets map[string]*Processor

	err := yaml.NewDecoder(reader).Decode(&presets)
	if err != nil {
		return nil, fmt.Errorf("preset store decode: %w", err)
	}

	return &PresetStore{presets: presets}, nil
}

// All returns all presets sorted by name.
func (s *PresetStore) All() []*Processor {
	processors := []*Processor{}
	for _, v := range s.presets {
		processors = append(processors, v)
	}
	sort.Slice(processors, func(i, j int) bool {
		return processors[i].Name < processors[j].Name
	})
	return processors
}

// Candidates returns all presets relevant to the given dir.
// A preset is relevant if there are files in the directory that match
// the preset's include patterns (and do not match any exclude pattern).
func (s *PresetStore) Candidates(dir string) []*Processor {
	// TODO: implement
	return s.All()
}

// Get returns the named processor.
func (s *PresetStore) Get(name string) (*Processor, error) {
	if preset, ok := s.presets[name]; ok {
		return preset, nil
	}
	return nil, fmt.Errorf("unknown preset: %s", name)
}
