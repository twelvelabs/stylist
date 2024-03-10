package lsp

import (
	"sync"
	"time"

	"github.com/jonboulle/clockwork"
)

// NewDebouncer returns a Debouncer for the given duration window.
func NewDebouncer(window time.Duration) *Debouncer {
	return &Debouncer{
		window: window,
		clock:  clockwork.NewRealClock(),
	}
}

type Debouncer struct {
	mu     sync.Mutex
	clock  clockwork.Clock
	timer  clockwork.Timer
	window time.Duration
}

// Run runs f in a goroutine after the debounce window elapses.
// Pending instances of f that have yet to be run are discarded.
func (d *Debouncer) Run(f func()) {
	d.mu.Lock()
	defer d.mu.Unlock()

	if d.timer != nil {
		d.timer.Stop()
	}
	d.timer = d.clock.AfterFunc(d.window, f)
}

// WithClock sets clock and returns the receiver.
func (d *Debouncer) WithClock(clock clockwork.Clock) *Debouncer {
	d.clock = clock
	return d
}
