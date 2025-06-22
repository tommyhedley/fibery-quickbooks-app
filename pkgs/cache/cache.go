package cache

import (
	"sync"
	"time"
)

type entry[V any] struct {
	value      V
	expiration time.Time
}

type Cache[I comparable, V any] struct {
	mu         sync.RWMutex
	entries    map[I]entry[V]
	defaultTTL time.Duration
}

func NewCache[I comparable, V any](defaultTTL time.Duration) *Cache[I, V] {
	return &Cache[I, V]{
		entries:    make(map[I]entry[V]),
		defaultTTL: defaultTTL,
	}
}

func (c *Cache[I, V]) Set(id I, value V) {
	c.SetWithTTL(id, value, c.defaultTTL)
}

func (c *Cache[I, V]) SetWithTTL(id I, value V, ttl time.Duration) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.entries[id] = entry[V]{
		value:      value,
		expiration: time.Now().Add(ttl),
	}
}

func (c *Cache[I, V]) Get(id I) (V, bool) {
	c.mu.RLock()
	e, ok := c.entries[id]
	c.mu.RUnlock()

	if !ok || time.Now().After(e.expiration) {
		c.mu.Lock()
		delete(c.entries, id)
		c.mu.Unlock()
		var zero V
		return zero, false
	}
	return e.value, true
}

func (c *Cache[I, V]) Cleanup() {
	now := time.Now()
	c.mu.Lock()
	defer c.mu.Unlock()

	for id, e := range c.entries {
		if now.After(e.expiration) {
			delete(c.entries, id)
		}
	}
}
