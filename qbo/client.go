// Copyright (c) 2018, Randy Westlund. All rights reserved.
// This code is under the BSD-2-Clause license.
package qbo

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"net/url"
	"os"
)

type RequestError struct {
	StatusCode int
	Err        error
	TryLater   bool
}

func (e *RequestError) Error() string {
	return e.Err.Error()
}

func NewRequestError(statusCode int, err error, tryLater bool) *RequestError {
	return &RequestError{
		StatusCode: statusCode,
		TryLater:   tryLater,
		Err:        err,
	}
}

// Client is your handle to the QuickBooks API.
type Client struct {
	// Get this from oauth2.NewClient().
	Client *http.Client
	// Set to ProductionEndpoint or SandboxEndpoint.
	endpoint *url.URL
	// The set of quickbooks APIs
	discoveryAPI *DiscoveryAPI
	// The client Id
	clientId string
	// The client Secret
	clientSecret string
	// The minor version of the QB API
	minorVersion string
	// The account Id you're connecting to.
	realmId string
}

// NewClient initializes a new QuickBooks client for interacting with their Online API
func NewClient(realmId string, token *BearerToken) (c *Client, err error) {
	var clientId, clientSecret string
	var clientEndpoint *url.URL
	mode := os.Getenv("MODE")
	switch mode {
	case "production":
		clientId = os.Getenv("OAUTH_CLIENT_ID_PRODUCTION")
		clientSecret = os.Getenv("OAUTH_CLIENT_SECRET_PRODUCTION")
		clientEndpoint, err = url.Parse(os.Getenv("PRODUCTION_ENDPOINT"))
		if err != nil {
			return nil, fmt.Errorf("failed to parse API endpoint: %v", err)
		}
	case "sandbox":
		clientId = os.Getenv("OAUTH_CLIENT_ID_SANDBOX")
		clientSecret = os.Getenv("OAUTH_CLIENT_SECRET_SANDBOX")
		clientEndpoint, err = url.Parse(os.Getenv("SANDBOX_ENDPOINT"))
		if err != nil {
			return nil, fmt.Errorf("failed to parse API endpoint: %v", err)
		}
	default:
		return nil, fmt.Errorf("invalid mode: %s", mode)
	}

	minorVersion := os.Getenv("MINOR_VERSION")
	if minorVersion == "" {
		return nil, errors.New("minor version is required")
	}

	client := Client{
		clientId:     clientId,
		clientSecret: clientSecret,
		endpoint:     clientEndpoint,
		discoveryAPI: DiscoveryAPIData,
		minorVersion: minorVersion,
		realmId:      realmId,
	}

	if token != nil {
		client.Client = getHttpClient(token)
	}

	return &client, nil
}

// FindAuthorizationUrl compiles the authorization url from the discovery api's auth endpoint.
//
// Example: qbClient.FindAuthorizationUrl("com.intuit.quickbooks.accounting", "security_token", "https://developer.intuit.com/v2/OAuth2Playground/RedirectUrl")
//
// You can find live examples from https://developer.intuit.com/app/developer/playground
func (c *Client) FindAuthorizationUrl(scope string, state string, redirectUri string) (string, error) {
	var authorizationUrl *url.URL

	authorizationUrl, err := url.Parse(c.discoveryAPI.AuthorizationEndpoint)
	if err != nil {
		return "", fmt.Errorf("failed to parse auth endpoint: %v", err)
	}

	urlValues := url.Values{}
	urlValues.Add("client_id", c.clientId)
	urlValues.Add("response_type", "code")
	urlValues.Add("scope", scope)
	urlValues.Add("redirect_uri", redirectUri)
	urlValues.Add("state", state)
	authorizationUrl.RawQuery = urlValues.Encode()

	return authorizationUrl.String(), nil
}

func (c *Client) req(method string, endpoint string, payloadData interface{}, responseObject interface{}, queryParameters map[string]string) error {
	endpointUrl := *c.endpoint
	endpointUrl.Path += "/v3/company/" + c.realmId + "/" + endpoint
	urlValues := url.Values{}

	if len(queryParameters) > 0 {
		for param, value := range queryParameters {
			urlValues.Add(param, value)
		}
	}

	urlValues.Set("minorversion", c.minorVersion)
	urlValues.Encode()
	endpointUrl.RawQuery = urlValues.Encode()

	var err error
	var marshalledJson []byte

	if payloadData != nil {
		marshalledJson, err = json.Marshal(payloadData)
		if err != nil {
			return fmt.Errorf("failed to marshal payload: %v", err)
		}
	}
	slog.Info("Requesting %s %s", method, endpointUrl.String())
	req, err := http.NewRequest(method, endpointUrl.String(), bytes.NewBuffer(marshalledJson))
	if err != nil {
		return fmt.Errorf("failed to create request: %v", err)
	}

	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/json")

	resp, err := c.Client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to make request: %v", err)
	}

	defer resp.Body.Close()

	switch resp.StatusCode {
	case http.StatusOK:
		break
	case http.StatusTooManyRequests:
		return RequestLimit{
			Message:      errors.New("rate limit exceeded"),
			ResetSeconds: 60,
		}
	default:
		return parseFailure(resp)
	}

	if responseObject != nil {
		if err = json.NewDecoder(resp.Body).Decode(&responseObject); err != nil {
			return fmt.Errorf("failed to unmarshal response into object: %v", err)
		}
	}

	return nil
}

func (c *Client) get(endpoint string, responseObject interface{}, queryParameters map[string]string) error {
	return c.req("GET", endpoint, nil, responseObject, queryParameters)
}

func (c *Client) post(endpoint string, payloadData interface{}, responseObject interface{}, queryParameters map[string]string) error {
	return c.req("POST", endpoint, payloadData, responseObject, queryParameters)
}

type QueryResponse[T any] struct {
	Items         []T
	StartPosition int
	MaxResults    int
}

// query makes the specified QBO `query` and unmarshals the result into `responseObject`
func (c *Client) query(query string, responseObject interface{}) error {
	return c.get("query", responseObject, map[string]string{"query": query})
}
