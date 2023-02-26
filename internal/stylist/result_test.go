package stylist

import (
	"math/rand"
	"sort"
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

func TestResultsByLocation(t *testing.T) {
	expected := []*Result{
		{
			Location: ResultLocation{
				Path:        "aaa.go",
				StartLine:   1,
				StartColumn: 1,
			},
		},
		{
			Location: ResultLocation{
				Path:        "aaa.go",
				StartLine:   1,
				StartColumn: 2,
			},
		},
		{
			Location: ResultLocation{
				Path:        "aaa.go",
				StartLine:   2,
				StartColumn: 1,
			},
		},
		{
			Location: ResultLocation{
				Path:        "aaa.go",
				StartLine:   2,
				StartColumn: 2,
			},
		},
		{
			Source: "aaa",
			Level:  ResultLevelError,
			Location: ResultLocation{
				Path:        "bbb.go",
				StartLine:   1,
				StartColumn: 1,
			},
		},
		{
			Source: "aaa",
			Level:  ResultLevelWarning,
			Location: ResultLocation{
				Path:        "bbb.go",
				StartLine:   1,
				StartColumn: 1,
			},
		},
		{
			Source: "bbb",
			Location: ResultLocation{
				Path:        "bbb.go",
				StartLine:   1,
				StartColumn: 1,
			},
		},
	}

	// Copy the `expected` slice and shuffle it
	actual := make([]*Result, len(expected))
	copy(actual, expected)
	rand.Shuffle(len(actual), func(i, j int) {
		actual[i], actual[j] = actual[j], actual[i]
	})
	// Ensure the two are now different
	assert.NotEqual(t, expected, actual)

	// Run the sorter
	sorter := ResultsByLocation{actual}
	sort.Sort(sorter)

	assert.Equal(t, expected, actual)
}

func TestResultsBySeverity(t *testing.T) {
	expected := []*Result{
		{
			Level: ResultLevelError,
			Location: ResultLocation{
				Path: "aaa.go",
			},
		},
		{
			Level: ResultLevelError,
			Location: ResultLocation{
				Path: "bbb.go",
			},
		},
		{
			Level: ResultLevelError,
			Location: ResultLocation{
				Path: "ccc.go",
			},
		},
		{
			Level: ResultLevelWarning,
		},
		{
			Level: ResultLevelInfo,
		},
		{
			Level: ResultLevelNone,
		},
	}

	// Copy the `expected` slice and shuffle it
	actual := make([]*Result, len(expected))
	copy(actual, expected)
	rand.Shuffle(len(actual), func(i, j int) {
		actual[i], actual[j] = actual[j], actual[i]
	})
	// Ensure the two are now different
	assert.NotEqual(t, expected, actual)

	// Run the sorter
	sorter := ResultsBySeverity{actual}
	sort.Sort(sorter)

	assert.Equal(t, expected, actual)
}

func TestResultsBySource(t *testing.T) {
	expected := []*Result{
		{
			Source: "aaa",
			Location: ResultLocation{
				Path: "aaa.go",
			},
		},
		{
			Source: "aaa",
			Location: ResultLocation{
				Path: "bbb.go",
			},
		},
		{
			Source: "aaa",
			Location: ResultLocation{
				Path: "ccc.go",
			},
		},
		{
			Source: "bbb",
		},
		{
			Source: "ccc",
		},
		{
			Source: "ddd",
		},
	}

	// Copy the `expected` slice and shuffle it
	actual := make([]*Result, len(expected))
	copy(actual, expected)
	rand.Shuffle(len(actual), func(i, j int) {
		actual[i], actual[j] = actual[j], actual[i]
	})
	// Ensure the two are now different
	assert.NotEqual(t, expected, actual)

	// Run the sorter
	sorter := ResultsBySource{actual}
	sort.Sort(sorter)

	assert.Equal(t, expected, actual)
}
