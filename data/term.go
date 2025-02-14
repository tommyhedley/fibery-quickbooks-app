package data

import "github.com/tommyhedley/fibery-quickbooks-app/pkgs/fibery"

var term = QuickBooksDualType{
	QuickBooksType: QuickBooksType{
		fiberyType: fiberyType{
			id:   "Term",
			name: "Term",
			schema: map[string]fibery.Field{
				"id": {
					Name: "id",
					Type: fibery.Text,
				},
				"name": {
					Name: "Name",
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
