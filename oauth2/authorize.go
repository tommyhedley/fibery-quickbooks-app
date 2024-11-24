package oauth2

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/tommyhedley/fibery/fibery-qbo-integration/internal/utils"
	"github.com/tommyhedley/fibery/fibery-qbo-integration/qbo"
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

	client, err := qbo.NewClient("", nil)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, fmt.Errorf("error creating new client: %w", err))
		return
	}

	redirectURI, err := client.FindAuthorizationUrl(os.Getenv("SCOPE"), reqBody.State, reqBody.CallbackURI)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, fmt.Errorf("error generating redirect uri: %w", err))
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, responseBody{
		RedirectURI: redirectURI,
	})
}
