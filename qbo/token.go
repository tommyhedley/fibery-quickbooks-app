package qbo

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"time"

	"golang.org/x/oauth2"
)

type BearerToken struct {
	RefreshToken           string      `json:"refresh_token"`
	AccessToken            string      `json:"access_token"`
	TokenType              string      `json:"token_type"`
	IdToken                string      `json:"id_token"`
	ExpiresIn              json.Number `json:"expires_in"`
	ExpiresOn              string      `json:"expires_on,omitempty"`
	XRefreshTokenExpiresIn json.Number `json:"x_refresh_token_expires_in"`
}

// RefreshToken
// Call the refresh endpoint to generate new tokens
func (c *Client) RefreshToken(refreshToken string) (*BearerToken, error) {
	client := &http.Client{}
	urlValues := url.Values{}
	urlValues.Set("grant_type", "refresh_token")
	urlValues.Add("refresh_token", refreshToken)

	req, err := http.NewRequest("POST", c.discoveryAPI.TokenEndpoint, bytes.NewBufferString(urlValues.Encode()))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded;charset=UTF-8")
	req.Header.Set("Authorization", "Basic "+basicAuth(c))

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New(string(body))
	}

	bearerTokenResponse, err := getBearerTokenResponse(body)
	c.Client = getHttpClient(bearerTokenResponse)

	return bearerTokenResponse, err
}

// RetrieveBearerToken
// Method to retrieve access token (bearer token).
// This method can only be called once
func (c *Client) RetrieveBearerToken(authorizationCode, redirectURI string) (*BearerToken, error) {
	client := &http.Client{}
	urlValues := url.Values{}
	// set parameters
	urlValues.Add("code", authorizationCode)
	urlValues.Set("grant_type", "authorization_code")
	urlValues.Add("redirect_uri", redirectURI)

	req, err := http.NewRequest("POST", c.discoveryAPI.TokenEndpoint, bytes.NewBufferString(urlValues.Encode()))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded;charset=UTF-8")
	req.Header.Set("Authorization", "Basic "+basicAuth(c))

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, parseFailure(resp)
	}

	bearerTokenResponse, err := getBearerTokenResponse(body)

	return bearerTokenResponse, err
}

// RevokeToken
// Call the revoke endpoint to revoke tokens
func (c *Client) RevokeToken(refreshToken string) error {
	client := &http.Client{}
	urlValues := url.Values{}
	urlValues.Add("token", refreshToken)

	req, err := http.NewRequest("POST", c.discoveryAPI.RevocationEndpoint, bytes.NewBufferString(urlValues.Encode()))
	if err != nil {
		return err
	}

	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded;charset=UTF-8")
	req.Header.Set("Authorization", "Basic "+basicAuth(c))

	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK {
		return errors.New(string(body))
	}

	c.Client = nil

	return nil
}

// RefreshNeeded
// Check if the token needs to be refreshed based on the TOKEN_REFRESH_BEFORE_EXPRIRATION environment variable (in seconds)
func (b *BearerToken) RefreshNeeded() (bool, error) {
	refreshBeforeDuration, err := strconv.Atoi(os.Getenv("TOKEN_REFRESH_BEFORE_EXPIRATION"))
	if err != nil {
		return false, fmt.Errorf("unable to parse TOKEN_REFRESH_BEFORE_EXPRIRATION to int: %w", err)
	}

	expiresOn, err := time.Parse(time.RFC3339, b.ExpiresOn)
	if err != nil {
		return false, fmt.Errorf("unable to parse token expiration time: %w", err)
	}

	refreshBefore := time.Now().Add(time.Duration(refreshBeforeDuration) * time.Second)

	return refreshBefore.After(expiresOn), nil
}
func basicAuth(c *Client) string {
	return base64.StdEncoding.EncodeToString([]byte(c.clientId + ":" + c.clientSecret))
}

func getBearerTokenResponse(body []byte) (*BearerToken, error) {
	token := BearerToken{}

	if err := json.Unmarshal(body, &token); err != nil {
		return nil, errors.New(string(body))
	}

	expiresIn, err := token.ExpiresIn.Int64()
	if err != nil {
		return nil, fmt.Errorf("unable to parse expires_in to int64: %w", err)
	}

	token.ExpiresOn = time.Now().UTC().Add(time.Duration(expiresIn) * time.Second).Format(time.RFC3339)

	return &token, nil
}

func getHttpClient(bearerToken *BearerToken) *http.Client {
	ctx := context.Background()
	token := oauth2.Token{
		AccessToken: bearerToken.AccessToken,
		TokenType:   "Bearer",
	}
	return oauth2.NewClient(ctx, oauth2.StaticTokenSource(&token))
}
