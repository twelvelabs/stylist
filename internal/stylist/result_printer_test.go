package stylist

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/twelvelabs/termite/ioutil"
)

func TestNewResultPrinter(t *testing.T) {
	ios := ioutil.Test()
	config := NewConfig()
	// Ensure a printer exists for each enum value.
	for _, name := range ResultFormatNames() {
		config.Output.Format = ResultFormat(name)
		assert.NotPanics(t, func() {
			_ = NewResultPrinter(ios, config)
		})
	}
	assert.PanicsWithValue(t, "unknown result format: unknown", func() {
		config.Output.Format = ResultFormat("unknown")
		_ = NewResultPrinter(ios, config)
	})
}
