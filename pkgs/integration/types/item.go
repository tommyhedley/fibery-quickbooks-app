package types

import (
	"github.com/tommyhedley/fibery-quickbooks-app/pkgs/fibery"
	"github.com/tommyhedley/fibery-quickbooks-app/pkgs/integration"
	"github.com/tommyhedley/quickbooks-go"
)

var item = integration.NewDualType(
	"item",
	"item",
	"Item",
	func(i quickbooks.Item) string {
		return i.Id
	},
	func(i quickbooks.Item) string {
		return i.Status
	},
	func(id string) quickbooks.Item {
		return quickbooks.Item{
			Id: id,
		}
	},
	func(bir quickbooks.BatchItemResponse) quickbooks.Item {
		return bir.Item
	},
	func(bqr quickbooks.BatchQueryResponse) []quickbooks.Item {
		return bqr.Item
	},
	func(cr quickbooks.CDCQueryResponse) []quickbooks.Item {
		return cr.Item
	},
	map[string]integration.FieldDef[quickbooks.Item]{
		"QBOId": {
			Params: fibery.Field{
				Name: "QBO ID",
				Type: fibery.Text,
			},
			Convert: func(sd integration.StandardData[quickbooks.Item]) (any, error) {
				return sd.Item.Id, nil
			},
		},
		"Name": {
			Params: fibery.Field{
				Name: "Base Name",
				Type: fibery.Text,
			},
			Convert: func(sd integration.StandardData[quickbooks.Item]) (any, error) {
				return sd.Item.Name, nil
			},
		},
		"FullyQualifiedName": {
			Params: fibery.Field{
				Name:    "Full Name",
				Type:    fibery.Text,
				SubType: fibery.Title,
			},
			Convert: func(sd integration.StandardData[quickbooks.Item]) (any, error) {
				return sd.Item.FullyQualifiedName, nil
			},
		},
		"SyncToken": {
			Params: fibery.Field{
				Name:     "Sync Token",
				Type:     fibery.Text,
				ReadOnly: true,
			},
			Convert: func(sd integration.StandardData[quickbooks.Item]) (any, error) {
				return sd.Item.SyncToken, nil
			},
		},
		"__syncAction": {
			Params: fibery.Field{
				Type: fibery.Text,
				Name: "Sync Action",
			},
			Convert: func(sd integration.StandardData[quickbooks.Item]) (any, error) {
				return fibery.SET, nil
			},
		},
		"Active": {
			Params: fibery.Field{
				Name:    "Active",
				Type:    fibery.Text,
				SubType: fibery.Boolean,
			},
			Convert: func(sd integration.StandardData[quickbooks.Item]) (any, error) {
				return sd.Item.Active, nil
			},
		},
		"Description": {
			Params: fibery.Field{
				Name:    "Description",
				Type:    fibery.Text,
				SubType: fibery.MD,
			},
			Convert: func(sd integration.StandardData[quickbooks.Item]) (any, error) {
				return sd.Item.Description, nil
			},
		},
		"PurchaseDesc": {
			Params: fibery.Field{
				Name: "Purchase Description",
				Type: fibery.Text,
			},
			Convert: func(sd integration.StandardData[quickbooks.Item]) (any, error) {
				return sd.Item.PurchaseDesc, nil
			},
		},
		"InvStartDate": {
			Params: fibery.Field{
				Name:    "Inventory Start",
				Type:    fibery.DateType,
				SubType: fibery.Day,
			},
			Convert: func(sd integration.StandardData[quickbooks.Item]) (any, error) {
				if !sd.Item.InvStartDate.IsZero() {
					return sd.Item.InvStartDate.Format(fibery.DateFormat), nil
				}
				return "", nil
			},
		},
		"Type": {
			Params: fibery.Field{
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
			Convert: func(sd integration.StandardData[quickbooks.Item]) (any, error) {
				switch sd.Item.Type {
				case "Inventory":
					return "Inventory", nil
				case "Service":
					return "Service", nil
				case "NonInventory":
					return "Non-Inventory", nil
				default:
					return sd.Item.Type, nil
				}
			},
		},
		"QtyOnHand": {
			Params: fibery.Field{
				Name: "Quantity On Hand",
				Type: fibery.Number,
				Format: map[string]any{
					"format":               "Number",
					"hasThousandSeparator": true,
					"precision":            2,
				},
			},
			Convert: func(sd integration.StandardData[quickbooks.Item]) (any, error) {
				return sd.Item.QtyOnHand, nil
			},
		},
		"ReorderPoint": {
			Params: fibery.Field{
				Name: "Reorder Quantity",
				Type: fibery.Number,
				Format: map[string]any{
					"format":               "Number",
					"hasThousandSeparator": true,
					"precision":            2,
				},
			},
			Convert: func(sd integration.StandardData[quickbooks.Item]) (any, error) {
				return sd.Item.ReorderPoint, nil
			},
		},
		"SKU": {
			Params: fibery.Field{
				Name: "SKU",
				Type: fibery.Text,
			},
			Convert: func(sd integration.StandardData[quickbooks.Item]) (any, error) {
				return sd.Item.SKU, nil
			},
		},
		"Taxable": {
			Params: fibery.Field{
				Name:    "Taxable",
				Type:    fibery.Text,
				SubType: fibery.Boolean,
			},
			Convert: func(sd integration.StandardData[quickbooks.Item]) (any, error) {
				return sd.Item.Taxable, nil
			},
		},
		"SalesTaxIncluded": {
			Params: fibery.Field{
				Name:    "Sales Tax Included",
				Type:    fibery.Text,
				SubType: fibery.Boolean,
			},
			Convert: func(sd integration.StandardData[quickbooks.Item]) (any, error) {
				return sd.Item.SalesTaxIncluded, nil
			},
		},
		"PurchaseTaxIncluded": {
			Params: fibery.Field{
				Name:    "Purchase Tax Included",
				Type:    fibery.Text,
				SubType: fibery.Boolean,
			},
			Convert: func(sd integration.StandardData[quickbooks.Item]) (any, error) {
				return sd.Item.PurchaseTaxIncluded, nil
			},
		},
		"SalesTaxCodeId": {
			Params: fibery.Field{
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
			Convert: func(sd integration.StandardData[quickbooks.Item]) (any, error) {
				if sd.Item.SalesTaxCodeRef != nil {
					return sd.Item.SalesTaxCodeRef.Value, nil
				}
				return "", nil
			},
		},
		"PurchaseTaxCodeId": {
			Params: fibery.Field{
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
			Convert: func(sd integration.StandardData[quickbooks.Item]) (any, error) {
				if sd.Item.PurchaseTaxCodeRef != nil {
					return sd.Item.PurchaseTaxCodeRef.Value, nil
				}
				return "", nil
			},
		},
		"ClassId": {
			Params: fibery.Field{
				Name: "Class ID",
				Type: fibery.Text,
				Relation: &fibery.Relation{
					Cardinality:   fibery.MTO,
					Name:          "Class",
					TargetName:    "Items",
					TargetType:    "Class",
					TargetFieldID: "id",
				},
			},
			Convert: func(sd integration.StandardData[quickbooks.Item]) (any, error) {
				if sd.Item.ClassRef != nil {
					return sd.Item.ClassRef.Value, nil
				}
				return "", nil
			},
		},
		"PrefVendorId": {
			Params: fibery.Field{
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
			Convert: func(sd integration.StandardData[quickbooks.Item]) (any, error) {
				if sd.Item.PrefVendorRef != nil {
					return sd.Item.PrefVendorRef.Value, nil
				}
				return "", nil
			},
		},
		"ParentId": {
			Params: fibery.Field{
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
			Convert: func(sd integration.StandardData[quickbooks.Item]) (any, error) {
				if sd.Item.ParentRef != nil {
					return sd.Item.ParentRef.Value, nil
				}
				return "", nil
			},
		},
		"PurchaseCost": {
			Params: fibery.Field{
				Name: "Purchase Cost",
				Type: fibery.Number,
				Format: map[string]any{
					"format":               "Money",
					"currencyCode":         "USD",
					"hasThousandSeperator": true,
					"precision":            2,
				},
			},
			Convert: func(sd integration.StandardData[quickbooks.Item]) (any, error) {
				return sd.Item.PurchaseCost, nil
			},
		},
		"UnitPrice": {
			Params: fibery.Field{
				Name: "Unit Price",
				Type: fibery.Number,
				Format: map[string]any{
					"format":               "Money",
					"currencyCode":         "USD",
					"hasThousandSeperator": true,
					"precision":            2,
				},
			},
			Convert: func(sd integration.StandardData[quickbooks.Item]) (any, error) {
				return sd.Item.UnitPrice, nil
			},
		},
		"AssetAccountId": {
			Params: fibery.Field{
				Name: "Asset Account ID",
				Type: fibery.Text,
				Relation: &fibery.Relation{
					Cardinality:   fibery.MTO,
					Name:          "Asset Account",
					TargetName:    "Inventory Items",
					TargetType:    "Account",
					TargetFieldID: "id",
				},
			},
			Convert: func(sd integration.StandardData[quickbooks.Item]) (any, error) {
				return sd.Item.AssetAccountRef.Value, nil
			},
		},
		"ExpenseAccountId": {
			Params: fibery.Field{
				Name: "Expense Account ID",
				Type: fibery.Text,
				Relation: &fibery.Relation{
					Cardinality:   fibery.MTO,
					Name:          "Expense Account",
					TargetName:    "Purchase Items",
					TargetType:    "Account",
					TargetFieldID: "id",
				},
			},
			Convert: func(sd integration.StandardData[quickbooks.Item]) (any, error) {
				if sd.Item.ExpenseAccountRef != nil {
					return sd.Item.ExpenseAccountRef.Value, nil
				}
				return "", nil
			},
		},
		"IncomeAccountId": {
			Params: fibery.Field{
				Name: "Income Account ID",
				Type: fibery.Text,
				Relation: &fibery.Relation{
					Cardinality:   fibery.MTO,
					Name:          "Income Account",
					TargetName:    "Sale Items",
					TargetType:    "Account",
					TargetFieldID: "id",
				},
			},
			Convert: func(sd integration.StandardData[quickbooks.Item]) (any, error) {
				return sd.Item.IncomeAccountRef.Value, nil
			},
		},
	},
	nil,
)

func init() {
	integration.Types.Register(item)
}

