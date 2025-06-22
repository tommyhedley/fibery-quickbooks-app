package types

import (
	"fmt"
	"regexp"

	"github.com/tommyhedley/fibery-quickbooks-app/pkgs/app"
	"github.com/tommyhedley/fibery-quickbooks-app/pkgs/fibery"
	"github.com/tommyhedley/quickbooks-go"
)

var account = app.NewDualType(
	"Account",
	"account",
	"Account",
	func(a quickbooks.Account) string {
		return a.Id
	},
	func(a quickbooks.Account) string {
		return a.Status
	},
	func(id string) quickbooks.Account {
		return quickbooks.Account{
			Id: id,
		}
	},
	func(bir quickbooks.BatchItemResponse) quickbooks.Account {
		return bir.Account
	},
	func(bqr quickbooks.BatchQueryResponse) []quickbooks.Account {
		return bqr.Account
	},
	func(cr quickbooks.CDCQueryResponse) []quickbooks.Account {
		return cr.Account
	},
	map[string]app.FieldDef[quickbooks.Account]{
		"qboId": {
			Params: fibery.Field{
				Name:     "QBO Id",
				Type:     fibery.Text,
				ReadOnly: true,
			},
			Convert: func(sd app.StandardData[quickbooks.Account]) (any, error) {
				return sd.Item.Id, nil
			},
		},
		"name": {
			Params: fibery.Field{
				Name: "Base Name",
				Type: fibery.Text,
			},
			Convert: func(sd app.StandardData[quickbooks.Account]) (any, error) {
				return sd.Item.Name, nil
			},
		},
		"fullyQualifiedName": {
			Params: fibery.Field{
				Name:    "Full Name",
				Type:    fibery.Text,
				SubType: fibery.Title,
			},
			Convert: func(sd app.StandardData[quickbooks.Account]) (any, error) {
				return sd.Item.FullyQualifiedName, nil
			},
		},
		"syncToken": {
			Params: fibery.Field{
				Name:     "Sync Token",
				Type:     fibery.Text,
				ReadOnly: true,
			},
			Convert: func(sd app.StandardData[quickbooks.Account]) (any, error) {
				return sd.Item.SyncToken, nil
			},
		},
		"__syncAction": {
			Params: fibery.Field{
				Type: fibery.Text,
				Name: "Sync Action",
			},
			Convert: func(sd app.StandardData[quickbooks.Account]) (any, error) {
				return fibery.SET, nil
			},
		},
		"active": {
			Params: fibery.Field{
				Name:    "Active",
				Type:    fibery.Text,
				SubType: fibery.Boolean,
			},
			Convert: func(sd app.StandardData[quickbooks.Account]) (any, error) {
				return sd.Item.Active, nil
			},
		},
		"description": {
			Params: fibery.Field{
				Name:    "Description",
				Type:    fibery.Text,
				SubType: fibery.MD,
			},
			Convert: func(sd app.StandardData[quickbooks.Account]) (any, error) {
				return sd.Item.Description, nil
			},
		},
		"acctNum": {
			Params: fibery.Field{
				Name: "Account Number",
				Type: fibery.Text,
			},
			Convert: func(sd app.StandardData[quickbooks.Account]) (any, error) {
				return sd.Item.AcctNum, nil
			},
		},
		"currentBalance": {
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
			Convert: func(sd app.StandardData[quickbooks.Account]) (any, error) {
				return sd.Item.CurrentBalance, nil
			},
		},
		"currentBalanceWithSubAccounts": {
			Params: fibery.Field{
				Name: "Balance With Sub-Accounts",
				Type: fibery.Number,
				Format: map[string]any{
					"format":               "Money",
					"currencyCode":         "USD",
					"hasThousandSeperator": true,
					"precision":            2,
				},
			},
			Convert: func(sd app.StandardData[quickbooks.Account]) (any, error) {
				return sd.Item.CurrentBalanceWithSubAccounts, nil
			},
		},
		"classification": {
			Params: fibery.Field{
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
			Convert: func(sd app.StandardData[quickbooks.Account]) (any, error) {
				return sd.Item.Classification, nil
			},
		},
		"accountType": {
			Params: fibery.Field{
				Name:     "Account Type",
				Type:     fibery.Text,
				SubType:  fibery.SingleSelect,
				ReadOnly: true,
			},
			Convert: func(sd app.StandardData[quickbooks.Account]) (any, error) {
				return sd.Item.AccountType, nil
			},
		},
		"accountSubType": {
			Params: fibery.Field{
				Name:     "Account Sub-Type",
				Type:     fibery.Text,
				SubType:  fibery.SingleSelect,
				ReadOnly: true,
			},
			Convert: func(sd app.StandardData[quickbooks.Account]) (any, error) {
				reg, err := regexp.Compile(`([a-z])([A-Z])`)
				if err != nil {
					return nil, fmt.Errorf("error creating regex: %w", err)
				}

				subtype := reg.ReplaceAllString(sd.Item.AccountSubType, `$1 $2`)
				return subtype, nil
			},
		},
		"parentAccountId": {
			Params: fibery.Field{
				Name: "Parent Account ID",
				Type: fibery.Text,
				Relation: &fibery.Relation{
					Cardinality:   fibery.MTO,
					Name:          "Parent Account",
					TargetName:    "Sub-Accounts",
					TargetType:    "account",
					TargetFieldID: "id",
				},
			},
			Convert: func(sd app.StandardData[quickbooks.Account]) (any, error) {
				var parentAccountId string
				if sd.Item.ParentRef != nil {
					parentAccountId = sd.Item.ParentRef.Value
				}
				return parentAccountId, nil
			},
		},
	},
	nil,
)

func init() {
	app.Types.Register(account)
}
