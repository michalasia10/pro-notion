package domain

import "time"

// Clock provides time-related operations for domain services
type Clock interface {
	Now() time.Time
}

// SystemClock implements Clock using system time
type SystemClock struct{}

func NewSystemClock() Clock {
	return &SystemClock{}
}

func (c *SystemClock) Now() time.Time {
	return time.Now()
}
