package app

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/tommyhedley/fibery-quickbooks-app/pkgs/fibery"
	"github.com/tommyhedley/quickbooks-go"
)

type Integration struct {
	appConfig  fibery.AppConfig
	syncConfig fibery.SyncConfig
	config     Config
	types      TypeRegistry
	client     *quickbooks.Client
	opManager  *OperationManager
	idStore    *IdStore
	ctx        context.Context
	cancel     context.CancelFunc
}

func New(parentCtx context.Context, version string) (*Integration, error) {
	ctx, cancel := context.WithCancel(parentCtx)
	config, err := NewConfig(version)
	if err != nil {
		cancel()
		return nil, fmt.Errorf("unable to build config: %w", err)
	}

	discoveryAPI, err := quickbooks.CallDiscoveryAPI(config.DiscoverURL())
	if err != nil {
		cancel()
		return nil, fmt.Errorf("error calling discovery API: %w", err)
	}

	clientReq := config.NewClientRequest(discoveryAPI, http.DefaultClient)

	client, err := quickbooks.NewClient(clientReq)
	if err != nil {
		cancel()
		return nil, fmt.Errorf("error creating quickbooks client: %w", err)
	}

	opManager := NewOperationManager(config.OperationTTL)
	idStore := NewIdStore(config.IdCacheTTL)
	integration := &Integration{
		appConfig: fibery.AppConfig{
			Id:          "qbo",
			Name:        "QuickBooks Online",
			Website:     "https://quickbooks.intuit.com",
			Version:     config.Version,
			Description: "Integrate QuickBooks Online data with Fibery",
			Authentication: []fibery.Authentication{
				{
					Id:          "oauth2",
					Name:        "OAuth v2 Authentication",
					Description: "OAuth v2-based authentication and authorization for access to QuickBooks Online",
					Fields: []fibery.AuthField{
						{
							Id:          "callback_uri",
							Title:       "callback_uri",
							Description: "OAuth post-auth redirect URI",
							Type:        "oauth",
						},
					},
				},
			},
			Sources: []string{},
			ResponsibleFor: fibery.ResponsibleFor{
				DataSynchronization: true,
				Automations:         true,
			},
			Actions: Actions.GetAll(),
		},
		syncConfig: fibery.SyncConfig{
			Types:   Types.GetAll(),
			Filters: []fibery.SyncFilter{},
			Webhooks: fibery.SyncConfigWebhook{
				Enabled: true,
				Type:    "ui",
			},
		},
		config:    config,
		types:     Types,
		client:    client,
		opManager: opManager,
		idStore:   idStore,
		ctx:       ctx,
		cancel:    cancel,
	}
	integration.StartCacheCleaner()
	slog.SetDefault(config.BuildLogger())
	return integration, nil
}

func (i *Integration) Cleanup() {
	i.opManager.CleanupExpired()
	i.idStore.CleanupExpired()
}

func (i *Integration) StartCacheCleaner() {
	ticker := time.NewTicker(5 * time.Second)
	go func() {
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				i.Cleanup()
			case <-i.ctx.Done():
				return
			}
		}
	}()
}

func (i *Integration) Port() string {
	return i.config.Port
}

type QuickBooksAccountInfo struct {
	Name    string `json:"name,omitempty"`
	RealmId string `json:"realmId,omitempty"`
	quickbooks.BearerToken
}
