package data

import "github.com/tommyhedley/fibery-quickbooks-app/pkgs/fibery"

var Item = QuickBooksDualType{
	QuickBooksType: QuickBooksType{
		fiberyType: fiberyType{
			id:   "Item",
			name: "Item",
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
				"name": {
					Name: "Name",
					Type: fibery.Text,
				},
				"fully_qualified_name": {
					Name:    "Full Name",
					Type:    fibery.Text,
					SubType: fibery.Title,
				},
				"description": {
					Name:    "Description",
					Type:    fibery.Text,
					SubType: fibery.MD,
				},
				"purchase description": {
					Name: "Purchase Description",
					Type: fibery.Text,
				},
				"inv_start_date": {
					Name:    "Inventory Start",
					Type:    fibery.DateType,
					SubType: fibery.Day,
				},
				"type": {
					Name:     "Type",
					Type:     fibery.Text,
					SubType:  fibery.SingleSelect,
					ReadOnly: true,
					Options: []map[string]any{
						{
							"name": "Inventory",
						},
						{
							"name": "Service",
						},
						{
							"name": "Non-Inventory",
						},
					},
				},
				"qty_on_hand": {
					Name: "Quantity On Hand",
					Type: fibery.Number,
					Format: map[string]any{
						"format":               "Number",
						"hasThousandSeparator": true,
						"precision":            2,
					},
				},
				"reorder_point": {
					Name: "Reorder Quantity",
					Type: fibery.Number,
					Format: map[string]any{
						"format":               "Number",
						"hasThousandSeparator": true,
						"precision":            2,
					},
				},
				"asset_account_id": {
					Name: "Asset Account ID",
					Type: fibery.Text,
					Relation: &fibery.Relation{
						Cardinality:   fibery.MTO,
						Name:          "Asset Account",
						TargetName:    "Items",
						TargetType:    "Account",
						TargetFieldID: "id",
					},
				},
				"expense_account_id": {
					Name: "Expense Account ID",
					Type: fibery.Text,
					Relation: &fibery.Relation{
						Cardinality:   fibery.MTO,
						Name:          "Expense Account",
						TargetName:    "Items",
						TargetType:    "Account",
						TargetFieldID: "id",
					},
				},
				"income_account_id": {
					Name: "Income Account ID",
					Type: fibery.Text,
					Relation: &fibery.Relation{
						Cardinality:   fibery.MTO,
						Name:          "Income Account",
						TargetName:    "Items",
						TargetType:    "Account",
						TargetFieldID: "id",
					},
				},
				"sku": {
					Name: "SKU",
					Type: fibery.Text,
				},
				"sales_tax_included": {
					Name:    "Sales Tax Included",
					Type:    fibery.Text,
					SubType: fibery.Boolean,
				},
				"purchase_tax_included": {
					Name:    "Purchase Tax Included",
					Type:    fibery.Text,
					SubType: fibery.Boolean,
				},
				"sales_tax_code_id": {
					Name: "Sales Tax Code ID",
					Type: fibery.Text,
					Relation: &fibery.Relation{
						Cardinality:   fibery.MTO,
						Name:          "Sales Tax",
						TargetName:    "Sales Tax On Items",
						TargetType:    "TaxCode",
						TargetFieldID: "id",
					},
				},
				"purchase_tax_code_id": {
					Name: "Purchase Tax Code ID",
					Type: fibery.Text,
					Relation: &fibery.Relation{
						Cardinality:   fibery.MTO,
						Name:          "Purchase Tax",
						TargetName:    "Purchase Tax On Items",
						TargetType:    "TaxCode",
						TargetFieldID: "id",
					},
				},
				"class": {
					Name: "Class ID",
					Type: fibery.Text,
					Relation: &fibery.Relation{
						Cardinality:   fibery.MTO,
						Name:          "Class",
						TargetName:    "Expense Account Line(s)",
						TargetType:    "Class",
						TargetFieldID: "id",
					},
				},
				"taxable": {
					Name:    "Taxable",
					Type:    fibery.Text,
					SubType: fibery.Boolean,
				},
				"preferred_vendor_id": {
					Name: "Preferred Vendor ID",
					Type: fibery.Text,
					Relation: &fibery.Relation{
						Cardinality:   fibery.MTO,
						Name:          "Preferred Vendor",
						TargetName:    "Primary Sale Items",
						TargetType:    "Vendor",
						TargetFieldID: "id",
					},
				},
				"parent_id": {
					Name: "Parent ID",
					Type: fibery.Text,
					Relation: &fibery.Relation{
						Cardinality:   fibery.MTO,
						Name:          "Parent",
						TargetName:    "Sub-Items",
						TargetType:    "Item",
						TargetFieldID: "id",
					},
				},
				"purchase_cost": {
					Name: "Purchase Cost",
					Type: fibery.Number,
					Format: map[string]any{
						"format":               "Money",
						"currencyCode":         "USD",
						"hasThousandSeperator": true,
						"precision":            2,
					},
				},
				"unit_price": {
					Name: "Unit Price",
					Type: fibery.Number,
					Format: map[string]any{
						"format":               "Money",
						"currencyCode":         "USD",
						"hasThousandSeperator": true,
						"precision":            2,
					},
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
