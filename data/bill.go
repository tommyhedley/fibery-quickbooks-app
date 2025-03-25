package data

import (
	"fmt"

	"github.com/tommyhedley/fibery-quickbooks-app/pkgs/fibery"
	"github.com/tommyhedley/quickbooks-go"
)

var Bill = QuickBooksDualType[quickbooks.Bill]{
	QuickBooksType: QuickBooksType[quickbooks.Bill]{
		BaseType: fibery.BaseType{
			TypeId:   "Bill",
			TypeName: "Bill",
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
					Name:    "Name",
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
						TargetFieldID: "id",
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
						TargetFieldID: "id",
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
						TargetFieldID: "id",
					},
				},
			},
		},
		schemaGen: func(b quickbooks.Bill) (map[string]any, error) {
			var apAccountId string
			if b.APAccountRef != nil {
				apAccountId = b.APAccountRef.Value
			}

			var salesTermId string
			if b.SalesTermRef != nil {
				salesTermId = b.SalesTermRef.Value
			}

			return map[string]any{
				"id":           b.Id,
				"QBOId":        b.Id,
				"Name":         b.PrivateNote,
				"SyncToken":    b.SyncToken,
				"__syncAction": fibery.SET,
				"DocNumber":    b.DocNumber,
				"TxnDate":      b.TxnDate.Format(fibery.DateFormat),
				"DueDate":      b.DueDate.Format(fibery.DateFormat),
				"PrivateNote":  b.PrivateNote,
				"TotalAmt":     b.TotalAmt,
				"Balance":      b.Balance,
				"VendorId":     b.VendorRef.Value,
				"APAccountId":  apAccountId,
				"SalesTermId":  salesTermId,
			}, nil
		},
		pageQuery: func(req Request) ([]quickbooks.Bill, error) {
			params := quickbooks.RequestParameters{
				Ctx:     req.Ctx,
				RealmId: req.RealmId,
				Token:   req.Token,
			}

			items, err := req.Client.FindBillsByPage(params, req.StartPosition, req.PageSize)
			if err != nil {
				return nil, err
			}

			return items, nil
		},
	},
	entityId: func(b quickbooks.Bill) string {
		return b.Id
	},
	entityStatus: func(b quickbooks.Bill) string {
		return b.Status
	},
}

var BillItemLine = DependentDualType[quickbooks.Bill]{
	dependentBaseType: dependentBaseType[quickbooks.Bill]{
		BaseType: fibery.BaseType{
			TypeId:   "BillItemLine",
			TypeName: "Bill Item Line",
			TypeSchema: map[string]fibery.Field{
				"id": {
					Name: "ID",
					Type: fibery.Id,
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
						TargetFieldID: "id",
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
						TargetFieldID: "id",
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
						TargetFieldID: "id",
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
						TargetFieldID: "id",
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
						TargetFieldID: "id",
					},
				},
			},
		},
		schemaGen: func(b quickbooks.Bill) ([]map[string]any, error) {
			items := []map[string]any{}
			for _, line := range b.Line {
				if line.DetailType == quickbooks.ItemExpenseLine {
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

					item := map[string]any{
						"id":              fmt.Sprintf("%s:i:%s", b.Id, line.Id),
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
						"BillId":          b.Id,
						"ItemId":          line.ItemBasedExpenseLineDetail.ItemRef.Value,
						"CustomerId":      line.AccountBasedExpenseLineDetail.CustomerRef.Value,
						"ClassId":         line.ItemBasedExpenseLineDetail.ClassRef.Value,
						"MarkupAccountId": line.ItemBasedExpenseLineDetail.MarkupInfo.MarkUpIncomeAccountRef.Value,
					}
					items = append(items, item)
				}
			}
			return items, nil
		},
	},
	sourceType: &Bill,
	sourceId: func(b quickbooks.Bill) string {
		return b.Id
	},
	sourceStatus: func(b quickbooks.Bill) string {
		return b.Status
	},
	sourceMapper: func(b quickbooks.Bill) map[string]struct{} {
		sourceMap := map[string]struct{}{}
		for _, line := range b.Line {
			if line.DetailType == quickbooks.ItemExpenseLine {
				sourceMap[fmt.Sprintf("%s:%s", b.Id, line.Id)] = struct{}{}
			}
		}
		return sourceMap
	},
}

var BillAccountLine = DependentDualType[quickbooks.Bill]{
	dependentBaseType: dependentBaseType[quickbooks.Bill]{
		BaseType: fibery.BaseType{
			TypeId:   "BillAccountLine",
			TypeName: "Bill Account Line",
			TypeSchema: map[string]fibery.Field{
				"id": {
					Name: "ID",
					Type: fibery.Id,
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
						TargetFieldID: "id",
					},
				},
				"AccountId": {
					Name: "Account ID",
					Type: fibery.Text,
					Relation: &fibery.Relation{
						Cardinality:   fibery.MTO,
						Name:          "Category",
						TargetName:    "Bill Account Lines",
						TargetType:    "Account",
						TargetFieldID: "id",
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
						TargetFieldID: "id",
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
						TargetFieldID: "id",
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
						TargetFieldID: "id",
					},
				},
			},
		},
		schemaGen: func(b quickbooks.Bill) ([]map[string]any, error) {
			items := []map[string]any{}
			for _, line := range b.Line {
				if line.DetailType == quickbooks.AccountExpenseLine {
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

					item := map[string]any{
						"id":              fmt.Sprintf("%s:a:%s", b.Id, line.Id),
						"QBOId":           line.Id,
						"Description":     line.Description,
						"__syncAction":    fibery.SET,
						"LineNum":         line.LineNum,
						"Tax":             tax,
						"Billable":        billable,
						"Billed":          billed,
						"MarkupPercent":   line.AccountBasedExpenseLineDetail.MarkupInfo.Percent,
						"Amount":          line.Amount,
						"BillId":          b.Id,
						"AccountId":       line.AccountBasedExpenseLineDetail.AccountRef.Value,
						"CustomerId":      line.AccountBasedExpenseLineDetail.CustomerRef.Value,
						"ClassId":         line.AccountBasedExpenseLineDetail.ClassRef.Value,
						"MarkupAccountId": line.AccountBasedExpenseLineDetail.MarkupInfo.MarkUpIncomeAccountRef.Value,
					}

					items = append(items, item)
				}
			}
			return items, nil
		},
	},
	sourceType: &Bill,
	sourceId: func(b quickbooks.Bill) string {
		return b.Id
	},
	sourceStatus: func(b quickbooks.Bill) string {
		return b.Status
	},
	sourceMapper: func(b quickbooks.Bill) map[string]struct{} {
		sourceMap := map[string]struct{}{}
		for _, line := range b.Line {
			if line.DetailType == quickbooks.AccountExpenseLine {
				sourceMap[fmt.Sprintf("%s:%s", b.Id, line.Id)] = struct{}{}
			}
		}
		return sourceMap
	},
}

func init() {
	registerType(&Bill)
	registerType(&BillItemLine)
	registerType(&BillAccountLine)
}
