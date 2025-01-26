package qbo

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

type DiscoveryAPI struct {
	Issuer                string `json:"issuer"`
	AuthorizationEndpoint string `json:"authorization_endpoint"`
	TokenEndpoint         string `json:"token_endpoint"`
	UserinfoEndpoint      string `json:"userinfo_endpoint"`
	RevocationEndpoint    string `json:"revocation_endpoint"`
	JwksUri               string `json:"jwks_uri"`
}

var DiscoveryAPIData *DiscoveryAPI

// CallDiscoveryAPI
// See https://developer.intuit.com/app/developer/qbo/docs/develop/authentication-and-authorization/openid-connect#discovery-document
func CallDiscoveryAPI(discoveryEndpoint EndpointUrl) (*DiscoveryAPI, error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", string(discoveryEndpoint), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create req: %v", err)
	}

	req.Header.Set("Accept", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to make req: %v", err)
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read body: %v", err)
	}

	respData := DiscoveryAPI{}
	if err = json.Unmarshal(body, &respData); err != nil {
		return nil, fmt.Errorf("error getting DiscoveryAPIResponse: %v", err)
	}

	return &respData, nil
}

func init() {
	godotenv.Load()
	mode := os.Getenv("MODE")
	var discoverURL string
	switch mode {
	case "production":
		discoverURL = os.Getenv("DISCOVERY_PRODUCTION_ENDPOINT")
	case "sandbox":
		discoverURL = os.Getenv("DISCOVERY_SANDBOX_ENDPOINT")
	default:
		log.Fatalf("invalid mode: %s", mode)
	}
	var err error
	DiscoveryAPIData, err = CallDiscoveryAPI(EndpointUrl(discoverURL))
	if err != nil {
		log.Fatalf("error calling production discovery api: %v", err)
	}
}
