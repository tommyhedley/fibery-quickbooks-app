package data

import (
	"github.com/tommyhedley/fibery-quickbooks-app/pkgs/fibery"
)

var Customer = QuickBooksDualType{
	QuickBooksType: QuickBooksType{
		FiberyType: FiberyType{
			id:   "customer",
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
				"primary_email": {
					Name:    "Email",
					Type:    fibery.Text,
					SubType: fibery.Email,
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
			},
		},
	},
}
