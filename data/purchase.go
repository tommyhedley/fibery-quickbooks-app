package data

import (
	"fmt"

	"github.com/tommyhedley/fibery-quickbooks-app/pkgs/fibery"
	"github.com/tommyhedley/quickbooks-go"
)

var Purchase = QuickBooksDualType[quickbooks.Purchase]{
	QuickBooksType: QuickBooksType[quickbooks.Purchase]{
		BaseType: fibery.BaseType{
			TypeId:   "Purchase",
			TypeName: "Expense",
			TypeSchema: map[string]fibery.Field{
				"id": {
					Name: "id",
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
				"PaymentType": {
					Name:     "Payment Type",
					Type:     fibery.Text,
					SubType:  fibery.SingleSelect,
					ReadOnly: true,
					Options: []map[string]any{
						{
							"name": "Cash",
						},
						{
							"name": "Check",
						},
						{
							"name": "Credit Card",
						},
					},
				},
				"PaymentAccountId": {
					Name: "Payment Account ID",
					Type: fibery.Text,
					Relation: &fibery.Relation{
						Cardinality:   fibery.MTO,
						Name:          "Payment Account",
						TargetName:    "Expenses",
						TargetType:    "Account",
						TargetFieldID: "id",
					},
				},
				"PaymentMethodId": {
					Name: "Payment Method ID",
					Type: fibery.Text,
					Relation: &fibery.Relation{
						Cardinality:   fibery.MTO,
						Name:          "Payment Method",
						TargetName:    "Expenses",
						TargetType:    "PaymentMethod",
						TargetFieldID: "id",
					},
				},
				"EntityId": {
					Name: "Entity ID",
					Type: fibery.Text,
					Relation: &fibery.Relation{
						Cardinality:   fibery.MTO,
						Name:          "Entity",
						TargetName:    "Expenses",
						TargetType:    "Entity",
						TargetFieldID: "id",
					},
				},
				"Total": {
					Name: "Total",
					Type: fibery.Number,
					Format: map[string]any{
						"format":               "Money",
						"currencyCode":         "USD",
						"hasThousandSeperator": true,
						"precision":            2,
					},
				},
				"DocNumber": {
					Name: "Reference Number",
					Type: fibery.Text,
				},
				"TxnDate": {
					Name: "Invoice Date",
					Type: fibery.DateType,
				},
				"Credit": {
					Name:    "Credit",
					Type:    fibery.Text,
					SubType: fibery.Boolean,
				},
				"PrivateNote": {
					Name:    "Memo",
					Type:    fibery.Text,
					SubType: fibery.MD,
				},
			},
		},
		schemaGen: func(p quickbooks.Purchase) (map[string]any, error) {
			var paymentType string
			switch p.PaymentType {
			case "Cash":
				paymentType = "Cash"
			case "Check":
				paymentType = "Check"
			case "CreditCard":
				paymentType = "Credit Card"
			}

			var entityId string
			var entityName string
			if p.EntityRef != nil {
				entityName = p.EntityRef.Name
				switch p.EntityRef.Type {
				case "Customer":
					entityId = "c:" + p.EntityRef.Value
				case "Vendor":
					entityId = "v:" + p.EntityRef.Value
				case "Employee":
					entityId = "e:" + p.EntityRef.Value
				}
			}

			var paymentMethodId string
			if p.PaymentMethodRef != nil {
				paymentMethodId = p.PaymentMethodRef.Value
			}

			var txnDate string
			if p.TxnDate != nil {
				txnDate = p.TxnDate.Format(fibery.DateFormat)
			}

			return map[string]any{
				"id":               p.Id,
				"QBOId":            p.Id,
				"Name":             entityName,
				"SyncToken":        p.SyncToken,
				"__syncToken":      fibery.SET,
				"PaymentType":      paymentType,
				"PaymentAccountId": p.AccountRef.Value,
				"PaymentMethodId":  paymentMethodId,
				"EntityId":         entityId,
				"Total":            p.TotalAmt,
				"DocNumber":        p.DocNumber,
				"TxnDate":          txnDate,
				"Credit":           p.Credit,
				"PrivateNote":      p.PrivateNote,
			}, nil
		},
		pageQuery: func(req Request) ([]quickbooks.Purchase, error) {
			params := quickbooks.RequestParameters{
				Ctx:     req.Ctx,
				RealmId: req.RealmId,
				Token:   req.Token,
			}

			items, err := req.Client.FindPurchasesByPage(params, req.StartPosition, req.PageSize)
			if err != nil {
				return nil, err
			}

			return items, nil
		},
	},
	entityId: func(p quickbooks.Purchase) string {
		return p.Id
	},
	entityStatus: func(p quickbooks.Purchase) string {
		return p.Status
	},
}
var PurchaseItemLine = DependentDualType[quickbooks.Purchase]{
	dependentBaseType: dependentBaseType[quickbooks.Purchase]{
		BaseType: fibery.BaseType{
			TypeId:   "PurchaseItemLine",
			TypeName: "Purchase Item Line",
			TypeSchema: map[string]fibery.Field{
				"id": {
					Name: "id",
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
				"PurchaseId": {
					Name: "Purchase ID",
					Type: fibery.Text,
					Relation: &fibery.Relation{
						Cardinality:   fibery.MTO,
						Name:          "Expense",
						TargetName:    "Item Lines",
						TargetType:    "Purchase",
						TargetFieldID: "id",
					},
				},
				"ItemId": {
					Name: "Item ID",
					Type: fibery.Text,
					Relation: &fibery.Relation{
						Cardinality:   fibery.MTO,
						Name:          "Item",
						TargetName:    "Expense Item Lines",
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
						TargetName:    "Expense Item Lines",
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
						TargetName:    "Expense Item Lines",
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
						TargetName:    "Expense Item Line Markup",
						TargetType:    "Account",
						TargetFieldID: "id",
					},
				},
			},
		},
		schemaGen: func(p quickbooks.Purchase) ([]map[string]any, error) {
			items := []map[string]any{}
			for _, line := range p.Line {
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
						"id":              fmt.Sprintf("%s:i:%s", p.Id, line.Id),
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
						"BillId":          p.Id,
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
	sourceType: &Purchase,
	sourceId: func(p quickbooks.Purchase) string {
		return p.Id
	},
	sourceStatus: func(p quickbooks.Purchase) string {
		return p.Status
	},
	sourceMapper: func(p quickbooks.Purchase) map[string]struct{} {
		sourceMap := map[string]struct{}{}
		for _, line := range p.Line {
			if line.DetailType == quickbooks.ItemExpenseLine {
				sourceMap[fmt.Sprintf("%s:%s", p.Id, line.Id)] = struct{}{}
			}
		}
		return sourceMap
	},
}

var PurchaseAccountLine = DependentDualType[quickbooks.Purchase]{
	dependentBaseType: dependentBaseType[quickbooks.Purchase]{
		BaseType: fibery.BaseType{
			TypeId:   "PurchaseAccountLine",
			TypeName: "Purchase Account Line",
			TypeSchema: map[string]fibery.Field{
				"id": {
					Name: "id",
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
				"PurchaseId": {
					Name: "Purchase ID",
					Type: fibery.Text,
					Relation: &fibery.Relation{
						Cardinality:   fibery.MTO,
						Name:          "Expense",
						TargetName:    "Item Lines",
						TargetType:    "Purchase",
						TargetFieldID: "id",
					},
				},
				"AccountId": {
					Name: "Account ID",
					Type: fibery.Text,
					Relation: &fibery.Relation{
						Cardinality:   fibery.MTO,
						Name:          "Category",
						TargetName:    "Expense Account Lines",
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
						TargetName:    "Expense Item Lines",
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
						TargetName:    "Expense Item Lines",
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
						TargetName:    "Expense Item Line Markup",
						TargetType:    "Account",
						TargetFieldID: "id",
					},
				},
			},
		},
		schemaGen: func(p quickbooks.Purchase) ([]map[string]any, error) {
			items := []map[string]any{}
			for _, line := range p.Line {
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
						"id":              fmt.Sprintf("%s:a:%s", p.Id, line.Id),
						"QBOId":           line.Id,
						"Description":     line.Description,
						"__syncAction":    fibery.SET,
						"LineNum":         line.LineNum,
						"Tax":             tax,
						"Billable":        billable,
						"Billed":          billed,
						"MarkupPercent":   line.AccountBasedExpenseLineDetail.MarkupInfo.Percent,
						"Amount":          line.Amount,
						"PurchaseId":      p.Id,
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
	sourceType: &Purchase,
	sourceId: func(p quickbooks.Purchase) string {
		return p.Id
	},
	sourceStatus: func(p quickbooks.Purchase) string {
		return p.Status
	},
	sourceMapper: func(p quickbooks.Purchase) map[string]struct{} {
		sourceMap := map[string]struct{}{}
		for _, line := range p.Line {
			if line.DetailType == quickbooks.AccountExpenseLine {
				sourceMap[fmt.Sprintf("%s:%s", p.Id, line.Id)] = struct{}{}
			}
		}
		return sourceMap
	},
}

func init() {
	registerType(&Purchase)
	registerType(&PurchaseItemLine)
	registerType(&PurchaseAccountLine)
}
