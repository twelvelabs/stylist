package stylist

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestProcessorList_All(t *testing.T) {
	p1 := &Processor{}
	p2 := &Processor{}
	p3 := &Processor{}
	list := ProcessorList{p1, p2, p3}.All()
	assert.Equal(t, 3, len(list))
}

func TestProcessorList_Named(t *testing.T) {
	p1 := &Processor{Name: "foo"}
	p2 := &Processor{Name: "bar"}
	p3 := &Processor{Name: "baz"}
	list := ProcessorList{p1, p2, p3}

	assert.Equal(t, []*Processor{p1}, list.Named("foo"))
	assert.Equal(t, []*Processor{p1, p3}, list.Named("foo", "baz"))
}

func TestProcessorList_Tagged(t *testing.T) {
	p1 := &Processor{
		Tags: []string{"foo", "bar"},
	}
	p2 := &Processor{}
	p3 := &Processor{
		Tags: []string{"foo"},
	}
	list := ProcessorList{p1, p2, p3}

	assert.Equal(t, []*Processor{p1, p3}, list.Tagged("foo"))
	assert.Equal(t, []*Processor{p1}, list.Tagged("foo", "bar"))
}
