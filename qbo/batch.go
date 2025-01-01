package qbo

type BatchItemRequest struct {
	BID                string `json:"bId"`
	OptionsData        string `json:"optionsData,omitempty"`
	Operation          string `json:"operation,omitempty"`
	Query              string `json:",omitempty"`
	QuickbooksDataType `json:",omitempty"`
}
