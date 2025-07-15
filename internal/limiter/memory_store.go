package limiter

import (
	"sync"
	"time"
)

type InMemoryStore struct {
	mu          sync.Mutex
	requests    map[string][]time.Time
	lockedUntil map[string]time.Time
}

func NewInMemoryStore() *InMemoryStore {
	return &InMemoryStore{
		requests:    make(map[string][]time.Time),
		lockedUntil: make(map[string]time.Time),
	}
}

func (s *InMemoryStore) Allow(id string, limit int, window time.Duration) bool {
	s.mu.Lock()
	defer s.mu.Unlock()

	now := time.Now()
	reqs := s.requests[id]

	var filtered []time.Time
	for _, t := range reqs {
		if now.Sub(t) <= window {
			filtered = append(filtered, t)
		}
	}

	if len(filtered) >= limit {
		s.requests[id] = filtered
		return false
	}

	filtered = append(filtered, now)
	s.requests[id] = filtered
	return true
}

func (s *InMemoryStore) IsLocked(id string) bool {
	s.mu.Lock()
	defer s.mu.Unlock()

	until, ok := s.lockedUntil[id]
	return ok && time.Now().Before(until)
}

func (s *InMemoryStore) SetLock(id string, until time.Time) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.lockedUntil[id] = until
}

func (s *InMemoryStore) GetLockUntil(id string) time.Time {
	s.mu.Lock()
	defer s.mu.Unlock()

	return s.lockedUntil[id]
}
