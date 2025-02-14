package data

import (
	"github.com/tommyhedley/fibery-quickbooks-app/pkgs/fibery"
)

var TaxRate = QuickBooksType{
	fiberyType: fiberyType{
		id:   "TaxRate",
		name: "Tax Rate",
		schema: map[string]fibery.Field{
			"id": {
				Name: "id",
				Type: fibery.ID,
			},
			"name": {
				Name: "Name",
				Type: fibery.Text,
			},
			"description": {
				Name: "Description",
				Type: fibery.Text,
			},
			"sync_token": {
				Name:     "Sync Token",
				Type:     fibery.Text,
				ReadOnly: true,
			},
			"rate_value": {
				Name: "Rate",
				Type: fibery.Number,
				Format: map[string]any{
					"format":    "Percent",
					"precision": 2,
				},
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
}
