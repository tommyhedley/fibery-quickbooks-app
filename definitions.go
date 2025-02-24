package main

import (
	"sync"

	"github.com/tommyhedley/quickbooks-go"
)

type integrationAccountInfo struct {
	Name    string `json:"name,omitempty"`
	RealmID string `json:"realmId,omitempty"`
	quickbooks.BearerToken
}

var globalOperationCache = make(map[string]*OperationCache)
var globalCacheMutex sync.Mutex

type CacheEntry struct {
	result chan any
}

type OperationCache struct {
	RequestedTypes map[string]bool
	Results        map[string]*CacheEntry
	sync.Mutex
}
