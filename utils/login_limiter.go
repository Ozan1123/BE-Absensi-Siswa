package utils

import (
	"sync"
	"time"
)

// LoginRateLimiter blocks repeated failed login attempts for a specific identifier.
type LoginRateLimiter struct {
	mu          sync.Mutex
	attempts    map[string]*loginAttempt
	maxAttempts int
	window      time.Duration
	lockout     time.Duration
}

type loginAttempt struct {
	failures     int
	firstFailure time.Time
	lockedUntil  time.Time
}

func NewLoginRateLimiter(maxAttempts int, window, lockout time.Duration) *LoginRateLimiter {
	if maxAttempts <= 0 {
		maxAttempts = 3
	}
	if window <= 0 {
		window = 5 * time.Minute
	}
	if lockout <= 0 {
		lockout = 5 * time.Minute
	}

	return &LoginRateLimiter{
		attempts:    make(map[string]*loginAttempt),
		maxAttempts: maxAttempts,
		window:      window,
		lockout:     lockout,
	}
}

func (l *LoginRateLimiter) RecordFailure(identifier string, now time.Time) (allowed bool, retryAfter time.Duration) {
	l.mu.Lock()
	defer l.mu.Unlock()

	entry, ok := l.attempts[identifier]
	if !ok || entry.firstFailure.IsZero() {
		entry = &loginAttempt{}
		l.attempts[identifier] = entry
	}

	if entry.lockedUntil.After(now) {
		return false, entry.lockedUntil.Sub(now)
	}

	if entry.firstFailure.IsZero() || now.Sub(entry.firstFailure) > l.window {
		entry.failures = 0
		entry.firstFailure = now
	}

	entry.failures++
	if entry.failures >= l.maxAttempts {
		entry.lockedUntil = now.Add(l.lockout)
		entry.firstFailure = now
		return false, l.lockout
	}

	return true, 0
}

func (l *LoginRateLimiter) Reset(identifier string) {
	l.mu.Lock()
	defer l.mu.Unlock()
	delete(l.attempts, identifier)
}
