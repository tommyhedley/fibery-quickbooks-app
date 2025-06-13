package types

import (
	"fmt"

	"github.com/tommyhedley/fibery-quickbooks-app/pkgs/fibery"
	"github.com/tommyhedley/fibery-quickbooks-app/pkgs/integration"
	"github.com/tommyhedley/quickbooks-go"
)

var bill = integration.NewDualType(
	"bill",
	"bill",
	"Bill",
	func(b quickbooks.Bill) string {
		return b.Id
	},
	func(b quickbooks.Bill) string {
		return b.Status
	},
	func(id string) quickbooks.Bill {
		return quickbooks.Bill{
			Id: id,
		}
	},
	func(bir quickbooks.BatchItemResponse) quickbooks.Bill {
		return bir.Bill
	},
	func(bqr quickbooks.BatchQueryResponse) []quickbooks.Bill {
		return bqr.Bill
	},
	func(cr quickbooks.CDCQueryResponse) []quickbooks.Bill {
		return cr.Bill
	},
	map[string]integration.FieldDef[quickbooks.Bill]{
		"QBOId": {
			Params: fibery.Field{
				Name: "QBO Id",
				Type: fibery.Text,
			},
			Convert: func(sd integration.StandardData[quickbooks.Bill]) (any, error) {
				return sd.Item.Id, nil
			},
		},
		"Name": {
			Params: fibery.Field{
				Name:    "Name",
				Type:    fibery.Text,
				SubType: fibery.Title,
			},
			Convert: func(sd integration.StandardData[quickbooks.Bill]) (any, error) {
				if sd.Item.PrivateNote == "" {
					return sd.Item.VendorRef.Name, nil
				}
				return sd.Item.VendorRef.Name + " – " + sd.Item.PrivateNote, nil
			},
		},
		"SyncToken": {
			Params: fibery.Field{
				Name:     "Sync Token",
				Type:     fibery.Text,
				ReadOnly: true,
			},
			Convert: func(sd integration.StandardData[quickbooks.Bill]) (any, error) {
				return sd.Item.SyncToken, nil
			},
		},
		"__syncAction": {
			Params: fibery.Field{
				Name: "Sync Action",
				Type: fibery.Text,
			},
			Convert: func(sd integration.StandardData[quickbooks.Bill]) (any, error) {
				return fibery.SET, nil
			},
		},
		"DocNumber": {
			Params: fibery.Field{
				Name: "Bill Number",
				Type: fibery.Text,
			},
			Convert: func(sd integration.StandardData[quickbooks.Bill]) (any, error) {
				return sd.Item.DocNumber, nil
			},
		},
		"TxnDate": {
			Params: fibery.Field{
				Name:    "Bill Date",
				Type:    fibery.DateType,
				SubType: fibery.Day,
			},
			Convert: func(sd integration.StandardData[quickbooks.Bill]) (any, error) {
				if sd.Item.TxnDate.IsZero() {
					return "", nil
				}
				return sd.Item.TxnDate.Format(fibery.DateFormat), nil
			},
		},
		"DueDate": {
			Params: fibery.Field{
				Name:    "Due Date",
				Type:    fibery.DateType,
				SubType: fibery.Day,
			},
			Convert: func(sd integration.StandardData[quickbooks.Bill]) (any, error) {
				if sd.Item.DueDate.IsZero() {
					return "", nil
				}
				return sd.Item.DueDate.Format(fibery.DateFormat), nil
			},
		},
		"PrivateNote": {
			Params: fibery.Field{
				Name:    "Memo",
				Type:    fibery.Text,
				SubType: fibery.MD,
			},
			Convert: func(sd integration.StandardData[quickbooks.Bill]) (any, error) {
				return sd.Item.PrivateNote, nil
			},
		},
		"TotalAmt": {
			Params: fibery.Field{
				Name: "Total",
				Type: fibery.Number,
				Format: map[string]any{
					"format":               "Money",
					"currencyCode":         "USD",
					"hasThousandSeperator": true,
					"precision":            2,
				},
			},
			Convert: func(sd integration.StandardData[quickbooks.Bill]) (any, error) {
				return sd.Item.TotalAmt, nil
			},
		},
		"Balance": {
			Params: fibery.Field{
				Name: "Balance",
				Type: fibery.Number,
				Format: map[string]any{
					"format":               "Money",
					"currencyCode":         "USD",
					"hasThousandSeperator": true,
					"precision":            2,
				},
			},
			Convert: func(sd integration.StandardData[quickbooks.Bill]) (any, error) {
				return sd.Item.Balance, nil
			},
		},
		"VendorId": {
			Params: fibery.Field{
				Name: "Vendor Id",
				Type: fibery.Text,
				Relation: &fibery.Relation{
					Cardinality:   fibery.MTO,
					Name:          "Vendor",
					TargetName:    "Bills",
					TargetType:    "Vendor",
					TargetFieldID: "id",
				},
			},
			Convert: func(sd integration.StandardData[quickbooks.Bill]) (any, error) {
				return sd.Item.VendorRef.Value, nil
			},
		},
		"APAccountId": {
			Params: fibery.Field{
				Name: "AP Account Id",
				Type: fibery.Text,
				Relation: &fibery.Relation{
					Cardinality:   fibery.MTO,
					Name:          "AP Account",
					TargetName:    "Bills",
					TargetType:    "Account",
					TargetFieldID: "id",
				},
			},
			Convert: func(sd integration.StandardData[quickbooks.Bill]) (any, error) {
				if sd.Item.APAccountRef != nil {
					return sd.Item.APAccountRef.Value, nil
				}
				return "", nil
			},
		},
		"SalesTermId": {
			Params: fibery.Field{
				Name: "Sales Term Id",
				Type: fibery.Text,
				Relation: &fibery.Relation{
					Cardinality:   fibery.MTO,
					Name:          "Terms",
					TargetName:    "Bills",
					TargetType:    "SalesTerm",
					TargetFieldID: "id",
				},
			},
			Convert: func(sd integration.StandardData[quickbooks.Bill]) (any, error) {
				if sd.Item.SalesTermRef != nil {
					return sd.Item.SalesTermRef.Value, nil
				}
				return "", nil
			},
		},
		"Files": {
			Params: fibery.Field{
				Name:    "Files",
				Type:    fibery.TextArray,
				SubType: fibery.File,
			},
			Convert: func(sd integration.StandardData[quickbooks.Bill]) (any, error) {
				id := sd.Item.Id
				if attachables, ok := sd.Attachables[id]; ok {
					output := make([]string, 0, len(attachables))
					for _, attachable := range attachables {
						url := integration.AttachableURL(attachable)
						output = append(output, url)
					}
					return output, nil
				}
				return nil, nil
			},
		},
	},
	[]integration.CDCType{reimburseCharge},
)

var billItemLine = integration.NewDependentDualType(
	"bill",
	"billItemLine",
	"Bill Item Line",
	func(b quickbooks.Bill, l quickbooks.Line) string {
		return fmt.Sprintf("%s:i:%s", b.Id, l.Id)
	},
	func(b quickbooks.Bill) []quickbooks.Line {
		return b.Line
	},
	func(b quickbooks.Bill) string {
		return b.Id
	},
	func(b quickbooks.Bill) string {
		return b.Status
	},
	func(id string) quickbooks.Bill {
		return quickbooks.Bill{
			Id: id,
		}
	},
	func(bir quickbooks.BatchItemResponse) quickbooks.Bill {
		return bir.Bill
	},
	func(bqr quickbooks.BatchQueryResponse) []quickbooks.Bill {
		return bqr.Bill
	},
	func(cr quickbooks.CDCQueryResponse) []quickbooks.Bill {
		return cr.Bill
	},
	map[string]integration.DependentFieldDef[quickbooks.Bill, quickbooks.Line]{
		"QBOId": {
			Params: fibery.Field{
				Name: "QBO ID",
				Type: fibery.Text,
			},
			Convert: func(dd integration.DependentData[quickbooks.Bill, quickbooks.Line]) (any, error) {
				return dd.Item.Id, nil
			},
		},
		"Name": {
			Params: fibery.Field{
				Name:    "Name",
				Type:    fibery.Text,
				SubType: fibery.Title,
			},
			Convert: func(dd integration.DependentData[quickbooks.Bill, quickbooks.Line]) (any, error) {
				var name string
				if dd.Item.Description == "" {
					name = dd.Item.ItemBasedExpenseLineDetail.ItemRef.Name
				} else {
					name = dd.Item.ItemBasedExpenseLineDetail.ItemRef.Name + dd.Item.Description
				}
				return name, nil
			},
		},
		"Description": {
			Params: fibery.Field{
				Name: "Description",
				Type: fibery.Text,
			},
			Convert: func(dd integration.DependentData[quickbooks.Bill, quickbooks.Line]) (any, error) {
				return dd.Item.Description, nil
			},
		},
		"__syncAction": {
			Params: fibery.Field{
				Type: fibery.Text,
				Name: "Sync Action",
			},
			Convert: func(dd integration.DependentData[quickbooks.Bill, quickbooks.Line]) (any, error) {
				return fibery.SET, nil
			},
		},
		"LineNum": {
			Params: fibery.Field{
				Name:    "Line",
				Type:    fibery.Number,
				SubType: fibery.Integer,
			},
			Convert: func(dd integration.DependentData[quickbooks.Bill, quickbooks.Line]) (any, error) {
				return dd.Item.LineNum, nil
			},
		},
		"Tax": {
			Params: fibery.Field{
				Name:    "Tax",
				Type:    fibery.Text,
				SubType: fibery.Boolean,
			},
			Convert: func(dd integration.DependentData[quickbooks.Bill, quickbooks.Line]) (any, error) {
				tax := false
				if dd.Item.ItemBasedExpenseLineDetail.TaxCodeRef.Value == "TAX" {
					tax = true
				}
				return tax, nil
			},
		},
		"Billable": {
			Params: fibery.Field{
				Name:    "Billable",
				Type:    fibery.Text,
				SubType: fibery.Boolean,
			},
			Convert: func(dd integration.DependentData[quickbooks.Bill, quickbooks.Line]) (any, error) {
				var billable bool
				switch dd.Item.ItemBasedExpenseLineDetail.BillableStatus {
				case quickbooks.BillableStatusType:
					billable = true
				case quickbooks.HasBeenBilledStatusType:
					billable = true
				default:
					billable = false
				}
				return billable, nil
			},
		},
		"Billed": {
			Params: fibery.Field{
				Name:    "Billed",
				Type:    fibery.Text,
				SubType: fibery.Boolean,
			},
			Convert: func(dd integration.DependentData[quickbooks.Bill, quickbooks.Line]) (any, error) {
				billed := false
				if dd.Item.ItemBasedExpenseLineDetail.BillableStatus == quickbooks.HasBeenBilledStatusType {
					billed = true
				}
				return billed, nil
			},
		},
		"Qty": {
			Params: fibery.Field{
				Name: "Quantity",
				Type: fibery.Number,
				Format: map[string]any{
					"format":               "Number",
					"hasThousandSeparator": true,
					"precision":            2,
				},
			},
			Convert: func(dd integration.DependentData[quickbooks.Bill, quickbooks.Line]) (any, error) {
				return dd.Item.ItemBasedExpenseLineDetail.Qty, nil
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
			Convert: func(dd integration.DependentData[quickbooks.Bill, quickbooks.Line]) (any, error) {
				return dd.Item.ItemBasedExpenseLineDetail.UnitPrice, nil
			},
		},
		"MarkupPercent": {
			Params: fibery.Field{
				Name: "Markup",
				Type: fibery.Number,
				Format: map[string]any{
					"format":    "Percent",
					"precision": 2,
				},
			},
			Convert: func(dd integration.DependentData[quickbooks.Bill, quickbooks.Line]) (any, error) {
				return dd.Item.ItemBasedExpenseLineDetail.MarkupInfo, nil
			},
		},
		"Amount": {
			Params: fibery.Field{
				Name: "Amount",
				Type: fibery.Number,
				Format: map[string]any{
					"format":               "Money",
					"currencyCode":         "USD",
					"hasThousandSeperator": true,
					"precision":            2,
				},
			},
			Convert: func(dd integration.DependentData[quickbooks.Bill, quickbooks.Line]) (any, error) {
				return dd.Item.Amount, nil
			},
		},
		"BillId": {
			Params: fibery.Field{
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
			Convert: func(dd integration.DependentData[quickbooks.Bill, quickbooks.Line]) (any, error) {
				return dd.SourceItem.Id, nil
			},
		},
		"ItemId": {
			Params: fibery.Field{
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
			Convert: func(dd integration.DependentData[quickbooks.Bill, quickbooks.Line]) (any, error) {
				return dd.Item.ItemBasedExpenseLineDetail.ItemRef.Value, nil
			},
		},
		"CustomerId": {
			Params: fibery.Field{
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
			Convert: func(dd integration.DependentData[quickbooks.Bill, quickbooks.Line]) (any, error) {
				return dd.Item.ItemBasedExpenseLineDetail.CustomerRef.Value, nil
			},
		},
		"ClassId": {
			Params: fibery.Field{
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
			Convert: func(dd integration.DependentData[quickbooks.Bill, quickbooks.Line]) (any, error) {
				return dd.Item.ItemBasedExpenseLineDetail.ClassRef.Value, nil
			},
		},
		"MarkupAccountId": {
			Params: fibery.Field{
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
			Convert: func(dd integration.DependentData[quickbooks.Bill, quickbooks.Line]) (any, error) {
				return dd.Item.ItemBasedExpenseLineDetail.MarkupInfo.MarkUpIncomeAccountRef.Value, nil
			},
		},
		"ReimburseChargeId": {
			Params: fibery.Field{
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
			Convert: func(dd integration.DependentData[quickbooks.Bill, quickbooks.Line]) (any, error) {
				var reimburseChargeId string
				for _, txn := range dd.Item.LinkedTxn {
					if txn.TxnType == "ReimburseCharge" {
						reimburseChargeId = txn.TxnId
					}
				}
				return reimburseChargeId, nil
			},
		},
	},
)

func init() {
	integration.Types.Register(bill)
	integration.Types.Register(billItemLine)
}
