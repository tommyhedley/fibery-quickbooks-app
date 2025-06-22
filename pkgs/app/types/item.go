package types

import (
	"github.com/tommyhedley/fibery-quickbooks-app/pkgs/app"
	"github.com/tommyhedley/fibery-quickbooks-app/pkgs/fibery"
	"github.com/tommyhedley/quickbooks-go"
)

var item = app.NewDualType(
	"Item",
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
	map[string]app.FieldDef[quickbooks.Item]{
		"qboId": {
			Params: fibery.Field{
				Name:     "QBO ID",
				Type:     fibery.Text,
				ReadOnly: true,
			},
			Convert: func(sd app.StandardData[quickbooks.Item]) (any, error) {
				return sd.Item.Id, nil
			},
		},
		"name": {
			Params: fibery.Field{
				Name: "Base Name",
				Type: fibery.Text,
			},
			Convert: func(sd app.StandardData[quickbooks.Item]) (any, error) {
				return sd.Item.Name, nil
			},
		},
		"fullyQualifiedName": {
			Params: fibery.Field{
				Name:    "Full Name",
				Type:    fibery.Text,
				SubType: fibery.Title,
			},
			Convert: func(sd app.StandardData[quickbooks.Item]) (any, error) {
				return sd.Item.FullyQualifiedName, nil
			},
		},
		"syncToken": {
			Params: fibery.Field{
				Name:     "Sync Token",
				Type:     fibery.Text,
				ReadOnly: true,
			},
			Convert: func(sd app.StandardData[quickbooks.Item]) (any, error) {
				return sd.Item.SyncToken, nil
			},
		},
		"__syncAction": {
			Params: fibery.Field{
				Type: fibery.Text,
				Name: "Sync Action",
			},
			Convert: func(sd app.StandardData[quickbooks.Item]) (any, error) {
				return fibery.SET, nil
			},
		},
		"active": {
			Params: fibery.Field{
				Name:    "Active",
				Type:    fibery.Text,
				SubType: fibery.Boolean,
			},
			Convert: func(sd app.StandardData[quickbooks.Item]) (any, error) {
				return sd.Item.Active, nil
			},
		},
		"description": {
			Params: fibery.Field{
				Name:    "Description",
				Type:    fibery.Text,
				SubType: fibery.MD,
			},
			Convert: func(sd app.StandardData[quickbooks.Item]) (any, error) {
				return sd.Item.Description, nil
			},
		},
		"purchaseDesc": {
			Params: fibery.Field{
				Name: "Purchase Description",
				Type: fibery.Text,
			},
			Convert: func(sd app.StandardData[quickbooks.Item]) (any, error) {
				return sd.Item.PurchaseDesc, nil
			},
		},
		"invStartDate": {
			Params: fibery.Field{
				Name:    "Inventory Start",
				Type:    fibery.DateType,
				SubType: fibery.Day,
			},
			Convert: func(sd app.StandardData[quickbooks.Item]) (any, error) {
				if !sd.Item.InvStartDate.IsZero() {
					return sd.Item.InvStartDate.Format(fibery.DateFormat), nil
				}
				return "", nil
			},
		},
		"type": {
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
					{
						"name": "Category",
					},
				},
			},
			Convert: func(sd app.StandardData[quickbooks.Item]) (any, error) {
				switch sd.Item.Type {
				case "Inventory":
					return "Inventory", nil
				case "Service":
					return "Service", nil
				case "NonInventory":
					return "Non-Inventory", nil
				case "Category":
					return "Category", nil
				default:
					return sd.Item.Type, nil
				}
			},
		},
		"qtyOnHand": {
			Params: fibery.Field{
				Name: "Quantity On Hand",
				Type: fibery.Number,
				Format: map[string]any{
					"format":               "Number",
					"hasThousandSeparator": true,
					"precision":            2,
				},
			},
			Convert: func(sd app.StandardData[quickbooks.Item]) (any, error) {
				return sd.Item.QtyOnHand, nil
			},
		},
		"reorderPoint": {
			Params: fibery.Field{
				Name: "Reorder Quantity",
				Type: fibery.Number,
				Format: map[string]any{
					"format":               "Number",
					"hasThousandSeparator": true,
					"precision":            2,
				},
			},
			Convert: func(sd app.StandardData[quickbooks.Item]) (any, error) {
				return sd.Item.ReorderPoint, nil
			},
		},
		"sku": {
			Params: fibery.Field{
				Name: "SKU",
				Type: fibery.Text,
			},
			Convert: func(sd app.StandardData[quickbooks.Item]) (any, error) {
				return sd.Item.SKU, nil
			},
		},
		"taxable": {
			Params: fibery.Field{
				Name:    "Taxable",
				Type:    fibery.Text,
				SubType: fibery.Boolean,
			},
			Convert: func(sd app.StandardData[quickbooks.Item]) (any, error) {
				return sd.Item.Taxable, nil
			},
		},
		"salesTaxIncluded": {
			Params: fibery.Field{
				Name:    "Sales Tax Included",
				Type:    fibery.Text,
				SubType: fibery.Boolean,
			},
			Convert: func(sd app.StandardData[quickbooks.Item]) (any, error) {
				return sd.Item.SalesTaxIncluded, nil
			},
		},
		"purchaseTaxIncluded": {
			Params: fibery.Field{
				Name:    "Purchase Tax Included",
				Type:    fibery.Text,
				SubType: fibery.Boolean,
			},
			Convert: func(sd app.StandardData[quickbooks.Item]) (any, error) {
				return sd.Item.PurchaseTaxIncluded, nil
			},
		},
		"salesTaxCodeId": {
			Params: fibery.Field{
				Name: "Sales Tax Code ID",
				Type: fibery.Text,
				Relation: &fibery.Relation{
					Cardinality:   fibery.MTO,
					Name:          "Sales Tax",
					TargetName:    "Sales Tax On Items",
					TargetType:    "taxCode",
					TargetFieldID: "id",
				},
			},
			Convert: func(sd app.StandardData[quickbooks.Item]) (any, error) {
				if sd.Item.SalesTaxCodeRef != nil {
					return sd.Item.SalesTaxCodeRef.Value, nil
				}
				return "", nil
			},
		},
		"purchaseTaxCodeId": {
			Params: fibery.Field{
				Name: "Purchase Tax Code ID",
				Type: fibery.Text,
				Relation: &fibery.Relation{
					Cardinality:   fibery.MTO,
					Name:          "Purchase Tax",
					TargetName:    "Purchase Tax On Items",
					TargetType:    "taxCode",
					TargetFieldID: "id",
				},
			},
			Convert: func(sd app.StandardData[quickbooks.Item]) (any, error) {
				if sd.Item.PurchaseTaxCodeRef != nil {
					return sd.Item.PurchaseTaxCodeRef.Value, nil
				}
				return "", nil
			},
		},
		"classId": {
			Params: fibery.Field{
				Name: "Class ID",
				Type: fibery.Text,
				Relation: &fibery.Relation{
					Cardinality:   fibery.MTO,
					Name:          "Class",
					TargetName:    "Items",
					TargetType:    "class",
					TargetFieldID: "id",
				},
			},
			Convert: func(sd app.StandardData[quickbooks.Item]) (any, error) {
				if sd.Item.ClassRef != nil {
					return sd.Item.ClassRef.Value, nil
				}
				return "", nil
			},
		},
		"prefVendorId": {
			Params: fibery.Field{
				Name: "Preferred Vendor ID",
				Type: fibery.Text,
				Relation: &fibery.Relation{
					Cardinality:   fibery.MTO,
					Name:          "Preferred Vendor",
					TargetName:    "Primary Sale Items",
					TargetType:    "vendor",
					TargetFieldID: "id",
				},
			},
			Convert: func(sd app.StandardData[quickbooks.Item]) (any, error) {
				if sd.Item.PrefVendorRef != nil {
					return sd.Item.PrefVendorRef.Value, nil
				}
				return "", nil
			},
		},
		"categoryId": {
			Params: fibery.Field{
				Name: "Parent ID",
				Type: fibery.Text,
				Relation: &fibery.Relation{
					Cardinality:   fibery.MTO,
					Name:          "Category",
					TargetName:    "Items",
					TargetType:    "item",
					TargetFieldID: "id",
				},
			},
			Convert: func(sd app.StandardData[quickbooks.Item]) (any, error) {
				if sd.Item.ParentRef != nil {
					return sd.Item.ParentRef.Value, nil
				}
				return "", nil
			},
		},
		"purchaseCost": {
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
			Convert: func(sd app.StandardData[quickbooks.Item]) (any, error) {
				return sd.Item.PurchaseCost, nil
			},
		},
		"unitPrice": {
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
			Convert: func(sd app.StandardData[quickbooks.Item]) (any, error) {
				return sd.Item.UnitPrice, nil
			},
		},
		"assetAccountId": {
			Params: fibery.Field{
				Name: "Asset Account ID",
				Type: fibery.Text,
				Relation: &fibery.Relation{
					Cardinality:   fibery.MTO,
					Name:          "Asset Account",
					TargetName:    "Inventory Items",
					TargetType:    "account",
					TargetFieldID: "id",
				},
			},
			Convert: func(sd app.StandardData[quickbooks.Item]) (any, error) {
				return sd.Item.AssetAccountRef.Value, nil
			},
		},
		"expenseAccountId": {
			Params: fibery.Field{
				Name: "Expense Account ID",
				Type: fibery.Text,
				Relation: &fibery.Relation{
					Cardinality:   fibery.MTO,
					Name:          "Expense Account",
					TargetName:    "Purchase Items",
					TargetType:    "account",
					TargetFieldID: "id",
				},
			},
			Convert: func(sd app.StandardData[quickbooks.Item]) (any, error) {
				if sd.Item.ExpenseAccountRef != nil {
					return sd.Item.ExpenseAccountRef.Value, nil
				}
				return "", nil
			},
		},
		"incomeAccountId": {
			Params: fibery.Field{
				Name: "Income Account ID",
				Type: fibery.Text,
				Relation: &fibery.Relation{
					Cardinality:   fibery.MTO,
					Name:          "Income Account",
					TargetName:    "Sale Items",
					TargetType:    "account",
					TargetFieldID: "id",
				},
			},
			Convert: func(sd app.StandardData[quickbooks.Item]) (any, error) {
				return sd.Item.IncomeAccountRef.Value, nil
			},
		},
	},
	nil,
)

func init() {
	app.Types.Register(item)
}
