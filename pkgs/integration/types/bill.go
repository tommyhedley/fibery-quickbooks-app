package types

import (
	"fmt"

	"github.com/tommyhedley/fibery-quickbooks-app/pkgs/fibery"
	"github.com/tommyhedley/fibery-quickbooks-app/pkgs/integration"
	"github.com/tommyhedley/quickbooks-go"
)

var bill = integration.NewDualType(
	"Bill",
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
		"qboId": {
			Params: fibery.Field{
				Name:     "QBO Id",
				Type:     fibery.Text,
				ReadOnly: true,
			},
			Convert: func(sd integration.StandardData[quickbooks.Bill]) (any, error) {
				return sd.Item.Id, nil
			},
		},
		"name": {
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
		"syncToken": {
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
		"docNumber": {
			Params: fibery.Field{
				Name: "Bill Number",
				Type: fibery.Text,
			},
			Convert: func(sd integration.StandardData[quickbooks.Bill]) (any, error) {
				return sd.Item.DocNumber, nil
			},
		},
		"txnDate": {
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
		"dueDate": {
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
		"privateNote": {
			Params: fibery.Field{
				Name:    "Memo",
				Type:    fibery.Text,
				SubType: fibery.MD,
			},
			Convert: func(sd integration.StandardData[quickbooks.Bill]) (any, error) {
				return sd.Item.PrivateNote, nil
			},
		},
		"totalAmt": {
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
		"balance": {
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
		"vendorId": {
			Params: fibery.Field{
				Name: "Vendor Id",
				Type: fibery.Text,
				Relation: &fibery.Relation{
					Cardinality:   fibery.MTO,
					Name:          "Vendor",
					TargetName:    "Bills",
					TargetType:    "vendor",
					TargetFieldID: "id",
				},
			},
			Convert: func(sd integration.StandardData[quickbooks.Bill]) (any, error) {
				return sd.Item.VendorRef.Value, nil
			},
		},
		"apAccountId": {
			Params: fibery.Field{
				Name: "AP Account Id",
				Type: fibery.Text,
				Relation: &fibery.Relation{
					Cardinality:   fibery.MTO,
					Name:          "AP Account",
					TargetName:    "Bills",
					TargetType:    "account",
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
		"salesTermId": {
			Params: fibery.Field{
				Name: "Sales Term Id",
				Type: fibery.Text,
				Relation: &fibery.Relation{
					Cardinality:   fibery.MTO,
					Name:          "Terms",
					TargetName:    "Bills",
					TargetType:    "salesTerm",
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
		"attachables": {
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
		items := make([]quickbooks.Line, 0)
		for _, line := range b.Line {
			if line.DetailType == quickbooks.ItemExpenseLine {
				items = append(items, line)
			}
		}
		return items
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
		"qboId": {
			Params: fibery.Field{
				Name:     "QBO ID",
				Type:     fibery.Text,
				ReadOnly: true,
			},
			Convert: func(dd integration.DependentData[quickbooks.Bill, quickbooks.Line]) (any, error) {
				return dd.Item.Id, nil
			},
		},
		"name": {
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
					name = dd.Item.ItemBasedExpenseLineDetail.ItemRef.Name + " - " + dd.Item.Description
				}
				return name, nil
			},
		},
		"description": {
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
		"lineNum": {
			Params: fibery.Field{
				Name:    "Line",
				Type:    fibery.Number,
				SubType: fibery.Integer,
			},
			Convert: func(dd integration.DependentData[quickbooks.Bill, quickbooks.Line]) (any, error) {
				return dd.Item.LineNum, nil
			},
		},
		"tax": {
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
		"billable": {
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
		"billed": {
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
		"qty": {
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
			Convert: func(dd integration.DependentData[quickbooks.Bill, quickbooks.Line]) (any, error) {
				return dd.Item.ItemBasedExpenseLineDetail.UnitPrice, nil
			},
		},
		"markupPercent": {
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
		"amount": {
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
		"billId": {
			Params: fibery.Field{
				Name: "Bill ID",
				Type: fibery.Text,
				Relation: &fibery.Relation{
					Cardinality:   fibery.MTO,
					Name:          "Bill",
					TargetName:    "Item Lines",
					TargetType:    "bill",
					TargetFieldID: "id",
				},
			},
			Convert: func(dd integration.DependentData[quickbooks.Bill, quickbooks.Line]) (any, error) {
				return dd.SourceItem.Id, nil
			},
		},
		"itemId": {
			Params: fibery.Field{
				Name: "Item ID",
				Type: fibery.Text,
				Relation: &fibery.Relation{
					Cardinality:   fibery.MTO,
					Name:          "Item",
					TargetName:    "Bill Item Lines",
					TargetType:    "item",
					TargetFieldID: "id",
				},
			},
			Convert: func(dd integration.DependentData[quickbooks.Bill, quickbooks.Line]) (any, error) {
				return dd.Item.ItemBasedExpenseLineDetail.ItemRef.Value, nil
			},
		},
		"customerId": {
			Params: fibery.Field{
				Name: "Customer ID",
				Type: fibery.Text,
				Relation: &fibery.Relation{
					Cardinality:   fibery.MTO,
					Name:          "Customer",
					TargetName:    "Bill Item Lines",
					TargetType:    "customer",
					TargetFieldID: "id",
				},
			},
			Convert: func(dd integration.DependentData[quickbooks.Bill, quickbooks.Line]) (any, error) {
				return dd.Item.ItemBasedExpenseLineDetail.CustomerRef.Value, nil
			},
		},
		"classId": {
			Params: fibery.Field{
				Name: "Class ID",
				Type: fibery.Text,
				Relation: &fibery.Relation{
					Cardinality:   fibery.MTO,
					Name:          "Class",
					TargetName:    "Bill Item Lines",
					TargetType:    "class",
					TargetFieldID: "id",
				},
			},
			Convert: func(dd integration.DependentData[quickbooks.Bill, quickbooks.Line]) (any, error) {
				return dd.Item.ItemBasedExpenseLineDetail.ClassRef.Value, nil
			},
		},
		"markupAccountId": {
			Params: fibery.Field{
				Name: "Markup Account ID",
				Type: fibery.Text,
				Relation: &fibery.Relation{
					Cardinality:   fibery.MTO,
					Name:          "Markup Income Account",
					TargetName:    "Bill Item Line Markup",
					TargetType:    "account",
					TargetFieldID: "id",
				},
			},
			Convert: func(dd integration.DependentData[quickbooks.Bill, quickbooks.Line]) (any, error) {
				return dd.Item.ItemBasedExpenseLineDetail.MarkupInfo.MarkUpIncomeAccountRef.Value, nil
			},
		},
		"reimburseChargeId": {
			Params: fibery.Field{
				Name: "Reimburse Charge ID",
				Type: fibery.Text,
				Relation: &fibery.Relation{
					Cardinality:   fibery.OTO,
					Name:          "Reimburse Charge",
					TargetName:    "Bill Item Line",
					TargetType:    "reimburseCharge",
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
