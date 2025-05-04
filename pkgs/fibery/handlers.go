package fibery

import (
	"net/http"
)

// CoreIntegration specifies the required endpoints for a Fibery integration
type IntegrationCore interface {
	AppConfigHandler(w http.ResponseWriter, r *http.Request)
	AccountValidateHandler(w http.ResponseWriter, r *http.Request)
	SyncConfigHandler(w http.ResponseWriter, r *http.Request)
	SyncSchemaHandler(w http.ResponseWriter, r *http.Request)
	SyncDataHandler(w http.ResponseWriter, r *http.Request)
}

// WebhookIntegration specifies the required endpoints for a Fibery integration that uses webhooks
type IntegrationWebhooks interface {
	WebhookInitHandler(w http.ResponseWriter, r *http.Request)
	WebhookPreProcessHandler(w http.ResponseWriter, r *http.Request)
	WebhookTransformHandler(w http.ResponseWriter, r *http.Request)
}

type IntegrationOauth1 interface {
	Oauth1AuthorizeHandler(w http.ResponseWriter, r *http.Request)
	Oauth1TokenHandler(w http.ResponseWriter, r *http.Request)
}

type IntegrationOauth2 interface {
	Oauth2AuthorizeHandler(w http.ResponseWriter, r *http.Request)
	Oauth2TokenHandler(w http.ResponseWriter, r *http.Request)
}

type IntegrationLogo interface {
	LogoHandler(w http.ResponseWriter, r *http.Request)
}

type IntegrationDatalist interface {
	SyncDatalistHandler(w http.ResponseWriter, r *http.Request)
}

type IntegrationFilter interface {
	SyncFilterValidateHandler(w http.ResponseWriter, r *http.Request)
}

type IntegrationResource interface {
	SyncResourceHandler(w http.ResponseWriter, r *http.Request)
}

type IntegrationActions interface {
	ActionHandler(w http.ResponseWriter, r *http.Request)
}
