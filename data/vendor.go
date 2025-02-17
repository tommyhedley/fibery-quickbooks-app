package data

import "github.com/tommyhedley/fibery-quickbooks-app/pkgs/fibery"

var Vendor = QuickBooksDualType{
	QuickBooksType: QuickBooksType{
		fiberyType: fiberyType{
			id:   "Vendor",
			name: "Vendor",
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
				"display_name": {
					Name:    "Display Name",
					Type:    fibery.Text,
					SubType: fibery.Title,
				},
				"title": {
					Name: "Title",
					Type: fibery.Text,
				},
				"given_name": {
					Name: "First Name",
					Type: fibery.Text,
				},
				"middle_name": {
					Name: "Middle Name",
					Type: fibery.Text,
				},
				"family_name": {
					Name: "Last Name",
					Type: fibery.Text,
				},
				"suffix": {
					Name: "Suffix",
					Type: fibery.Text,
				},
				"company_name": {
					Name: "Company Name",
					Type: fibery.Text,
				},
				"primary_email": {
					Name:    "Email",
					Type:    fibery.Text,
					SubType: fibery.Email,
				},
				"sales_term_id": {
					Name: "Sales Term ID",
					Type: fibery.Text,
					Relation: &fibery.Relation{
						Cardinality:   fibery.MTO,
						Name:          "Sales Term",
						TargetName:    "Customers",
						TargetType:    "SalesTerm",
						TargetFieldID: "id",
					},
				},
				"primary_phone": {
					Name: "Phone",
					Type: fibery.Text,
					Format: map[string]any{
						"format": "phone",
					},
				},
				"alt_phone": {
					Name: "Alternate Phone",
					Type: fibery.Text,
					Format: map[string]any{
						"format": "phone",
					},
				},
				"mobile": {
					Name: "Mobile",
					Type: fibery.Text,
					Format: map[string]any{
						"format": "phone",
					},
				},
				"fax": {
					Name: "Fax",
					Type: fibery.Text,
					Format: map[string]any{
						"format": "phone",
					},
				},
				"1099": {
					Name:        "1099",
					Type:        fibery.Text,
					SubType:     fibery.Boolean,
					Description: "Is the Vendor a 1099 contractor?",
				},
				"cost_rate": {
					Name:        "Cost Rate",
					Type:        fibery.Number,
					Description: "Default cost rate of the Vendor",
					Format: map[string]any{
						"format":               "Money",
						"currencyCode":         "USD",
						"hasThousandSeperator": true,
						"precision":            2,
					},
				},
				"bill_rate": {
					Name:        "Bill Rate",
					Type:        fibery.Number,
					Description: "Default billing rate of the Vendor",
					Format: map[string]any{
						"format":               "Money",
						"currencyCode":         "USD",
						"hasThousandSeperator": true,
						"precision":            2,
					},
				},
				"website": {
					Name:    "Website",
					Type:    fibery.Text,
					SubType: fibery.URL,
				},
				"account_number": {
					Name:        "Account Number",
					Type:        fibery.Text,
					Description: "Name or number of the account associated with this vendor",
				},
				"balance": {
					Name: "Balance",
					Type: fibery.Number,
					Format: map[string]any{
						"format":               "Money",
						"currencyCode":         "USD",
						"hasThousandSeperator": true,
						"precision":            2,
					},
				},
				"billing_address": {
					Name:    "Billing Address",
					Type:    fibery.Text,
					SubType: fibery.MD,
				},
				"billing_line_1": {
					Name: "Billing Line 1",
					Type: fibery.Text,
				},
				"billing_line_2": {
					Name: "Billing Line 2",
					Type: fibery.Text,
				},
				"billing_line_3": {
					Name: "Billing Line 3",
					Type: fibery.Text,
				},
				"billing_line_4": {
					Name: "Billing Line 4",
					Type: fibery.Text,
				},
				"billing_line_5": {
					Name: "Billing Line 5",
					Type: fibery.Text,
				},
				"billing_city": {
					Name: "Billing City",
					Type: fibery.Text,
				},
				"billing_state": {
					Name: "Billing State",
					Type: fibery.Text,
				},
				"billing_postal_code": {
					Name: "Billing Postal Code",
					Type: fibery.Text,
				},
				"billing_country": {
					Name: "Billing Country",
					Type: fibery.Text,
				},
				"billing_lat": {
					Name: "Billing Latitude",
					Type: fibery.Text,
				},
				"billing_long": {
					Name: "Billing Longitude",
					Type: fibery.Text,
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
