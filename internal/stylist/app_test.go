package stylist

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewApp(t *testing.T) {
	app, err := NewApp()
	assert.NoError(t, err)
	assert.IsType(t, &App{}, app)
	assert.Equal(t, false, app.CmdClient.IsStubbed())
}

func TestNewTestApp(t *testing.T) {
	app := NewTestApp()
	assert.IsType(t, &App{}, app)
	assert.Equal(t, true, app.CmdClient.IsStubbed())
}
