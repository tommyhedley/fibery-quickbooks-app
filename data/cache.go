package data

import (
	"context"
	"sync"
	"time"

	"github.com/tommyhedley/fibery-quickbooks-app/pkgs/fibery"
)

type CacheEntry struct {
	Data chan any
}

type OperationCache struct {
	sync.RWMutex
	RequestedTypes map[string]bool
	Results        map[string]*CacheEntry
	SyncTypes      map[string]fibery.SyncType
	lastAccess     time.Time
	ctx            context.Context
	cancel         context.CancelFunc
}

func NewOperationCache(requestedTypes map[string]bool, syncTypes map[string]fibery.SyncType) *OperationCache {
	ctx, cancel := context.WithCancel(context.Background())
	return &OperationCache{
		RequestedTypes: requestedTypes,
		Results:        make(map[string]*CacheEntry),
		SyncTypes:      syncTypes,
		lastAccess:     time.Now(),
		ctx:            ctx,
		cancel:         cancel,
	}
}

func (c *OperationCache) MarkTypeComplete(typeId string) {
	c.Lock()
	c.RequestedTypes[typeId] = true
	c.Unlock()
}

func (c *OperationCache) RefreshAccessTime() {
	c.Lock()
	c.lastAccess = time.Now()
	c.Unlock()
}

type Cache struct {
	sync.Mutex
	operations map[string]*OperationCache
	ttl        time.Duration
}

func NewDataCache(ttl time.Duration) *Cache {
	c := &Cache{
		operations: make(map[string]*OperationCache),
		ttl:        ttl,
	}
	return c
}

func (c *Cache) CleanupExpired() {
	c.Lock()
	defer c.Unlock()
	now := time.Now()
	for opId, opCache := range c.operations {
		opCache.Lock()
		if now.Sub(opCache.lastAccess) > c.ttl {
			opCache.cancel()
			delete(c.operations, opId)
		}
		opCache.Unlock()
	}
}

func (c *Cache) Get(operationId string) (*OperationCache, bool) {
	c.Lock()
	defer c.Unlock()
	opCache, ok := c.operations[operationId]
	return opCache, ok
}

func (c *Cache) Set(operationId string, opCache *OperationCache) {
	c.Lock()
	defer c.Unlock()
	c.operations[operationId] = opCache
}

func (c *Cache) Delete(operationId string) {
	c.Lock()
	defer c.Unlock()
	delete(c.operations, operationId)
}

// IdSet is intented to cache QuickBooks dependent/subtype IDs to use for delta syncing.
// The intended structure is sourceId:entityId.
type IdSet map[string]map[string]struct{}

type idEntry struct {
	value      IdSet
	expiration time.Time
}

// IdCache is intended to be used as a memcache for storing data for delta sync comparison.
// The indended structure of Entries is realmId:dataType:idEntry
type IdCache struct {
	sync.RWMutex
	data map[string]map[string]idEntry
	ttl  time.Duration
}

func NewIdCache(ttl time.Duration) *IdCache {
	c := &IdCache{
		data: make(map[string]map[string]idEntry),
		ttl:  ttl,
	}
	return c
}

func (c *IdCache) CleanupExpired() {
	c.Lock()
	defer c.Unlock()
	now := time.Now()
	for realmId, realmData := range c.data {
		for dataType, entry := range realmData {
			if !entry.expiration.IsZero() && now.After(entry.expiration) {
				delete(realmData, dataType)
			}
		}
		if len(realmData) == 0 {
			delete(c.data, realmId)
		}
	}
}

// Get retrieves the idEntry for a key.
func (c *IdCache) Get(realmId, dataType string) (IdSet, bool) {
	c.RLock()
	defer c.RUnlock()
	realmData, ok := c.data[realmId]
	if !ok {
		return nil, false
	}
	entry, ok := realmData[dataType]
	if !ok || (!entry.expiration.IsZero() && time.Now().After(entry.expiration)) {
		return nil, false
	}
	return entry.value, ok
}

// Set stores an idEntry under the given key.
func (c *IdCache) Set(realmId, dataType string, value IdSet) {
	c.Lock()
	defer c.Unlock()
	if _, exists := c.data[realmId]; !exists {
		c.data[realmId] = make(map[string]idEntry)
	}
	c.data[realmId][dataType] = idEntry{
		value:      value,
		expiration: time.Now().Add(c.ttl),
	}
}

// Delete removes a key from the cache.
func (c *IdCache) Delete(realmId, dataType string) {
	c.Lock()
	defer c.Unlock()
	if realmData, exists := c.data[realmId]; exists {
		delete(realmData, dataType)
		if len(realmData) == 0 {
			delete(c.data, realmId)
		}
	}
}

func (c *IdCache) AddID(realmId, dataType, sourceId, id string) {
	c.Lock()
	defer c.Unlock()
	now := time.Now()
	if _, exists := c.data[realmId]; !exists {
		c.data[realmId] = make(map[string]idEntry)
	}
	entry, exists := c.data[realmId][dataType]
	if !exists || (exists && now.After(entry.expiration)) {
		entry = idEntry{
			value:      make(IdSet),
			expiration: now.Add(c.ttl),
		}
	} else {
		entry.expiration = now.Add(c.ttl)
	}
	if entry.value[sourceId] == nil {
		entry.value[sourceId] = make(map[string]struct{})
	}
	entry.value[sourceId][id] = struct{}{}
	c.data[realmId][dataType] = entry
}

// RemoveID removes an id from the set at the given key.
func (c *IdCache) RemoveID(realmId, dataType, sourceId, id string) {
	c.Lock()
	defer c.Unlock()
	if realmData, exists := c.data[realmId]; exists {
		if entry, exists := realmData[dataType]; exists {
			if set, exists := entry.value[sourceId]; exists {
				delete(set, id)
				if len(set) == 0 {
					delete(entry.value, sourceId)
				}
			}
			if len(entry.value) == 0 {
				delete(realmData, dataType)
			} else {
				realmData[dataType] = entry
			}
		}
		if len(realmData) == 0 {
			delete(c.data, realmId)
		}
	}
}

// RemoveSource removes an id from the set at the given key.
func (c *IdCache) RemoveSource(realmId, dataType, sourceId string) {
	c.Lock()
	defer c.Unlock()
	if realmData, exists := c.data[realmId]; exists {
		if entry, exists := realmData[dataType]; exists {
			delete(entry.value, sourceId)
			if len(entry.value) == 0 {
				delete(realmData, dataType)
			} else {
				realmData[dataType] = entry
			}
		}
		if len(realmData) == 0 {
			delete(c.data, realmId)
		}
	}
}
