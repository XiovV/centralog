package main

import (
	"time"
)

type Limiter struct {
	allowedAttempts int
	tries           int
	timeout         time.Duration
	unlockAt        int64
}

func NewLimiter(attempts int, timeout time.Duration) *Limiter {
	return &Limiter{
		allowedAttempts: attempts,
		timeout:         timeout,
	}
}

func (l *Limiter) Fail() {
	l.tries++
}

func (l *Limiter) reset() {
	l.tries = 0
	l.unlockAt = 0
}

func (l *Limiter) Try() bool {
	currTime := time.Now().Unix()

	if currTime != 0 && l.unlockAt-currTime > 0 {
		return false
	}

	if l.unlockAt > 0 && l.unlockAt-currTime <= 0 && l.tries == l.allowedAttempts {
		l.reset()
	}

	if l.tries == l.allowedAttempts {
		l.unlockAt = currTime + int64(l.timeout.Seconds())

		return false
	}

	return true
}
