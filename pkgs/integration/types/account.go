package types

import (
	"github.com/tommyhedley/fibery-quickbooks-app/pkgs/fibery"
	"github.com/tommyhedley/fibery-quickbooks-app/pkgs/integration"
	"github.com/tommyhedley/quickbooks-go"
)

var account = integration.NewDualType(
	"account",
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
	map[string]integration.FieldDef[quickbooks.Account]{
		"QBOId": {
			Params: fibery.Field{
				Name: "QBO Id",
				Type: fibery.Text,
			},
			Convert: func(sd integration.StandardData[quickbooks.Account]) (any, error) {
				return sd.Item.Id, nil
			},
		},
		"Name": {
			Params: fibery.Field{
				Name: "Base Name",
				Type: fibery.Text,
			},
			Convert: func(sd integration.StandardData[quickbooks.Account]) (any, error) {
				return sd.Item.Name, nil
			},
		},
		"FullyQualifiedName": {
			Params: fibery.Field{
				Name:    "Full Name",
				Type:    fibery.Text,
				SubType: fibery.Title,
			},
			Convert: func(sd integration.StandardData[quickbooks.Account]) (any, error) {
				return sd.Item.FullyQualifiedName, nil
			},
		},
		"SyncToken": {
			Params: fibery.Field{
				Name:     "Sync Token",
				Type:     fibery.Text,
				ReadOnly: true,
			},
			Convert: func(sd integration.StandardData[quickbooks.Account]) (any, error) {
				return sd.Item.SyncToken, nil
			},
		},
		"__syncAction": {
			Params: fibery.Field{
				Type: fibery.Text,
				Name: "Sync Action",
			},
			Convert: func(sd integration.StandardData[quickbooks.Account]) (any, error) {
				return fibery.SET, nil
			},
		},
		"Active": {
			Params: fibery.Field{
				Name:    "Active",
				Type:    fibery.Text,
				SubType: fibery.Boolean,
			},
			Convert: func(sd integration.StandardData[quickbooks.Account]) (any, error) {
				return sd.Item.Active, nil
			},
		},
		"Description": {
			Params: fibery.Field{
				Name:    "Description",
				Type:    fibery.Text,
				SubType: fibery.MD,
			},
			Convert: func(sd integration.StandardData[quickbooks.Account]) (any, error) {
				return sd.Item.Description, nil
			},
		},
		"AcctNum": {
			Params: fibery.Field{
				Name: "Account Number",
				Type: fibery.Text,
			},
			Convert: func(sd integration.StandardData[quickbooks.Account]) (any, error) {
				return sd.Item.AcctNum, nil
			},
		},
		"CurrentBalance": {
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
			Convert: func(sd integration.StandardData[quickbooks.Account]) (any, error) {
				return sd.Item.CurrentBalance, nil
			},
		},
		"CurrentBalanceWithSubAccounts": {
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
			Convert: func(sd integration.StandardData[quickbooks.Account]) (any, error) {
				return sd.Item.CurrentBalanceWithSubAccounts, nil
			},
		},
		"Classification": {
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
			Convert: func(sd integration.StandardData[quickbooks.Account]) (any, error) {
				return sd.Item.Classification, nil
			},
		},
		"AccountType": {
			Params: fibery.Field{
				Name:     "Account Type",
				Type:     fibery.Text,
				SubType:  fibery.SingleSelect,
				ReadOnly: true,
			},
			Convert: func(sd integration.StandardData[quickbooks.Account]) (any, error) {
				return sd.Item.AccountType, nil
			},
		},
		"AccountSubType": {
			Params: fibery.Field{
				Name:     "Account Sub-Type",
				Type:     fibery.Text,
				SubType:  fibery.SingleSelect,
				ReadOnly: true,
			},
			Convert: func(sd integration.StandardData[quickbooks.Account]) (any, error) {
				return sd.Item.AccountSubType, nil
			},
		},
		"ParentAccountId": {
			Params: fibery.Field{
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
			Convert: func(sd integration.StandardData[quickbooks.Account]) (any, error) {
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
	integration.Types.Register(account)
}
