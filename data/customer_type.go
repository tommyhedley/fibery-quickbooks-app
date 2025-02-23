package data

import (
	"fmt"

	"github.com/tommyhedley/fibery-quickbooks-app/pkgs/fibery"
	"github.com/tommyhedley/quickbooks-go"
)

var CustomerType = QuickBooksCDCType{
	QuickBooksType: QuickBooksType{
		fiberyType: fiberyType{
			id:   "CustomerType",
			name: "Customer Type",
			schema: map[string]fibery.Field{
				"Id": {
					Name: "ID",
					Type: fibery.Text,
				},
				"QBOId": {
					Name: "QBO ID",
					Type: fibery.Text,
				},
				"Name": {
					Name: "Name",
					Type: fibery.Text,
				},
				"SyncToken": {
					Name:     "Sync Token",
					Type:     fibery.Text,
					ReadOnly: true,
				},
				"__syncAction": {
					Type: fibery.Text,
					Name: "Sync Action",
				},
				"Active": {
					Name:    "Active",
					Type:    fibery.Text,
					SubType: fibery.Boolean,
				},
			},
		},
		schemaGen: func(entity any) (map[string]any, error) {
			customerType, ok := entity.(quickbooks.CustomerType)
			if !ok {
				return nil, fmt.Errorf("unable to convert entity to CustomerType")
			}

			return map[string]any{
				"Id":           customerType.Id,
				"QBOId":        customerType.Id,
				"Name":         customerType.Name,
				"SyncToken":    customerType.SyncToken,
				"__syncAction": fibery.SET,
				"Active":       customerType.Active,
			}, nil
		},
		query:          func(req Request) (Response, error) {},
		queryProcessor: func(entityArray any, schemaGen schemaGenFunc) ([]map[string]any, error) {},
	},
	cdcProcessor: func(cdc quickbooks.ChangeDataCapture, schemaGen schemaGenFunc) ([]map[string]any, error) {},
}
