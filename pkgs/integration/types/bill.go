package types

import (
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
)

func init() {
	integration.Types.Register(bill)
}
