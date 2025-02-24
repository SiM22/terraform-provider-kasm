package client

import (
	"fmt"
	"sync"
	"time"
)

type RateLimiter struct {
	tokens     int
	capacity   int
	refillRate time.Duration
	lastRefill time.Time
	mu         sync.Mutex
}

func NewRateLimiter(capacity int, refillRate time.Duration) *RateLimiter {
	return &RateLimiter{
		tokens:     capacity,
		capacity:   capacity,
		refillRate: refillRate,
		lastRefill: time.Now(),
	}
}

func (r *RateLimiter) Take() error {
	r.mu.Lock()
	defer r.mu.Unlock()

	now := time.Now()
	elapsedTime := now.Sub(r.lastRefill)
	newTokens := int(elapsedTime / r.refillRate)

	if newTokens > 0 {
		r.tokens = min(r.capacity, r.tokens+newTokens)
		r.lastRefill = now
	}

	if r.tokens <= 0 {
		return fmt.Errorf("rate limit exceeded, try again in %v", r.refillRate)
	}

	r.tokens--
	return nil
}

// min returns the smaller of x or y
func min(x, y int) int {
	if x < y {
		return x
	}
	return y
}
