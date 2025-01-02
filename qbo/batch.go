package qbo

import "fmt"

type BatchFault struct {
	Message string
	Code    string `json:"code"`
	Detail  string
	Element string `json:"element"`
}

type BatchFaultResponse struct {
	FaultType string       `json:"type"`
	Faults    []BatchFault `json:"Error"`
}

type BatchItemRequest struct {
	BID                string `json:"bId"`
	OptionsData        string `json:"optionsData,omitempty"`
	Operation          string `json:"operation,omitempty"`
	Query              string `json:",omitempty"`
	QuickbooksDataType `json:",omitempty"`
}

type BatchItemResponse struct {
	BID                string `json:"bId"`
	QuickbooksDataType `json:",omitempty"`
	Fault              BatchFaultResponse `json:",omitempty"`
	QueryResponse
}

func (c *Client) BatchRequest(items []BatchItemRequest) ([]BatchItemResponse, error) {
	if len(items) == 0 {
		return nil, nil
	}

	var allResponses []BatchItemResponse

	// each BatchRequest is limited to 30 items
	chunkSize := 30
	for start := 0; start < len(items); start += chunkSize {
		end := start + chunkSize
		if end > len(items) {
			end = len(items)
		}
		batch := items[start:end]

		var req struct {
			BatchItemRequest []BatchItemRequest `json:"BatchItemRequest"`
		}

		var res struct {
			BatchItemResponse []BatchItemResponse `json:"BatchItemResponse"`
		}

		req.BatchItemRequest = batch

		err := c.req("POST", "/batch", req, &res, nil)
		if err != nil {
			return nil, fmt.Errorf("failed to make batch request: %w", err)
		}

		allResponses = append(allResponses, res.BatchItemResponse...)
	}

	return allResponses, nil
}
