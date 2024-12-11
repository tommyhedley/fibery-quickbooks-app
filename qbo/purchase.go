package qbo

type Purchase struct {
}

type AccountPurchaseLine struct {
}

type ItemPurchaseLine struct {
}

func (Purchase) TypeInfo() TypeArray {
	return TypeArray{
		ID:   "purchase",
		Name: "Expense",
	}
}

func (Purchase) Schema() map[string]Field {
	return map[string]Field{
		"id": {
			Name: "id",
			Type: ID,
		},
		"name": {
			Name: "Name",
			Type: Text,
		},
		"payment_type": {
			Name:     "Payment Type",
			Type:     Text,
			SubType:  SingleSelect,
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
			Type: Text,
			Relation: &Relation{
				Cardinality:   MTO,
				Name:          "Payment Method",
				TargetName:    "Expenses",
				TargetType:    "payment_method",
				TargetFieldID: "id",
			},
		},
		"account_id": {
			Name: "Account ID",
			Type: Text,
			Relation: &Relation{
				Cardinality:   MTO,
				Name:          "Account",
				TargetName:    "Expenses",
				TargetType:    "account",
				TargetFieldID: "id",
			},
		},
		"sync_token": {
			Name:     "Sync Token",
			Type:     Text,
			Ignore:   true,
			ReadOnly: true,
		},
		"date": {
			Name:        "Date",
			Description: "The date that the transaction occured",
			Type:        DateType,
			SubType:     Day,
		},
		"reference_number": {
			Name: "Reference Number",
			Type: Text,
		},
		"private_memo": {
			Name:        "Memo",
			Type:        Text,
			Description: "The private memo line on the QB expense form",
		},
		"credit": {
			Name:        "Credit",
			Type:        Text,
			SubType:     Boolean,
			Description: "Valid only for credit card charges",
		},
		"entity_id": {
			Name: "Entity ID",
			Type: Text,
			Relation: &Relation{
				Cardinality:   MTO,
				Name:          "Payee",
				TargetName:    "Expenses",
				TargetType:    "entity",
				TargetFieldID: "id",
			},
		},
		"created_qbo": {
			Name: "Created (QBO)",
			Type: DateType,
		},
		"last_updated_qbo": {
			Name: "Last Updated (QBO)",
			Type: DateType,
		},
		"total": {
			Name:        "Total (QBO)",
			Description: "Calculated total from QBO",
			Type:        Number,
			Format: map[string]any{
				"format":               "Money",
				"currencyCode":         "USD",
				"hasThousandSeperator": true,
				"precision":            2,
			},
		},
	}
}

func (P Purchase) TransformData(params ...any) (any, error) {
	return nil, nil
}

func (Purchase) FullSync(*FullSyncRequest) ([]map[string]any, bool, error) {
	return nil, false, nil
}

func (Purchase) DeltaSync(*DeltaSyncRequest) ([]map[string]any, error) {
	return nil, nil
}

func (Purchase) Webhook(*WebhookRequest) ([]map[string]any, error) {
	return nil, nil
}

func (AccountPurchaseLine) TypeInfo() TypeArray {
	return TypeArray{
		ID:   "purchase_account_line",
		Name: "Purchase Account Line",
	}
}

func (AccountPurchaseLine) Schema() map[string]Field {
	return map[string]Field{
		"id": {
			Name: "id",
			Type: ID,
		},
		"name": {
			Name: "Name",
			Type: Text,
		},
		"purchase_id": {
			Name: "Purchase ID",
			Type: Text,
			Relation: &Relation{
				Cardinality:   MTO,
				Name:          "Expense",
				TargetName:    "Account Line(s)",
				TargetType:    "purchase",
				TargetFieldID: "id",
			},
		},
		"amount": {
			Name: "Amount",
			Type: Number,
			Format: map[string]any{
				"format":               "Money",
				"currencyCode":         "USD",
				"hasThousandSeperator": true,
				"precision":            2,
			},
		},
		"description": {
			Name: "Description",
			Type: Text,
		},
		"line": {
			Name:    "Line Number",
			Type:    Number,
			SubType: Integer,
		},
		"account_id": {
			Name: "Account ID",
			Type: Text,
			Relation: &Relation{
				Cardinality:   MTO,
				Name:          "Account or Item",
				TargetName:    "Expense Line(s)",
				TargetType:    "account",
				TargetFieldID: "id",
			},
		},
		"customer_id": {
			Name: "Customer ID",
			Type: Text,
			Relation: &Relation{
				Cardinality:   MTO,
				Name:          "Customer",
				TargetName:    "Expense Account Line(s)",
				TargetType:    "customer",
				TargetFieldID: "id",
			},
		},
		"class_id": {
			Name: "Class ID",
			Type: Text,
			Relation: &Relation{
				Cardinality:   MTO,
				Name:          "Class",
				TargetName:    "Expense Account Line(s)",
				TargetType:    "class",
				TargetFieldID: "id",
			},
		},
		"tax_code_id": {
			Name: "Tax Code ID",
			Type: Text,
			Relation: &Relation{
				Cardinality:   MTO,
				Name:          "Tax Code",
				TargetName:    "Expense Account Line(s)",
				TargetType:    "tax_code",
				TargetFieldID: "id",
			},
		},
		"billable_status": {
			Name:     "Billlable Status",
			Type:     Text,
			SubType:  SingleSelect,
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
			Type: Number,
			Format: map[string]any{
				"format":    "Percent",
				"precision": 2,
			},
		},
	}
}

func (Apl AccountPurchaseLine) TransformData(params ...any) (any, error) {
	return nil, nil
}

func (AccountPurchaseLine) FullSync(*FullSyncRequest) ([]map[string]any, bool, error) {
	return nil, false, nil
}

func (AccountPurchaseLine) DeltaSync(*DeltaSyncRequest) ([]map[string]any, error) {
	return nil, nil
}

func (AccountPurchaseLine) Webhook(*WebhookRequest) ([]map[string]any, error) {
	return nil, nil
}

func (ItemPurchaseLine) TypeInfo() TypeArray {
	return TypeArray{
		ID:   "purchase_item_line",
		Name: "Purchase Item Line",
	}
}

func (ItemPurchaseLine) Schema() map[string]Field {
	return map[string]Field{}
}

func (Ipl ItemPurchaseLine) TransformData(params ...any) (any, error) {
	return nil, nil
}

func (ItemPurchaseLine) FullSync(*FullSyncRequest) ([]map[string]any, bool, error) {
	return nil, false, nil
}

func (ItemPurchaseLine) DeltaSync(*DeltaSyncRequest) ([]map[string]any, error) {
	return nil, nil
}

func (ItemPurchaseLine) Webhook(*WebhookRequest) ([]map[string]any, error) {
	return nil, nil
}

func init() {
	RegisterType(Purchase{})
	RegisterType(AccountPurchaseLine{})
	RegisterType(ItemPurchaseLine{})
}
