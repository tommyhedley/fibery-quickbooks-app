package qbo

import (
	"fmt"
	"strings"
	"time"
)

type ChangeDataCapture struct {
	CDCResponse []struct {
		QueryResponse []struct {
			Invoice       []CDCInvoice `json:"Invoice,omitempty"`
			StartPosition int          `json:"startPosition"`
			MaxResults    int          `json:"maxResults"`
			TotalCount    int          `json:"totalCount,omitempty"`
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
