package data

import "github.com/tommyhedley/fibery-quickbooks-app/pkgs/fibery"

var BillPayment = QuickBooksDualType{
	QuickBooksType: QuickBooksType{
		fiberyType: fiberyType{
			id:   "BillPayment",
			name: "Bill Payement",
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
						TargetFieldID: "Id",
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
						TargetFieldID: "Id",
					},
				},
			},
		},
		schemaGen:      func(entity any) (map[string]any, error) {},
		query:          func(req Request) (Response, error) {},
		queryProcessor: func(entityArray any, schemaGen schemaGenFunc) ([]map[string]any, error) {},
	},
	cdcProcessor: func(cdc quickbooks.ChangeDataCapture, schemaGen schemaGenFunc) ([]map[string]any, error) {},
	whBatchProcessor: func(itemResponse quickbooks.BatchItemResponse, response *map[string][]map[string]any, cache *cache.Cache, realmId string, queryProcessor queryProcessorFunc, schemaGen schemaGenFunc, typeId string) error {
	},
}

var BillPaymentLine = DependentDualType{
	dependentBaseType: dependentBaseType{
		fiberyType: fiberyType{
			id:   "BillPaymentLine",
			name: "Bill Payment Line",
			schema: map[string]fibery.Field{
				"Id": {
					Name: "ID",
					Type: fibery.Text,
				},
				"QBOId": {
					Name: "QBO ID",
					Type: fibery.Text,
				},
				"Name": {
					Name: "Name",
					Type: fibery.Text,
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
						TargetFieldID: "Id",
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
						TargetFieldID: "Id",
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
						TargetFieldID: "Id",
					},
				},
			},
		},
		schemaGen:      func(entity, source any) (map[string]any, error) {},
		queryProcessor: func(sourceArray any, schemaGen depSchemaGenFunc) ([]map[string]any, error) {},
	},
	source:       BillPayment,
	sourceMapper: func(source any) (map[string]bool, error) {},
	typeMapper:   func(sourceArray any, sourceMapper sourceMapperFunc) (map[string]map[string]bool, error) {},
	cdcProcessor: func(cdc quickbooks.ChangeDataCapture, cacheEntry *IdCache, sourceMapper sourceMapperFunc, schemaGen depSchemaGenFunc) ([]map[string]any, error) {
	},
	whBatchProcessor: func(sourceArray any, cacheEntry *IdCache, sourceMapper sourceMapperFunc, schemaGen depSchemaGenFunc) ([]map[string]any, error) {
	},
}
