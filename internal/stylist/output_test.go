package stylist

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewOutputParser(t *testing.T) {
	// Ensure a parser exists for each enum value.
	for _, name := range OutputTypeNames() {
		assert.NotPanics(t, func() {
			_ = NewOutputParser(OutputType(name))
		})
	}
	assert.PanicsWithValue(t, "unknown output type: unknown", func() {
		_ = NewOutputParser(OutputType("unknown"))
	})
}
