package stylist

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewApp(t *testing.T) {
	app, err := NewApp(nil)
	assert.NoError(t, err)
	assert.IsType(t, &App{}, app)
	assert.Equal(t, false, app.CmdClient.IsStubbed())
}

func TestNewTestApp(t *testing.T) {
	app := NewTestApp()
	assert.IsType(t, &App{}, app)
	assert.Equal(t, true, app.CmdClient.IsStubbed())
}

func TestNewAppMeta(t *testing.T) {
	meta := NewAppMeta("1.2.3", "9b11774", "1676781982")
	assert.Equal(t, "9b11774", meta.BuildCommit)
	assert.Equal(t, "2023-02-18 22:46", meta.BuildTime.Format("2006-01-02 15:04"))
	assert.Equal(t, "1.2.3", meta.Version)
}
