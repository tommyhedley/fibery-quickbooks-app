package integration

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/tommyhedley/fibery-quickbooks-app/pkgs/fibery"
	"github.com/tommyhedley/quickbooks-go"
)

type TypeRegistry map[string]fibery.Type

func (tr TypeRegistry) Register(t fibery.Type) {
	tr[t.Id()] = t
}

func (tr TypeRegistry) Get(id string) (fibery.Type, bool) {
	if typ, exists := tr[id]; exists {
		return typ, true
	}
	return nil, false
}

func (tr TypeRegistry) GetAll() []fibery.SyncConfigTypes {
	types := make([]fibery.SyncConfigTypes, 0, len(tr))
	for _, typ := range tr {
		types = append(types, fibery.SyncConfigTypes{
			Id:   typ.Id(),
			Name: typ.Name(),
		})
	}
	return types
}

var Types = make(TypeRegistry)

type ActionRegistry map[string]fibery.Action

func (ar ActionRegistry) Register(a fibery.Action) {
	ar[a.ActionId] = a
}

func (ar ActionRegistry) Get(id string) (fibery.Action, bool) {
	if action, exists := ar[id]; exists {
		return action, true
	}
	return fibery.Action{}, false
}

func (ar ActionRegistry) GetAll() []fibery.Action {
	actions := make([]fibery.Action, 0, len(ar))
	for _, action := range ar {
		actions = append(actions, action)
	}
	return actions
}

var Actions = make(ActionRegistry)

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

func New(parentCtx context.Context, params Parameters) (*Integration, error) {
	ctx, cancel := context.WithCancel(parentCtx)
	config, err := BuildConfig(params)
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

	opManager := NewOperationManager(params.OperationTTL)
	idStore := NewIdStore(params.IdCacheTTL)
	integration := &Integration{
		appConfig: fibery.AppConfig{
			Id:          "qbo",
			Name:        "QuickBooks Online",
			Website:     "https://quickbooks.intuit.com",
			Version:     params.Version,
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
