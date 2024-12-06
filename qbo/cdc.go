package qbo

type ChangeDataCapture struct {
	CDCResponse []struct {
		QueryResponse []struct {
			Customer []struct {
				ID string `json:"Id"`
			} `json:"Customer,omitempty"`
			StartPosition int `json:"startPosition"`
			MaxResults    int `json:"maxResults"`
			TotalCount    int `json:"totalCount,omitempty"`
			Estimate      []struct {
				ID       string `json:"Id"`
				Status   string `json:"status,omitempty"`
				Domain   string `json:"domain,omitempty"`
				MetaData struct {
					LastUpdatedTime string `json:"LastUpdatedTime"`
				} `json:"MetaData,omitempty"`
			} `json:"Estimate,omitempty"`
		} `json:"QueryResponse"`
	} `json:"CDCResponse"`
	Time string `json:"time"`
}

func (c *Client) ChangeDataCapture() (*ChangeDataCapture, error) {
	return &ChangeDataCapture{}, nil
}
