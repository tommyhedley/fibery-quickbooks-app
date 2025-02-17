package data

import (
	"fmt"

	"github.com/tommyhedley/fibery-quickbooks-app/pkgs/fibery"
	"github.com/tommyhedley/quickbooks-go"
)

var Bill = QuickBooksDualType{
	QuickBooksType: QuickBooksType{
		fiberyType: fiberyType{
			id:   "Bill",
			name: "Bill",
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
				"SyncToken": {
					Name:     "Sync Token",
					Type:     fibery.Text,
					ReadOnly: true,
				},
				"__syncAction": {
					Type: fibery.Text,
					Name: "Sync Action",
				},
				"DocNumber": {
					Name: "Bill Number",
					Type: fibery.Text,
				},
				"TxnDate": {
					Name:    "Bill Date",
					Type:    fibery.DateType,
					SubType: fibery.Day,
				},
				"DueDate": {
					Name:    "Due Date",
					Type:    fibery.DateType,
					SubType: fibery.Day,
				},
				"PrivateNote": {
					Name:    "Memo",
					Type:    fibery.Text,
					SubType: fibery.MD,
				},
				"TotalAmt": {
					Name: "Total",
					Type: fibery.Number,
					Format: map[string]any{
						"format":               "Money",
						"currencyCode":         "USD",
						"hasThousandSeperator": true,
						"precision":            2,
					},
				},
				"Balance": {
					Name: "Balance",
					Type: fibery.Number,
					Format: map[string]any{
						"format":               "Money",
						"currencyCode":         "USD",
						"hasThousandSeperator": true,
						"precision":            2,
					},
				},
				"VendorId": {
					Name: "Vendor ID",
					Type: fibery.Text,
					Relation: &fibery.Relation{
						Cardinality:   fibery.MTO,
						Name:          "Vendor",
						TargetName:    "Bills",
						TargetType:    "Vendor",
						TargetFieldID: "Id",
					},
				},
				"APAccountId": {
					Name: "AP Account ID",
					Type: fibery.Text,
					Relation: &fibery.Relation{
						Cardinality:   fibery.MTO,
						Name:          "AP Account",
						TargetName:    "Bills",
						TargetType:    "Account",
						TargetFieldID: "Id",
					},
				},
				"SalesTermId": {
					Name: "Sales Term ID",
					Type: fibery.Text,
					Relation: &fibery.Relation{
						Cardinality:   fibery.MTO,
						Name:          "Terms",
						TargetName:    "Bills",
						TargetType:    "SalesTerm",
						TargetFieldID: "Id",
					},
				},
			},
		},
		schemaGen: func(entity any) (map[string]any, error) {
			bill, ok := entity.(quickbooks.Bill)
			if !ok {
				return nil, fmt.Errorf("unable to convert entity to bill")
			}

			var name string

			if bill.VendorRef.Name != "" {
				if !bill.TxnDate.IsZero() {
					name = bill.VendorRef.Name + " - " + bill.TxnDate.Format(fibery.DateFormat)
				} else {
					name = bill.VendorRef.Name
				}
			}
			data := map[string]any{
				"Id":           bill.Id,
				"QBOId":        bill.Id,
				"Name":         name,
				"SyncToken":    bill.SyncToken,
				"__syncAction": fibery.SET,
				"DocNumber":    bill.DocNumber,
				"TxnDate":      bill.TxnDate.Format(fibery.DateFormat),
				"DueDate":      bill.DueDate.Format(fibery.DateFormat),
				"PrivateNote":  bill.PrivateNote,
				"TotalAmt":     bill.TotalAmt,
				"Balance":      bill.Balance,
				"VendorId":     bill.VendorRef.Value,
				"APAccountId":  bill.APAccountRef.Value,
				"SalesTermId":  bill.SalesTermRef.Value,
			}
			return data, nil
		},
		query:          func(req Request) (Response, error) {},
		queryProcessor: func(entityArray any, schemaGen schemaGenFunc) ([]map[string]any, error) {},
	},
	cdcProcessor: func(cdc quickbooks.ChangeDataCapture, schemaGen schemaGenFunc) ([]map[string]any, error) {},
	whBatchProcessor: func(itemResponse quickbooks.BatchItemResponse, response *map[string][]map[string]any, cache *cache.Cache, realmId string, queryProcessor queryProcessorFunc, schemaGen schemaGenFunc, typeId string) error {
	},
}

var BillItemLine = DependentDualType{
	dependentBaseType: dependentBaseType{
		fiberyType: fiberyType{
			id:   "BillItemLine",
			name: "Bill Item Line",
			schema: map[string]fibery.Field{
				"Id": {
					Name: "ID",
					Type: fibery.ID,
				},
				"QBOId": {
					Name: "QBO ID",
					Type: fibery.Text,
				},
				"Description": {
					Name:    "Description",
					Type:    fibery.Text,
					SubType: fibery.Title,
				},
				"__syncAction": {
					Type: fibery.Text,
					Name: "Sync Action",
				},
				"LineNum": {
					Name: "Line",
					Type: fibery.Number,
				},
				"Tax": {
					Name:    "Tax",
					Type:    fibery.Text,
					SubType: fibery.Boolean,
				},
				"Billable": {
					Name:    "Billable",
					Type:    fibery.Text,
					SubType: fibery.Boolean,
				},
				"Billed": {
					Name:    "Billed",
					Type:    fibery.Text,
					SubType: fibery.Boolean,
				},
				"Qty": {
					Name: "Quantity",
					Type: fibery.Number,
					Format: map[string]any{
						"format":               "Number",
						"hasThousandSeparator": true,
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
				"MarkupPercent": {
					Name: "Markup",
					Type: fibery.Number,
					Format: map[string]any{
						"format":    "Percent",
						"precision": 2,
					},
				},
				"Amount": {
					Name: "Amount",
					Type: fibery.Number,
					Format: map[string]any{
						"format":               "Money",
						"currencyCode":         "USD",
						"hasThousandSeperator": true,
						"precision":            2,
					},
				},
				"BillId": {
					Name: "Bill ID",
					Type: fibery.Text,
					Relation: &fibery.Relation{
						Cardinality:   fibery.MTO,
						Name:          "Bill",
						TargetName:    "Item Lines",
						TargetType:    "Bill",
						TargetFieldID: "Id",
					},
				},
				"ItemId": {
					Name: "Item ID",
					Type: fibery.Text,
					Relation: &fibery.Relation{
						Cardinality:   fibery.MTO,
						Name:          "Item",
						TargetName:    "Bill Item Lines",
						TargetType:    "Item",
						TargetFieldID: "Id",
					},
				},
				"CustomerId": {
					Name: "Customer ID",
					Type: fibery.Text,
					Relation: &fibery.Relation{
						Cardinality:   fibery.MTO,
						Name:          "Customer",
						TargetName:    "Bill Item Lines",
						TargetType:    "Customer",
						TargetFieldID: "Id",
					},
				},
				"ClassId": {
					Name: "Class ID",
					Type: fibery.Text,
					Relation: &fibery.Relation{
						Cardinality:   fibery.MTO,
						Name:          "Class",
						TargetName:    "Bill Item Lines",
						TargetType:    "Class",
						TargetFieldID: "Id",
					},
				},
				"MarkupAccountId": {
					Name: "Markup Account ID",
					Type: fibery.Text,
					Relation: &fibery.Relation{
						Cardinality:   fibery.MTO,
						Name:          "Markup Income Account",
						TargetName:    "Bill Item Line Markup",
						TargetType:    "Account",
						TargetFieldID: "Id",
					},
				},
			},
		},
		schemaGen: func(entity, source any) (map[string]any, error) {
			line, ok := entity.(quickbooks.Line)
			if !ok {
				return nil, fmt.Errorf("unable to convert entity to Line")
			}

			bill, ok := source.(quickbooks.Bill)
			if !ok {
				return nil, fmt.Errorf("unable to convert source to Bill")
			}

			tax := false
			if line.ItemBasedExpenseLineDetail.TaxCodeRef.Value == "TAX" {
				tax = true
			}

			var billable bool
			switch line.ItemBasedExpenseLineDetail.BillableStatus {
			case quickbooks.BillableStatusType:
				billable = true
			case quickbooks.HasBeenBilledStatusType:
				billable = true
			default:
				billable = false
			}

			billed := false
			if line.ItemBasedExpenseLineDetail.BillableStatus == quickbooks.HasBeenBilledStatusType {
				billed = true
			}

			data := map[string]any{
				"Id":              fmt.Sprintf("%s:i:%s", bill.Id, line.Id),
				"QBOId":           line.Id,
				"Description":     line.Description,
				"__syncAction":    fibery.SET,
				"LineNum":         line.LineNum,
				"Tax":             tax,
				"Billable":        billable,
				"Billed":          billed,
				"Qty":             line.ItemBasedExpenseLineDetail.Qty,
				"UnitPrice":       line.ItemBasedExpenseLineDetail.UnitPrice,
				"MarkupPercent":   line.ItemBasedExpenseLineDetail.MarkupInfo.Percent,
				"Amount":          line.Amount,
				"BillId":          bill.Id,
				"ItemId":          line.ItemBasedExpenseLineDetail.ItemRef.Value,
				"CustomerId":      line.AccountBasedExpenseLineDetail.CustomerRef.Value,
				"ClassId":         line.ItemBasedExpenseLineDetail.ClassRef.Value,
				"MarkupAccountId": line.ItemBasedExpenseLineDetail.MarkupInfo.MarkUpIncomeAccountRef.Value,
			}

			return data, nil
		},
		queryProcessor: func(sourceArray any, schemaGen depSchemaGenFunc) ([]map[string]any, error) {},
	},
	source:       Bill,
	sourceMapper: func(source any) (map[string]bool, error) {},
	typeMapper:   func(sourceArray any, sourceMapper sourceMapperFunc) (map[string]map[string]bool, error) {},
	cdcProcessor: func(cdc quickbooks.ChangeDataCapture, cacheEntry *IdCache, sourceMapper sourceMapperFunc, schemaGen depSchemaGenFunc) ([]map[string]any, error) {
	},
	whBatchProcessor: func(sourceArray any, cacheEntry *IdCache, sourceMapper sourceMapperFunc, schemaGen depSchemaGenFunc) ([]map[string]any, error) {
	},
}

var BillAccountLine = DependentDualType{
	dependentBaseType: dependentBaseType{
		fiberyType: fiberyType{
			id:   "BillAccountLine",
			name: "Bill Account Line",
			schema: map[string]fibery.Field{
				"Id": {
					Name: "ID",
					Type: fibery.ID,
				},
				"QBOId": {
					Name: "QBO ID",
					Type: fibery.Text,
				},
				"Description": {
					Name:    "Description",
					Type:    fibery.Text,
					SubType: fibery.Title,
				},
				"__syncAction": {
					Type: fibery.Text,
					Name: "Sync Action",
				},
				"LineNum": {
					Name: "Line",
					Type: fibery.Number,
				},
				"Tax": {
					Name:    "Tax",
					Type:    fibery.Text,
					SubType: fibery.Boolean,
				},
				"Billable": {
					Name:    "Billable",
					Type:    fibery.Text,
					SubType: fibery.Boolean,
				},
				"Billed": {
					Name:    "Billed",
					Type:    fibery.Text,
					SubType: fibery.Boolean,
				},
				"MarkupPercent": {
					Name: "Markup",
					Type: fibery.Number,
					Format: map[string]any{
						"format":    "Percent",
						"precision": 2,
					},
				},
				"Amount": {
					Name: "Amount",
					Type: fibery.Number,
					Format: map[string]any{
						"format":               "Money",
						"currencyCode":         "USD",
						"hasThousandSeperator": true,
						"precision":            2,
					},
				},
				"BillId": {
					Name: "Bill ID",
					Type: fibery.Text,
					Relation: &fibery.Relation{
						Cardinality:   fibery.MTO,
						Name:          "Bill",
						TargetName:    "Item Lines",
						TargetType:    "Bill",
						TargetFieldID: "Id",
					},
				},
				"AccountId": {
					Name: "Item ID",
					Type: fibery.Text,
					Relation: &fibery.Relation{
						Cardinality:   fibery.MTO,
						Name:          "Category",
						TargetName:    "Bill Account Lines",
						TargetType:    "Account",
						TargetFieldID: "Id",
					},
				},
				"CustomerId": {
					Name: "Customer ID",
					Type: fibery.Text,
					Relation: &fibery.Relation{
						Cardinality:   fibery.MTO,
						Name:          "Customer",
						TargetName:    "Bill Item Lines",
						TargetType:    "Customer",
						TargetFieldID: "Id",
					},
				},
				"ClassId": {
					Name: "Class ID",
					Type: fibery.Text,
					Relation: &fibery.Relation{
						Cardinality:   fibery.MTO,
						Name:          "Class",
						TargetName:    "Bill Item Lines",
						TargetType:    "Class",
						TargetFieldID: "Id",
					},
				},
				"MarkupAccountId": {
					Name: "Markup Account ID",
					Type: fibery.Text,
					Relation: &fibery.Relation{
						Cardinality:   fibery.MTO,
						Name:          "Markup Income Account",
						TargetName:    "Bill Item Line Markup",
						TargetType:    "Account",
						TargetFieldID: "Id",
					},
				},
			},
		},
		schemaGen: func(entity, source any) (map[string]any, error) {
			line, ok := entity.(quickbooks.Line)
			if !ok {
				return nil, fmt.Errorf("unable to convert entity to Line")
			}

			bill, ok := source.(quickbooks.Bill)
			if !ok {
				return nil, fmt.Errorf("unable to convert source to Bill")
			}

			tax := false
			if line.AccountBasedExpenseLineDetail.TaxCodeRef.Value == "TAX" {
				tax = true
			}

			var billable bool
			switch line.AccountBasedExpenseLineDetail.BillableStatus {
			case quickbooks.BillableStatusType:
				billable = true
			case quickbooks.HasBeenBilledStatusType:
				billable = true
			default:
				billable = false
			}

			billed := false
			if line.AccountBasedExpenseLineDetail.BillableStatus == quickbooks.HasBeenBilledStatusType {
				billed = true
			}

			data := map[string]any{
				"Id":              fmt.Sprintf("%s:a:%s", bill.Id, line.Id),
				"QBOId":           line.Id,
				"Description":     line.Description,
				"__syncAction":    fibery.SET,
				"LineNum":         line.LineNum,
				"Tax":             tax,
				"Billable":        billable,
				"Billed":          billed,
				"MarkupPercent":   line.AccountBasedExpenseLineDetail.MarkupInfo.Percent,
				"Amount":          line.Amount,
				"BillId":          bill.Id,
				"AccountId":       line.AccountBasedExpenseLineDetail.AccountRef.Value,
				"CustomerId":      line.AccountBasedExpenseLineDetail.CustomerRef.Value,
				"ClassId":         line.AccountBasedExpenseLineDetail.ClassRef.Value,
				"MarkupAccountId": line.AccountBasedExpenseLineDetail.MarkupInfo.MarkUpIncomeAccountRef.Value,
			}

			return data, nil
		},
		queryProcessor: func(sourceArray any, schemaGen depSchemaGenFunc) ([]map[string]any, error) {},
	},
	source:       Bill,
	sourceMapper: func(source any) (map[string]bool, error) {},
	typeMapper:   func(sourceArray any, sourceMapper sourceMapperFunc) (map[string]map[string]bool, error) {},
	cdcProcessor: func(cdc quickbooks.ChangeDataCapture, cacheEntry *IdCache, sourceMapper sourceMapperFunc, schemaGen depSchemaGenFunc) ([]map[string]any, error) {
	},
	whBatchProcessor: func(sourceArray any, cacheEntry *IdCache, sourceMapper sourceMapperFunc, schemaGen depSchemaGenFunc) ([]map[string]any, error) {
	},
}
