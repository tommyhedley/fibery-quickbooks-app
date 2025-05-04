package fibery

import (
	"net/http"
)

func RegisterFiberyRoutes(mux *http.ServeMux, integration any) {
	if core, ok := integration.(IntegrationCore); ok {
		mux.HandleFunc("GET /", core.AppConfigHandler)
		mux.HandleFunc("POST /validate", core.AccountValidateHandler)
		mux.HandleFunc("POST /api/v1/synchronizer/config", core.SyncConfigHandler)
		mux.HandleFunc("POST /api/v1/synchronizer/schema", core.SyncSchemaHandler)
		mux.HandleFunc("POST /api/v1/synchronizer/data", core.SyncDataHandler)
	} else {
		panic("Integration does not implement CoreIntegration interface")
	}

	if webhooks, ok := integration.(IntegrationWebhooks); ok {
		mux.HandleFunc("POST /api/v1/synchronizer/webhooks", webhooks.WebhookInitHandler)
		mux.HandleFunc("POST /api/v1/synchronizer/webhooks/pre-process", webhooks.WebhookPreProcessHandler)
		mux.HandleFunc("POST /api/v1/synchronizer/webhooks/transform", webhooks.WebhookTransformHandler)
	}

	if oauth1, ok := integration.(IntegrationOauth1); ok {
		mux.HandleFunc("POST /oauth1/v1/authorize", oauth1.Oauth1AuthorizeHandler)
		mux.HandleFunc("POST /oauth1/v1/access_token", oauth1.Oauth1TokenHandler)
	}

	if oauth2, ok := integration.(IntegrationOauth2); ok {
		mux.HandleFunc("POST /oauth2/v1/authorize", oauth2.Oauth2AuthorizeHandler)
		mux.HandleFunc("POST /oauth2/v1/access_token", oauth2.Oauth2TokenHandler)
	}

	if logo, ok := integration.(IntegrationLogo); ok {
		mux.HandleFunc("GET /logo", logo.LogoHandler)
	}

	if datalist, ok := integration.(IntegrationDatalist); ok {
		mux.HandleFunc("POST /api/v1/synchronizer/datalist", datalist.SyncDatalistHandler)
	}

	if filter, ok := integration.(IntegrationFilter); ok {
		mux.HandleFunc("POST /api/v1/synchronizer/filter/validate", filter.SyncFilterValidateHandler)
	}

	if resource, ok := integration.(IntegrationResource); ok {
		mux.HandleFunc("POST /api/v1/synchronizer/resource", resource.SyncResourceHandler)
	}

	if actions, ok := integration.(IntegrationActions); ok {
		mux.HandleFunc("POST /api/v1/automations/action/execute", actions.ActionHandler)
	}
}
