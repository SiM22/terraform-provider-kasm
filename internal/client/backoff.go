// internal/client/backoff.go

package client

import (
	"math/rand"
	"time"
)

// ExponentialBackoff implements exponential backoff with jitter
type ExponentialBackoff struct {
	InitialInterval     time.Duration
	MaxInterval         time.Duration
	Multiplier          float64
	RandomizationFactor float64
	currentInterval     time.Duration
}

func NewExponentialBackoff(config *RetryConfig) *ExponentialBackoff {
	return &ExponentialBackoff{
		InitialInterval:     config.InitialInterval,
		MaxInterval:         config.MaxInterval,
		Multiplier:          config.Multiplier,
		RandomizationFactor: config.RandomizationFactor,
		currentInterval:     0,
	}
}

// NextBackOff calculates the next backoff duration
func (b *ExponentialBackoff) NextBackOff() time.Duration {
	if b.currentInterval == 0 {
		b.currentInterval = b.InitialInterval
	} else {
		b.currentInterval = time.Duration(float64(b.currentInterval) * b.Multiplier)
		if b.currentInterval > b.MaxInterval {
			b.currentInterval = b.MaxInterval
		}
	}

	// Add jitter
	delta := b.RandomizationFactor * float64(b.currentInterval)
	minInterval := float64(b.currentInterval) - delta
	maxInterval := float64(b.currentInterval) + delta

	return time.Duration(minInterval + (rand.Float64() * (maxInterval - minInterval)))
}

// Reset resets the backoff to its initial state
func (b *ExponentialBackoff) Reset() {
	b.currentInterval = 0
}
