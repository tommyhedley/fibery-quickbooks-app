package data

var Purchase = QuickBooksDualType{
	QuickBooksType: QuickBooksType{
		fiberyType: fiberyType{
			id:   "Purchase",
			name: "Expense",
			schema: map[string]fibery.Field{
				"id": {
					Name: "id",
					Type: fibery.ID,
				},
				"qbo_id": {
					Name: "QBO ID",
					Type: fibery.Text,
				},
				"name": {
					Name: "Name",
					Type: fibery.Text,
				},
				"sync_token": {
					Name:     "Sync Token",
					Type:     fibery.Text,
					ReadOnly: true,
				},
				"payment_type": {
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
				"payment_account_id": {
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
				"payment_method_id": {
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
				"entity_id": {
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
				"total": {
					Name: "Total",
					Type: fibery.Number,
					Format: map[string]any{
						"format":               "Money",
						"currencyCode":         "USD",
						"hasThousandSeperator": true,
						"precision":            2,
					},
				},
				"doc_number": {
					Name: "Reference Number",
					Type: fibery.Text,
				},
				"txn_date": {
					Name: "Invoice Date",
					Type: fibery.DateType,
				},
				"credit": {
					Name:    "Credit",
					Type:    fibery.Text,
					SubType: fibery.Boolean,
				},
				"private_note": {
					Name:    "Memo",
					Type:    fibery.Text,
					SubType: fibery.MD,
				},
				"__syncAction": {
					Type: fibery.Text,
					Name: "Sync Action",
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
var PurchaseItemLine = DependentDualType{
	dependentBaseType: dependentBaseType{
		fiberyType: fiberyType{
			id:   "PurchaseItemLine",
			name: "Purchase Item Line",
			schema: map[string]fibery.Field{
				"id": {
					Name: "id",
					Type: fibery.ID,
				},
				"qbo_id": {
					Name: "QBO ID",
					Type: fibery.Text,
				},
				"sync_token": {
					Name:     "Sync Token",
					Type:     fibery.Text,
					ReadOnly: true,
				},
				"description": {
					Name:    "Description",
					Type:    fibery.Text,
					SubType: fibery.Title,
				},
				"line_num": {
					Name: "Line",
					Type: fibery.Number,
				},
				"bill_id": {
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
				"item_id": {
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
				"customer_id": {
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
				"class_id": {
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
				"tax": {
					Name:    "Tax",
					Type:    fibery.Text,
					SubType: fibery.Boolean,
				},
				"markup_percent": {
					Name: "Markup",
					Type: fibery.Number,
					Format: map[string]any{
						"format":    "Percent",
						"precision": 2,
					},
				},
				"markup_account_id": {
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
				"billable": {
					Name:    "Billable",
					Type:    fibery.Text,
					SubType: fibery.Boolean,
				},
				"billed": {
					Name:    "Billed",
					Type:    fibery.Text,
					SubType: fibery.Boolean,
				},
				"quantity": {
					Name: "Quantity",
					Type: fibery.Number,
					Format: map[string]any{
						"format":               "Number",
						"hasThousandSeparator": true,
						"precision":            2,
					},
				},
				"unit_price": {
					Name: "Unit Price",
					Type: fibery.Number,
					Format: map[string]any{
						"format":               "Money",
						"currencyCode":         "USD",
						"hasThousandSeperator": true,
						"precision":            2,
					},
				},
				"amount": {
					Name: "Amount",
					Type: fibery.Number,
					Format: map[string]any{
						"format":               "Money",
						"currencyCode":         "USD",
						"hasThousandSeperator": true,
						"precision":            2,
					},
				},
				"__syncAction": {
					Type: fibery.Text,
					Name: "Sync Action",
				},
			},
		},
		schemaGen:      func(entity, source any) (map[string]any, error) {},
		queryProcessor: func(sourceArray any, schemaGen depSchemaGenFunc) ([]map[string]any, error) {},
	},
	source:       Bill,
	sourceMapper: func(source any) (map[string]bool, error) {},
	typeMapper:   func(sourceArray any, sourceMapper sourceMapperFunc) (map[string]map[string]bool, error) {},
	cdcProcessor: func(cdc quickbooks.ChangeDataCapture, cacheEntry *IdCache, sourceMapper sourceMapperFunc, schemaGen depSchemaGenFunc) ([]map[string]any, error) {
	},
	whBatchProcessor: func(sourceArray any, cacheEntry *IdCache, sourceMapper sourceMapperFunc, schemaGen depSchemaGenFunc) ([]map[string]any, error) {
	},
}

var PurchaseAccountLine = DependentDualType{
	dependentBaseType: dependentBaseType{
		fiberyType: fiberyType{
			id:   "PurchaseAccountLine",
			name: "Purchase Account Line",
			schema: map[string]fibery.Field{
				"id": {
					Name: "id",
					Type: fibery.ID,
				},
				"qbo_id": {
					Name: "QBO ID",
					Type: fibery.Text,
				},
				"sync_token": {
					Name:     "Sync Token",
					Type:     fibery.Text,
					ReadOnly: true,
				},
				"description": {
					Name:    "Description",
					Type:    fibery.Text,
					SubType: fibery.Title,
				},
				"line_num": {
					Name: "Line",
					Type: fibery.Number,
				},
				"bill_id": {
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
				"account_id": {
					Name: "Item ID",
					Type: fibery.Text,
					Relation: &fibery.Relation{
						Cardinality:   fibery.MTO,
						Name:          "Category",
						TargetName:    "Bill Account Lines",
						TargetType:    "Account",
						TargetFieldID: "id",
					},
				},
				"customer_id": {
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
				"class_id": {
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
				"tax": {
					Name:    "Tax",
					Type:    fibery.Text,
					SubType: fibery.Boolean,
				},
				"markup_percent": {
					Name: "Markup",
					Type: fibery.Number,
					Format: map[string]any{
						"format":    "Percent",
						"precision": 2,
					},
				},
				"markup_account_id": {
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
				"billable": {
					Name:    "Billable",
					Type:    fibery.Text,
					SubType: fibery.Boolean,
				},
				"billed": {
					Name:    "Billed",
					Type:    fibery.Text,
					SubType: fibery.Boolean,
				},
				"amount": {
					Name: "Amount",
					Type: fibery.Number,
					Format: map[string]any{
						"format":               "Money",
						"currencyCode":         "USD",
						"hasThousandSeperator": true,
						"precision":            2,
					},
				},
				"__syncAction": {
					Type: fibery.Text,
					Name: "Sync Action",
				},
			},
		},
		schemaGen:      func(entity, source any) (map[string]any, error) {},
		queryProcessor: func(sourceArray any, schemaGen depSchemaGenFunc) ([]map[string]any, error) {},
	},
	source:       Bill,
	sourceMapper: func(source any) (map[string]bool, error) {},
	typeMapper:   func(sourceArray any, sourceMapper sourceMapperFunc) (map[string]map[string]bool, error) {},
	cdcProcessor: func(cdc quickbooks.ChangeDataCapture, cacheEntry *IdCache, sourceMapper sourceMapperFunc, schemaGen depSchemaGenFunc) ([]map[string]any, error) {
	},
	whBatchProcessor: func(sourceArray any, cacheEntry *IdCache, sourceMapper sourceMapperFunc, schemaGen depSchemaGenFunc) ([]map[string]any, error) {
	},
}
