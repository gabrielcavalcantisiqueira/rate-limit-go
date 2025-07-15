package limiter

import "time"

type Store interface {
	Allow(id string, limit int, window time.Duration) bool
	IsLocked(id string) bool
	SetLock(id string, until time.Time)
	GetLockUntil(id string) time.Time
}
