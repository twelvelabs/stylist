package stylist

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestPathAdjuster_WhenAbsolute(t *testing.T) {
	require := require.New(t)

	adjuster := NewPathAdjuster("/foo", ResultPathAbsolute)

	var (
		path string
		err  error
	)

	path, err = adjuster.Convert("bar")
	require.Equal("/foo/bar", path)
	require.NoError(err)

	path, err = adjuster.Convert("./bar")
	require.Equal("/foo/bar", path)
	require.NoError(err)

	path, err = adjuster.Convert("/foo/bar")
	require.Equal("/foo/bar", path)
	require.NoError(err)

	path, err = adjuster.Convert("/other/path")
	require.Equal("/other/path", path)
	require.NoError(err)
}

func TestPathAdjuster_WhenRelative(t *testing.T) {
	require := require.New(t)

	adjuster := NewPathAdjuster("/foo", ResultPathRelative)

	var (
		path string
		err  error
	)

	path, err = adjuster.Convert("bar")
	require.Equal("bar", path)
	require.NoError(err)

	path, err = adjuster.Convert("./bar")
	require.Equal("bar", path)
	require.NoError(err)

	path, err = adjuster.Convert("/foo/bar")
	require.Equal("bar", path)
	require.NoError(err)

	path, err = adjuster.Convert("/other/path")
	require.Equal("../other/path", path)
	require.NoError(err)
}

func TestPathAdjuster_WhenRelative_WhenError(t *testing.T) {
	require := require.New(t)

	adjuster := NewPathAdjuster("./foo", ResultPathRelative)

	_, err := adjuster.Convert("/other/path")
	require.Error(err)
}
