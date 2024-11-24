package sync

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/patrickmn/go-cache"
	"github.com/tommyhedley/fibery/fibery-qbo-integration/internal/utils"
	"golang.org/x/sync/singleflight"
)

type SyncType string

const (
	DeltaSync SyncType = "delta"
	FullSync  SyncType = "full"
)

type DataHandlerRequest struct {
	RequestedType string         `json:"requestedType"`
	OperationID   string         `json:"operationId"`
	Types         []string       `json:"types"`
	Filter        map[string]any `json:"filter"`
	Account       struct {
		AccessToken string `json:"access_token"`
	} `json:"account"`
	LastSyncronized string                               `json:"lastSynchronizedAt"`
	Pagination      NextPageConfig                       `json:"pagination"`
	Schema          map[string]map[string]map[string]any `json:"schema"`
}

type NextPageConfig struct {
	StartPosition int `json:"startPosition"`
	MaxResults    int `json:"maxResults"`
}

type Pagination struct {
	HasNext        bool           `json:"hasNext"`
	NextPageConfig NextPageConfig `json:"nextPageConfig"`
}

type DataHandlerResponse struct {
	Items               []map[string]any `json:"items"`
	Pagination          Pagination       `json:"pagination"`
	SynchronizationType SyncType         `json:"synchronizationType"`
}

type RequestParameters struct {
	Cache         *cache.Cache
	Group         *singleflight.Group
	StartPosition int
	MaxResults    int
	Token         string
	LastSynced    string
	Filter        map[string]any
}

type DataRequest func(RequestParameters) ([]map[string]any, bool, error)

type DataType struct {
	ID     string
	Name   string
	Schema map[string]Field
	DataRequest
}

var Types = []TypeArray{}
var Schema = make(map[string]map[string]Field)
var DataRequests = make(map[string]*DataRequest)

func (dt *DataType) Register() {
	Types = append(Types, TypeArray{
		ID:   dt.ID,
		Name: dt.Name,
	})
	Schema[dt.ID] = dt.Schema
	DataRequests[dt.ID] = &dt.DataRequest
}

func GetRequestFunctions(id string) (*DataRequest, bool) {
	dr, exists := DataRequests[id]
	return dr, exists
}

func DataHandler(c *cache.Cache, group *singleflight.Group) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		decoder := json.NewDecoder(r.Body)
		params := DataHandlerRequest{}
		err := decoder.Decode(&params)
		if err != nil {
			utils.RespondWithError(w, http.StatusInternalServerError, fmt.Errorf("unable to decode request parameters: %w", err))
			return
		}

		startPosition := params.Pagination.StartPosition
		maxResults := params.Pagination.MaxResults
		reqType := params.RequestedType
		opID := params.OperationID

		if reqType == "" || opID == "" {
			utils.RespondWithError(w, http.StatusBadRequest, fmt.Errorf("request parameters missing type: %s and/or operation ID (%s", reqType, opID))
			return
		}

		rf, exists := GetRequestFunctions(reqType)
		if !exists {
			utils.RespondWithError(w, http.StatusBadRequest, fmt.Errorf("requested type was not found: %s", reqType))
			return
		}

		if startPosition == 0 {
			startPosition = 1
		}

		if maxResults == 0 {
			maxResults = 1000
		}

		req := RequestParameters{
			Cache:         c,
			Group:         group,
			StartPosition: startPosition,
			MaxResults:    maxResults,
			Token:         params.Account.AccessToken,
			LastSynced:    params.LastSyncronized,
			Filter:        params.Filter,
		}

		items, more, err := (*rf)(req)
		if err != nil {
			utils.RespondWithError(w, http.StatusBadRequest, fmt.Errorf("unable to retreive full sync data"))
			return
		}
		resp := DataHandlerResponse{
			Items: items,
			Pagination: Pagination{
				HasNext: more,
				NextPageConfig: NextPageConfig{
					StartPosition: startPosition + maxResults,
					MaxResults:    maxResults,
				},
			},
			SynchronizationType: FullSync,
		}
		utils.RespondWithJSON(w, http.StatusOK, resp)
	}
}
