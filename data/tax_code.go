package data

import (
	"github.com/tommyhedley/fibery-quickbooks-app/pkgs/fibery"
)

var TaxCode = QuickBooksType{
	fiberyType: fiberyType{
		id:   "TaxCode",
		name: "Tax Code",
		schema: map[string]fibery.Field{
			"id": {
				Name: "id",
				Type: fibery.Text,
			},
			"name": {
				Name: "Name",
				Type: fibery.Text,
			},
			"description": {
				Name: "Description",
				Type: fibery.Text,
			},
			"sync_token": {
				Name: "Sync Token",
				Type: fibery.Text,
			},
			"tax_group": {
				Name:    "Tax Group",
				Type:    fibery.Text,
				SubType: fibery.Boolean,
			},
			"taxable": {
				Name:    "Taxable",
				Type:    fibery.Text,
				SubType: fibery.Boolean,
			},
			"active": {
				Name:    "Active",
				Type:    fibery.Text,
				SubType: fibery.Boolean,
			},
			"hidden": {
				Name:    "Hidden",
				Type:    fibery.Text,
				SubType: fibery.Boolean,
			},
			"tax_code_type": {
				Name:     "Type",
				Type:     fibery.Text,
				SubType:  fibery.SingleSelect,
				ReadOnly: true,
				Options: []map[string]any{
					{
						"name": "User Defined",
					},
					{
						"name": "System Generated",
					},
				},
			},
		},
	},
	schemaGen:      func(entity any) (map[string]any, error) {},
	query:          func(req Request) (Response, error) {},
	queryProcessor: func(entityArray any, schemaGen schemaGenFunc) ([]map[string]any, error) {},
}

var TaxCodeLine = DependentDataType{
	dependentBaseType: dependentBaseType{
		fiberyType: fiberyType{
			id:   "TaxCodeLine",
			name: "Tax Code Line",
			schema: map[string]fibery.Field{
				"id": {
					Name: "id",
					Type: fibery.Text,
				},
				"name": {
					Name: "Name",
					Type: fibery.Text,
				},
				"tax_rate_id": {
					Name: "Tax Rate ID",
					Type: fibery.Text,
					Relation: &fibery.Relation{
						Cardinality:   fibery.MTO,
						Name:          "Tax Rate",
						TargetName:    "Tax Code Lines",
						TargetType:    "TaxRate",
						TargetFieldID: "id",
					},
				},
				"tax_type": {
					Name:     "Tax Type",
					Type:     fibery.Text,
					SubType:  fibery.SingleSelect,
					ReadOnly: true,
					Options: []map[string]any{
						{
							"name": "Tax On Amount",
						},
						{
							"name": "Tax On Amount Plus Tax",
						},
						{
							"name": "Tax On Tax",
						},
					},
				},
				"tax_order": {
					Name: "Tax Order",
					Type: fibery.Number,
				},
				"tax_code_id_purchase": {
					Name: "Purchase Tax Code ID",
					Type: fibery.Text,
					Relation: &fibery.Relation{
						Cardinality:   fibery.MTO,
						Name:          "Purchase Tax Code",
						TargetName:    "Purchase Tax Rates",
						TargetType:    "TaxCode",
						TargetFieldID: "id",
					},
				},
				"tax_code_id_sale": {
					Name: "Sale Tax Code ID",
					Type: fibery.Text,
					Relation: &fibery.Relation{
						Cardinality:   fibery.MTO,
						Name:          "Sales Tax Code",
						TargetName:    "Sales Tax Rates",
						TargetType:    "TaxCode",
						TargetFieldID: "id",
					},
				},
			},
		},
		schemaGen:      func(entity, source any) (map[string]any, error) {},
		queryProcessor: func(sourceArray any, schemaGen depSchemaGenFunc) ([]map[string]any, error) {},
	},
	source: TaxCode,
}
