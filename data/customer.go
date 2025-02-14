package data

import (
	"github.com/tommyhedley/fibery-quickbooks-app/pkgs/fibery"
	"github.com/tommyhedley/quickbooks-go"
)

var Customer = QuickBooksDualType{
	QuickBooksType: QuickBooksType{
		fiberyType: fiberyType{
			id:   "Customer",
			name: "Customer",
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
				"active": {
					Name:    "Active",
					Type:    fibery.Text,
					SubType: fibery.Boolean,
				},
				"display_name": {
					Name:    "Display Name",
					Type:    fibery.Text,
					SubType: fibery.Title,
				},
				"title": {
					Name: "Title",
					Type: fibery.Text,
				},
				"given_name": {
					Name: "First Name",
					Type: fibery.Text,
				},
				"middle_name": {
					Name: "Middle Name",
					Type: fibery.Text,
				},
				"family_name": {
					Name: "Last Name",
					Type: fibery.Text,
				},
				"suffix": {
					Name: "Suffix",
					Type: fibery.Text,
				},
				"company_name": {
					Name: "Company Name",
					Type: fibery.Text,
				},
				"primary_email": {
					Name:    "Email",
					Type:    fibery.Text,
					SubType: fibery.Email,
				},
				"taxable": {
					Name:    "Taxable",
					Type:    fibery.Text,
					SubType: fibery.Boolean,
				},
				"tax_exemption_id": {
					Name: "Tax Exemption ID",
					Type: fibery.Text,
					Relation: &fibery.Relation{
						Cardinality:   fibery.MTO,
						Name:          "Tax Exemption",
						TargetName:    "Customers",
						TargetType:    "TaxExemption",
						TargetFieldID: "id",
					},
				},
				"resale_num": {
					Name: "Resale ID",
					Type: fibery.Text,
				},
				"default_tax_code_id": {
					Name: "Default Tax Code ID",
					Type: fibery.Text,
					Relation: &fibery.Relation{
						Cardinality:   fibery.MTO,
						Name:          "Default Tax Code",
						TargetName:    "Customers",
						TargetType:    "TaxCode",
						TargetFieldID: "id",
					},
				},
				"customer_type_id": {
					Name: "Customer Type ID",
					Type: fibery.Text,
					Relation: &fibery.Relation{
						Cardinality:   fibery.MTO,
						Name:          "Customer Type",
						TargetName:    "Customers",
						TargetType:    "CustomerType",
						TargetFieldID: "id",
					},
				},
				"delivery_method": {
					Name:     "Delivery Method",
					Type:     fibery.Text,
					SubType:  fibery.SingleSelect,
					ReadOnly: true,
					Options: []map[string]any{
						{
							"name": "Print",
						},
						{
							"name": "Email",
						},
					},
				},
				"sales_term_id": {
					Name: "Sales Term ID",
					Type: fibery.Text,
					Relation: &fibery.Relation{
						Cardinality:   fibery.MTO,
						Name:          "Sales Term",
						TargetName:    "Customers",
						TargetType:    "SalesTerm",
						TargetFieldID: "id",
					},
				},
				"phone": {
					Name: "Phone",
					Type: fibery.Text,
					Format: map[string]any{
						"format": "phone",
					},
				},
				"alt_phone": {
					Name: "Alternate Phone",
					Type: fibery.Text,
					Format: map[string]any{
						"format": "phone",
					},
				},
				"mobile": {
					Name: "Mobile",
					Type: fibery.Text,
					Format: map[string]any{
						"format": "phone",
					},
				},
				"fax": {
					Name: "Fax",
					Type: fibery.Text,
					Format: map[string]any{
						"format": "phone",
					},
				},
				"parent_id": {
					Name: "Parent ID",
					Type: fibery.Text,
					Relation: &fibery.Relation{
						Cardinality:   fibery.MTO,
						Name:          "Parent",
						TargetName:    "Jobs",
						TargetType:    "Customer",
						TargetFieldID: "id",
					},
				},
				"job": {
					Name:    "Job",
					Type:    fibery.Text,
					SubType: fibery.Boolean,
				},
				"bill_with_parent": {
					Name:    "Bill With Parent",
					Type:    fibery.Text,
					SubType: fibery.Boolean,
				},
				"notes": {
					Name:    "Notes",
					Type:    fibery.Text,
					SubType: fibery.MD,
				},
				"website": {
					Name:    "Website",
					Type:    fibery.Text,
					SubType: fibery.URL,
				},
				"balance": {
					Name: "Balance",
					Type: fibery.Number,
					Format: map[string]any{
						"format":               "Money",
						"currencyCode":         "USD",
						"hasThousandSeperator": true,
						"precision":            2,
					},
				},
				"balance_with_jobs": {
					Name: "Balance With Jobs",
					Type: fibery.Number,
					Format: map[string]any{
						"format":               "Money",
						"currencyCode":         "USD",
						"hasThousandSeperator": true,
						"precision":            2,
					},
				},
				"payment_method_id": {
					Name: "Payment Method ID",
					Type: fibery.Text,
					Relation: &fibery.Relation{
						Cardinality:   fibery.MTO,
						Name:          "Payment Method",
						TargetName:    "Customers",
						TargetType:    "PaymentMethod",
						TargetFieldID: "id",
					},
				},
				"shipping_address": {
					Name:    "Shipping Address",
					Type:    fibery.Text,
					SubType: fibery.MD,
				},
				"shipping_line_1": {
					Name: "Shipping Line 1",
					Type: fibery.Text,
				},
				"shipping_line_2": {
					Name: "Shipping Line 2",
					Type: fibery.Text,
				},
				"shipping_line_3": {
					Name: "Shipping Line 3",
					Type: fibery.Text,
				},
				"shipping_line_4": {
					Name: "Shipping Line 4",
					Type: fibery.Text,
				},
				"shipping_line_5": {
					Name: "Shipping Line 5",
					Type: fibery.Text,
				},
				"shipping_city": {
					Name: "Shipping City",
					Type: fibery.Text,
				},
				"shipping_state": {
					Name: "Shipping State",
					Type: fibery.Text,
				},
				"shipping_postal_code": {
					Name: "Shipping Postal Code",
					Type: fibery.Text,
				},
				"shipping_country": {
					Name: "Shipping Country",
					Type: fibery.Text,
				},
				"shipping_lat": {
					Name: "Shipping Latitude",
					Type: fibery.Text,
				},
				"shipping_long": {
					Name: "Shipping Longitude",
					Type: fibery.Text,
				},
				"billing_address": {
					Name:    "Billing Address",
					Type:    fibery.Text,
					SubType: fibery.MD,
				},
				"billing_line_1": {
					Name: "Billing Line 1",
					Type: fibery.Text,
				},
				"billing_line_2": {
					Name: "Billing Line 2",
					Type: fibery.Text,
				},
				"billing_line_3": {
					Name: "Billing Line 3",
					Type: fibery.Text,
				},
				"billing_line_4": {
					Name: "Billing Line 4",
					Type: fibery.Text,
				},
				"billing_line_5": {
					Name: "Billing Line 5",
					Type: fibery.Text,
				},
				"billing_city": {
					Name: "Billing City",
					Type: fibery.Text,
				},
				"billing_state": {
					Name: "Billing State",
					Type: fibery.Text,
				},
				"billing_postal_code": {
					Name: "Billing Postal Code",
					Type: fibery.Text,
				},
				"billing_country": {
					Name: "Billing Country",
					Type: fibery.Text,
				},
				"billing_lat": {
					Name: "Billing Latitude",
					Type: fibery.Text,
				},
				"billing_long": {
					Name: "Billing Longitude",
					Type: fibery.Text,
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
