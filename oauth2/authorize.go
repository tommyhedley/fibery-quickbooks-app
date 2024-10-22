package oauth2

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"

	"github.com/tommyhedley/fibery/fibery-tsheets-integration/internal/utils"
)

func AuthorizeHandler(w http.ResponseWriter, r *http.Request) {
	type requestBody struct {
		CallbackURI string `json:"callback_uri"`
		State       string `json:"state"`
	}
	type responseBody struct {
		RedirectURI string `json:"redirect_uri"`
	}
	decoder := json.NewDecoder(r.Body)
	reqBody := requestBody{}
	err := decoder.Decode(&reqBody)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, fmt.Errorf("unable to decode request parameters: %w", err))
		return
	}

	redirectURI, err := url.Parse(discoveryParams.AuthorizationEndpoint)
	if err != nil {
		utils.RespondWithError(w, http.StatusUnauthorized, fmt.Errorf("error parsing base url: %w", err))
		return
	}

	parameters := url.Values{}
	parameters.Add("client_id", os.Getenv("OAUTH_CLIENT_ID"))
	parameters.Add("response_type", "code")
	parameters.Add("scope", "com.intuit.quickbooks.accounting openid email")
	parameters.Add("redirect_uri", reqBody.CallbackURI)
	parameters.Add("state", reqBody.State)

	redirectURI.RawQuery = parameters.Encode()

	utils.RespondWithJSON(w, http.StatusOK, responseBody{
		RedirectURI: redirectURI.String(),
	})
}
