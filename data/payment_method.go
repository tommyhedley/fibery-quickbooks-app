package data

var PaymentMethod = QuickBooksDualType{
	QuickBooksType: QuickBooksType{
		fiberyType: fiberyType{
			id:   "PaymentMethod",
			name: "Payment Method",
			schema: map[string]fibery.Field{
				"id": {
					Name: "id",
					Type: fibery.ID,
				},
				"qbo_id": {
					Name: "QBO ID",
					Type: fibery.Text,
				},
				"sync_token": {
					Name:     "Sync Token",
					Type:     fibery.Text,
					ReadOnly: true,
				},
				"active": {
					Name:    "Active",
					Type:    fibery.Text,
					SubType: fibery.Boolean,
				},
				"type": {
					Name:     "Type",
					Type:     fibery.Text,
					SubType:  fibery.SingleSelect,
					ReadOnly: true,
					Options: []map[string]any{
						{
							"name": "Credit Card",
						},
						{
							"name": "Non Credit Card",
						},
					},
				},
				"__syncAction": {
					Type: fibery.Text,
					Name: "Sync Action",
				},
			},
		},
		schemaGen:      func(entity any) (map[string]any, error) {},
		query:          func(req Request) (Response, error) {},
		queryProcessor: func(entityArray any, schemaGen schemaGenFunc) ([]map[string]any, error) {},
	},
	cdcProcessor: func(cdc quickbooks.ChangeDataCapture, schemaGen schemaGenFunc) ([]map[string]any, error) {},
	whBatchProcessor: func(itemResponse quickbooks.BatchItemResponse, response *map[string][]map[string]any, cache *cache.Cache, realmId string, queryProcessor queryProcessorFunc, schemaGen schemaGenFunc, typeId string) error {
	},
}
