package qbo

import "github.com/tommyhedley/fibery/fibery-qbo-integration/sync"

type Purchase struct {
}

type AccountPurchaseLine struct {
}

type ItemPurchaseLine struct {
}

var PurchaseType = sync.DataType{
	ID:   "purchase",
	Name: "Expense",
	Schema: map[string]sync.Field{
		"id": {
			Name: "id",
			Type: sync.ID,
		},
		"name": {
			Name: "Name",
			Type: sync.Text,
		},
		"payment_type": {
			Name:     "Payment Type",
			Type:     sync.Text,
			SubType:  sync.SingleSelect,
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
		"payment_method": {
			Name: "Payment Method ID",
			Type: sync.Text,
			Relation: &sync.Relation{
				Cardinality:   sync.MTO,
				Name:          "Payment Method",
				TargetName:    "Expenses",
				TargetType:    "payment_method",
				TargetFieldID: "id",
			},
		},
		"account_id": {
			Name: "Account ID",
			Type: sync.Text,
			Relation: &sync.Relation{
				Cardinality:   sync.MTO,
				Name:          "Account",
				TargetName:    "Expenses",
				TargetType:    "account",
				TargetFieldID: "id",
			},
		},
		"sync_token": {
			Name:     "Sync Token",
			Type:     sync.Text,
			Ignore:   true,
			ReadOnly: true,
		},
		"date": {
			Name:        "Date",
			Description: "The date that the transaction occured",
			Type:        sync.Date,
			SubType:     sync.Day,
		},
		"reference_number": {
			Name: "Reference Number",
			Type: sync.Text,
		},
		"private_memo": {
			Name:        "Memo",
			Type:        sync.Text,
			Description: "The private memo line on the QB expense form",
		},
		"credit": {
			Name:        "Credit",
			Type:        sync.Text,
			SubType:     sync.Boolean,
			Description: "Valid only for credit card charges",
		},
		"entity_id": {
			Name: "Entity ID",
			Type: sync.Text,
			Relation: &sync.Relation{
				Cardinality:   sync.MTO,
				Name:          "Payee",
				TargetName:    "Expenses",
				TargetType:    "entity",
				TargetFieldID: "id",
			},
		},
		"created_qbo": {
			Name: "Created (QBO)",
			Type: sync.Date,
		},
		"last_updated_qbo": {
			Name: "Last Updated (QBO)",
			Type: sync.Date,
		},
		"total": {
			Name:        "Total (QBO)",
			Description: "Calculated total from QBO",
			Type:        sync.Number,
			Format: map[string]any{
				"format":               "Money",
				"currencyCode":         "USD",
				"hasThousandSeperator": true,
				"precision":            2,
			},
		},
	},
	DataRequest: func(params sync.RequestParameters) ([]map[string]any, bool, error) {
		return nil, false, nil
	},
}

var PurchaseAccountLineType = sync.DataType{
	ID:   "purchase_account_line",
	Name: "Purchase Account Line",
	Schema: map[string]sync.Field{
		"id": {
			Name: "id",
			Type: sync.ID,
		},
		"name": {
			Name: "Name",
			Type: sync.Text,
		},
		"purchase_id": {
			Name: "Purchase ID",
			Type: sync.Text,
			Relation: &sync.Relation{
				Cardinality:   sync.MTO,
				Name:          "Expense",
				TargetName:    "Account Line(s)",
				TargetType:    "purchase",
				TargetFieldID: "id",
			},
		},
		"amount": {
			Name: "Amount",
			Type: sync.Number,
			Format: map[string]any{
				"format":               "Money",
				"currencyCode":         "USD",
				"hasThousandSeperator": true,
				"precision":            2,
			},
		},
		"description": {
			Name: "Description",
			Type: sync.Text,
		},
		"line": {
			Name:    "Line Number",
			Type:    sync.Number,
			SubType: sync.Integer,
		},
		"account_id": {
			Name: "Account ID",
			Type: sync.Text,
			Relation: &sync.Relation{
				Cardinality:   sync.MTO,
				Name:          "Account or Item",
				TargetName:    "Expense Line(s)",
				TargetType:    "account",
				TargetFieldID: "id",
			},
		},
		"customer_id": {
			Name: "Customer ID",
			Type: sync.Text,
			Relation: &sync.Relation{
				Cardinality:   sync.MTO,
				Name:          "Customer",
				TargetName:    "Expense Account Line(s)",
				TargetType:    "customer",
				TargetFieldID: "id",
			},
		},
		"class_id": {
			Name: "Class ID",
			Type: sync.Text,
			Relation: &sync.Relation{
				Cardinality:   sync.MTO,
				Name:          "Class",
				TargetName:    "Expense Account Line(s)",
				TargetType:    "class",
				TargetFieldID: "id",
			},
		},
		"tax_code_id": {
			Name: "Tax Code ID",
			Type: sync.Text,
			Relation: &sync.Relation{
				Cardinality:   sync.MTO,
				Name:          "Tax Code",
				TargetName:    "Expense Account Line(s)",
				TargetType:    "tax_code",
				TargetFieldID: "id",
			},
		},
		"billable_status": {
			Name:     "Billlable Status",
			Type:     sync.Text,
			SubType:  sync.SingleSelect,
			ReadOnly: true,
			Options: []map[string]any{
				{
					"name": "Billable",
				},
				{
					"name": "Not Billable",
				},
				{
					"name": "Billed",
				},
			},
		},
		"markup_info": {
			Name: "Markup",
			Type: sync.Number,
			Format: map[string]any{
				"format":    "Percent",
				"precision": 2,
			},
		},
	},
}

var PurchaseItemLineType = sync.DataType{
	ID:   "purchase_item_line",
	Name: "Purchase Item Line",
}

func init() {
	PurchaseType.Register()
	PurchaseAccountLineType.Register()
	PurchaseItemLineType.Register()
}
