package utils

import (
	"testing"
	"time"
)

func TestLoginRateLimiterBlocksAfterMaxAttempts(t *testing.T) {
	limiter := NewLoginRateLimiter(3, 5*time.Minute, 5*time.Minute)

	for i := 0; i < 2; i++ {
		allowed, _ := limiter.RecordFailure("nisn-1", time.Now())
		if !allowed {
			t.Fatalf("expected attempt %d to be allowed", i+1)
		}
	}

	allowed, _ := limiter.RecordFailure("nisn-1", time.Now())
	if allowed {
		t.Fatal("expected third failure to be blocked")
	}
}

func TestLoginRateLimiterResetsAfterWindowExpires(t *testing.T) {
	limiter := NewLoginRateLimiter(2, 5*time.Minute, 5*time.Minute)

	allowed, _ := limiter.RecordFailure("nisn-2", time.Now())
	if !allowed {
		t.Fatal("expected first failure to be allowed")
	}

	allowed, _ = limiter.RecordFailure("nisn-2", time.Now().Add(6*time.Minute))
	if !allowed {
		t.Fatal("expected failure after window expiry to be allowed")
	}
}
