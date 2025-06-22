package types

import (
	"fmt"
	"math"
	"strings"

	"github.com/tommyhedley/fibery-quickbooks-app/pkgs/app"
	"github.com/tommyhedley/fibery-quickbooks-app/pkgs/fibery"
	"github.com/tommyhedley/quickbooks-go"
)

var reimburseCharge = app.NewCDCType(
	"Reimbursecharge",
	"reimbursecharge",
	"Reimburse Charge",
	func(r quickbooks.ReimburseCharge) string {
		return r.Id
	},
	func(r quickbooks.ReimburseCharge) string {
		return r.Status
	},
	func(bir quickbooks.BatchItemResponse) quickbooks.ReimburseCharge {
		return bir.ReimburseCharge
	},
	func(bqr quickbooks.BatchQueryResponse) []quickbooks.ReimburseCharge {
		return bqr.ReimburseCharge
	},
	func(cr quickbooks.CDCQueryResponse) []quickbooks.ReimburseCharge {
		return cr.ReimburseCharge
	},
	map[string]app.FieldDef[quickbooks.ReimburseCharge]{
		"qboId": {
			Params: fibery.Field{
				Name:     "QBO ID",
				Type:     fibery.Text,
				ReadOnly: true,
			},
			Convert: func(sd app.StandardData[quickbooks.ReimburseCharge]) (any, error) {
				return sd.Item.Id, nil
			},
		},
		"name": {
			Params: fibery.Field{
				Name:    "Name",
				Type:    fibery.Text,
				SubType: fibery.Title,
			},
			Convert: func(sd app.StandardData[quickbooks.ReimburseCharge]) (any, error) {
				var name string
				if sd.Item.PrivateNote != "" {
					name = sd.Item.CustomerRef.Name + " - " + sd.Item.PrivateNote
				} else {
					name = sd.Item.CustomerRef.Name
				}
				return name, nil
			},
		},
		"description": {
			Params: fibery.Field{
				Name: "Description",
				Type: fibery.Text,
			},
			Convert: func(sd app.StandardData[quickbooks.ReimburseCharge]) (any, error) {
				return sd.Item.PrivateNote, nil
			},
		},
		"__syncAction": {
			Params: fibery.Field{
				Type: fibery.Text,
				Name: "Sync Action",
			},
			Convert: func(sd app.StandardData[quickbooks.ReimburseCharge]) (any, error) {
				return fibery.SET, nil
			},
		},
		"syncToken": {
			Params: fibery.Field{
				Name:     "Sync Token",
				Type:     fibery.Text,
				ReadOnly: true,
			},
			Convert: func(sd app.StandardData[quickbooks.ReimburseCharge]) (any, error) {
				return sd.Item.SyncToken, nil
			},
		},
		"txnDate": {
			Params: fibery.Field{
				Name:    "Date",
				Type:    fibery.DateType,
				SubType: fibery.Day,
			},
			Convert: func(sd app.StandardData[quickbooks.ReimburseCharge]) (any, error) {
				if sd.Item.TxnDate != nil {
					return sd.Item.TxnDate.Format(fibery.DateFormat), nil
				}
				return "", nil
			},
		},
		"customerId": {
			Params: fibery.Field{
				Name: "Customer ID",
				Type: fibery.Text,
				Relation: &fibery.Relation{
					Cardinality:   fibery.MTO,
					Name:          "Customer",
					TargetName:    "Reimburse Charges",
					TargetType:    "customer",
					TargetFieldID: "id",
				},
			},
			Convert: func(sd app.StandardData[quickbooks.ReimburseCharge]) (any, error) {
				return sd.Item.CustomerRef.Value, nil
			},
		},
		"totalAmount": {
			Params: fibery.Field{
				Name: "Total Amount",
				Type: fibery.Number,
				Format: map[string]any{
					"format":               "Money",
					"currencyCode":         "USD",
					"hasThousandSeperator": true,
					"precision":            2,
				},
			},
			Convert: func(sd app.StandardData[quickbooks.ReimburseCharge]) (any, error) {
				return sd.Item.Amount, nil
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
			Convert: func(sd app.StandardData[quickbooks.ReimburseCharge]) (any, error) {
				for _, line := range sd.Item.Line {
					if line.LineNum == 1 {
						return line.Amount, nil
					}
				}
				return 0, nil
			},
		},
		"accountId": {
			Params: fibery.Field{
				Name: "Account ID",
				Type: fibery.Text,
				Relation: &fibery.Relation{
					Cardinality:   fibery.MTO,
					Name:          "Account",
					TargetName:    "Reimburse Charges",
					TargetType:    "account",
					TargetFieldID: "id",
				},
			},
			Convert: func(sd app.StandardData[quickbooks.ReimburseCharge]) (any, error) {
				for _, line := range sd.Item.Line {
					if line.LineNum == 1 {
						return line.ReimburseLineDetail.ItemAccountRef.Value, nil
					}
				}
				return "", nil
			},
		},
		"markup": {
			Params: fibery.Field{
				Name: "Markup",
				Type: fibery.Text,
			},
			Convert: func(sd app.StandardData[quickbooks.ReimburseCharge]) (any, error) {
				for _, line := range sd.Item.Line {
					if line.LineNum == 2 {
						markupAmount, err := line.Amount.Float64()
						if err != nil {
							return "", fmt.Errorf("error converting markup amount json.Number to float: %w", err)
						}

						markupPercent, err := line.ReimburseLineDetail.MarkupInfo.Percent.Float64()
						if err != nil {
							return "", fmt.Errorf("error converting markup percent json.Number to float: %w", err)
						}

						var markupStr string
						if math.Mod(markupPercent, 1.0) == 0 {
							markupStr = fmt.Sprintf("%.0f%%", markupPercent)
						} else {
							markupStr = fmt.Sprintf("%.5f%%", markupPercent)
							markupStr = strings.TrimRight(markupStr, "0")
							markupStr = strings.TrimRight(markupStr, ".")
						}

						return fmt.Sprintf("$%.2f (%s)", markupAmount, markupStr), nil
					}
				}
				return "", nil
			},
		},
		"markupAccountId": {
			Params: fibery.Field{
				Name: "Markup Account ID",
				Type: fibery.Text,
				Relation: &fibery.Relation{
					Cardinality:   fibery.MTO,
					Name:          "Markup Account",
					TargetName:    "Reimburse Charge Markup",
					TargetType:    "account",
					TargetFieldID: "id",
				},
			},
			Convert: func(sd app.StandardData[quickbooks.ReimburseCharge]) (any, error) {
				for _, line := range sd.Item.Line {
					if line.LineNum == 2 {
						return line.ReimburseLineDetail.ItemRef.Value, nil
					}
				}
				return "", nil
			},
		},
		"taxable": {
			Params: fibery.Field{
				Name:    "Taxable",
				Type:    fibery.Text,
				SubType: fibery.Boolean,
			},
			Convert: func(sd app.StandardData[quickbooks.ReimburseCharge]) (any, error) {
				for _, line := range sd.Item.Line {
					if line.LineNum == 1 {
						return line.ReimburseLineDetail.TaxCodeRef.Value == "TAX", nil
					}
				}
				return false, nil
			},
		},
		"linkedInvoiceId": {
			Params: fibery.Field{
				Name: "Linked Invoice ID",
				Type: fibery.Text,
				Relation: &fibery.Relation{
					Cardinality:   fibery.MTO,
					Name:          "Linked Invoice",
					TargetName:    "Reimburse Charges",
					TargetType:    "invoice",
					TargetFieldID: "id",
				},
			},
			Convert: func(sd app.StandardData[quickbooks.ReimburseCharge]) (any, error) {
				for _, txn := range sd.Item.LinkedTxn {
					if txn.TxnType == "Invoice" {
						return txn.TxnId, nil
					}
				}
				return "", nil
			},
		},
	},
)

func init() {
	app.Types.Register(reimburseCharge)
}
