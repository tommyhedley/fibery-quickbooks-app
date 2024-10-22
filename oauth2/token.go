package oauth2

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/tommyhedley/fibery/fibery-tsheets-integration/internal/utils"
)

type baseTokenRequest struct {
	GrantType    string `url:"grant_type"`
	ClientId     string `url:"client_id"`
	ClientSecret string `url:"client_secret"`
}

type accessTokenRequest struct {
	baseTokenRequest
	Code        string `url:"code"`
	RedirectURI string `url:"redirect_uri"`
}

type refreshTokenRequest struct {
	baseTokenRequest
	RefreshToken string `url:"refresh_token"`
}

type accessTokenResponse struct {
	AccessToken            string `json:"access_token"`
	RefreshToken           string `json:"refresh_token"`
	TokenType              string `json:"token_type"`
	IdToken                string `json:"id_token"`
	ExpiresIn              int    `json:"expires_in"`
	XRefreshTokenExpiresIn int    `json:"x_refresh_token_expires_in"`
}

type tokenHandlerResponse struct {
	accessTokenResponse
	ExpiresOn string `json:"expires_on"`
}

type address struct {
	StreetAddress string `json:"streetAddress"`
	Locality      string `json:"locality"`
	Region        string `json:"region"`
	PostalCode    string `json:"postalCode"`
	Country       string `json:"country"`
}

type userInfoResponse struct {
	Sub                 string  `json:"sub"`
	Email               string  `json:"email"`
	EmailVerified       bool    `json:"emailVerified"`
	GivenName           string  `json:"givenName"`
	FamilyName          string  `json:"familyName"`
	PhoneNumber         string  `json:"phoneNumber"`
	PhoneNumberVerified bool    `json:"phoneNumberVerified"`
	Address             address `json:"address"`
}

var currentAccessToken string

func (params *accessTokenRequest) get(URL string) (accessTokenResponse, error) {
	baseURL, err := url.Parse(URL)
	if err != nil {
		return accessTokenResponse{}, fmt.Errorf("error parsing base url: %w", err)
	}

	body := url.Values{}
	body.Set("grant_type", params.GrantType)
	body.Add("code", params.Code)
	body.Add("redirect_uri", params.RedirectURI)

	req, err := http.NewRequest("POST", baseURL.String(), strings.NewReader(body.Encode()))
	if err != nil {
		return accessTokenResponse{}, fmt.Errorf("error creating request: %w", err)
	}

	auth := base64.StdEncoding.EncodeToString([]byte(params.ClientId + ":" + params.ClientSecret))

	req.Header.Set("accept", "application/json")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded;charset=UTF-8")
	req.Header.Set("Authorization", "Basic "+auth)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return accessTokenResponse{}, fmt.Errorf("error executing request: %w", err)
	}

	defer res.Body.Close()

	if res.StatusCode > 299 {
		return accessTokenResponse{}, fmt.Errorf("request error: %d", res.StatusCode)
	}

	decoder := json.NewDecoder(res.Body)
	var resp accessTokenResponse
	err = decoder.Decode(&resp)
	if err != nil {
		return accessTokenResponse{}, fmt.Errorf("unable to decode response: %w", err)
	}
	return resp, nil
}

func (params *refreshTokenRequest) refresh(URL string) (accessTokenResponse, error) {
	baseURL, err := url.Parse(URL)
	if err != nil {
		return accessTokenResponse{}, fmt.Errorf("error parsing base url: %w", err)
	}

	body := url.Values{}
	body.Set("grant_type", params.GrantType)
	body.Add("refresh_token", params.RefreshToken)

	req, err := http.NewRequest("POST", baseURL.String(), strings.NewReader(body.Encode()))
	if err != nil {
		return accessTokenResponse{}, fmt.Errorf("error creating request: %w", err)
	}

	auth := base64.StdEncoding.EncodeToString([]byte(params.ClientId + ":" + params.ClientSecret))

	req.Header.Set("accept", "application/json")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded;charset=UTF-8")
	req.Header.Set("Authorization", "Basic "+auth)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return accessTokenResponse{}, fmt.Errorf("error executing request: %w", err)
	}

	defer res.Body.Close()

	if res.StatusCode > 299 {
		return accessTokenResponse{}, fmt.Errorf("request error: %d", res.StatusCode)
	}

	decoder := json.NewDecoder(res.Body)
	var resp accessTokenResponse
	err = decoder.Decode(&resp)
	if err != nil {
		return accessTokenResponse{}, fmt.Errorf("unable to decode response: %w", err)
	}
	return resp, nil
}

func TokenHandler(w http.ResponseWriter, r *http.Request) {
	type requestBody struct {
		Fields struct {
			CallbackURI string `json:"callback_uri"`
		} `json:"fields"`
		Code string `json:"code"`
	}

	decoder := json.NewDecoder(r.Body)
	reqBody := requestBody{}
	err := decoder.Decode(&reqBody)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, fmt.Errorf("unable to decode request parameters: %w", err))
		return
	}

	tokenRequest := accessTokenRequest{
		baseTokenRequest: baseTokenRequest{
			GrantType:    "authorization_code",
			ClientId:     os.Getenv("OAUTH_CLIENT_ID"),
			ClientSecret: os.Getenv("OAUTH_CLIENT_SECRET"),
		},
		Code:        reqBody.Code,
		RedirectURI: reqBody.Fields.CallbackURI,
	}

	tokenResponse, err := tokenRequest.get(discoveryParams.TokenEndpoint)
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, fmt.Errorf("error with access token request: %w", err))
		return
	}

	currentAccessToken = tokenResponse.AccessToken

	utils.RespondWithJSON(w, http.StatusOK, tokenHandlerResponse{
		accessTokenResponse: tokenResponse,
		ExpiresOn:           time.Now().UTC().Add(time.Duration(tokenResponse.ExpiresIn) * time.Second).Format(time.RFC3339),
	})
}

func ValidateHandler(w http.ResponseWriter, r *http.Request) {
	type requestBody struct {
		Id     string `json:"id"`
		Fields struct {
			Name string `json:"name"`
			tokenHandlerResponse
		} `json:"fields"`
	}
	type responseBody struct {
		Name string `json:"name"`
		tokenHandlerResponse
	}

	decoder := json.NewDecoder(r.Body)
	reqBody := requestBody{}
	err := decoder.Decode(&reqBody)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, fmt.Errorf("unable to decode request parameters: %w", err))
		return
	}

	refreshNeeded, err := refreshNeeded(reqBody.Fields.ExpiresOn, 2)
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, fmt.Errorf("error checking token expiration: %w", err))
		return
	}

	if refreshNeeded {
		requestParams := refreshTokenRequest{
			baseTokenRequest: baseTokenRequest{
				GrantType:    "refresh_token",
				ClientId:     os.Getenv("OAUTH_CLIENT_ID"),
				ClientSecret: os.Getenv("OAUTH_CLIENT_SECRET"),
			},
			RefreshToken: reqBody.Fields.RefreshToken,
		}
		refreshResponse, err := requestParams.refresh(discoveryParams.TokenEndpoint)
		if err != nil {
			utils.RespondWithError(w, http.StatusBadRequest, fmt.Errorf("error with refresh token request: %w", err))
		}
		currentUser, err := validate(discoveryParams.UserinfoEndpoint, refreshResponse.AccessToken)
		if err != nil {
			utils.RespondWithError(w, http.StatusBadRequest, fmt.Errorf("token validation error: %w", err))
		}

		currentAccessToken = refreshResponse.AccessToken

		utils.RespondWithJSON(w, http.StatusOK, responseBody{
			Name: currentUser.Email,
			tokenHandlerResponse: tokenHandlerResponse{
				accessTokenResponse: refreshResponse,
				ExpiresOn:           time.Now().UTC().Add(time.Duration(refreshResponse.ExpiresIn) * time.Second).Format(time.RFC3339),
			},
		})
		return
	}

	currentUser, err := validate(discoveryParams.UserinfoEndpoint, reqBody.Fields.AccessToken)
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, fmt.Errorf("token validation error: %w", err))
	}
	utils.RespondWithJSON(w, http.StatusOK, responseBody{
		Name: currentUser.Email,
	})
}

func validate(URL, token string) (userInfoResponse, error) {
	baseURL, err := url.Parse(URL)
	if err != nil {
		return userInfoResponse{}, fmt.Errorf("error parsing base url: %w", err)
	}

	req, err := http.NewRequest("GET", baseURL.String(), nil)
	if err != nil {
		return userInfoResponse{}, fmt.Errorf("error creating request: %w", err)
	}

	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", token))

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return userInfoResponse{}, fmt.Errorf("error executing request: %w", err)
	}

	defer res.Body.Close()

	decoder := json.NewDecoder(res.Body)
	var resp userInfoResponse
	err = decoder.Decode(&resp)
	if err != nil {
		return userInfoResponse{}, fmt.Errorf("unable to decode response: %w", err)
	}

	return resp, nil
}

func refreshNeeded(expiresOn string, hoursToRefresh int) (bool, error) {
	expiration, err := time.Parse(time.RFC3339, expiresOn)
	if err != nil {
		return false, fmt.Errorf("unable to parse token expiration time: %w", err)
	}
	deadline := expiration.Add(time.Duration(hoursToRefresh) * time.Hour)
	return time.Now().UTC().After(deadline), nil
}
