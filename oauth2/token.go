package oauth2

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/tommyhedley/fibery/fibery-qbo-integration/internal/utils"
	"github.com/tommyhedley/fibery/fibery-qbo-integration/qbo"
)

func TokenHandler(w http.ResponseWriter, r *http.Request) {
	type requestBody struct {
		Fields struct {
			ID          string `json:"_id"`
			App         string `json:"app"`
			Owner       string `json:"owner"`
			ExpireOn    string `json:"expireOn"`
			CallbackURI string `json:"callback_uri"`
			State       string `json:"state"`
		} `json:"fields"`
		Code    string `json:"code"`
		State   string `json:"state"`
		RealmID string `json:"realmId"`
	}

	type responseBody struct {
		RealmID string          `json:"realmId"`
		Token   qbo.BearerToken `json:"token"`
	}

	decoder := json.NewDecoder(r.Body)
	reqBody := requestBody{}
	err := decoder.Decode(&reqBody)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, fmt.Errorf("unable to decode request parameters: %w", err))
		return
	}

	// Temporary workaround until Fibery authorization response supports dynamic parameters
	realmId := reqBody.RealmID
	if realmId == "" {
		mode := os.Getenv("MODE")
		switch mode {
		case "production":
			realmId = os.Getenv("REALM_ID_PRODUCTION")
		case "sandbox":
			realmId = os.Getenv("REALM_ID_SANDBOX")
		default:
			utils.RespondWithError(w, http.StatusInternalServerError, fmt.Errorf("invalid mode: %s", mode))
		}
	}
	client, err := qbo.NewClient(realmId, nil)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, fmt.Errorf("unable to create new client: %w", err))
		return
	}

	token, err := client.RetrieveBearerToken(reqBody.Code, reqBody.Fields.CallbackURI)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, fmt.Errorf("unable to retreive bearer token: %w", err))
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, responseBody{
		RealmID: realmId,
		Token:   (*token),
	})
}

func ValidateHandler(w http.ResponseWriter, r *http.Request) {
	type requestBody struct {
		Id     string `json:"id"`
		Fields struct {
			Name    string          `json:"name"`
			RealmID string          `json:"realmId"`
			Token   qbo.BearerToken `json:"token"`
		} `json:"fields"`
	}

	type responseBody struct {
		Name    string          `json:"name"`
		RealmID string          `json:"realmId"`
		Token   qbo.BearerToken `json:"token"`
	}

	decoder := json.NewDecoder(r.Body)
	reqBody := requestBody{}
	err := decoder.Decode(&reqBody)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, fmt.Errorf("unable to decode request parameters: %w", err))
		return
	}

	client, err := qbo.NewClient(reqBody.Fields.RealmID, &reqBody.Fields.Token)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, fmt.Errorf("unable to create new client: %w", err))
		return
	}

	token := reqBody.Fields.Token

	refreshNeeded, err := token.RefreshNeeded()
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, fmt.Errorf("unable to determine if token refresh is needed: %w", err))
		return
	}

	if refreshNeeded {
		newToken, err := client.RefreshToken(token.RefreshToken)
		if err != nil {
			utils.RespondWithError(w, http.StatusInternalServerError, fmt.Errorf("unable to refresh token: %w", err))
			return
		}
		token = *newToken
	}

	info, err := client.FindCompanyInfo()
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, fmt.Errorf("unable to find company info: %w", err))
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, responseBody{
		Name:    info.CompanyName,
		RealmID: reqBody.Fields.RealmID,
		Token:   token,
	})
}
