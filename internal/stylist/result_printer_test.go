package stylist

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewResultPrinter(t *testing.T) {
	// Ensure a printer exists for each enum value.
	for _, name := range ResultFormatNames() {
		assert.NotPanics(t, func() {
			_ = NewResultPrinter(ResultFormat(name))
		})
	}
	assert.PanicsWithValue(t, "unknown result format: unknown", func() {
		_ = NewResultPrinter(ResultFormat("unknown"))
	})
}
