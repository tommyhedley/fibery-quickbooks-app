package data

import (
	"github.com/tommyhedley/fibery-quickbooks-app/pkgs/fibery"
	"github.com/tommyhedley/quickbooks-go"
)

var Item = QuickBooksDualType[quickbooks.Item]{
	QuickBooksType: QuickBooksType[quickbooks.Item]{
		BaseType: fibery.BaseType{
			TypeId:   "Item",
			TypeName: "Item",
			TypeSchema: map[string]fibery.Field{
				"id": {
					Name: "ID",
					Type: fibery.Id,
				},
				"QBOId": {
					Name: "QBO ID",
					Type: fibery.Text,
				},
				"Name": {
					Name: "Base Name",
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
						TargetFieldID: "id",
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
						TargetFieldID: "id",
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
						TargetFieldID: "id",
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
						TargetFieldID: "id",
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
						TargetFieldID: "id",
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
						TargetFieldID: "id",
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
						TargetFieldID: "id",
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
						TargetFieldID: "id",
					},
				},
			},
		},
		schemaGen: func(i quickbooks.Item) (map[string]any, error) {
			var itemType string
			switch i.Type {
			case "Inventory":
				itemType = "Inventory"
			case "Service":
				itemType = "Service"
			case "NonInventory":
				itemType = "Non-Inventory"
			}

			var salesTaxCodeId string
			if i.SalesTaxCodeRef != nil {
				salesTaxCodeId = i.SalesTaxCodeRef.Value
			}

			var purchaseTaxCodeId string
			if i.PurchaseTaxCodeRef != nil {
				purchaseTaxCodeId = i.PurchaseTaxCodeRef.Value
			}

			var classId string
			if i.ClassRef != nil {
				classId = i.ClassRef.Value
			}

			var vendorId string
			if i.PrefVendorRef != nil {
				vendorId = i.PrefVendorRef.Value
			}

			var parentId string
			if i.ParentRef != nil {
				parentId = i.ParentRef.Value
			}

			var expenseAccountId string
			if i.ExpenseAccountRef != nil {
				expenseAccountId = i.ExpenseAccountRef.Value
			}

			return map[string]any{
				"id":                  i.Id,
				"QBOId":               i.Id,
				"Name":                i.Name,
				"FullyQualifiedName":  i.FullyQualifiedName,
				"SyncToken":           i.SyncToken,
				"__syncAction":        fibery.SET,
				"Active":              i.Active,
				"Description":         i.Description,
				"PurchaseDesc":        i.PurchaseDesc,
				"InvStartDate":        i.InvStartDate.Format(fibery.DateFormat),
				"Type":                itemType,
				"QtyOnHand":           i.QtyOnHand,
				"ReorderPoint":        i.ReorderPoint,
				"SKU":                 i.SKU,
				"Taxable":             i.Taxable,
				"SalesTaxIncluded":    i.SalesTaxIncluded,
				"PurchaseTaxIncluded": i.PurchaseTaxIncluded,
				"SalesTaxCodeId":      salesTaxCodeId,
				"PurchaseTaxCodeId":   purchaseTaxCodeId,
				"ClassId":             classId,
				"PrefVendorId":        vendorId,
				"ParentId":            parentId,
				"PurchaseCost":        i.PurchaseCost,
				"UnitPrice":           i.UnitPrice,
				"AssetAccountId":      i.AssetAccountRef.Value,
				"ExpenseAccountId":    expenseAccountId,
				"IncomeAccountId":     i.IncomeAccountRef.Value,
			}, nil
		},
		pageQuery: func(req Request) ([]quickbooks.Item, error) {
			params := quickbooks.RequestParameters{
				Ctx:     req.Ctx,
				RealmId: req.RealmId,
				Token:   req.Token,
			}

			items, err := req.Client.FindItemsByPage(params, req.StartPosition, req.PageSize)
			if err != nil {
				return nil, err
			}

			return items, nil
		},
	},
	entityId: func(i quickbooks.Item) string {
		return i.Id
	},
	entityStatus: func(i quickbooks.Item) string {
		return i.Status
	},
}

func init() {
	registerType(&Item)
}
