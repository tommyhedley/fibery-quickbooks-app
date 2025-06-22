package app

import (
	"time"

	"github.com/tommyhedley/fibery-quickbooks-app/pkgs/cache"
)

type DataType string

type SourceID string

type ItemID string

type ItemSet map[ItemID]struct{}

type SourceMap map[SourceID]ItemSet

type TypeMap map[DataType]SourceMap

type IDCache struct {
	cache *cache.Cache[string, TypeMap]
}

func NewIDCache(defaultTTL time.Duration) *IDCache {
	return &IDCache{
		cache: cache.NewCache[string, TypeMap](defaultTTL),
	}
}

func (ic *IDCache) Add(realmID string, dt DataType, src SourceID, id ItemID) {
	dm, ok := ic.cache.Get(realmID)
	if !ok {
		dm = make(TypeMap)
	}

	sm, ok := dm[dt]
	if !ok {
		sm = make(SourceMap)
		dm[dt] = sm
	}

	is, ok := sm[src]
	if !ok {
		is = make(ItemSet)
		sm[src] = is
	}

	is[id] = struct{}{}

	ic.cache.Set(realmID, dm)
}
