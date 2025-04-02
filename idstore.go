package main

import (
	"sync"
	"time"
)

type IdKey struct {
	EntityType string
	EntityId   string
}

type IdStore struct {
	sync.Mutex
	idCaches map[string]*IdCache
	ttl      time.Duration
}

type IdCache struct {
	sync.RWMutex
	ids        map[IdKey]map[string]map[string]struct{}
	expiration time.Time
}

func NewIdStore(ttl time.Duration) *IdStore {
	return &IdStore{
		idCaches: make(map[string]*IdCache),
		ttl:      ttl,
	}
}

func (s *IdStore) CleanupExpired() {
	s.Lock()
	defer s.Unlock()
	now := time.Now()
	for companyId, idCache := range s.idCaches {
		if !idCache.expiration.IsZero() && now.After(idCache.expiration) {
			delete(s.idCaches, companyId)
		}
	}
}

func (s *IdStore) GetOrCreateIdCache(realmId string) (*IdCache, bool) {
	s.Lock()
	defer s.Unlock()

	cache, ok := s.idCaches[realmId]
	if ok && time.Now().Before(cache.expiration) {
		return cache, true
	}

	cache = &IdCache{
		ids:        make(map[IdKey]map[string]map[string]struct{}),
		expiration: time.Now().Add(s.ttl),
	}
	s.idCaches[realmId] = cache
	return cache, false
}

func (c *IdCache) SetIds(source IdKey, entityType string, newIds map[string]struct{}) {
	c.Lock()
	defer c.Unlock()

	if _, exists := c.ids[source]; !exists {
		c.ids[source] = make(map[string]map[string]struct{})
	}
	c.ids[source][entityType] = newIds
}

func (c *IdCache) AddIds(source IdKey, entityType string, newIds map[string]struct{}) {
	c.Lock()
	defer c.Unlock()

	if _, exists := c.ids[source]; !exists {
		c.ids[source] = make(map[string]map[string]struct{})
	}
	if existing, exists := c.ids[source][entityType]; exists {
		for id := range newIds {
			existing[id] = struct{}{}
		}
	} else {
		c.ids[source][entityType] = newIds
	}
}

func (c *IdCache) AddId(source IdKey, entityType string, entityId string) {
	c.Lock()
	defer c.Unlock()

	if _, exists := c.ids[source]; !exists {
		c.ids[source] = make(map[string]map[string]struct{})
	}
	if _, exists := c.ids[source][entityType]; !exists {
		c.ids[source][entityType] = make(map[string]struct{})
	}
	c.ids[source][entityType][entityId] = struct{}{}
}

func (c *IdCache) GetSourceMap(source IdKey) (map[string]map[string]struct{}, bool) {
	c.RLock()
	defer c.RUnlock()

	ids, exists := c.ids[source]
	return ids, exists
}

func (c *IdCache) GetIdsByType(source IdKey, entityType string) (map[string]struct{}, bool) {
	c.RLock()
	defer c.RUnlock()

	ids, exists := c.ids[source][entityType]
	return ids, exists
}

func (c *IdCache) CheckId(source IdKey, entityType string, entityId string) bool {
	c.RLock()
	defer c.RUnlock()

	if entityMap, exists := c.ids[source]; exists {
		if idSet, exists := entityMap[entityType]; exists {
			if _, exists := idSet[entityId]; exists {
				return true
			}
		}
	}
	return false
}

func (c *IdCache) RemoveEntityType(source IdKey, entityType string) bool {
	c.Lock()
	defer c.Unlock()

	if entityMap, exists := c.ids[source]; exists {
		if _, exists := entityMap[entityType]; exists {
			delete(entityMap, entityType)
			if len(entityMap) == 0 {
				delete(c.ids, source)
			}
			return true
		}
	}
	return false
}

func (c *IdCache) RemoveSource(source IdKey) bool {
	c.Lock()
	defer c.Unlock()

	if _, exists := c.ids[source]; exists {
		delete(c.ids, source)
		return true
	}
	return false
}

