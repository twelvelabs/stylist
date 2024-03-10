package lsp

import (
	"testing"
	"time"

	"github.com/jonboulle/clockwork"
	"github.com/stretchr/testify/require"
)

func TestDebouncer_Run(t *testing.T) {
	clock := clockwork.NewFakeClock()
	debouncer := NewDebouncer(5 * time.Second).WithClock(clock)

	// Invoke 5 times
	var calls int
	for i := 0; i < 5; i++ {
		debouncer.Run(func() {
			calls++
		})
		clock.BlockUntil(1)
		clock.Advance(1 * time.Second)
	}

	// Advance past the debounce window
	// and ensure the callback was only run once.
	clock.BlockUntil(1)
	clock.Advance(5 * time.Second)
	// The func passed to debouncer.Run() is run in a goroutine.
	// Give it time to execute.
	time.Sleep(1 * time.Millisecond)

	require.Equal(t, 1, calls)
}
