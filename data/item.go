package data

import (
	"fmt"

	"github.com/tommyhedley/fibery-quickbooks-app/pkgs/fibery"
	"github.com/tommyhedley/quickbooks-go"
)

var Item = QuickBooksDualType{
	QuickBooksType: QuickBooksType{
		fiberyType: fiberyType{
			id:   "Item",
			name: "Item",
			schema: map[string]fibery.Field{
				"Id": {
					Name: "ID",
					Type: fibery.ID,
				},
				"QBOId": {
					Name: "QBO ID",
					Type: fibery.Text,
				},
				"Name": {
					Name: "Name",
					Type: fibery.Text,
				},
				"FullyQualifiedName": {
					Name:    "Full Name",
					Type:    fibery.Text,
					SubType: fibery.Title,
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
				"Description": {
					Name:    "Description",
					Type:    fibery.Text,
					SubType: fibery.MD,
				},
				"PurchaseDesc": {
					Name: "Purchase Description",
					Type: fibery.Text,
				},
				"InvStartDate": {
					Name:    "Inventory Start",
					Type:    fibery.DateType,
					SubType: fibery.Day,
				},
				"Type": {
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
				"QtyOnHand": {
					Name: "Quantity On Hand",
					Type: fibery.Number,
					Format: map[string]any{
						"format":               "Number",
						"hasThousandSeparator": true,
						"precision":            2,
					},
				},
				"ReorderPoint": {
					Name: "Reorder Quantity",
					Type: fibery.Number,
					Format: map[string]any{
						"format":               "Number",
						"hasThousandSeparator": true,
						"precision":            2,
					},
				},
				"SKU": {
					Name: "SKU",
					Type: fibery.Text,
				},
				"Taxable": {
					Name:    "Taxable",
					Type:    fibery.Text,
					SubType: fibery.Boolean,
				},
				"SalesTaxIncluded": {
					Name:    "Sales Tax Included",
					Type:    fibery.Text,
					SubType: fibery.Boolean,
				},
				"PurchaseTaxIncluded": {
					Name:    "Purchase Tax Included",
					Type:    fibery.Text,
					SubType: fibery.Boolean,
				},
				"SalesTaxCodeId": {
					Name: "Sales Tax Code ID",
					Type: fibery.Text,
					Relation: &fibery.Relation{
						Cardinality:   fibery.MTO,
						Name:          "Sales Tax",
						TargetName:    "Sales Tax On Items",
						TargetType:    "TaxCode",
						TargetFieldID: "Id",
					},
				},
				"PurchaseTaxCodeId": {
					Name: "Purchase Tax Code ID",
					Type: fibery.Text,
					Relation: &fibery.Relation{
						Cardinality:   fibery.MTO,
						Name:          "Purchase Tax",
						TargetName:    "Purchase Tax On Items",
						TargetType:    "TaxCode",
						TargetFieldID: "Id",
					},
				},
				"ClassId": {
					Name: "Class ID",
					Type: fibery.Text,
					Relation: &fibery.Relation{
						Cardinality:   fibery.MTO,
						Name:          "Class",
						TargetName:    "Expense Account Line(s)",
						TargetType:    "Class",
						TargetFieldID: "Id",
					},
				},
				"PrefVendorId": {
					Name: "Preferred Vendor ID",
					Type: fibery.Text,
					Relation: &fibery.Relation{
						Cardinality:   fibery.MTO,
						Name:          "Preferred Vendor",
						TargetName:    "Primary Sale Items",
						TargetType:    "Vendor",
						TargetFieldID: "Id",
					},
				},
				"ParentId": {
					Name: "Parent ID",
					Type: fibery.Text,
					Relation: &fibery.Relation{
						Cardinality:   fibery.MTO,
						Name:          "Parent",
						TargetName:    "Sub-Items",
						TargetType:    "Item",
						TargetFieldID: "Id",
					},
				},
				"PurchaseCost": {
					Name: "Purchase Cost",
					Type: fibery.Number,
					Format: map[string]any{
						"format":               "Money",
						"currencyCode":         "USD",
						"hasThousandSeperator": true,
						"precision":            2,
					},
				},
				"UnitPrice": {
					Name: "Unit Price",
					Type: fibery.Number,
					Format: map[string]any{
						"format":               "Money",
						"currencyCode":         "USD",
						"hasThousandSeperator": true,
						"precision":            2,
					},
				},
				"AssetAccountId": {
					Name: "Asset Account ID",
					Type: fibery.Text,
					Relation: &fibery.Relation{
						Cardinality:   fibery.MTO,
						Name:          "Asset Account",
						TargetName:    "Items",
						TargetType:    "Account",
						TargetFieldID: "Id",
					},
				},
				"ExpenseAccountId": {
					Name: "Expense Account ID",
					Type: fibery.Text,
					Relation: &fibery.Relation{
						Cardinality:   fibery.MTO,
						Name:          "Expense Account",
						TargetName:    "Items",
						TargetType:    "Account",
						TargetFieldID: "Id",
					},
				},
				"IncomeAccountId": {
					Name: "Income Account ID",
					Type: fibery.Text,
					Relation: &fibery.Relation{
						Cardinality:   fibery.MTO,
						Name:          "Income Account",
						TargetName:    "Items",
						TargetType:    "Account",
						TargetFieldID: "Id",
					},
				},
			},
		},
		schemaGen: func(entity any) (map[string]any, error) {
			item, ok := entity.(quickbooks.Item)
			if !ok {
				return nil, fmt.Errorf("unable to convert entity to Item")
			}

			return map[string]any{
				"Id": item.Id,
				"QBOId": item.Id,
				"Name": item.Name,
				"FullyQualifiedName": 
			}, nil
		},
		query:          func(req Request) (Response, error) {},
		queryProcessor: func(entityArray any, schemaGen schemaGenFunc) ([]map[string]any, error) {},
	},
	cdcProcessor: func(cdc quickbooks.ChangeDataCapture, schemaGen schemaGenFunc) ([]map[string]any, error) {},
	whBatchProcessor: func(itemResponse quickbooks.BatchItemResponse, response *map[string][]map[string]any, cache *cache.Cache, realmId string, queryProcessor queryProcessorFunc, schemaGen schemaGenFunc, typeId string) error {
	},
}
