package stylist

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewPresetStore(t *testing.T) {
	prevPresetsPath := presetsPath
	t.Cleanup(func() {
		presetsPath = prevPresetsPath
	})

	_, err := NewPresetStore()
	assert.NoError(t, err)

	presetsPath = "missing.yml"
	_, err = NewPresetStore()
	assert.ErrorContains(t, err, "file does not exist")
}

func TestNewPresetStoreFromReader(t *testing.T) {
	_, err := NewPresetStoreFromReader(strings.NewReader("---\n{}"))
	assert.NoError(t, err)

	_, err = NewPresetStoreFromReader(strings.NewReader("[]"))
	assert.ErrorContains(t, err, "cannot unmarshal")
}

func TestPresetStore_All(t *testing.T) {
	yaml := `---
bbb:
  name: bbb
ccc:
  name: ccc
aaa:
  name: aaa
`
	store, err := NewPresetStoreFromReader(strings.NewReader(yaml))
	require.NoError(t, err)

	presets := store.All()
	assert.Equal(t, 3, len(presets))
	assert.Equal(t, "aaa", presets[0].Name)
	assert.Equal(t, "bbb", presets[1].Name)
	assert.Equal(t, "ccc", presets[2].Name)

	var preset *Processor

	preset, err = store.Get("aaa")
	require.NoError(t, err)
	require.Equal(t, presets[0], preset)

	preset, err = store.Get("bbb")
	require.NoError(t, err)
	require.Equal(t, presets[1], preset)

	preset, err = store.Get("ccc")
	require.NoError(t, err)
	require.Equal(t, presets[2], preset)

	preset, err = store.Get("unknown")
	require.ErrorContains(t, err, "unknown preset")
	assert.Nil(t, preset)
}

func TestResolvePresets(t *testing.T) {
	processors := []*Processor{
		{
			Name: "foo",
		},
		{
			Name:   "markdown",
			Preset: "markdownlint",
			Includes: []string{
				"foo.md",
				"bar.md",
			},
		},
	}

	processors, err := ResolvePresets(processors)
	require.NoError(t, err)
	require.Equal(t, 2, len(processors))

	// Non-preset using processor should be untouched.
	require.Equal(t, &Processor{Name: "foo"}, processors[0])

	require.Equal(t, "markdown", processors[1].Name)
	require.Equal(t, "markdownlint", processors[1].Preset)
	require.Equal(t, []string{
		"foo.md",
		"bar.md",
	}, processors[1].Includes)
	require.Equal(t, "markdownlint --json", processors[1].CheckCommand.Template)
}

func TestResolvePresets_WhenPresetStoreError(t *testing.T) {
	prevPresetsPath := presetsPath
	t.Cleanup(func() {
		presetsPath = prevPresetsPath
	})

	presetsPath = "missing.yml"
	_, err := ResolvePresets([]*Processor{})
	assert.ErrorContains(t, err, "file does not exist")
}

func TestResolvePresets_WhenUnknownPreset(t *testing.T) {
	processors := []*Processor{
		{Preset: "unknown"},
	}
	processors, err := ResolvePresets(processors)
	require.ErrorContains(t, err, "unknown preset")
	require.Equal(t, 0, len(processors))
}
