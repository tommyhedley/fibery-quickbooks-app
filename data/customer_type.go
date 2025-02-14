package data

import (
	"github.com/tommyhedley/fibery-quickbooks-app/pkgs/fibery"
	"github.com/tommyhedley/quickbooks-go"
)

var CustomerType = QuickBooksCDCType{
	QuickBooksType: QuickBooksType{
		fiberyType: fiberyType{
			id:   "CustomerType",
			name: "Customer Type",
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
}
