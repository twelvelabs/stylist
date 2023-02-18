package stylist

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestResultLocation_String(t *testing.T) {
	var loc ResultLocation

	loc = ResultLocation{}
	assert.Equal(t, "<none>:0:0", loc.String())

	loc = ResultLocation{
		Path: "foo/bar.go",
	}
	assert.Equal(t, "foo/bar.go:0:0", loc.String())

	loc = ResultLocation{
		Path:        "foo/bar.go",
		StartLine:   10,
		StartColumn: 12,
	}
	assert.Equal(t, "foo/bar.go:10:12", loc.String())
}

func TestNewResultsError(t *testing.T) {
	var err error

	// `nil` slice should return a nil error
	err = NewResultsError(nil)
	assert.NoError(t, err)

	// empty slice should return a nil error
	err = NewResultsError([]*Result{})
	assert.NoError(t, err)

	// non-empty slice should return an error
	err = NewResultsError([]*Result{
		{Source: "test-linter"},
	})
	assert.Error(t, err)
}

func TestResultsError_Error(t *testing.T) {
	var err error

	err = &ResultsError{}
	assert.Equal(t, "0 issue(s)", err.Error())

	err = &ResultsError{
		results: []*Result{
			{Source: "test-linter"},
			{Source: "test-linter"},
		},
	}
	assert.Equal(t, "2 issue(s)", err.Error())
}
