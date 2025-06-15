package temp

import (
	"encoding/json"
	"fmt"
	"math"
	"strings"

	"github.com/tommyhedley/fibery-quickbooks-app/pkgs/fibery"
	"github.com/tommyhedley/quickbooks-go"
)

type TypeRegistry map[string]Type

func NewTypeRegistry() TypeRegistry {
	tr := TypeRegistry{}
	BuildTypes(tr)
	return tr
}

func (tr TypeRegistry) Register(t Type) {
	tr[t.Id()] = t
}

func (tr TypeRegistry) GetType(id string) (Type, bool) {
	if typ, exists := tr[id]; exists {
		return typ, true
	}
	return nil, false
}

func BuildTypes(tr TypeRegistry) {
	account := NewQuickBooksDualType(
		"Account",
		"Account",
		map[string]fibery.Field{
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
			"AcctNum": {
				Name: "Account Number",
				Type: fibery.Text,
			},
			"CurrentBalance": {
				Name: "Balance",
				Type: fibery.Number,
				Format: map[string]any{
					"format":               "Money",
					"currencyCode":         "USD",
					"hasThousandSeperator": true,
					"precision":            2,
				},
			},
			"CurrentBalanceWithSubAccounts": {
				Name: "Balance With Sub-Accounts",
				Type: fibery.Number,
				Format: map[string]any{
					"format":               "Money",
					"currencyCode":         "USD",
					"hasThousandSeperator": true,
					"precision":            2,
				},
			},
			"Classification": {
				Name:     "Classification",
				Type:     fibery.Text,
				SubType:  fibery.SingleSelect,
				ReadOnly: true,
				Options: []map[string]any{
					{
						"name": "Asset",
					},
					{
						"name": "Equity",
					},
					{
						"name": "Expense",
					},
					{
						"name": "Liability",
					},
					{
						"name": "Revenue",
					},
				},
			},
			"AccountType": {
				Name:     "Account Type",
				Type:     fibery.Text,
				SubType:  fibery.SingleSelect,
				ReadOnly: true,
			},
			"AccountSubType": {
				Name:     "Account Sub-Type",
				Type:     fibery.Text,
				SubType:  fibery.SingleSelect,
				ReadOnly: true,
			},
			"ParentAccountId": {
				Name: "Parent Account ID",
				Type: fibery.Text,
				Relation: &fibery.Relation{
					Cardinality:   fibery.MTO,
					Name:          "Parent Account",
					TargetName:    "Sub-Accounts",
					TargetType:    "Account",
					TargetFieldID: "id",
				},
			},
		},
		func(a quickbooks.Account) (map[string]any, error) {
			var parentAccountId string
			if a.ParentRef != nil {
				parentAccountId = a.ParentRef.Value
			}

			return map[string]any{
				"id":                            a.Id,
				"QBOId":                         a.Id,
				"Name":                          a.Name,
				"FullyQualifiedName":            a.FullyQualifiedName,
				"SyncToken":                     a.SyncToken,
				"__syncAction":                  fibery.SET,
				"Active":                        a.Active,
				"Description":                   a.Description,
				"AcctNum":                       a.AcctNum,
				"CurrentBalance":                a.CurrentBalance,
				"CurrentBalanceWithSubAccounts": a.CurrentBalanceWithSubAccounts,
				"Classification":                a.Classification,
				"AccountType":                   a.AccountType,
				"AccountSubType":                a.AccountSubType,
				"ParentAccountId":               parentAccountId,
			}, nil
		},
		func(client *quickbooks.Client, requestParams quickbooks.RequestParameters, startPosition, pageSize int) ([]quickbooks.Account, error) {
			items, err := client.FindAccountsByPage(requestParams, startPosition, pageSize)
			if err != nil {
				return nil, err
			}

			return items, nil
		},
		func(a quickbooks.Account) string {
			return a.Id
		},
		func(a quickbooks.Account) string {
			return a.Status
		},
	)

	tr.Register(account)

	bill := NewQuickBooksDualType(
		"Bill",
		"Bill",
		map[string]fibery.Field{
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
			"Files": {
				Name:    "Files",
				Type:    fibery.TextArray,
				SubType: fibery.File,
			},
		},
		func(b quickbooks.Bill) (map[string]any, error) {
			var apAccountId string
			if b.APAccountRef != nil {
				apAccountId = b.APAccountRef.Value
			}

			var salesTermId string
			if b.SalesTermRef != nil {
				salesTermId = b.SalesTermRef.Value
			}

			var txnDate string
			if !b.TxnDate.IsZero() {
				txnDate = b.TxnDate.Format(fibery.DateFormat)
			}

			var dueDate string
			if !b.DueDate.IsZero() {
				dueDate = b.DueDate.Format(fibery.DateFormat)
			}

			var name string
			if b.PrivateNote == "" {
				name = b.VendorRef.Name
			} else {
				name = b.VendorRef.Name + " - " + b.PrivateNote
			}

			return map[string]any{
				"id":           b.Id,
				"QBOId":        b.Id,
				"Name":         name,
				"SyncToken":    b.SyncToken,
				"__syncAction": fibery.SET,
				"DocNumber":    b.DocNumber,
				"TxnDate":      txnDate,
				"DueDate":      dueDate,
				"PrivateNote":  b.PrivateNote,
				"TotalAmt":     b.TotalAmt,
				"Balance":      b.Balance,
				"VendorId":     b.VendorRef.Value,
				"APAccountId":  apAccountId,
				"SalesTermId":  salesTermId,
			}, nil
		},
		func(client *quickbooks.Client, requestParams quickbooks.RequestParameters, startPosition, pageSize int) ([]quickbooks.Bill, error) {
			items, err := client.FindBillsByPage(requestParams, startPosition, pageSize)
			if err != nil {
				return nil, err
			}

			return items, nil
		},
		func(b quickbooks.Bill) string {
			return b.Id
		},
		func(b quickbooks.Bill) string {
			return b.Status
		},
	)

	tr.Register(bill)

	billItemLine := NewDependentDualType(
		"BillItemLine",
		"Bill Item Line",
		map[string]fibery.Field{
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
			"Description": {
				Name: "Description",
				Type: fibery.Text,
			},
			"__syncAction": {
				Type: fibery.Text,
				Name: "Sync Action",
			},
			"LineNum": {
				Name:    "Line",
				Type:    fibery.Number,
				SubType: fibery.Integer,
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
			"ReimburseChargeId": {
				Name: "Reimburse Charge ID",
				Type: fibery.Text,
				Relation: &fibery.Relation{
					Cardinality:   fibery.OTO,
					Name:          "Reimburse Charge",
					TargetName:    "Bill Item Line",
					TargetType:    "ReimburseCharge",
					TargetFieldID: "id",
				},
			},
		},
		func(b quickbooks.Bill) ([]map[string]any, error) {
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

					var reimburseChargeId string
					for _, txn := range line.LinkedTxn {
						if txn.TxnType == "ReimburseCharge" {
							reimburseChargeId = txn.TxnId
						}
					}

					var name string
					if line.Description == "" {
						name = line.ItemBasedExpenseLineDetail.ItemRef.Name
					} else {
						name = line.ItemBasedExpenseLineDetail.ItemRef.Name + line.Description
					}

					item := map[string]any{
						"id":                fmt.Sprintf("%s:i:%s", b.Id, line.Id),
						"QBOId":             line.Id,
						"Name":              name,
						"Description":       line.Description,
						"__syncAction":      fibery.SET,
						"LineNum":           line.LineNum,
						"Tax":               tax,
						"Billable":          billable,
						"Billed":            billed,
						"Qty":               line.ItemBasedExpenseLineDetail.Qty,
						"UnitPrice":         line.ItemBasedExpenseLineDetail.UnitPrice,
						"MarkupPercent":     line.ItemBasedExpenseLineDetail.MarkupInfo.Percent,
						"Amount":            line.Amount,
						"BillId":            b.Id,
						"ItemId":            line.ItemBasedExpenseLineDetail.ItemRef.Value,
						"CustomerId":        line.AccountBasedExpenseLineDetail.CustomerRef.Value,
						"ClassId":           line.ItemBasedExpenseLineDetail.ClassRef.Value,
						"MarkupAccountId":   line.ItemBasedExpenseLineDetail.MarkupInfo.MarkUpIncomeAccountRef.Value,
						"ReimburseChargeId": reimburseChargeId,
					}
					items = append(items, item)
				}
			}
			return items, nil
		},
		bill,
		func(b quickbooks.Bill) string {
			return b.Id
		},
		func(b quickbooks.Bill) string {
			return b.Status
		},
		func(b quickbooks.Bill) map[string]struct{} {
			sourceMap := map[string]struct{}{}
			for _, line := range b.Line {
				if line.DetailType == quickbooks.ItemExpenseLine {
					sourceMap[line.Id] = struct{}{}
				}
			}
			return sourceMap
		},
	)

	tr.Register(billItemLine)

	billAccountLine := NewDependentDualType(
		"BillAccountLine",
		"Bill Account Line",
		map[string]fibery.Field{
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
			"Description": {
				Name: "Description",
				Type: fibery.Text,
			},
			"__syncAction": {
				Type: fibery.Text,
				Name: "Sync Action",
			},
			"LineNum": {
				Name:    "Line",
				Type:    fibery.Number,
				SubType: fibery.Integer,
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
					TargetName:    "Account Lines",
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
					TargetName:    "Bill Account Lines",
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
					TargetName:    "Bill Account Lines",
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
					TargetName:    "Bill Account Line Markup",
					TargetType:    "Account",
					TargetFieldID: "id",
				},
			},
			"ReimburseChargeId": {
				Name: "Reimburse Charge ID",
				Type: fibery.Text,
				Relation: &fibery.Relation{
					Cardinality:   fibery.OTO,
					Name:          "Reimburse Charge",
					TargetName:    "Bill Account Line",
					TargetType:    "ReimburseCharge",
					TargetFieldID: "id",
				},
			},
		},
		func(b quickbooks.Bill) ([]map[string]any, error) {
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

					var reimburseChargeId string
					for _, txn := range line.LinkedTxn {
						if txn.TxnType == "ReimburseCharge" {
							reimburseChargeId = txn.TxnId
						}
					}

					var name string
					if line.Description == "" {
						name = line.AccountBasedExpenseLineDetail.AccountRef.Name
					} else {
						name = line.AccountBasedExpenseLineDetail.AccountRef.Name + " - " + line.Description
					}

					item := map[string]any{
						"id":                fmt.Sprintf("%s:a:%s", b.Id, line.Id),
						"QBOId":             line.Id,
						"Name":              name,
						"Description":       line.Description,
						"__syncAction":      fibery.SET,
						"LineNum":           line.LineNum,
						"Tax":               tax,
						"Billable":          billable,
						"Billed":            billed,
						"MarkupPercent":     line.AccountBasedExpenseLineDetail.MarkupInfo.Percent,
						"Amount":            line.Amount,
						"BillId":            b.Id,
						"AccountId":         line.AccountBasedExpenseLineDetail.AccountRef.Value,
						"CustomerId":        line.AccountBasedExpenseLineDetail.CustomerRef.Value,
						"ClassId":           line.AccountBasedExpenseLineDetail.ClassRef.Value,
						"MarkupAccountId":   line.AccountBasedExpenseLineDetail.MarkupInfo.MarkUpIncomeAccountRef.Value,
						"ReimburseChargeId": reimburseChargeId,
					}

					items = append(items, item)
				}
			}
			return items, nil
		},
		bill,
		func(b quickbooks.Bill) string {
			return b.Id
		},
		func(b quickbooks.Bill) string {
			return b.Status
		},
		func(b quickbooks.Bill) map[string]struct{} {
			sourceMap := map[string]struct{}{}
			for _, line := range b.Line {
				if line.DetailType == quickbooks.AccountExpenseLine {
					sourceMap[line.Id] = struct{}{}
				}
			}
			return sourceMap
		},
	)

	tr.Register(billAccountLine)

	billPayment := NewQuickBooksDualType(
		"BillPayment",
		"Bill Payment",
		map[string]fibery.Field{
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
				Name: "Reference Number",
				Type: fibery.Text,
			},
			"TxnDate": {
				Name:    "Payment Date",
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
			"PayType": {
				Name:     "Payment Type",
				Type:     fibery.Text,
				SubType:  fibery.SingleSelect,
				ReadOnly: true,
				Options: []map[string]any{
					{
						"name": "Check",
					},
					{
						"name": "Credit Card",
					},
				},
			},
			"VendorId": {
				Name: "Vendor ID",
				Type: fibery.Text,
				Relation: &fibery.Relation{
					Cardinality:   fibery.MTO,
					Name:          "Vendor",
					TargetName:    "Bill Payments",
					TargetType:    "Vendor",
					TargetFieldID: "id",
				},
			},
			"PaymentAccountId": {
				Name: "Payment Account ID",
				Type: fibery.Text,
				Relation: &fibery.Relation{
					Cardinality:   fibery.MTO,
					Name:          "Payment Account",
					TargetName:    "Bill Payments",
					TargetType:    "Account",
					TargetFieldID: "id",
				},
			},
		},
		func(bp quickbooks.BillPayment) (map[string]any, error) {
			var paymentAccountId string
			if bp.APAccountRef != nil {
				paymentAccountId = bp.APAccountRef.Value
			}

			var payType string
			switch bp.PayType {
			case quickbooks.CreditCardPaymentType:
				payType = "Credit Card"
				paymentAccountId = bp.CreditCardPayment.CCAccountRef.Value
			case quickbooks.CheckPaymentType:
				payType = "Check"
				paymentAccountId = bp.CheckPayment.BankAccountRef.Value
			}

			var txnDate string
			if !bp.TxnDate.IsZero() {
				txnDate = bp.TxnDate.Format(fibery.DateFormat)
			}

			return map[string]any{
				"id":               bp.Id,
				"QBOId":            bp.Id,
				"Name":             bp.VendorRef.Name,
				"SyncToken":        bp.SyncToken,
				"__syncAction":     fibery.SET,
				"DocNumber":        bp.DocNumber,
				"TxnDate":          txnDate,
				"PrivateNote":      bp.PrivateNote,
				"TotalAmt":         bp.TotalAmt,
				"PayType":          payType,
				"VendorId":         bp.VendorRef.Value,
				"PaymentAccountId": paymentAccountId,
			}, nil
		},
		func(client *quickbooks.Client, requestParams quickbooks.RequestParameters, startPosition, pageSize int) ([]quickbooks.BillPayment, error) {
			items, err := client.FindBillPaymentsByPage(requestParams, startPosition, pageSize)
			if err != nil {
				return nil, err
			}

			return items, nil
		},
		func(bp quickbooks.BillPayment) string {
			return bp.Id
		},
		func(bp quickbooks.BillPayment) string {
			return bp.Status
		},
	)

	tr.Register(billPayment)

	billPaymentLine := NewDependentDualType(
		"BillPaymentLine",
		"Bill Payment Line",
		map[string]fibery.Field{
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
					TargetName:    "Bill Payment Lines",
					TargetType:    "Bill",
					TargetFieldID: "id",
				},
			},
			"VendorCreditId": {
				Name: "Vendor Credit ID",
				Type: fibery.Text,
				Relation: &fibery.Relation{
					Cardinality:   fibery.MTO,
					Name:          "Vendor Credit",
					TargetName:    "Bill Payment Lines",
					TargetType:    "VendorCredit",
					TargetFieldID: "id",
				},
			},
			"DepositId": {
				Name: "Deposit ID",
				Type: fibery.Text,
				Relation: &fibery.Relation{
					Cardinality:   fibery.MTO,
					Name:          "Deposit",
					TargetName:    "Bill Payment Lines",
					TargetType:    "Deposit",
					TargetFieldID: "id",
				},
			},
		},
		func(bp quickbooks.BillPayment) ([]map[string]any, error) {
			items := []map[string]any{}
			for _, line := range bp.Line {
				var description string
				var billId string
				var vendorCreditId string
				switch line.LinkedTxn[0].TxnType {
				case "Bill":
					description = "Bill Payment"
					billId = line.LinkedTxn[0].TxnId
				case "VendorCredit":
					description = "Vendor Credit"
					vendorCreditId = line.LinkedTxn[0].TxnId
				}

				item := map[string]any{
					"id":             fmt.Sprintf("%s:%s", bp.Id, line.Id),
					"QBOId":          line.Id,
					"Description":    description,
					"__syncAction":   fibery.SET,
					"Amount":         line.Amount,
					"BillId":         billId,
					"VendorCreditId": vendorCreditId,
				}

				items = append(items, item)
			}
			return items, nil
		},
		billPayment,
		func(bp quickbooks.BillPayment) string {
			return bp.Id
		},
		func(bp quickbooks.BillPayment) string {
			return bp.Status
		},
		func(bp quickbooks.BillPayment) map[string]struct{} {
			sourceMap := map[string]struct{}{}
			for _, line := range bp.Line {
				sourceMap[line.Id] = struct{}{}
			}
			return sourceMap
		},
	)

	tr.Register(billPaymentLine)

	class := NewQuickBooksDualType(
		"Class",
		"Class",
		map[string]fibery.Field{
			"id": {
				Name: "ID",
				Type: fibery.Id,
			},
			"QBOId": {
				Name: "QBOId",
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
			"ParentClassId": {
				Name: "Parent Class ID",
				Type: fibery.Text,
				Relation: &fibery.Relation{
					Cardinality:   fibery.MTO,
					Name:          "Parent Class",
					TargetName:    "Sub-Classes",
					TargetType:    "Class",
					TargetFieldID: "id",
				},
			},
		},
		func(c quickbooks.Class) (map[string]any, error) {
			return map[string]any{
				"id":                 c.Id,
				"QBOId":              c.Id,
				"Name":               c.Name,
				"FullyQualifiedName": c.FullyQualifiedName,
				"SyncToken":          c.SyncToken,
				"__syncAction":       fibery.SET,
				"Active":             c.Active,
				"ParentClassId":      c.ParentRef.Value,
			}, nil
		},
		func(client *quickbooks.Client, requestParams quickbooks.RequestParameters, startPosition, pageSize int) ([]quickbooks.Class, error) {
			items, err := client.FindClassesByPage(requestParams, startPosition, pageSize)
			if err != nil {
				return nil, err
			}

			return items, nil
		},
		func(c quickbooks.Class) string {
			return c.Id
		},
		func(c quickbooks.Class) string {
			return c.Status
		},
	)

	tr.Register(class)

	customer := NewQuickBooksDualType(
		"Customer",
		"Customer",
		map[string]fibery.Field{
			"id": {
				Name: "ID",
				Type: fibery.Id,
			},
			"QBOId": {
				Name: "QBO ID",
				Type: fibery.Text,
			},
			"DisplayName": {
				Name:    "Display Name",
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
			"Title": {
				Name: "Title",
				Type: fibery.Text,
			},
			"GivenName": {
				Name: "First Name",
				Type: fibery.Text,
			},
			"MiddleName": {
				Name: "Middle Name",
				Type: fibery.Text,
			},
			"FamilyName": {
				Name: "Last Name",
				Type: fibery.Text,
			},
			"Suffix": {
				Name: "Suffix",
				Type: fibery.Text,
			},
			"CompanyName": {
				Name: "Company Name",
				Type: fibery.Text,
			},
			"PrimaryEmail": {
				Name:    "Email",
				Type:    fibery.Text,
				SubType: fibery.Email,
			},
			"Taxable": {
				Name:    "Taxable",
				Type:    fibery.Text,
				SubType: fibery.Boolean,
			},
			"ResaleNum": {
				Name: "Resale ID",
				Type: fibery.Text,
			},
			"PrimaryPhone": {
				Name: "Phone",
				Type: fibery.Text,
				Format: map[string]any{
					"format": "phone",
				},
			},
			"AlternatePhone": {
				Name: "Alternate Phone",
				Type: fibery.Text,
				Format: map[string]any{
					"format": "phone",
				},
			},
			"Mobile": {
				Name: "Mobile",
				Type: fibery.Text,
				Format: map[string]any{
					"format": "phone",
				},
			},
			"Fax": {
				Name: "Fax",
				Type: fibery.Text,
				Format: map[string]any{
					"format": "phone",
				},
			},
			"Job": {
				Name:    "Job",
				Type:    fibery.Text,
				SubType: fibery.Boolean,
			},
			"BillWithParent": {
				Name:    "Bill With Parent",
				Type:    fibery.Text,
				SubType: fibery.Boolean,
			},
			"Notes": {
				Name:    "Notes",
				Type:    fibery.Text,
				SubType: fibery.MD,
			},
			"Website": {
				Name:    "Website",
				Type:    fibery.Text,
				SubType: fibery.URL,
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
			"BalanceWithJobs": {
				Name: "Balance With Jobs",
				Type: fibery.Number,
				Format: map[string]any{
					"format":               "Money",
					"currencyCode":         "USD",
					"hasThousandSeperator": true,
					"precision":            2,
				},
			},
			"ShippingLine1": {
				Name: "Shipping Line 1",
				Type: fibery.Text,
			},
			"ShippingLine2": {
				Name: "Shipping Line 2",
				Type: fibery.Text,
			},
			"ShippingLine3": {
				Name: "Shipping Line 3",
				Type: fibery.Text,
			},
			"ShippingLine4": {
				Name: "Shipping Line 4",
				Type: fibery.Text,
			},
			"ShippingLine5": {
				Name: "Shipping Line 5",
				Type: fibery.Text,
			},
			"ShippingCity": {
				Name: "Shipping City",
				Type: fibery.Text,
			},
			"ShippingState": {
				Name: "Shipping State",
				Type: fibery.Text,
			},
			"ShippingPostalCode": {
				Name: "Shipping Postal Code",
				Type: fibery.Text,
			},
			"ShippingCountry": {
				Name: "Shipping Country",
				Type: fibery.Text,
			},
			"ShippingLat": {
				Name: "Shipping Latitude",
				Type: fibery.Text,
			},
			"ShippingLong": {
				Name: "Shipping Longitude",
				Type: fibery.Text,
			},
			"BillingLine1": {
				Name: "Billing Line 1",
				Type: fibery.Text,
			},
			"BillingLine2": {
				Name: "Billing Line 2",
				Type: fibery.Text,
			},
			"BillingLine3": {
				Name: "Billing Line 3",
				Type: fibery.Text,
			},
			"BillingLine4": {
				Name: "Billing Line 4",
				Type: fibery.Text,
			},
			"BillingLine5": {
				Name: "Billing Line 5",
				Type: fibery.Text,
			},
			"BillingCity": {
				Name: "Billing City",
				Type: fibery.Text,
			},
			"BillingState": {
				Name: "Billing State",
				Type: fibery.Text,
			},
			"BillingPostalCode": {
				Name: "Billing Postal Code",
				Type: fibery.Text,
			},
			"BillingCountry": {
				Name: "Billing Country",
				Type: fibery.Text,
			},
			"BillingLat": {
				Name: "Billing Latitude",
				Type: fibery.Text,
			},
			"BillingLong": {
				Name: "Billing Longitude",
				Type: fibery.Text,
			},
			"TaxExemptionId": {
				Name: "Tax Exemption ID",
				Type: fibery.Text,
				Relation: &fibery.Relation{
					Cardinality:   fibery.MTO,
					Name:          "Tax Exemption",
					TargetName:    "Customers",
					TargetType:    "TaxExemption",
					TargetFieldID: "id",
				},
			},
			"DefaultTaxCodeId": {
				Name: "Default Tax Code ID",
				Type: fibery.Text,
				Relation: &fibery.Relation{
					Cardinality:   fibery.MTO,
					Name:          "Default Tax Code",
					TargetName:    "Customers",
					TargetType:    "TaxCode",
					TargetFieldID: "id",
				},
			},
			"CustomerTypeId": {
				Name: "Customer Type ID",
				Type: fibery.Text,
				Relation: &fibery.Relation{
					Cardinality:   fibery.MTO,
					Name:          "Customer Type",
					TargetName:    "Customers",
					TargetType:    "CustomerType",
					TargetFieldID: "id",
				},
			},
			"SalesTermId": {
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
			"PaymentMethodId": {
				Name: "Payment Method ID",
				Type: fibery.Text,
				Relation: &fibery.Relation{
					Cardinality:   fibery.MTO,
					Name:          "Payment Method",
					TargetName:    "Customers",
					TargetType:    "PaymentMethod",
					TargetFieldID: "id",
				},
			},
			"ParentId": {
				Name: "Parent ID",
				Type: fibery.Text,
				Relation: &fibery.Relation{
					Cardinality:   fibery.MTO,
					Name:          "Parent",
					TargetName:    "Jobs",
					TargetType:    "Customer",
					TargetFieldID: "id",
				},
			},
		},
		func(c quickbooks.Customer) (map[string]any, error) {
			var email string
			if c.PrimaryEmailAddr != nil {
				email = c.PrimaryEmailAddr.Address
			}

			var primaryPhone string
			if c.PrimaryPhone != nil {
				primaryPhone = c.PrimaryPhone.FreeFormNumber
			}

			var alternatePhone string
			if c.AlternatePhone != nil {
				alternatePhone = c.AlternatePhone.FreeFormNumber
			}

			var mobile string
			if c.Mobile != nil {
				mobile = c.Mobile.FreeFormNumber
			}

			var fax string
			if c.Fax != nil {
				fax = c.Fax.FreeFormNumber
			}

			var website string
			if c.WebAddr != nil {
				website = c.WebAddr.URI
			}

			var shipAddr quickbooks.PhysicalAddress
			if c.ShipAddr != nil {
				shipAddr = *c.ShipAddr
			}

			var billAddr quickbooks.PhysicalAddress
			if c.BillAddr != nil {
				billAddr = *c.BillAddr
			}

			job := false
			if c.Job.Valid {
				job = c.Job.Bool
			}

			var defaultTaxCodeId string
			if c.DefaultTaxCodeRef != nil {
				defaultTaxCodeId = c.DefaultTaxCodeRef.Value
			}

			var customerTypeId string
			if c.CustomerTypeRef != nil {
				customerTypeId = c.CustomerTypeRef.Value
			}

			var salesTermId string
			if c.SalesTermRef != nil {
				salesTermId = c.SalesTermRef.Value
			}

			var paymentMethodId string
			if c.PaymentMethodRef != nil {
				paymentMethodId = c.PaymentMethodRef.Value
			}

			var parentId string
			if c.ParentRef != nil {
				parentId = c.ParentRef.Value
			}

			return map[string]any{
				"id":                 c.Id,
				"QBOId":              c.Id,
				"DisplayName":        c.DisplayName,
				"SyncToken":          c.SyncToken,
				"__syncAction":       fibery.SET,
				"Active":             c.Active,
				"Title":              c.Title,
				"GivenName":          c.GivenName,
				"MiddleName":         c.MiddleName,
				"FamilyName":         c.FamilyName,
				"Suffix":             c.Suffix,
				"CompanyName":        c.CompanyName,
				"PrimaryEmail":       email,
				"Taxable":            c.Taxable,
				"ResaleNum":          c.ResaleNum,
				"PrimaryPhone":       primaryPhone,
				"AlternatePhone":     alternatePhone,
				"Mobile":             mobile,
				"Fax":                fax,
				"Job":                job,
				"BillWithParent":     c.BillWithParent,
				"Notes":              c.Notes,
				"Website":            website,
				"Balance":            c.Balance,
				"BalanceWithJobs":    c.BalanceWithJobs,
				"ShippingLine1":      shipAddr.Line1,
				"ShippingLine2":      shipAddr.Line2,
				"ShippingLine3":      shipAddr.Line3,
				"ShippingLine4":      shipAddr.Line4,
				"ShippingLine5":      shipAddr.Line5,
				"ShippingCity":       shipAddr.City,
				"ShippingState":      shipAddr.CountrySubDivisionCode,
				"ShippingPostalCode": shipAddr.PostalCode,
				"ShippingCountry":    shipAddr.Country,
				"ShippingLat":        shipAddr.Lat,
				"ShippingLong":       shipAddr.Long,
				"BillingLine1":       billAddr.Line1,
				"BillingLine2":       billAddr.Line2,
				"BillingLine3":       billAddr.Line3,
				"BillingLine4":       billAddr.Line4,
				"BillingLine5":       billAddr.Line5,
				"BillingCity":        billAddr.City,
				"BillingState":       billAddr.CountrySubDivisionCode,
				"BillingPostalCode":  billAddr.PostalCode,
				"BillingCountry":     billAddr.Country,
				"BillingLat":         billAddr.Lat,
				"BillingLong":        billAddr.Long,
				"TaxExemptionId":     c.TaxExemptionReasonId,
				"DefaultTaxCodeId":   defaultTaxCodeId,
				"CustomerTypeId":     customerTypeId,
				"SalesTermId":        salesTermId,
				"PaymentMethodId":    paymentMethodId,
				"ParentId":           parentId,
			}, nil
		},
		func(client *quickbooks.Client, requestParams quickbooks.RequestParameters, startPosition, pageSize int) ([]quickbooks.Customer, error) {
			items, err := client.FindCustomersByPage(requestParams, startPosition, pageSize)
			if err != nil {
				return nil, err
			}

			return items, nil
		},
		func(c quickbooks.Customer) string {
			return c.Id
		},
		func(c quickbooks.Customer) string {
			return c.Status
		},
	)

	tr.Register(customer)

	customerType := NewQuickBooksCDCType(
		"CustomerType",
		"Customer Type",
		map[string]fibery.Field{
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
			"Active": {
				Name:    "Active",
				Type:    fibery.Text,
				SubType: fibery.Boolean,
			},
		},
		func(ct quickbooks.CustomerType) (map[string]any, error) {
			return map[string]any{
				"id":           ct.Id,
				"QBOId":        ct.Id,
				"Name":         ct.Name,
				"SyncToken":    ct.SyncToken,
				"__syncAction": fibery.SET,
				"Active":       ct.Active,
			}, nil
		},
		func(client *quickbooks.Client, requestParams quickbooks.RequestParameters, startPosition, pageSize int) ([]quickbooks.CustomerType, error) {
			items, err := client.FindCustomerTypesByPage(requestParams, startPosition, pageSize)
			if err != nil {
				return nil, err
			}

			return items, nil

		},
		func(ct quickbooks.CustomerType) string {
			return ct.Id
		},
		func(ct quickbooks.CustomerType) string {
			return ct.Status
		},
	)

	tr.Register(customerType)

	employee := NewQuickBooksDualType(
		"Employee",
		"Employee",
		map[string]fibery.Field{
			"id": {
				Name: "ID",
				Type: fibery.Id,
			},
			"QBOId": {
				Name: "QBO ID",
				Type: fibery.Text,
			},
			"DisplayName": {
				Name:    "Display Name",
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
			"Title": {
				Name: "Title",
				Type: fibery.Text,
			},
			"GivenName": {
				Name: "First Name",
				Type: fibery.Text,
			},
			"MiddleName": {
				Name: "Middle Name",
				Type: fibery.Text,
			},
			"FamilyName": {
				Name: "Last Name",
				Type: fibery.Text,
			},
			"Suffix": {
				Name: "Suffix",
				Type: fibery.Text,
			},
			"PrimaryEmailAddr": {
				Name:    "Email",
				Type:    fibery.Text,
				SubType: fibery.Email,
			},
			"BillableTime": {
				Name:        "Billable",
				Type:        fibery.Text,
				SubType:     fibery.Boolean,
				Description: "Is the entity enabled for use in QuickBooks?",
			},
			"BirthDate": {
				Name:    "Date of Birth",
				Type:    fibery.DateType,
				SubType: fibery.Day,
			},
			"PrimaryPhone": {
				Name: "Phone",
				Type: fibery.Text,
				Format: map[string]any{
					"format": "phone",
				},
			},
			"Mobile": {
				Name: "Mobile",
				Type: fibery.Text,
				Format: map[string]any{
					"format": "phone",
				},
			},
			"CostRate": {
				Name: "Cost Rate",
				Type: fibery.Number,
				Format: map[string]any{
					"format":               "Money",
					"currencyCode":         "USD",
					"hasThousandSeperator": true,
					"precision":            2,
				},
			},
			"BillRate": {
				Name: "Bill Rate",
				Type: fibery.Number,
				Format: map[string]any{
					"format":               "Money",
					"currencyCode":         "USD",
					"hasThousandSeperator": true,
					"precision":            2,
				},
			},
			"EmployeeNumber": {
				Name: "Employee ID",
				Type: fibery.Text,
			},
			"AddressLine1": {
				Name: "Address Line 1",
				Type: fibery.Text,
			},
			"AddressLine2": {
				Name: "Address Line 2",
				Type: fibery.Text,
			},
			"AddressLine3": {
				Name: "Address Line 3",
				Type: fibery.Text,
			},
			"AddressLine4": {
				Name: "Address Line 4",
				Type: fibery.Text,
			},
			"AddressLine5": {
				Name: "Address Line 5",
				Type: fibery.Text,
			},
			"AddressCity": {
				Name: "Address City",
				Type: fibery.Text,
			},
			"AddressState": {
				Name: "Address State",
				Type: fibery.Text,
			},
			"AddressPostalCode": {
				Name: "Address Postal Code",
				Type: fibery.Text,
			},
			"AddressCountry": {
				Name: "Address Country",
				Type: fibery.Text,
			},
			"AddressLat": {
				Name: "Address Latitude",
				Type: fibery.Text,
			},
			"AddressLong": {
				Name: "Address Longitude",
				Type: fibery.Text,
			},
		},
		func(e quickbooks.Employee) (map[string]any, error) {
			var email string
			if e.PrimaryEmailAddr != nil {
				email = e.PrimaryEmailAddr.Address
			}

			var primaryPhone string
			if e.PrimaryPhone != nil {
				primaryPhone = e.PrimaryPhone.FreeFormNumber
			}

			var mobile string
			if e.Mobile != nil {
				mobile = e.Mobile.FreeFormNumber
			}

			var birthDate string
			if e.BirthDate != nil {
				birthDate = e.BirthDate.Format(fibery.DateFormat)
			}

			return map[string]any{
				"id":                e.Id,
				"QBOId":             e.Id,
				"DisplayName":       e.DisplayName,
				"SyncToken":         e.SyncToken,
				"__syncAction":      fibery.SET,
				"Active":            e.Active,
				"Title":             e.Title,
				"GivenName":         e.GivenName,
				"MiddleName":        e.MiddleName,
				"FamilyName":        e.FamilyName,
				"Suffix":            e.Suffix,
				"PrimaryEmailAddr":  email,
				"BillableTime":      e.BillableTime,
				"BirthDate":         birthDate,
				"PrimaryPhone":      primaryPhone,
				"Mobile":            mobile,
				"CostRate":          e.CostRate,
				"BillRate":          e.BillRate,
				"EmployeeNumber":    e.EmployeeNumber,
				"AddressLine1":      e.PrimaryAddr.Line1,
				"AddressLine2":      e.PrimaryAddr.Line2,
				"AddressLine3":      e.PrimaryAddr.Line3,
				"AddressLine4":      e.PrimaryAddr.Line4,
				"AddressLine5":      e.PrimaryAddr.Line1,
				"AddressCity":       e.PrimaryAddr.City,
				"AddressState":      e.PrimaryAddr.CountrySubDivisionCode,
				"AddressPostalCode": e.PrimaryAddr.PostalCode,
				"AddressCountry":    e.PrimaryAddr.Country,
				"AddressLat":        e.PrimaryAddr.Lat,
				"AddressLong":       e.PrimaryAddr.Long,
			}, nil
		},
		func(client *quickbooks.Client, requestParams quickbooks.RequestParameters, startPosition, pageSize int) ([]quickbooks.Employee, error) {
			items, err := client.FindEmployeesByPage(requestParams, startPosition, pageSize)
			if err != nil {
				return nil, err
			}

			return items, nil
		},
		func(e quickbooks.Employee) string {
			return e.Id
		},
		func(e quickbooks.Employee) string {
			return e.Status
		},
	)

	tr.Register(employee)

	estimate := NewQuickBooksDualType(
		"Estimate",
		"Estimate",
		map[string]fibery.Field{
			"id": {
				Name: "ID",
				Type: fibery.Id,
			},
			"QBOId": {
				Name: "QBO ID",
				Type: fibery.Text,
			},
			"InvoiceNum": {
				Name:    "Invoice Number",
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
			"ShippingFromLine1": {
				Name: "Sale Line 1",
				Type: fibery.Text,
			},
			"ShippingFromLine2": {
				Name: "Sale Line 2",
				Type: fibery.Text,
			},
			"ShippingFromLine3": {
				Name: "Sale Line 3",
				Type: fibery.Text,
			},
			"ShippingFromLine4": {
				Name: "Sale Line 4",
				Type: fibery.Text,
			},
			"ShippingFromLine5": {
				Name: "Sale Line 5",
				Type: fibery.Text,
			},
			"ShippingFromCity": {
				Name: "Sale City",
				Type: fibery.Text,
			},
			"ShippingFromState": {
				Name: "Sale State",
				Type: fibery.Text,
			},
			"ShippingFromPostalCode": {
				Name: "Sale Postal Code",
				Type: fibery.Text,
			},
			"ShippingFromCountry": {
				Name: "Sale Country",
				Type: fibery.Text,
			},
			"shippingFromLat": {
				Name: "Sale Latitude",
				Type: fibery.Text,
			},
			"ShippingFromLong": {
				Name: "Sale Longitude",
				Type: fibery.Text,
			},
			"ShippingLine1": {
				Name: "Shipping Line 1",
				Type: fibery.Text,
			},
			"ShippingLine2": {
				Name: "Shipping Line 2",
				Type: fibery.Text,
			},
			"ShippingLine3": {
				Name: "Shipping Line 3",
				Type: fibery.Text,
			},
			"ShippingLine4": {
				Name: "Shipping Line 4",
				Type: fibery.Text,
			},
			"ShippingLine5": {
				Name: "Shipping Line 5",
				Type: fibery.Text,
			},
			"ShippingCity": {
				Name: "Shipping City",
				Type: fibery.Text,
			},
			"ShippingState": {
				Name: "Shipping State",
				Type: fibery.Text,
			},
			"ShippingPostalCode": {
				Name: "Shipping Postal Code",
				Type: fibery.Text,
			},
			"ShippingCountry": {
				Name: "Shipping Country",
				Type: fibery.Text,
			},
			"ShippingLat": {
				Name: "Shipping Latitude",
				Type: fibery.Text,
			},
			"ShippingLong": {
				Name: "Shipping Longitude",
				Type: fibery.Text,
			},
			"BillingLine1": {
				Name: "Billing Line 1",
				Type: fibery.Text,
			},
			"BillingLine2": {
				Name: "Billing Line 2",
				Type: fibery.Text,
			},
			"BillingLine3": {
				Name: "Billing Line 3",
				Type: fibery.Text,
			},
			"BillingLine4": {
				Name: "Billing Line 4",
				Type: fibery.Text,
			},
			"BillingLine5": {
				Name: "Billing Line 5",
				Type: fibery.Text,
			},
			"BillingCity": {
				Name: "Billing City",
				Type: fibery.Text,
			},
			"BillingState": {
				Name: "Billing State",
				Type: fibery.Text,
			},
			"BillingPostalCode": {
				Name: "Billing Postal Code",
				Type: fibery.Text,
			},
			"BillingCountry": {
				Name: "Billing Country",
				Type: fibery.Text,
			},
			"BillingLat": {
				Name: "Billing Latitude",
				Type: fibery.Text,
			},
			"BillingLong": {
				Name: "Billing Longitude",
				Type: fibery.Text,
			},
			"DocNumber": {
				Name: "Number",
				Type: fibery.Text,
			},
			"TxnStatus": {
				Name:     "Status",
				Type:     fibery.Text,
				SubType:  fibery.SingleSelect,
				ReadOnly: true,
				Options: []map[string]any{
					{
						"name":    "Pending",
						"default": true,
					},
					{
						"name": "Accepted",
					},
					{
						"name": "Closed",
					},
					{
						"name": "Rejected",
					},
				},
			},
			"AcceptedBy": {
				Name: "Accepted By",
				Type: fibery.Text,
			},
			"AcceptedDate": {
				Name:    "Accepted Date",
				Type:    fibery.DateType,
				SubType: fibery.Day,
			},
			"ExpirationDate": {
				Name:    "Expiration Date",
				Type:    fibery.DateType,
				SubType: fibery.Day,
			},
			"Email": {
				Name:    "To",
				Type:    fibery.Text,
				SubType: fibery.Email,
			},
			"EmailCC": {
				Name:    "CC",
				Type:    fibery.Text,
				SubType: fibery.Email,
			},
			"EmailBCC": {
				Name:    "BCC",
				Type:    fibery.Text,
				SubType: fibery.Email,
			},
			"EmailSendLater": {
				Name:    "Send Later",
				Type:    fibery.Text,
				SubType: fibery.Boolean,
			},
			"EmailSent": {
				Name:    "Sent",
				Type:    fibery.Text,
				SubType: fibery.Boolean,
			},
			"EmailSendTime": {
				Name: "Send Time",
				Type: fibery.DateType,
			},
			"TxnDate": {
				Name:    "Date",
				Type:    fibery.DateType,
				SubType: fibery.Day,
			},
			"PrivateNote": {
				Name:    "Message on Statement",
				Type:    fibery.Text,
				SubType: fibery.MD,
			},
			"CustomerMemo": {
				Name:    "Message on Estimate",
				Type:    fibery.Text,
				SubType: fibery.MD,
			},
			"DiscountPosition": {
				Name:        "Discount Before Tax",
				Type:        fibery.Text,
				SubType:     fibery.Boolean,
				Description: "Should the discount be applied before or after the sales tax calculation? Default is false as tax should generally be calculated first before a discount is given.",
			},
			"DiscountType": {
				Name:     "Discount Type",
				Type:     fibery.Text,
				SubType:  fibery.SingleSelect,
				ReadOnly: true,
				Options: []map[string]any{
					{
						"name": "Percent",
					},
					{
						"name": "Amount",
					},
				},
			},
			"DiscountPercent": {
				Name: "Discount Percent",
				Type: fibery.Number,
				Format: map[string]any{
					"format":    "Percent",
					"precision": 2,
				},
			},
			"DiscountAmount": {
				Name: "Discount Amount",
				Type: fibery.Number,
				Format: map[string]any{
					"format":               "Money",
					"currencyCode":         "USD",
					"hasThousandSeperator": true,
					"precision":            2,
				},
			},
			"Tax": {
				Name: "Tax",
				Type: fibery.Number,
				Format: map[string]any{
					"format":               "Money",
					"currencyCode":         "USD",
					"hasThousandSeperator": true,
					"precision":            2,
				},
			},
			"SubtotalAmt": {
				Name: "Subtotal",
				Type: fibery.Number,
				Format: map[string]any{
					"format":               "Money",
					"currencyCode":         "USD",
					"hasThousandSeperator": true,
					"precision":            2,
				},
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
			"ClassId": {
				Name: "Class ID",
				Type: fibery.Text,
				Relation: &fibery.Relation{
					Cardinality:   fibery.MTO,
					Name:          "Class",
					TargetName:    "Estimates",
					TargetType:    "Class",
					TargetFieldID: "id",
				},
			},
			"TaxCodeId": {
				Name: "Tax Code ID",
				Type: fibery.Text,
				Relation: &fibery.Relation{
					Cardinality:   fibery.MTO,
					Name:          "Sales Tax Rate",
					TargetName:    "Estimates",
					TargetType:    "TaxCode",
					TargetFieldID: "id",
				},
			},
			"TaxExemptionId": {
				Name: "Tax Exemption ID",
				Type: fibery.Text,
				Relation: &fibery.Relation{
					Cardinality:   fibery.MTO,
					Name:          "Tax Exemption",
					TargetName:    "Estimates",
					TargetType:    "TaxExemption",
					TargetFieldID: "id",
				},
			},
			"CustomerId": {
				Name: "Customer ID",
				Type: fibery.Text,
				Relation: &fibery.Relation{
					Cardinality:   fibery.MTO,
					Name:          "Customer",
					TargetName:    "Estimates",
					TargetType:    "Customer",
					TargetFieldID: "id",
				},
			},
			"LinkedInvoiceId": {
				Name: "Linked Invoice ID",
				Type: fibery.Text,
				Relation: &fibery.Relation{
					Cardinality:   fibery.OTO,
					Name:          "Linked Invoice",
					TargetName:    "Linked Estimate",
					TargetType:    "Invoice",
					TargetFieldID: "id",
				},
			},
		},
		func(e quickbooks.Estimate) (map[string]any, error) {
			var discountType = map[bool]string{
				true:  "Percentage",
				false: "Amount",
			}
			var discountLine *quickbooks.Line
			var subtotalLine *quickbooks.Line
			for _, line := range e.Line {
				if line.DetailType == "DiscountLineDetail" {
					if discountLine != nil {
						return nil, fmt.Errorf("estimate %s has more than one discount line", e.Id)
					}
					discountLine = &line
				}
				if line.DetailType == "SubTotalLineDetail" {
					if subtotalLine != nil {
						return nil, fmt.Errorf("estimate %s has more than one subtotal line", e.Id)
					}
					subtotalLine = &line
				}
			}

			if subtotalLine == nil {
				return nil, fmt.Errorf("estimate %s has no subtotal lines", e.Id)
			}

			var discountTypeValue string
			var discountPercent json.Number
			var discountAmount json.Number

			if discountLine != nil {
				discountTypeValue = discountType[discountLine.DiscountLineDetail.PercentBased]
				discountPercent = discountLine.DiscountLineDetail.DiscountPercent
				discountAmount = discountLine.Amount
			}

			var emailSendTime string
			if e.DeliveryInfo != nil && !e.DeliveryInfo.DeliveryTime.IsZero() {
				emailSendTime = e.DeliveryInfo.DeliveryTime.Format(fibery.DateFormat)
			}

			var name string
			if e.CustomerRef.Name == "" {
				name = e.DocNumber
			} else {
				name = e.DocNumber + " - " + e.CustomerRef.Name
			}

			var shipAddr quickbooks.PhysicalAddress
			if e.ShipAddr != nil {
				shipAddr = *e.ShipAddr
			}

			var billAddr quickbooks.PhysicalAddress
			if e.BillAddr != nil {
				billAddr = *e.BillAddr
			}

			var acceptedDate string
			if e.AcceptedDate != nil {
				acceptedDate = e.AcceptedDate.Format(fibery.DateFormat)
			}

			var expirationDate string
			if e.ExpirationDate != nil {
				expirationDate = e.ExpirationDate.Format(fibery.DateFormat)
			}

			var billEmailCC string
			if e.BillEmailCC != nil {
				billEmailCC = e.BillEmailCC.Address
			}

			var billEmailBCC string
			if e.BillEmailBCC != nil {
				billEmailBCC = e.BillEmailCC.Address
			}

			var txnDate string
			if e.TxnDate != nil {
				txnDate = e.TxnDate.Format(fibery.DateFormat)
			}

			var totalTax json.Number
			var taxCodeId string
			if e.TxnTaxDetail != nil {
				totalTax = e.TxnTaxDetail.TotalTax
				taxCodeId = e.TxnTaxDetail.TxnTaxCodeRef.Value
			}

			var classId string
			if e.ClassRef != nil {
				classId = e.ClassRef.Value
			}

			var taxExemptionId string
			if e.TaxExemptionRef != nil {
				taxExemptionId = e.TaxExemptionRef.Value
			}

			var linkedInvoiceId string
			for _, txn := range e.LinkedTxn {
				if txn.TxnType == "Invoice" {
					linkedInvoiceId = txn.TxnId
				}
			}

			return map[string]any{
				"id":                     e.Id,
				"QBOId":                  e.Id,
				"Name":                   name,
				"SyncToken":              e.SyncToken,
				"__syncAction":           fibery.SET,
				"ShippingLine1":          shipAddr.Line1,
				"ShippingLine2":          shipAddr.Line2,
				"ShippingLine3":          shipAddr.Line3,
				"ShippingLine4":          shipAddr.Line4,
				"ShippingLine5":          shipAddr.Line5,
				"ShippingCity":           shipAddr.City,
				"ShippingState":          shipAddr.CountrySubDivisionCode,
				"ShippingPostalCode":     shipAddr.PostalCode,
				"ShippingCountry":        shipAddr.Country,
				"ShippingLat":            shipAddr.Lat,
				"ShippingLong":           shipAddr.Long,
				"ShippingFromLine1":      e.ShipFromAddr.Line1,
				"ShippingFromLine2":      e.ShipFromAddr.Line2,
				"ShippingFromLine3":      e.ShipFromAddr.Line3,
				"ShippingFromLine4":      e.ShipFromAddr.Line4,
				"ShippingFromLine5":      e.ShipFromAddr.Line5,
				"ShippingFromCity":       e.ShipFromAddr.City,
				"ShippingFromState":      e.ShipFromAddr.CountrySubDivisionCode,
				"ShippingFromPostalCode": e.ShipFromAddr.PostalCode,
				"ShippingFromCountry":    e.ShipFromAddr.Country,
				"ShippingFromLat":        e.ShipFromAddr.Lat,
				"ShippingFromLong":       e.ShipFromAddr.Long,
				"BillingLine1":           billAddr.Line1,
				"BillingLine2":           billAddr.Line2,
				"BillingLine3":           billAddr.Line3,
				"BillingLine4":           billAddr.Line4,
				"BillingLine5":           billAddr.Line5,
				"BillingCity":            billAddr.City,
				"BillingState":           billAddr.CountrySubDivisionCode,
				"BillingPostalCode":      billAddr.PostalCode,
				"BillingCountry":         billAddr.Country,
				"BillingLat":             billAddr.Lat,
				"BillingLong":            billAddr.Long,
				"DocNumber":              e.DocNumber,
				"TxnStatus":              e.TxnStatus,
				"AcceptedBy":             e.AcceptedBy,
				"AcceptedDate":           acceptedDate,
				"ExpirationDate":         expirationDate,
				"Email":                  e.BillEmail.Address,
				"EmailCC":                billEmailCC,
				"EmailBCC":               billEmailBCC,
				"EmailSendLater":         e.EmailStatus == "NeedToSend",
				"EmailSent":              e.EmailStatus == "EmailSent",
				"EmailSendTime":          emailSendTime,
				"TxnDate":                txnDate,
				"PrivateNote":            e.PrivateNote,
				"CustomerMemo":           e.CustomerMemo.Value,
				"DiscountPosition":       e.ApplyTaxAfterDiscount,
				"DiscountType":           discountTypeValue,
				"DiscountPercent":        discountPercent,
				"DiscountAmount":         discountAmount,
				"Tax":                    totalTax,
				"SubtotalAmt":            subtotalLine.Amount,
				"TotalAmt":               e.TotalAmt,
				"ClassId":                classId,
				"TaxCodeId":              taxCodeId,
				"TaxExemptionId":         taxExemptionId,
				"CustomerId":             e.CustomerRef.Value,
				"LinkedInvoiceId":        linkedInvoiceId,
			}, nil
		},
		func(client *quickbooks.Client, requestParams quickbooks.RequestParameters, startPosition, pageSize int) ([]quickbooks.Estimate, error) {
			items, err := client.FindEstimatesByPage(requestParams, startPosition, pageSize)
			if err != nil {
				return nil, err
			}

			return items, nil
		},
		func(e quickbooks.Estimate) string {
			return e.Id
		},
		func(e quickbooks.Estimate) string {
			return e.Status
		},
	)

	tr.Register(estimate)

	estimateLine := NewDependentDualType(
		"EstimateLine",
		"Estimate Line",
		map[string]fibery.Field{
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
			"Description": {
				Name: "Description",
				Type: fibery.Text,
			},
			"__syncAction": {
				Type: fibery.Text,
				Name: "Sync Action",
			},
			"LineNum": {
				Name:    "Line",
				Type:    fibery.Number,
				SubType: fibery.Integer,
			},
			"Taxed": {
				Name:    "Taxed",
				Type:    fibery.Text,
				SubType: fibery.Boolean,
			},
			"ServiceDate": {
				Name:    "Service Date",
				Type:    fibery.DateType,
				SubType: fibery.Day,
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
			"GroupLineId": {
				Name: "Group Line ID",
				Type: fibery.Text,
				Relation: &fibery.Relation{
					Cardinality:   fibery.MTO,
					Name:          "Group",
					TargetName:    "Lines",
					TargetType:    "EstimateLine",
					TargetFieldID: "id",
				},
			},
			"ItemId": {
				Name: "Item",
				Type: fibery.Text,
				Relation: &fibery.Relation{
					Cardinality:   fibery.MTO,
					Name:          "Item",
					TargetName:    "Estimate Lines",
					TargetType:    "Item",
					TargetFieldID: "id",
				},
			},
			"ClassId": {
				Name: "Class ID",
				Type: fibery.Text,
				Relation: &fibery.Relation{
					Cardinality:   fibery.MTO,
					Name:          "Class",
					TargetName:    "Expense Lines",
					TargetType:    "Class",
					TargetFieldID: "id",
				},
			},
			"EstimateId": {
				Name: "Estimate ID",
				Type: fibery.Text,
				Relation: &fibery.Relation{
					Cardinality:   fibery.MTO,
					Name:          "Estimate",
					TargetName:    "Lines",
					TargetType:    "Estimate",
					TargetFieldID: "id",
				},
			},
		},
		func(e quickbooks.Estimate) ([]map[string]any, error) {
			items := []map[string]any{}
			for _, line := range e.Line {
				if line.DetailType == quickbooks.DescriptionLine || line.DetailType == quickbooks.SalesItemLine {
					var name string
					if line.SalesItemLineDetail.ItemRef.Name == "" {
						name = line.Description
					} else {
						if line.Description == "" {
							name = line.SalesItemLineDetail.ItemRef.Name
						} else {
							name = line.SalesItemLineDetail.ItemRef.Name + " - " + line.Description
						}
					}

					item := map[string]any{
						"id":           fmt.Sprintf("%s:%s", e.Id, line.Id),
						"QBOId":        line.Id,
						"Name":         name,
						"Description":  line.Description,
						"__syncAction": fibery.SET,
						"LineNum":      line.LineNum,
						"Taxed":        line.SalesItemLineDetail.TaxCodeRef.Value == "TAX",
						"ServiceDate":  line.SalesItemLineDetail.ServiceDate.Format(fibery.DateFormat),
						"Qty":          line.GroupLineDetail.Quantity,
						"UnitPrice":    line.SalesItemLineDetail.UnitPrice,
						"Amount":       line.Amount,
						"ItemId":       line.SalesItemLineDetail.ItemRef.Value,
						"ClassId":      line.SalesItemLineDetail.ClassRef.Value,
						"EstimateId":   e.Id,
					}
					items = append(items, item)
				}
				if line.DetailType == quickbooks.GroupLine {
					for _, groupLine := range line.GroupLineDetail.Line {
						var name string
						if groupLine.SalesItemLineDetail.ItemRef.Name == "" {
							name = groupLine.Description
						} else {
							if groupLine.Description == "" {
								name = groupLine.SalesItemLineDetail.ItemRef.Name
							} else {
								name = groupLine.SalesItemLineDetail.ItemRef.Name + " - " + groupLine.Description
							}
						}

						item := map[string]any{
							"id":           fmt.Sprintf("%s:%s:%s", e.Id, line.Id, groupLine.Id),
							"GroupLineId":  line.Id,
							"QBOId":        line.Id,
							"Name":         name,
							"Description":  line.Description,
							"__syncAction": fibery.SET,
							"LineNum":      groupLine.LineNum,
							"Taxed":        groupLine.SalesItemLineDetail.TaxCodeRef.Value == "TAX",
							"ServiceDate":  groupLine.SalesItemLineDetail.ServiceDate.Format(fibery.DateFormat),
							"Qty":          groupLine.SalesItemLineDetail.Qty,
							"UnitPrice":    groupLine.SalesItemLineDetail.UnitPrice,
							"Amount":       groupLine.Amount,
							"ItemId":       line.SalesItemLineDetail.ItemRef.Value,
							"ClassId":      groupLine.SalesItemLineDetail.ClassRef.Value,
							"EstimateId":   e.Id,
						}
						items = append(items, item)
					}

					var name string
					if line.Description == "" {
						name = line.GroupLineDetail.GroupItemRef.Name
					} else {
						name = line.GroupLineDetail.GroupItemRef.Name + " - " + line.Description
					}

					item := map[string]any{
						"id":           fmt.Sprintf("%s:%s", e.Id, line.Id),
						"QBOId":        line.Id,
						"Name":         name,
						"Description":  line.Description,
						"__syncAction": fibery.SET,
						"Qty":          line.GroupLineDetail.Quantity,
						"LineNum":      line.LineNum,
						"ItemId":       line.GroupLineDetail.GroupItemRef.Value,
						"EstimateId":   e.Id,
					}
					items = append(items, item)
				}
			}
			return items, nil
		},
		estimate,
		func(e quickbooks.Estimate) string {
			return e.Id
		},
		func(e quickbooks.Estimate) string {
			return e.Status
		},
		func(e quickbooks.Estimate) map[string]struct{} {
			sourceMap := map[string]struct{}{}
			for _, line := range e.Line {
				if line.DetailType == quickbooks.GroupLine {
					for _, groupLine := range line.GroupLineDetail.Line {
						sourceMap[fmt.Sprintf("%s:%s", line.Id, groupLine.Id)] = struct{}{}
					}
					sourceMap[line.Id] = struct{}{}
				}
				if line.DetailType == quickbooks.DescriptionLine || line.DetailType == quickbooks.SalesItemLine {
					sourceMap[line.Id] = struct{}{}
				}
			}
			return sourceMap
		},
	)

	tr.Register(estimateLine)

	invoice := NewQuickBooksDualType(
		"Invoice",
		"Invoice",
		map[string]fibery.Field{
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
			"ShippingFromLine1": {
				Name: "Sale Line 1",
				Type: fibery.Text,
			},
			"ShippingFromLine2": {
				Name: "Sale Line 2",
				Type: fibery.Text,
			},
			"ShippingFromLine3": {
				Name: "Sale Line 3",
				Type: fibery.Text,
			},
			"ShippingFromLine4": {
				Name: "Sale Line 4",
				Type: fibery.Text,
			},
			"ShippingFromLine5": {
				Name: "Sale Line 5",
				Type: fibery.Text,
			},
			"ShippingFromCity": {
				Name: "Sale City",
				Type: fibery.Text,
			},
			"ShippingFromState": {
				Name: "Sale State",
				Type: fibery.Text,
			},
			"ShippingFromPostalCode": {
				Name: "Sale Postal Code",
				Type: fibery.Text,
			},
			"ShippingFromCountry": {
				Name: "Sale Country",
				Type: fibery.Text,
			},
			"ShippingFromLat": {
				Name: "Sale Latitude",
				Type: fibery.Text,
			},
			"ShippingFromLong": {
				Name: "Sale Longitude",
				Type: fibery.Text,
			},
			"ShippingLine1": {
				Name: "Shipping Line 1",
				Type: fibery.Text,
			},
			"ShippingLine2": {
				Name: "Shipping Line 2",
				Type: fibery.Text,
			},
			"ShippingLine3": {
				Name: "Shipping Line 3",
				Type: fibery.Text,
			},
			"ShippingLine4": {
				Name: "Shipping Line 4",
				Type: fibery.Text,
			},
			"ShippingLine5": {
				Name: "Shipping Line 5",
				Type: fibery.Text,
			},
			"ShippingCity": {
				Name: "Shipping City",
				Type: fibery.Text,
			},
			"ShippingState": {
				Name: "Shipping State",
				Type: fibery.Text,
			},
			"ShippingPostalCode": {
				Name: "Shipping Postal Code",
				Type: fibery.Text,
			},
			"ShippingCountry": {
				Name: "Shipping Country",
				Type: fibery.Text,
			},
			"ShippingLat": {
				Name: "Shipping Latitude",
				Type: fibery.Text,
			},
			"ShippingLong": {
				Name: "Shipping Longitude",
				Type: fibery.Text,
			},
			"BillingLine1": {
				Name: "Billing Line 1",
				Type: fibery.Text,
			},
			"BillingLine2": {
				Name: "Billing Line 2",
				Type: fibery.Text,
			},
			"BillingLine3": {
				Name: "Billing Line 3",
				Type: fibery.Text,
			},
			"BillingLine4": {
				Name: "Billing Line 4",
				Type: fibery.Text,
			},
			"BillingLine5": {
				Name: "Billing Line 5",
				Type: fibery.Text,
			},
			"BillingCity": {
				Name: "Billing City",
				Type: fibery.Text,
			},
			"BillingState": {
				Name: "Billing State",
				Type: fibery.Text,
			},
			"BillingPostalCode": {
				Name: "Billing Postal Code",
				Type: fibery.Text,
			},
			"BillingCountry": {
				Name: "Billing Country",
				Type: fibery.Text,
			},
			"BillingLat": {
				Name: "Billing Latitude",
				Type: fibery.Text,
			},
			"BillingLong": {
				Name: "Billing Longitude",
				Type: fibery.Text,
			},
			"DocNumber": {
				Name: "Number",
				Type: fibery.Text,
			},
			"Email": {
				Name:    "To",
				Type:    fibery.Text,
				SubType: fibery.Email,
			},
			"EmailCC": {
				Name:    "CC",
				Type:    fibery.Text,
				SubType: fibery.Email,
			},
			"EmailBCC": {
				Name:    "BCC",
				Type:    fibery.Text,
				SubType: fibery.Email,
			},
			"EmailSendLater": {
				Name:    "Send Later",
				Type:    fibery.Text,
				SubType: fibery.Boolean,
			},
			"EmailSent": {
				Name:    "Sent",
				Type:    fibery.Text,
				SubType: fibery.Boolean,
			},
			"EmailSendTime": {
				Name: "Send Time",
				Type: fibery.DateType,
			},
			"TxnDate": {
				Name:    "Date",
				Type:    fibery.DateType,
				SubType: fibery.Day,
			},
			"DueDate": {
				Name:    "Due Date",
				Type:    fibery.DateType,
				SubType: fibery.Day,
			},
			"PrivateNote": {
				Name:    "Message on Statement",
				Type:    fibery.Text,
				SubType: fibery.MD,
			},
			"CustomerMemo": {
				Name:    "Message on Invoice",
				Type:    fibery.Text,
				SubType: fibery.MD,
			},
			"DiscountPosition": {
				Name:        "Discount Before Tax",
				Type:        fibery.Text,
				SubType:     fibery.Boolean,
				Description: "Should the discount be applied before or after the sales tax calculation? Default is false as tax should generally be calculated first before a discount is given.",
			},
			"DepositField": {
				Name: "Deposit Amount",
				Type: fibery.Number,
				Format: map[string]any{
					"format":               "Money",
					"currencyCode":         "USD",
					"hasThousandSeperator": true,
					"precision":            2,
				},
			},
			"DiscountType": {
				Name:     "Discount Type",
				Type:     fibery.Text,
				SubType:  fibery.SingleSelect,
				ReadOnly: true,
				Options: []map[string]any{
					{
						"name": "Percent",
					},
					{
						"name": "Amount",
					},
				},
			},
			"DiscountPercent": {
				Name: "Discount Percent",
				Type: fibery.Number,
				Format: map[string]any{
					"format":    "Percent",
					"precision": 2,
				},
			},
			"DiscountAmount": {
				Name: "Discount Amount",
				Type: fibery.Number,
				Format: map[string]any{
					"format":               "Money",
					"currencyCode":         "USD",
					"hasThousandSeperator": true,
					"precision":            2,
				},
			},
			"Tax": {
				Name: "Tax",
				Type: fibery.Number,
				Format: map[string]any{
					"format":               "Money",
					"currencyCode":         "USD",
					"hasThousandSeperator": true,
					"precision":            2,
				},
			},
			"TaxPercent": {
				Name: "Tax Percent",
				Type: fibery.Number,
				Format: map[string]any{
					"format":    "Percent",
					"precision": 2,
				},
			},

			"SubtotalAmt": {
				Name: "Subtotal",
				Type: fibery.Number,
				Format: map[string]any{
					"format":               "Money",
					"currencyCode":         "USD",
					"hasThousandSeperator": true,
					"precision":            2,
				},
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
			"AllowACH": {
				Name:    "ACH Payments",
				Type:    fibery.Text,
				SubType: fibery.Boolean,
			},
			"AllowCC": {
				Name:    "Credit Card Payments",
				Type:    fibery.Text,
				SubType: fibery.Boolean,
			},
			"ClassId": {
				Name: "Class ID",
				Type: fibery.Text,
				Relation: &fibery.Relation{
					Cardinality:   fibery.MTO,
					Name:          "Class",
					TargetName:    "Invoices",
					TargetType:    "Class",
					TargetFieldID: "id",
				},
			},
			"SalesTermId": {
				Name: "Sales Term ID",
				Type: fibery.Text,
				Relation: &fibery.Relation{
					Cardinality:   fibery.MTO,
					Name:          "Terms",
					TargetName:    "Invoices",
					TargetType:    "Term",
					TargetFieldID: "id",
				},
			},
			"TaxCodeId": {
				Name: "Tax Code ID",
				Type: fibery.Text,
				Relation: &fibery.Relation{
					Cardinality:   fibery.MTO,
					Name:          "Sales Tax Rate",
					TargetName:    "Invoices",
					TargetType:    "TaxCode",
					TargetFieldID: "id",
				},
			},
			"TaxExemptionId": {
				Name: "Tax Exemption ID",
				Type: fibery.Text,
				Relation: &fibery.Relation{
					Cardinality:   fibery.MTO,
					Name:          "Tax Exemption",
					TargetName:    "Invoices",
					TargetType:    "TaxExemption",
					TargetFieldID: "id",
				},
			},
			"CustomerId": {
				Name: "Customer ID",
				Type: fibery.Text,
				Relation: &fibery.Relation{
					Cardinality:   fibery.MTO,
					Name:          "Customer",
					TargetName:    "Invoices",
					TargetType:    "Customer",
					TargetFieldID: "id",
				},
			},
			"LinkedTimeActivityIds": {
				Name: "Linked Time Activity IDs",
				Type: fibery.TextArray,
				Relation: &fibery.Relation{
					Cardinality:   fibery.MTM,
					Name:          "Time Activities",
					TargetName:    "Linked Invoice",
					TargetType:    "TimeActivity",
					TargetFieldID: "id",
				},
			},
		},
		func(i quickbooks.Invoice) (map[string]any, error) {
			var discountType = map[bool]string{
				true:  "Percent",
				false: "Amount",
			}
			var discountLine *quickbooks.Line
			var subtotalLine *quickbooks.Line
			for _, line := range i.Line {
				if line.DetailType == "DiscountLineDetail" {
					if discountLine != nil {
						return nil, fmt.Errorf("invoice %s has more than one discount line", i.Id)
					}
					discountLine = &line
				}
				if line.DetailType == "SubTotalLineDetail" {
					if subtotalLine != nil {
						return nil, fmt.Errorf("invoice %s has more than one subtotal line", i.Id)
					}
					subtotalLine = &line
				}
			}

			if subtotalLine == nil {
				return nil, fmt.Errorf("estimate %s has no subtotal lines", i.Id)
			}

			var discountTypeValue string
			var discountPercent float64
			var discountAmount json.Number

			if discountLine != nil {
				var err error
				var discountNumber float64
				if discountLine.DiscountLineDetail.DiscountPercent == "" {
					discountNumber = 0.0
				} else {
					discountNumber, err = discountLine.DiscountLineDetail.DiscountPercent.Float64()
					if err != nil {
						return nil, err
					}
				}
				discountTypeValue = discountType[discountLine.DiscountLineDetail.PercentBased]
				discountPercent = discountNumber / 100
				discountAmount = discountLine.Amount
			}

			var emailSendTime string
			if i.DeliveryInfo != nil && !i.DeliveryInfo.DeliveryTime.IsZero() {
				emailSendTime = i.DeliveryInfo.DeliveryTime.Format(fibery.DateFormat)
			}

			var name string
			if i.CustomerRef.Name == "" {
				name = i.DocNumber
			} else {
				name = i.DocNumber + " - " + i.CustomerRef.Name
			}

			var shipAddr quickbooks.PhysicalAddress
			if i.ShipAddr != nil {
				shipAddr = *i.ShipAddr
			}

			var billAddr quickbooks.PhysicalAddress
			if i.BillAddr != nil {
				billAddr = *i.BillAddr
			}

			var billEmailCC string
			if i.BillEmailCC != nil {
				billEmailCC = i.BillEmailCC.Address
			}

			var billEmailBCC string
			if i.BillEmailBCC != nil {
				billEmailBCC = i.BillEmailCC.Address
			}

			var txnDate string
			if i.TxnDate != nil {
				txnDate = i.TxnDate.Format(fibery.DateFormat)
			}

			var dueDate string
			if i.DueDate != nil {
				dueDate = i.DueDate.Format(fibery.DateFormat)
			}

			var taxPercent float64
			var totalTax json.Number
			var taxCodeId string
			if i.TxnTaxDetail != nil {
				for _, taxLine := range i.TxnTaxDetail.TaxLine {
					if taxLine.TaxLineDetail.PercentBased {
						if taxLine.TaxLineDetail.TaxPercent == "" {
							continue
						}
						percentNumber, err := taxLine.TaxLineDetail.TaxPercent.Float64()
						if err != nil {
							return nil, err
						}
						taxPercent += percentNumber / 100
					}
				}
				totalTax = i.TxnTaxDetail.TotalTax
				taxCodeId = i.TxnTaxDetail.TxnTaxCodeRef.Value
			}

			var classId string
			if i.ClassRef != nil {
				classId = i.ClassRef.Value
			}

			var taxExemptionId string
			if i.TaxExemptionRef != nil {
				taxExemptionId = i.TaxExemptionRef.Value
			}

			linkedTimeActivityIds := []string{}
			for _, txn := range i.LinkedTxn {
				if txn.TxnType == "TimeActivity" {
					linkedTimeActivityIds = append(linkedTimeActivityIds, txn.TxnId)
				}
			}

			return map[string]any{
				"id":                     i.Id,
				"QBOId":                  i.Id,
				"Name":                   name,
				"SyncToken":              i.SyncToken,
				"__syncAction":           fibery.SET,
				"ShippingLine1":          shipAddr.Line1,
				"ShippingLine2":          shipAddr.Line2,
				"ShippingLine3":          shipAddr.Line3,
				"ShippingLine4":          shipAddr.Line4,
				"ShippingLine5":          shipAddr.Line5,
				"ShippingCity":           shipAddr.City,
				"ShippingState":          shipAddr.CountrySubDivisionCode,
				"ShippingPostalCode":     shipAddr.PostalCode,
				"ShippingCountry":        shipAddr.Country,
				"ShippingLat":            shipAddr.Lat,
				"ShippingLong":           shipAddr.Long,
				"ShippingFromLine1":      i.ShipFromAddr.Line1,
				"ShippingFromLine2":      i.ShipFromAddr.Line2,
				"ShippingFromLine3":      i.ShipFromAddr.Line3,
				"ShippingFromLine4":      i.ShipFromAddr.Line4,
				"ShippingFromLine5":      i.ShipFromAddr.Line5,
				"ShippingFromCity":       i.ShipFromAddr.City,
				"ShippingFromState":      i.ShipFromAddr.CountrySubDivisionCode,
				"ShippingFromPostalCode": i.ShipFromAddr.PostalCode,
				"ShippingFromCountry":    i.ShipFromAddr.Country,
				"ShippingFromLat":        i.ShipFromAddr.Lat,
				"ShippingFromLong":       i.ShipFromAddr.Long,
				"BillingLine1":           billAddr.Line1,
				"BillingLine2":           billAddr.Line2,
				"BillingLine3":           billAddr.Line3,
				"BillingLine4":           billAddr.Line4,
				"BillingLine5":           billAddr.Line5,
				"BillingCity":            billAddr.City,
				"BillingState":           billAddr.CountrySubDivisionCode,
				"BillingPostalCode":      billAddr.PostalCode,
				"BillingCountry":         billAddr.Country,
				"BillingLat":             billAddr.Lat,
				"BillingLong":            billAddr.Long,
				"DocNumber":              i.DocNumber,
				"Email":                  i.BillEmail.Address,
				"EmailCC":                billEmailCC,
				"EmailBCC":               billEmailBCC,
				"EmailSendLater":         i.EmailStatus == "NeedToSend",
				"EmailSent":              i.EmailStatus == "EmailSent",
				"EmailSendTime":          emailSendTime,
				"TxnDate":                txnDate,
				"DueDate":                dueDate,
				"PrivateNote":            i.PrivateNote,
				"CustomerMemo":           i.CustomerMemo.Value,
				"DiscountPosition":       i.ApplyTaxAfterDiscount,
				"DepositField":           i.Deposit,
				"DiscountType":           discountTypeValue,
				"DiscountPercent":        discountPercent,
				"DiscountAmount":         discountAmount,
				"Tax":                    totalTax,
				"SubtotalAmt":            subtotalLine.Amount,
				"TotalAmt":               i.TotalAmt,
				"Balance":                i.Balance,
				"AllowACH":               i.AllowOnlineACHPayment,
				"AllowCC":                i.AllowOnlineCreditCardPayment,
				"ClassId":                classId,
				"TaxCodeId":              taxCodeId,
				"TaxExemptionId":         taxExemptionId,
				"CustomerId":             i.CustomerRef.Value,
				"LinkedTimeActivityIds":  linkedTimeActivityIds,
			}, nil
		},
		func(client *quickbooks.Client, requestParams quickbooks.RequestParameters, startPosition, pageSize int) ([]quickbooks.Invoice, error) {
			items, err := client.FindInvoicesByPage(requestParams, startPosition, pageSize)
			if err != nil {
				return nil, err
			}

			return items, nil
		},
		func(i quickbooks.Invoice) string {
			return i.Id
		},
		func(i quickbooks.Invoice) string {
			return i.Status
		},
	)

	tr.Register(invoice)

	invoiceLine := NewDependentDualType(
		"InvoiceLine",
		"Invoice Line",
		map[string]fibery.Field{
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
			"Description": {
				Name: "Description",
				Type: fibery.Text,
			},
			"__syncAction": {
				Type: fibery.Text,
				Name: "Sync Action",
			},
			"LineNum": {
				Name:    "Line",
				Type:    fibery.Number,
				SubType: fibery.Integer,
			},
			"Taxed": {
				Name:    "Taxed",
				Type:    fibery.Text,
				SubType: fibery.Boolean,
			},
			"ServiceDate": {
				Name:    "Service Date",
				Type:    fibery.DateType,
				SubType: fibery.Day,
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
			"GroupLineId": {
				Name: "Group Line ID",
				Type: fibery.Text,
				Relation: &fibery.Relation{
					Cardinality:   fibery.MTO,
					Name:          "Group",
					TargetName:    "Lines",
					TargetType:    "InvoiceLine",
					TargetFieldID: "id",
				},
			},
			"ItemId": {
				Name: "Item",
				Type: fibery.Text,
				Relation: &fibery.Relation{
					Cardinality:   fibery.MTO,
					Name:          "Item",
					TargetName:    "Invoice Lines",
					TargetType:    "Item",
					TargetFieldID: "id",
				},
			},
			"ClassId": {
				Name: "Class ID",
				Type: fibery.Text,
				Relation: &fibery.Relation{
					Cardinality:   fibery.MTO,
					Name:          "Class",
					TargetName:    "Invoice Lines",
					TargetType:    "Class",
					TargetFieldID: "id",
				},
			},
			"ReimburseChargeId": {
				Name: "Reimburse Charge ID",
				Type: fibery.Text,
				Relation: &fibery.Relation{
					Cardinality:   fibery.OTO,
					Name:          "Reimburse Charge",
					TargetName:    "Linked Invoice Line",
					TargetType:    "ReimburseCharge",
					TargetFieldID: "id",
				},
			},
			"InvoiceId": {
				Name: "Invoice ID",
				Type: fibery.Text,
				Relation: &fibery.Relation{
					Cardinality:   fibery.MTO,
					Name:          "Invoice",
					TargetName:    "Lines",
					TargetType:    "Invoice",
					TargetFieldID: "id",
				},
			},
		},
		func(i quickbooks.Invoice) ([]map[string]any, error) {
			items := []map[string]any{}
			for _, line := range i.Line {
				if line.DetailType == quickbooks.DescriptionLine || line.DetailType == quickbooks.SalesItemLine {
					var reimburseChargeId string
					for _, txn := range line.LinkedTxn {
						if txn.TxnType == "ReimburseCharge" {
							reimburseChargeId = txn.TxnId
						}
					}

					var name string
					if line.SalesItemLineDetail.ItemRef.Name == "" {
						name = line.Description
					} else {
						if line.Description == "" {
							name = line.SalesItemLineDetail.ItemRef.Name
						} else {
							name = line.SalesItemLineDetail.ItemRef.Name + " - " + line.Description
						}
					}

					item := map[string]any{
						"id":                fmt.Sprintf("%s:%s", i.Id, line.Id),
						"QBOId":             line.Id,
						"Name":              name,
						"Description":       line.Description,
						"__syncAction":      fibery.SET,
						"LineNum":           line.LineNum,
						"Taxed":             line.SalesItemLineDetail.TaxCodeRef.Value == "TAX",
						"ServiceDate":       line.SalesItemLineDetail.ServiceDate.Format(fibery.DateFormat),
						"Qty":               line.SalesItemLineDetail.Qty,
						"UnitPrice":         line.SalesItemLineDetail.UnitPrice,
						"Amount":            line.Amount,
						"ItemId":            line.SalesItemLineDetail.ItemRef.Value,
						"ClassId":           line.SalesItemLineDetail.ClassRef.Value,
						"ReimburseChargeId": reimburseChargeId,
						"InvoiceId":         i.Id,
					}
					items = append(items, item)
				}
				if line.DetailType == quickbooks.GroupLine {
					for _, groupLine := range line.GroupLineDetail.Line {
						var name string
						if groupLine.SalesItemLineDetail.ItemRef.Name == "" {
							name = groupLine.Description
						} else {
							if groupLine.Description == "" {
								name = groupLine.SalesItemLineDetail.ItemRef.Name
							} else {
								name = groupLine.SalesItemLineDetail.ItemRef.Name + " - " + groupLine.Description
							}
						}

						item := map[string]any{
							"id":           fmt.Sprintf("%s:%s:%s", i.Id, line.Id, groupLine.Id),
							"GroupLineId":  line.Id,
							"QBOId":        groupLine.Id,
							"Name":         name,
							"Description":  groupLine.Description,
							"__syncAction": fibery.SET,
							"LineNum":      groupLine.LineNum,
							"Taxed":        groupLine.SalesItemLineDetail.TaxCodeRef.Value == "TAX",
							"ServiceDate":  groupLine.SalesItemLineDetail.ServiceDate.Format(fibery.DateFormat),
							"Qty":          groupLine.SalesItemLineDetail.Qty,
							"UnitPrice":    groupLine.SalesItemLineDetail.UnitPrice,
							"Amount":       groupLine.Amount,
							"ItemId":       line.SalesItemLineDetail.ItemRef.Value,
							"ClassId":      groupLine.SalesItemLineDetail.ClassRef.Value,
							"InvoiceId":    i.Id,
						}
						items = append(items, item)
					}

					var name string
					if line.Description == "" {
						name = line.GroupLineDetail.GroupItemRef.Name
					} else {
						name = line.GroupLineDetail.GroupItemRef.Name + " - " + line.Description
					}

					item := map[string]any{
						"id":           fmt.Sprintf("%s:%s", i.Id, line.Id),
						"QBOId":        line.Id,
						"Name":         name,
						"Description":  line.Description,
						"__syncAction": fibery.SET,
						"Qty":          line.GroupLineDetail.Quantity,
						"LineNum":      line.LineNum,
						"ItemId":       line.GroupLineDetail.GroupItemRef.Value,
						"InvoiceId":    i.Id,
					}
					items = append(items, item)
				}
			}
			return items, nil
		},
		invoice,
		func(i quickbooks.Invoice) string {
			return i.Id
		},
		func(i quickbooks.Invoice) string {
			return i.Status
		},
		func(i quickbooks.Invoice) map[string]struct{} {
			sourceMap := map[string]struct{}{}
			for _, line := range i.Line {
				if line.DetailType == "GroupLineDetail" {
					for _, groupLine := range line.GroupLineDetail.Line {
						sourceMap[fmt.Sprintf("%s:%s", line.Id, groupLine.Id)] = struct{}{}
					}
					sourceMap[line.Id] = struct{}{}
				}
				if line.DetailType == "DescriptionOnly" || line.DetailType == "SalesItemLineDetail" {
					sourceMap[line.Id] = struct{}{}
				}
			}
			return sourceMap
		},
	)

	tr.Register(invoiceLine)

	item := NewQuickBooksDualType(
		"Item",
		"Item",
		map[string]fibery.Field{
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
					TargetName:    "Items",
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
					TargetName:    "Inventory Items",
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
					TargetName:    "Purchase Items",
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
					TargetName:    "Sale Items",
					TargetType:    "Account",
					TargetFieldID: "id",
				},
			},
		},
		func(i quickbooks.Item) (map[string]any, error) {
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

			var invStart string
			if !i.InvStartDate.IsZero() {
				i.InvStartDate.Format(fibery.DateFormat)
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
				"InvStartDate":        invStart,
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
		func(client *quickbooks.Client, requestParams quickbooks.RequestParameters, startPosition, pageSize int) ([]quickbooks.Item, error) {
			items, err := client.FindItemsByPage(requestParams, startPosition, pageSize)
			if err != nil {
				return nil, err
			}

			return items, nil
		},
		func(i quickbooks.Item) string {
			return i.Id
		},
		func(i quickbooks.Item) string {
			return i.Status
		},
	)

	tr.Register(item)

	paymentMethod := NewQuickBooksDualType(
		"PaymentMethod",
		"Payment Method",
		map[string]fibery.Field{
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
				Name: "Reference Number",
				Type: fibery.Text,
			},
			"TxnDate": {
				Name:    "Payment Date",
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
			"PayType": {
				Name:     "Payment Type",
				Type:     fibery.Text,
				SubType:  fibery.SingleSelect,
				ReadOnly: true,
				Options: []map[string]any{
					{
						"name": "Check",
					},
					{
						"name": "Credit Card",
					},
				},
			},
			"VendorId": {
				Name: "Vendor ID",
				Type: fibery.Text,
				Relation: &fibery.Relation{
					Cardinality:   fibery.MTO,
					Name:          "Vendor",
					TargetName:    "Payment Methods",
					TargetType:    "Vendor",
					TargetFieldID: "id",
				},
			},
			"PaymentAccountId": {
				Name: "Payment Account ID",
				Type: fibery.Text,
				Relation: &fibery.Relation{
					Cardinality:   fibery.MTO,
					Name:          "Payment Account",
					TargetName:    "Payment Methods",
					TargetType:    "Account",
					TargetFieldID: "id",
				},
			},
		},
		func(pm quickbooks.PaymentMethod) (map[string]any, error) {
			var paymentType string
			switch pm.Type {
			case "CREDIT_CARD":
				paymentType = "Credit Card"
			case "NON_CREDIT_CARD":
				paymentType = "Non-Credit Card"
			}
			return map[string]any{
				"id":          pm.Id,
				"QBOId":       pm.Id,
				"Name":        pm.Name,
				"SyncToken":   pm.SyncToken,
				"__syncToken": fibery.SET,
				"Active":      pm.Active,
				"PayType":     paymentType,
			}, nil
		},
		func(client *quickbooks.Client, requestParams quickbooks.RequestParameters, startPosition, pageSize int) ([]quickbooks.PaymentMethod, error) {
			items, err := client.FindPaymentMethodsByPage(requestParams, startPosition, pageSize)
			if err != nil {
				return nil, err
			}

			return items, nil
		},
		func(pm quickbooks.PaymentMethod) string {
			return pm.Id
		},
		func(pm quickbooks.PaymentMethod) string {
			return pm.Status
		},
	)

	tr.Register(paymentMethod)

	purchase := NewQuickBooksDualType(
		"Purchase",
		"Expense",
		map[string]fibery.Field{
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
			"CustomerId": {
				Name: "Customer ID",
				Type: fibery.Text,
				Relation: &fibery.Relation{
					Cardinality:   fibery.MTO,
					Name:          "Customer Payee",
					TargetName:    "Expenses",
					TargetType:    "Customer",
					TargetFieldID: "id",
				},
			},
			"VendorId": {
				Name: "Vendor ID",
				Type: fibery.Text,
				Relation: &fibery.Relation{
					Cardinality:   fibery.MTO,
					Name:          "Vendor Payee",
					TargetName:    "Expenses",
					TargetType:    "Vendor",
					TargetFieldID: "id",
				},
			},
			"EmployeeId": {
				Name: "Employee ID",
				Type: fibery.Text,
				Relation: &fibery.Relation{
					Cardinality:   fibery.MTO,
					Name:          "Employee Payee",
					TargetName:    "Expenses",
					TargetType:    "Employee",
					TargetFieldID: "id",
				},
			},
			"EntityId": {
				Name: "Entity ID",
				Type: fibery.Text,
				Relation: &fibery.Relation{
					Cardinality:   fibery.MTO,
					Name:          "Entity Payee",
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
				Name:    "Date",
				Type:    fibery.DateType,
				SubType: fibery.Day,
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
			"Files": {
				Name:    "Files",
				Type:    fibery.TextArray,
				SubType: fibery.File,
			},
		},
		func(p quickbooks.Purchase) (map[string]any, error) {
			var paymentType string
			switch p.PaymentType {
			case "Cash":
				paymentType = "Cash"
			case "Check":
				paymentType = "Check"
			case "CreditCard":
				paymentType = "Credit Card"
			}

			var entityId, entityName, customerId, vendorId, employeeId string
			if p.EntityRef != nil {
				entityName = p.EntityRef.Name
				switch p.EntityRef.Type {
				case "Customer":
					entityId = "c:" + p.EntityRef.Value
					customerId = p.EntityRef.Value
				case "Vendor":
					entityId = "v:" + p.EntityRef.Value
					vendorId = p.EntityRef.Value
				case "Employee":
					entityId = "e:" + p.EntityRef.Value
					employeeId = p.EntityRef.Value
				}
			}

			var name string
			if p.PrivateNote == "" {
				name = entityName
			} else {
				name = entityName + " - " + p.PrivateNote
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
				"Name":             name,
				"SyncToken":        p.SyncToken,
				"__syncAction":     fibery.SET,
				"PaymentType":      paymentType,
				"PaymentAccountId": p.AccountRef.Value,
				"PaymentMethodId":  paymentMethodId,
				"CustomerId":       customerId,
				"VendorId":         vendorId,
				"EmployeeId":       employeeId,
				"EntityId":         entityId,
				"Total":            p.TotalAmt,
				"DocNumber":        p.DocNumber,
				"TxnDate":          txnDate,
				"Credit":           p.Credit,
				"PrivateNote":      p.PrivateNote,
			}, nil
		},
		func(client *quickbooks.Client, requestParams quickbooks.RequestParameters, startPosition, pageSize int) ([]quickbooks.Purchase, error) {
			items, err := client.FindPurchasesByPage(requestParams, startPosition, pageSize)
			if err != nil {
				return nil, err
			}

			return items, nil
		},
		func(p quickbooks.Purchase) string {
			return p.Id
		},
		func(p quickbooks.Purchase) string {
			return p.Status
		},
	)

	tr.Register(purchase)

	purchaseItemLine := NewDependentDualType(
		"PurchaseItemLine",
		"Expense Item Line",
		map[string]fibery.Field{
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
			"Description": {
				Name: "Description",
				Type: fibery.Text,
			},
			"__syncAction": {
				Type: fibery.Text,
				Name: "Sync Action",
			},
			"LineNum": {
				Name:    "Line",
				Type:    fibery.Number,
				SubType: fibery.Integer,
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
			"ReimburseChargeId": {
				Name: "Reimburse Charge ID",
				Type: fibery.Text,
				Relation: &fibery.Relation{
					Cardinality:   fibery.OTO,
					Name:          "Reimburse Charge",
					TargetName:    "Expense Item Line",
					TargetType:    "ReimburseCharge",
					TargetFieldID: "id",
				},
			},
		},
		func(p quickbooks.Purchase) ([]map[string]any, error) {
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

					var reimburseChargeId string
					for _, txn := range line.LinkedTxn {
						if txn.TxnType == "ReimburseCharge" {
							reimburseChargeId = txn.TxnId
						}
					}

					var name string
					if line.Description == "" {
						name = line.ItemBasedExpenseLineDetail.ItemRef.Name
					} else {
						name = line.ItemBasedExpenseLineDetail.ItemRef.Name + line.Description
					}

					item := map[string]any{
						"id":                fmt.Sprintf("%s:i:%s", p.Id, line.Id),
						"QBOId":             line.Id,
						"Name":              name,
						"Description":       line.Description,
						"__syncAction":      fibery.SET,
						"LineNum":           line.LineNum,
						"Tax":               tax,
						"Billable":          billable,
						"Billed":            billed,
						"Qty":               line.ItemBasedExpenseLineDetail.Qty,
						"UnitPrice":         line.ItemBasedExpenseLineDetail.UnitPrice,
						"MarkupPercent":     line.ItemBasedExpenseLineDetail.MarkupInfo.Percent,
						"Amount":            line.Amount,
						"PurchaseId":        p.Id,
						"ItemId":            line.ItemBasedExpenseLineDetail.ItemRef.Value,
						"CustomerId":        line.AccountBasedExpenseLineDetail.CustomerRef.Value,
						"ClassId":           line.ItemBasedExpenseLineDetail.ClassRef.Value,
						"MarkupAccountId":   line.ItemBasedExpenseLineDetail.MarkupInfo.MarkUpIncomeAccountRef.Value,
						"ReimburseChargeId": reimburseChargeId,
					}
					items = append(items, item)
				}
			}
			return items, nil
		},
		purchase,
		func(p quickbooks.Purchase) string {
			return p.Id
		},
		func(p quickbooks.Purchase) string {
			return p.Status
		},
		func(p quickbooks.Purchase) map[string]struct{} {
			sourceMap := map[string]struct{}{}
			for _, line := range p.Line {
				if line.DetailType == quickbooks.ItemExpenseLine {
					sourceMap[line.Id] = struct{}{}
				}
			}
			return sourceMap
		},
	)

	tr.Register(purchaseItemLine)

	purchaseAccountLine := NewDependentDualType(
		"PurchaseAccountLine",
		"Expense Account Line",
		map[string]fibery.Field{
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
			"Description": {
				Name: "Description",
				Type: fibery.Text,
			},
			"__syncAction": {
				Type: fibery.Text,
				Name: "Sync Action",
			},
			"LineNum": {
				Name:    "Line",
				Type:    fibery.Number,
				SubType: fibery.Integer,
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
					TargetName:    "Account Lines",
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
					TargetName:    "Expense Account Lines",
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
					TargetName:    "Expense Account Lines",
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
					TargetName:    "Expense Account Line Markup",
					TargetType:    "Account",
					TargetFieldID: "id",
				},
			},
			"ReimburseChargeId": {
				Name: "Reimburse Charge ID",
				Type: fibery.Text,
				Relation: &fibery.Relation{
					Cardinality:   fibery.OTO,
					Name:          "Reimburse Charge",
					TargetName:    "Expense Account Line",
					TargetType:    "ReimburseCharge",
					TargetFieldID: "id",
				},
			},
		},
		func(p quickbooks.Purchase) ([]map[string]any, error) {
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

					var reimburseChargeId string
					for _, txn := range line.LinkedTxn {
						if txn.TxnType == "ReimburseCharge" {
							reimburseChargeId = txn.TxnId
						}
					}

					var name string
					if line.Description == "" {
						name = line.AccountBasedExpenseLineDetail.AccountRef.Name
					} else {
						name = line.AccountBasedExpenseLineDetail.AccountRef.Name + " - " + line.Description
					}

					item := map[string]any{
						"id":                fmt.Sprintf("%s:a:%s", p.Id, line.Id),
						"QBOId":             line.Id,
						"Name":              name,
						"Description":       line.Description,
						"__syncAction":      fibery.SET,
						"LineNum":           line.LineNum,
						"Tax":               tax,
						"Billable":          billable,
						"Billed":            billed,
						"MarkupPercent":     line.AccountBasedExpenseLineDetail.MarkupInfo.Percent,
						"Amount":            line.Amount,
						"PurchaseId":        p.Id,
						"AccountId":         line.AccountBasedExpenseLineDetail.AccountRef.Value,
						"CustomerId":        line.AccountBasedExpenseLineDetail.CustomerRef.Value,
						"ClassId":           line.AccountBasedExpenseLineDetail.ClassRef.Value,
						"MarkupAccountId":   line.AccountBasedExpenseLineDetail.MarkupInfo.MarkUpIncomeAccountRef.Value,
						"ReimburseChargeId": reimburseChargeId,
					}

					items = append(items, item)
				}
			}
			return items, nil
		},
		purchase,
		func(p quickbooks.Purchase) string {
			return p.Id
		},
		func(p quickbooks.Purchase) string {
			return p.Status
		},
		func(p quickbooks.Purchase) map[string]struct{} {
			sourceMap := map[string]struct{}{}
			for _, line := range p.Line {
				if line.DetailType == quickbooks.AccountExpenseLine {
					sourceMap[line.Id] = struct{}{}
				}
			}
			return sourceMap
		},
	)

	tr.Register(purchaseAccountLine)

	reimburseCharge := NewQuickBooksCDCType(
		"ReimburseCharge",
		"Reimburse Charge",
		map[string]fibery.Field{
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
			"Description": {
				Name: "Description",
				Type: fibery.Text,
			},
			"__syncAction": {
				Type: fibery.Text,
				Name: "Sync Action",
			},
			"SyncToken": {
				Name: "Sync Token",
				Type: fibery.Text,
			},
			"TxnDate": {
				Name:    "Date",
				Type:    fibery.DateType,
				SubType: fibery.Day,
			},
			"CustomerId": {
				Name: "Customer ID",
				Type: fibery.Text,
				Relation: &fibery.Relation{
					Cardinality:   fibery.MTO,
					Name:          "Customer",
					TargetName:    "Reimburse Charges",
					TargetType:    "Customer",
					TargetFieldID: "id",
				},
			},
			"TotalAmount": {
				Name: "Total Amount",
				Type: fibery.Number,
				Format: map[string]any{
					"format":               "Money",
					"currencyCode":         "USD",
					"hasThousandSeperator": true,
					"precision":            2,
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
			"AccountId": {
				Name: "Account ID",
				Type: fibery.Text,
				Relation: &fibery.Relation{
					Cardinality:   fibery.MTO,
					Name:          "Account",
					TargetName:    "Reimburse Charges",
					TargetType:    "Account",
					TargetFieldID: "id",
				},
			},
			"Markup": {
				Name: "Markup",
				Type: fibery.Text,
			},
			"MarkupAccountId": {
				Name: "Markup Account ID",
				Type: fibery.Text,
				Relation: &fibery.Relation{
					Cardinality:   fibery.MTO,
					Name:          "Markup Account",
					TargetName:    "Reimburse Charge Markup",
					TargetType:    "Account",
					TargetFieldID: "id",
				},
			},
			"Taxable": {
				Name:    "Taxable",
				Type:    fibery.Text,
				SubType: fibery.Boolean,
			},
			"LinkedInvoiceId": {
				Name: "Linked Invoice ID",
				Type: fibery.Text,
				Relation: &fibery.Relation{
					Cardinality:   fibery.MTO,
					Name:          "Linked Invoice",
					TargetName:    "Reimburse Charges",
					TargetType:    "Invoice",
					TargetFieldID: "id",
				},
			},
		},
		func(r quickbooks.ReimburseCharge) (map[string]any, error) {
			var name string
			if r.PrivateNote != "" {
				name = r.CustomerRef.Name + " - " + r.PrivateNote
			}

			var date string
			if r.TxnDate != nil {
				date = r.TxnDate.Format(fibery.DateFormat)
			}

			var amount json.Number
			var accountId string
			var markup string
			var markupAccountId string

			taxable := false

			for _, line := range r.Line {
				if line.LineNum == 1 {
					amount = line.Amount
					accountId = line.ReimburseLineDetail.ItemAccountRef.Value
					if line.ReimburseLineDetail.TaxCodeRef.Value == "TAX" {
						taxable = true
					}
				}
				if line.LineNum == 2 {
					markupAmount, err := line.Amount.Float64()
					if err != nil {
						return nil, fmt.Errorf("error converting markup amount json.Number to float: %w", err)
					}

					markupPercent, err := line.ReimburseLineDetail.MarkupInfo.Percent.Float64()
					if err != nil {
						return nil, fmt.Errorf("error converting markup percent json.Number to float: %w", err)
					}

					var markupStr string

					if math.Mod(markupPercent, 1.0) == 0 {
						markupStr = fmt.Sprintf("%.0f%%", markupPercent)
					} else {
						markupStr = fmt.Sprintf("%.5f%%", markupPercent)

						markupStr = strings.TrimRight(markupStr, "0")
						markupStr = strings.TrimRight(markupStr, ".")
					}

					markup = fmt.Sprintf("$%.2f (%s)", markupAmount, markupStr)
					markupAccountId = line.ReimburseLineDetail.ItemRef.Value
				}
			}

			var invoiceId string
			for _, txn := range r.LinkedTxn {
				if txn.TxnType == "Invoice" {
					invoiceId = txn.TxnId
				}
			}

			return map[string]any{
				"id":              r.Id,
				"QBOId":           r.Id,
				"Name":            name,
				"Description":     r.PrivateNote,
				"__syncAction":    fibery.SET,
				"SyncToken":       r.SyncToken,
				"TxnDate":         date,
				"CustomerId":      r.CustomerRef.Value,
				"TotalAmount":     r.Amount,
				"Amount":          amount,
				"AccountId":       accountId,
				"Markup":          markup,
				"MarkupAccountId": markupAccountId,
				"Taxable":         taxable,
				"LinkedInvoiceId": invoiceId,
			}, nil
		},
		func(client *quickbooks.Client, requestParams quickbooks.RequestParameters, startPosition, pageSize int) ([]quickbooks.ReimburseCharge, error) {
			items, err := client.FindReimburseChargesByPage(requestParams, startPosition, pageSize)
			if err != nil {
				return nil, err
			}

			return items, nil
		},
		func(r quickbooks.ReimburseCharge) string {
			return r.Id
		},
		func(r quickbooks.ReimburseCharge) string {
			return r.Status
		},
	)

	tr.Register(reimburseCharge)

	taxCode := NewQuickBooksType(
		"TaxCode",
		"Tax Code",
		map[string]fibery.Field{
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
			"Description": {
				Name: "Description",
				Type: fibery.Text,
			},
			"SyncToken": {
				Name: "Sync Token",
				Type: fibery.Text,
			},
			"TaxGroup": {
				Name:    "Tax Group",
				Type:    fibery.Text,
				SubType: fibery.Boolean,
			},
			"Taxable": {
				Name:    "Taxable",
				Type:    fibery.Text,
				SubType: fibery.Boolean,
			},
			"Active": {
				Name:    "Active",
				Type:    fibery.Text,
				SubType: fibery.Boolean,
			},
			"Hidden": {
				Name:    "Hidden",
				Type:    fibery.Text,
				SubType: fibery.Boolean,
			},
			"TaxCodeType": {
				Name:     "Type",
				Type:     fibery.Text,
				SubType:  fibery.SingleSelect,
				ReadOnly: true,
				Options: []map[string]any{
					{
						"name": "User Defined",
					},
					{
						"name": "System Generated",
					},
				},
			},
		},
		func(tc quickbooks.TaxCode) (map[string]any, error) {
			var taxCodeType string
			switch tc.TaxCodeConfigType {
			case "SYSTEM_GENERATED":
				taxCodeType = "System Generated"
			case "USER_DEFINED":
				taxCodeType = "User Defined"
			}
			return map[string]any{
				"id":          tc.Id,
				"QBOId":       tc.Id,
				"Name":        tc.Name,
				"Description": tc.Description,
				"SyncToken":   tc.SyncToken,
				"TaxGroup":    tc.TaxGroup,
				"Taxable":     tc.Taxable,
				"Active":      tc.Active,
				"Hidden":      tc.Hidden,
				"TaxCodeType": taxCodeType,
			}, nil
		},
		func(client *quickbooks.Client, requestParams quickbooks.RequestParameters, startPosition, pageSize int) ([]quickbooks.TaxCode, error) {
			items, err := client.FindTaxCodesByPage(requestParams, startPosition, pageSize)
			if err != nil {
				return nil, err
			}

			return items, nil
		},
	)

	tr.Register(taxCode)

	taxCodeLine := NewDependentDataType(
		"TaxCodeLine",
		"Tax Code Line",
		map[string]fibery.Field{
			"id": {
				Name: "ID",
				Type: fibery.Id,
			},
			"Name": {
				Name:    "Name",
				Type:    fibery.Text,
				SubType: fibery.Title,
			},
			"TaxRateId": {
				Name: "Tax Rate ID",
				Type: fibery.Text,
				Relation: &fibery.Relation{
					Cardinality:   fibery.MTO,
					Name:          "Tax Rate",
					TargetName:    "Tax Code Lines",
					TargetType:    "TaxRate",
					TargetFieldID: "id",
				},
			},
			"TaxType": {
				Name:     "Tax Type",
				Type:     fibery.Text,
				SubType:  fibery.SingleSelect,
				ReadOnly: true,
				Options: []map[string]any{
					{
						"name": "Tax On Amount",
					},
					{
						"name": "Tax On Amount Plus Tax",
					},
					{
						"name": "Tax On Tax",
					},
				},
			},
			"TaxOrder": {
				Name: "Tax Order",
				Type: fibery.Number,
			},
			"TaxCodeIdPurchase": {
				Name: "Purchase Tax Code ID",
				Type: fibery.Text,
				Relation: &fibery.Relation{
					Cardinality:   fibery.MTO,
					Name:          "Purchase Tax Code",
					TargetName:    "Purchase Tax Rates",
					TargetType:    "TaxCode",
					TargetFieldID: "id",
				},
			},
			"TaxCodeIdSale": {
				Name: "Sale Tax Code ID",
				Type: fibery.Text,
				Relation: &fibery.Relation{
					Cardinality:   fibery.MTO,
					Name:          "Sales Tax Code",
					TargetName:    "Sales Tax Rates",
					TargetType:    "TaxCode",
					TargetFieldID: "id",
				},
			},
		},
		func(tc quickbooks.TaxCode) ([]map[string]any, error) {
			items := []map[string]any{}
			for _, ptRate := range tc.PurchaseTaxRateList.TaxRateDetail {
				var taxType string
				switch ptRate.TaxTypeApplicable {
				case "TaxOnAmount":
					taxType = "Tax On Amount"
				case "TaxOnAmountPlusTax":
					taxType = "Tax On Amount Plus Tax"
				case "TaxOnTax":
					taxType = "Tax On Tax"
				}
				item := map[string]any{
					"id":                fmt.Sprintf("pt:%s", ptRate.TaxOrder.String()),
					"Name":              ptRate.TaxRateRef.Name,
					"TaxRateId":         ptRate.TaxRateRef.Value,
					"TaxType":           taxType,
					"TaxOrder":          ptRate.TaxOrder,
					"TaxCodeIdPurchase": tc.Id,
				}
				items = append(items, item)
			}
			for _, stRate := range tc.SalesTaxRateList.TaxRateDetail {
				var taxType string
				switch stRate.TaxTypeApplicable {
				case "TaxOnAmount":
					taxType = "Tax On Amount"
				case "TaxOnAmountPlusTax":
					taxType = "Tax On Amount Plus Tax"
				case "TaxOnTax":
					taxType = "Tax On Tax"
				}
				item := map[string]any{
					"id":            fmt.Sprintf("st:%s", stRate.TaxOrder.String()),
					"Name":          stRate.TaxRateRef.Name,
					"TaxRateId":     stRate.TaxRateRef.Value,
					"TaxType":       taxType,
					"TaxOrder":      stRate.TaxOrder,
					"TaxCodeIdSale": tc.Id,
				}
				items = append(items, item)
			}
			return items, nil
		},
		taxCode,
	)

	tr.Register(taxCodeLine)

	type TaxExemptionEntity struct {
		Id   string
		Name string
	}

	taxExemption := NewQuickBooksType(
		"TaxExemption",
		"Tax Exemption",
		map[string]fibery.Field{
			"id": {
				Name: "ID",
				Type: fibery.Id,
			},
			"Name": {
				Name:    "Name",
				Type:    fibery.Text,
				SubType: fibery.Title,
			},
		},
		func(te TaxExemptionEntity) (map[string]any, error) {
			return map[string]any{
				"id":   te.Id,
				"Name": te.Name,
			}, nil
		},
		func(client *quickbooks.Client, requestParams quickbooks.RequestParameters, startPosition, pageSize int) ([]TaxExemptionEntity, error) {
			return []TaxExemptionEntity{
				{
					Id:   "1",
					Name: "Federal government",
				},
				{
					Id:   "2",
					Name: "State government",
				},
				{
					Id:   "3",
					Name: "Local government",
				},
				{
					Id:   "4",
					Name: "Tribal government",
				},
				{
					Id:   "5",
					Name: "Charitable organization",
				},
				{
					Id:   "6",
					Name: "Religious organization",
				},
				{
					Id:   "7",
					Name: "Educational organization",
				},
				{
					Id:   "8",
					Name: "Hospital",
				},
				{
					Id:   "9",
					Name: "Resale",
				},
				{
					Id:   "10",
					Name: "Direct pay permit",
				},
				{
					Id:   "11",
					Name: "Multiple points of use",
				},
				{
					Id:   "12",
					Name: "Direct mail",
				},
				{
					Id:   "13",
					Name: "Agricultural production",
				},
				{
					Id:   "14",
					Name: "Industrial production / manufacturing",
				},
				{
					Id:   "15",
					Name: "Foreign diplomat",
				},
			}, nil
		},
	)

	tr.Register(taxExemption)

	term := NewQuickBooksDualType(
		"Term",
		"Term",
		map[string]fibery.Field{
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
			"Active": {
				Name:    "Active",
				Type:    fibery.Text,
				SubType: fibery.Boolean,
			},
		},
		func(t quickbooks.Term) (map[string]any, error) {
			return map[string]any{
				"id":           t.Id,
				"QBOId":        t.Id,
				"Name":         t.Name,
				"SyncToken":    t.SyncToken,
				"__syncAction": fibery.SET,
				"Active":       t.Active,
			}, nil
		},
		func(client *quickbooks.Client, requestParams quickbooks.RequestParameters, startPosition, pageSize int) ([]quickbooks.Term, error) {
			items, err := client.FindTermsByPage(requestParams, startPosition, pageSize)
			if err != nil {
				return nil, err
			}

			return items, nil
		},
		func(t quickbooks.Term) string {
			return t.Id
		},
		func(t quickbooks.Term) string {
			return t.Status
		},
	)

	tr.Register(term)

	timeActivity := NewQuickBooksWHType(
		"TimeActivity",
		"Time Activity",
		map[string]fibery.Field{
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
			"Description": {
				Name:    "Description",
				Type:    fibery.Text,
				SubType: fibery.MD,
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
			"ActivityType": {
				Name:     "Activity Type",
				Type:     fibery.Text,
				SubType:  fibery.SingleSelect,
				ReadOnly: true,
				Options: []map[string]any{
					{
						"name": "Employee",
					},
					{
						"name": "Vendor",
					},
				},
			},
			"TxnDate": {
				Name:    "Date",
				Type:    fibery.DateType,
				SubType: fibery.Day,
			},
			"Hours": {
				Name: "Hours",
				Type: fibery.Number,
				Format: map[string]any{
					"format":               "Number",
					"currencyCode":         "h",
					"hasThousandSeperator": true,
					"precision":            0,
				},
			},
			"Minutes": {
				Name: "Minutes",
				Type: fibery.Number,
				Format: map[string]any{
					"format":               "Number",
					"currencyCode":         "m",
					"hasThousandSeperator": true,
					"precision":            0,
				},
			},
			"BreakHours": {
				Name: "Break Hours",
				Type: fibery.Number,
				Format: map[string]any{
					"format":               "Number",
					"currencyCode":         "h",
					"hasThousandSeperator": true,
					"precision":            0,
				},
			},
			"BreakMinutes": {
				Name: "Break Minutes",
				Type: fibery.Number,
				Format: map[string]any{
					"format":               "Number",
					"currencyCode":         "m",
					"hasThousandSeperator": true,
					"precision":            0,
				},
			},
			"StartTime": {
				Name: "Start Time",
				Type: fibery.DateType,
			},
			"EndTime": {
				Name: "End Time",
				Type: fibery.DateType,
			},
			"HourlyRate": {
				Name: "Rate",
				Type: fibery.Number,
				Format: map[string]any{
					"format":               "Money",
					"currencyCode":         "USD",
					"hasThousandSeperator": true,
					"precision":            2,
				},
			},
			"CostRate": {
				Name: "Cost Rate",
				Type: fibery.Number,
				Format: map[string]any{
					"format":               "Money",
					"currencyCode":         "USD",
					"hasThousandSeperator": true,
					"precision":            2,
				},
			},
			"Taxable": {
				Name:    "Taxable",
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
			"VendorId": {
				Name: "Vendor ID",
				Type: fibery.Text,
				Relation: &fibery.Relation{
					Cardinality:   fibery.MTO,
					Name:          "Vendor",
					TargetName:    "Time Activity",
					TargetType:    "Vendor",
					TargetFieldID: "id",
				},
			},
			"EmployeeId": {
				Name: "Employee ID",
				Type: fibery.Text,
				Relation: &fibery.Relation{
					Cardinality:   fibery.MTO,
					Name:          "Employee",
					TargetName:    "Time Activity",
					TargetType:    "Employee",
					TargetFieldID: "id",
				},
			},
			"CustomerId": {
				Name: "Customer ID",
				Type: fibery.Text,
				Relation: &fibery.Relation{
					Cardinality:   fibery.MTO,
					Name:          "Customer",
					TargetName:    "Time Activity",
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
					TargetName:    "Time Activity",
					TargetType:    "Class",
					TargetFieldID: "id",
				},
			},
			"ItemId": {
				Name: "Item ID",
				Type: fibery.Text,
				Relation: &fibery.Relation{
					Cardinality:   fibery.MTO,
					Name:          "Item",
					TargetName:    "Time Activity",
					TargetType:    "Item",
					TargetFieldID: "id",
				},
			},
		},
		func(ta quickbooks.TimeActivity) (map[string]any, error) {
			var billable bool
			switch ta.BillableStatus {
			case quickbooks.BillableStatusType:
				billable = true
			case quickbooks.HasBeenBilledStatusType:
				billable = true
			default:
				billable = false
			}

			billed := false
			if ta.BillableStatus == quickbooks.HasBeenBilledStatusType {
				billed = true
			}

			var classId string
			if ta.ClassRef != nil {
				classId = ta.ClassRef.Value
			}

			var itemId string
			if ta.ItemRef != nil {
				itemId = ta.ItemRef.Value
			}

			var startTime string
			if ta.StartTime != nil {
				startTime = ta.StartTime.Format(fibery.DateFormat)
			}

			var endTime string
			if ta.EndTime != nil {
				endTime = ta.EndTime.Format(fibery.DateFormat)
			}

			var txnDate string
			if !ta.TxnDate.IsZero() {
				txnDate = ta.TxnDate.Format(fibery.DateFormat)
			}

			var name string
			if ta.ItemRef.Name != "" {
				if ta.EmployeeRef.Name != "" {
					name = ta.EmployeeRef.Name + " - " + ta.ItemRef.Name
				} else {
					name = ta.VendorRef.Name + " - " + ta.ItemRef.Name
				}
			} else {
				if ta.EmployeeRef.Name != "" {
					name = ta.EmployeeRef.Name
				} else {
					name = ta.VendorRef.Name
				}
			}

			return map[string]any{
				"id":           ta.Id,
				"QBOId":        ta.Id,
				"Name":         name,
				"Description":  ta.Description,
				"SyncToken":    ta.SyncToken,
				"__syncAction": fibery.SET,
				"ActivityType": ta.NameOf,
				"TxnDate":      txnDate,
				"Hours":        ta.Hours,
				"Minutes":      ta.Minutes,
				"BreakHours":   ta.BreakHours,
				"BreakMinutes": ta.BreakMinutes,
				"StartTime":    startTime,
				"EndTime":      endTime,
				"HourlyRate":   ta.HourlyRate,
				"CostRate":     ta.CostRate,
				"Taxable":      ta.Taxable,
				"Billable":     billable,
				"Billed":       billed,
				"VendorId":     ta.VendorRef.Value,
				"EmployeeId":   ta.EmployeeRef.Value,
				"CustomerId":   ta.CustomerRef.Value,
				"ClassId":      classId,
				"ItemId":       itemId,
			}, nil
		},
		func(client *quickbooks.Client, requestParams quickbooks.RequestParameters, startPosition, pageSize int) ([]quickbooks.TimeActivity, error) {
			items, err := client.FindTimeActivitiesByPage(requestParams, startPosition, pageSize)
			if err != nil {
				return nil, err
			}

			return items, nil
		},
		func(ta quickbooks.TimeActivity) string {
			return ta.Id
		},
	)

	tr.Register(timeActivity)

	vendor := NewQuickBooksDualType(
		"Vendor",
		"Vendor",
		map[string]fibery.Field{
			"id": {
				Name: "ID",
				Type: fibery.Id,
			},
			"QBOId": {
				Name: "QBO ID",
				Type: fibery.Text,
			},
			"DisplayName": {
				Name:    "Display Name",
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
			"Title": {
				Name: "Title",
				Type: fibery.Text,
			},
			"GivenName": {
				Name: "First Name",
				Type: fibery.Text,
			},
			"MiddleName": {
				Name: "Middle Name",
				Type: fibery.Text,
			},
			"FamilyName": {
				Name: "Last Name",
				Type: fibery.Text,
			},
			"Suffix": {
				Name: "Suffix",
				Type: fibery.Text,
			},
			"CompanyName": {
				Name: "Company Name",
				Type: fibery.Text,
			},
			"PrimaryEmail": {
				Name:    "Email",
				Type:    fibery.Text,
				SubType: fibery.Email,
			},
			"SalesTermId": {
				Name: "Sales Term ID",
				Type: fibery.Text,
				Relation: &fibery.Relation{
					Cardinality:   fibery.MTO,
					Name:          "Sales Term",
					TargetName:    "Vendors",
					TargetType:    "SalesTerm",
					TargetFieldID: "id",
				},
			},
			"PrimaryPhone": {
				Name: "Phone",
				Type: fibery.Text,
				Format: map[string]any{
					"format": "phone",
				},
			},
			"AlternatePhone": {
				Name: "Alternate Phone",
				Type: fibery.Text,
				Format: map[string]any{
					"format": "phone",
				},
			},
			"Mobile": {
				Name: "Mobile",
				Type: fibery.Text,
				Format: map[string]any{
					"format": "phone",
				},
			},
			"Fax": {
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
			"CostRate": {
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
			"BillRate": {
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
			"Website": {
				Name:    "Website",
				Type:    fibery.Text,
				SubType: fibery.URL,
			},
			"AccountNumber": {
				Name:        "Account Number",
				Type:        fibery.Text,
				Description: "Name or number of the account associated with this vendor",
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
			"BillingAddress": {
				Name:    "Billing Address",
				Type:    fibery.Text,
				SubType: fibery.MD,
			},
			"BillingLine1": {
				Name: "Billing Line 1",
				Type: fibery.Text,
			},
			"BillingLine2": {
				Name: "Billing Line 2",
				Type: fibery.Text,
			},
			"BillingLine3": {
				Name: "Billing Line 3",
				Type: fibery.Text,
			},
			"BillingLine4": {
				Name: "Billing Line 4",
				Type: fibery.Text,
			},
			"BillingLine5": {
				Name: "Billing Line 5",
				Type: fibery.Text,
			},
			"BillingCity": {
				Name: "Billing City",
				Type: fibery.Text,
			},
			"BillingState": {
				Name: "Billing State",
				Type: fibery.Text,
			},
			"BillingPostalCode": {
				Name: "Billing Postal Code",
				Type: fibery.Text,
			},
			"BillingCountry": {
				Name: "Billing Country",
				Type: fibery.Text,
			},
			"BillingLat": {
				Name: "Billing Latitude",
				Type: fibery.Text,
			},
			"BillingLong": {
				Name: "Billing Longitude",
				Type: fibery.Text,
			},
		},
		func(v quickbooks.Vendor) (map[string]any, error) {
			var email string
			if v.PrimaryEmailAddr != nil {
				email = v.PrimaryEmailAddr.Address
			}

			var primaryPhone string
			if v.PrimaryPhone != nil {
				primaryPhone = v.PrimaryPhone.FreeFormNumber
			}

			var alternatePhone string
			if v.AlternatePhone != nil {
				alternatePhone = v.AlternatePhone.FreeFormNumber
			}

			var mobile string
			if v.Mobile != nil {
				mobile = v.Mobile.FreeFormNumber
			}

			var fax string
			if v.Fax != nil {
				fax = v.Fax.FreeFormNumber
			}

			var website string
			if v.WebAddr != nil {
				website = v.WebAddr.URI
			}

			var billAddr quickbooks.PhysicalAddress
			if v.BillAddr != nil {
				billAddr = *v.BillAddr
			}

			var termId string
			if v.TermRef != nil {
				termId = v.TermRef.Value
			}

			return map[string]any{
				"id":                v.Id,
				"QBOId":             v.Id,
				"DisplayName":       v.DisplayName,
				"SyncToken":         v.SyncToken,
				"__syncAction":      fibery.SET,
				"Active":            v.Active,
				"Title":             v.Title,
				"GivenName":         v.GivenName,
				"MiddleName":        v.MiddleName,
				"FamilyName":        v.FamilyName,
				"Suffix":            v.Suffix,
				"CompanyName":       v.CompanyName,
				"PrimaryEmail":      email,
				"SalesTermId":       termId,
				"PrimaryPhone":      primaryPhone,
				"AlternatePhone":    alternatePhone,
				"Mobile":            mobile,
				"Fax":               fax,
				"1099":              v.Vendor1099,
				"CostRate":          v.CostRate,
				"BillRate":          v.BillRate,
				"Website":           website,
				"Balance":           v.Balance,
				"BillingLine1":      billAddr.Line1,
				"BillingLine2":      billAddr.Line2,
				"BillingLine3":      billAddr.Line3,
				"BillingLine4":      billAddr.Line4,
				"BillingLine5":      billAddr.Line5,
				"BillingCity":       billAddr.City,
				"BillingState":      billAddr.CountrySubDivisionCode,
				"BillingPostalCode": billAddr.PostalCode,
				"BillingCountry":    billAddr.Country,
				"BillingLat":        billAddr.Lat,
				"BillingLong":       billAddr.Long,
			}, nil
		},
		func(client *quickbooks.Client, requestParams quickbooks.RequestParameters, startPosition, pageSize int) ([]quickbooks.Vendor, error) {
			items, err := client.FindVendorsByPage(requestParams, startPosition, pageSize)
			if err != nil {
				return nil, err
			}

			return items, nil
		},
		func(v quickbooks.Vendor) string {
			return v.Id
		},
		func(v quickbooks.Vendor) string {
			return v.Status
		},
	)

	tr.Register(vendor)

	// Build Union Types

	entity := NewUnionDataType(
		"Entity",
		"Entity",
		map[string]fibery.Field{
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
			"__syncAction": {
				Type: fibery.Text,
				Name: "Sync Action",
			},
			"CustomerId": {
				Name: "Customer ID",
				Type: fibery.Text,
				Relation: &fibery.Relation{
					Cardinality:   fibery.OTO,
					Name:          "Customer",
					TargetName:    "Entity",
					TargetType:    "Customer",
					TargetFieldID: "id",
				},
			},
			"EmployeeId": {
				Name: "Employee ID",
				Type: fibery.Text,
				Relation: &fibery.Relation{
					Cardinality:   fibery.OTO,
					Name:          "Employee",
					TargetName:    "Entity",
					TargetType:    "Employee",
					TargetFieldID: "id",
				},
			},
			"VendorId": {
				Name: "Vendor ID",
				Type: fibery.Text,
				Relation: &fibery.Relation{
					Cardinality:   fibery.OTO,
					Name:          "Vendor",
					TargetName:    "Entity",
					TargetType:    "Vendor",
					TargetFieldID: "id",
				},
			},
		},
		[]Type{customer, employee, vendor},
		func(typeId string, input []map[string]any) ([]map[string]any, error) {
			var items []map[string]any
			switch typeId {
			case "Customer":
				for _, inputItem := range input {
					var id, name string
					var syncAction fibery.SyncAction
					var ok bool
					if id, ok = inputItem["id"].(string); !ok {
						return nil, fmt.Errorf("unable to extract 'id' from Customer item")
					}
					if name, ok = inputItem["DisplayName"].(string); !ok {
						return nil, fmt.Errorf("unable to extract 'DisplayName' from Customer item")

					}
					if syncAction, ok = inputItem["__syncAction"].(fibery.SyncAction); !ok {
						return nil, fmt.Errorf("unable to extract '__syncAction' from Customer item")
					}

					item := map[string]any{
						"id":           "c:" + id,
						"QBOId":        id,
						"Name":         name + " (Customer)",
						"__syncAction": syncAction,
						"CustomerId":   id,
					}

					items = append(items, item)
				}
			case "Employee":
				for _, inputItem := range input {
					var id, name string
					var syncAction fibery.SyncAction
					var ok bool
					if id, ok = inputItem["id"].(string); !ok {
						return nil, fmt.Errorf("unable to extract 'id' from Employee item")
					}
					if name, ok = inputItem["DisplayName"].(string); !ok {
						return nil, fmt.Errorf("unable to extract 'DisplayName' from Employee item")

					}
					if syncAction, ok = inputItem["__syncAction"].(fibery.SyncAction); !ok {
						return nil, fmt.Errorf("unable to extract '__syncAction' from Employee item")
					}

					item := map[string]any{
						"id":           "e:" + id,
						"QBOId":        id,
						"Name":         name + " (Employee)",
						"__syncAction": syncAction,
						"EmployeeId":   id,
					}

					items = append(items, item)
				}
			case "Vendor":
				for _, inputItem := range input {
					var id, name string
					var syncAction fibery.SyncAction
					var ok bool
					if id, ok = inputItem["id"].(string); !ok {
						return nil, fmt.Errorf("unable to extract 'id' from Vendor item")
					}
					if name, ok = inputItem["DisplayName"].(string); !ok {
						return nil, fmt.Errorf("unable to extract 'DisplayName' from Vendor item")

					}
					if syncAction, ok = inputItem["__syncAction"].(fibery.SyncAction); !ok {
						return nil, fmt.Errorf("unable to extract '__syncAction' from Vendor item")
					}

					item := map[string]any{
						"id":           "v:" + id,
						"QBOId":        id,
						"Name":         name + " (Vendor)",
						"__syncAction": syncAction,
						"VendorId":     id,
					}

					items = append(items, item)
				}
			default:
				return nil, fmt.Errorf("invalid typeId: %s", typeId)
			}
			return items, nil
		},
	)

	tr.Register(entity)

	// Set related types

	purchase.relatedTypes = []CDCType{reimburseCharge}
	bill.relatedTypes = []CDCType{reimburseCharge}
	invoice.relatedTypes = []CDCType{reimburseCharge}
}
