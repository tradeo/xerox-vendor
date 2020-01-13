package util

import "time"

// Ticker delivers `ticks` of a clock at intervals.
//
// It is used to adapt the time.Ticker for easier testing.
// One can use the TickerFactory to abstract away the tickers production.
type Ticker interface {
	// Chan returns a channel with the produced ticks.
	Chan() <-chan time.Time

	// Stop stops the ticker.
	Stop()
}

// TimeTicker implements the ticker interface by using golang time.Ticker.
type TimeTicker struct {
	*time.Ticker
}

// Chan returns the ticker channel.
func (tt *TimeTicker) Chan() <-chan time.Time {
	return tt.C
}

// TickerFactory creates Tickers.
type TickerFactory interface {
	Create(period time.Duration) Ticker
}

// DefaultTickerFactory creates TimeTicker instances.
type DefaultTickerFactory struct {
}

// Create makes TimeTicker instance.
func (f *DefaultTickerFactory) Create(period time.Duration) Ticker {
	return &TimeTicker{time.NewTicker(period)}
}

// NewTickerFactory create DefaultTickerFactory instance.
func NewTickerFactory() *DefaultTickerFactory {
	return &DefaultTickerFactory{}
}

// TestTicker creates non active ticker for testing purposes.
// One can use the TestTicker.C channel to produce ticks.
type TestTicker struct {
	C chan time.Time
}

// Chan implements the Ticker interface.
func (tt *TestTicker) Chan() <-chan time.Time {
	return tt.C
}

// Stop implements the Ticker interface.
func (tt *TestTicker) Stop() {
}

// TestTickerFactory creates TestTicker instances.
type TestTickerFactory struct {
	C chan time.Time
}

// Create implements the TickerFactory interface.
func (f *TestTickerFactory) Create(period time.Duration) Ticker {
	return &TestTicker{
		C: f.C,
	}
}

// NewTestTickerFactory creates TestTickerFactory instance.
func NewTestTickerFactory(ch chan time.Time) *TestTickerFactory {
	return &TestTickerFactory{
		C: ch,
	}
}
