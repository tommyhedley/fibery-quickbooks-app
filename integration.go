package main

import (
	"time"

	"github.com/tommyhedley/fibery-quickbooks-app/pkgs/fibery"
	"github.com/tommyhedley/quickbooks-go"
)

type Integration struct {
	appConfig  fibery.AppConfig
	syncConfig fibery.SyncConfig
	config     ProgramConfig
	types      TypeRegistry
	client     *quickbooks.Client
	opManager  *OperationManager
	idStore    *IdStore
}

func NewIntegration(appConfig fibery.AppConfig, syncConfig fibery.SyncConfig, config ProgramConfig, types TypeRegistry, client *quickbooks.Client, opttl, idttl time.Duration) *Integration {
	opManager := NewOperationManager(opttl)
	idStore := NewIdStore(idttl)
	integration := &Integration{
		appConfig:  appConfig,
		syncConfig: syncConfig,
		config:     config,
		types:      types,
		client:     client,
		opManager:  opManager,
		idStore:    idStore,
	}
	integration.StartCacheCleaner()
	return integration
}

func (i *Integration) Cleanup() {
	i.opManager.CleanupExpired()
	i.idStore.CleanupExpired()
}

func (i *Integration) StartCacheCleaner() {
	ticker := time.NewTicker(5 * time.Second)
	go func() {
		defer ticker.Stop()
		for range ticker.C {
			i.Cleanup()
		}
	}()
}
