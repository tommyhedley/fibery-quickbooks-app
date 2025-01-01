package qbo

import (
	"fmt"
	"strings"
	"time"
)

type ChangeDataCapture struct {
	CDCResponse []struct {
		QueryResponse []struct {
			Customer      []Customer           `json:",omitempty"`
			Invoice       []CDCInvoiceResponse `json:",omitempty"`
			Purchase      []Purchase           `json:",omitempty"`
			StartPosition int                  `json:"startPosition"`
			MaxResults    int                  `json:"maxResults"`
			TotalCount    int                  `json:"totalCount,omitempty"`
		} `json:"QueryResponse"`
	} `json:"CDCResponse"`
	Time string `json:"time"`
}

func (c *Client) ChangeDataCapture(entities []string, changedSince time.Time) (ChangeDataCapture, error) {
	var res ChangeDataCapture

	queryParams := map[string]string{
		"entities":     strings.Join(entities, ","),
		"changedSince": changedSince.Format(qboDateFormat),
	}

	err := c.req("GET", "/cdc", nil, &res, queryParams)
	if err != nil {
		return ChangeDataCapture{}, fmt.Errorf("failed to make cdc request: %w", err)
	}
	return res, nil
}

func getChangeDataCapture(req *DataRequest) (ChangeDataCapture, error) {
	lastSyncedTime, err := time.Parse(time.RFC3339, req.LastSynced)
	if err != nil {
		return ChangeDataCapture{}, fmt.Errorf("unable to parse last synced time: %w", err)
	}

	client, err := NewClient(req.RealmID, req.Token)
	if err != nil {
		return ChangeDataCapture{}, fmt.Errorf("unable to create new client: %w", err)
	}

	cdc, err := client.ChangeDataCapture(req.CDCTypes, lastSyncedTime)
	if err != nil {
		return ChangeDataCapture{}, fmt.Errorf("unable to get change data capture: %w", err)
	}

	// consider checking if response count meets the limit (1000 entities) and returning an error or forcing a full sync

	return cdc, nil
}
