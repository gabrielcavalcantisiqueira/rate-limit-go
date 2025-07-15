package limiter

import (
	"testing"
	"time"
)

func TestAllowWithinLimit(t *testing.T) {
	store := NewInMemoryStore()
	id := "user1"
	limit := 3
	window := 1 * time.Second

	for i := 0; i < limit; i++ {
		if !store.Allow(id, limit, window) {
			t.Errorf("Request %d should be allowed", i+1)
		}
	}
}

func TestBlockAfterLimit(t *testing.T) {
	store := NewInMemoryStore()
	id := "user2"
	limit := 2
	window := 1 * time.Second

	store.Allow(id, limit, window)
	store.Allow(id, limit, window)
	if store.Allow(id, limit, window) {
		t.Error("Third request should be blocked")
	}
}

func TestAllowAfterWindowExpires(t *testing.T) {
	store := NewInMemoryStore()
	id := "user3"
	limit := 1
	window := 500 * time.Millisecond

	store.Allow(id, limit, window)

	time.Sleep(window + 100*time.Millisecond)

	if !store.Allow(id, limit, window) {
		t.Error("Request should be allowed after window expires")
	}
}

func TestLockAndIsLocked(t *testing.T) {
	store := NewInMemoryStore()
	id := "user4"
	lockDuration := 500 * time.Millisecond
	lockUntil := time.Now().Add(lockDuration)

	store.SetLock(id, lockUntil)

	if !store.IsLocked(id) {
		t.Error("User should be locked")
	}

	// Espera passar o tempo de lock
	time.Sleep(lockDuration + 100*time.Millisecond)

	if store.IsLocked(id) {
		t.Error("User should no longer be locked")
	}
}

func TestGetLockUntil(t *testing.T) {
	store := NewInMemoryStore()
	id := "user5"
	lockUntil := time.Now().Add(1 * time.Minute)
	store.SetLock(id, lockUntil)

	returned := store.GetLockUntil(id)

	if !returned.Equal(lockUntil) {
		t.Errorf("Expected lockUntil %v, got %v", lockUntil, returned)
	}
}
