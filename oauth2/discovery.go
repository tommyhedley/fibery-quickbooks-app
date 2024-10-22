package oauth2

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
)

type discoveryResponse struct {
	Issuer                string `json:"issuer"`
	AuthorizationEndpoint string `json:"authorization_endpoint"`
	TokenEndpoint         string `json:"token_endpoint"`
	UserinfoEndpoint      string `json:"userinfo_endpoint"`
	RevocationEndpoint    string `json:"revocation_endpoint"`
	JwksUri               string `json:"jwks_uri"`
}

var discoveryParams discoveryResponse

func GetDiscovery() error {
	baseURL, err := url.Parse(os.Getenv("DISCOVERY_URL"))
	if err != nil {
		return fmt.Errorf("error parsing base url: %w", err)
	}

	req, err := http.NewRequest("GET", baseURL.String(), nil)
	if err != nil {
		return fmt.Errorf("error creating request: %w", err)
	}

	req.Header.Set("accept", "application/json")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("error executing request: %w", err)
	}

	defer res.Body.Close()

	if res.StatusCode > 299 {
		return fmt.Errorf("request error: %d", res.StatusCode)
	}

	decoder := json.NewDecoder(res.Body)
	err = decoder.Decode(&discoveryParams)
	if err != nil {
		return fmt.Errorf("unable to decode response: %w", err)
	}
	return nil
}
