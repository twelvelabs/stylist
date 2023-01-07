package stylist

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestProcessorFilter_Filter(t *testing.T) {
	tests := []struct {
		desc       string
		processors []*Processor
		filter     *ProcessorFilter
		expected   []*Processor
		err        string
	}{
		{
			desc:       "returns an error when passed an empty processor list",
			processors: []*Processor{},
			filter:     &ProcessorFilter{},
			expected:   nil,
			err:        "no processors defined",
		},
		{
			desc: "returns an error when passed unnamed processors",
			processors: []*Processor{
				{Name: "p1"},
				{Name: " "},
			},
			filter:   &ProcessorFilter{},
			expected: nil,
			err:      "processor at index 1 is unnamed",
		},
		{
			desc: "returns an error when passed processors with duplicate names",
			processors: []*Processor{
				{Name: "p1"},
				{Name: "p1"},
			},
			filter:   &ProcessorFilter{},
			expected: nil,
			err:      "processor at index 1 has a duplicate name",
		},
		{
			desc: "returns an error when filtering by an unknown name",
			processors: []*Processor{
				{Name: "p1"},
			},
			filter: &ProcessorFilter{
				Names: []string{"p2"},
			},
			expected: nil,
			err:      "no processor named p2",
		},
		{
			desc: "returns an error when filtering by an unknown tag",
			processors: []*Processor{
				{Name: "p1"},
			},
			filter: &ProcessorFilter{
				Tags: []string{"foo"},
			},
			expected: nil,
			err:      "no processor tagged foo",
		},

		{
			desc: "returns the given processor list when no filters defined",
			processors: []*Processor{
				{Name: "p1"},
				{Name: "p2"},
				{Name: "p3"},
			},
			filter: &ProcessorFilter{},
			expected: []*Processor{
				{Name: "p1"},
				{Name: "p2"},
				{Name: "p3"},
			},
		},
		{
			desc: "returns processors filtered by a single name",
			processors: []*Processor{
				{Name: "p1"},
				{Name: "p2"},
				{Name: "p3"},
			},
			filter: &ProcessorFilter{
				Names: []string{"p1"},
			},
			expected: []*Processor{
				{Name: "p1"},
			},
		},
		{
			desc: "returns processors filtered by multiple names",
			processors: []*Processor{
				{Name: "p1"},
				{Name: "p2"},
				{Name: "p3"},
			},
			filter: &ProcessorFilter{
				Names: []string{"p1", "p3"},
			},
			expected: []*Processor{
				{Name: "p1"},
				{Name: "p3"},
			},
		},
		{
			desc: "returns processors filtered by a single tag",
			processors: []*Processor{
				{Name: "p1", Tags: []string{"foo"}},
				{Name: "p2", Tags: []string{"bar"}},
				{Name: "p3", Tags: []string{"foo", "bar"}},
			},
			filter: &ProcessorFilter{
				Tags: []string{"foo"},
			},
			expected: []*Processor{
				{Name: "p1", Tags: []string{"foo"}},
				{Name: "p3", Tags: []string{"foo", "bar"}},
			},
		},
		{
			desc: "returns processors filtered by multiple tags",
			processors: []*Processor{
				{Name: "p1", Tags: []string{"foo"}},
				{Name: "p2", Tags: []string{"bar"}},
				{Name: "p3", Tags: []string{"foo", "bar"}},
			},
			filter: &ProcessorFilter{
				Tags: []string{"foo", "bar"},
			},
			expected: []*Processor{
				{Name: "p1", Tags: []string{"foo"}},
				{Name: "p2", Tags: []string{"bar"}},
				{Name: "p3", Tags: []string{"foo", "bar"}},
			},
		},
		{
			desc: "returns processors filtered by both name and tag",
			processors: []*Processor{
				{Name: "p1", Tags: []string{"foo"}},
				{Name: "p2", Tags: []string{"bar"}},
				{Name: "p3", Tags: []string{"foo", "bar"}},
			},
			filter: &ProcessorFilter{
				Names: []string{"p1"},
				Tags:  []string{"bar"},
			},
			expected: []*Processor{
				{Name: "p1", Tags: []string{"foo"}},
				{Name: "p2", Tags: []string{"bar"}},
				{Name: "p3", Tags: []string{"foo", "bar"}},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			actual, err := tt.filter.Filter(tt.processors)

			if tt.err == "" {
				assert.NoError(t, err)
			} else {
				assert.ErrorContains(t, err, tt.err)
			}

			assert.ElementsMatch(t, tt.expected, actual)
		})
	}
}
