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
				"Id": {
					Name: "ID",
					Type: fibery.ID,
				},
				"QBOId": {
					Name: "QBO ID",
					Type: fibery.Text,
				},
				"DisplayName": {
					Name:    "Display Name",
					Type:    fibery.Text,
					SubType: fibery.Title,
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
				"Active": {
					Name:    "Active",
					Type:    fibery.Text,
					SubType: fibery.Boolean,
				},
				"Title": {
					Name: "Title",
					Type: fibery.Text,
				},
				"GivenName": {
					Name: "First Name",
					Type: fibery.Text,
				},
				"MiddleName": {
					Name: "Middle Name",
					Type: fibery.Text,
				},
				"FamilyName": {
					Name: "Last Name",
					Type: fibery.Text,
				},
				"Suffix": {
					Name: "Suffix",
					Type: fibery.Text,
				},
				"CompanyName": {
					Name: "Company Name",
					Type: fibery.Text,
				},
				"PrimaryEmail": {
					Name:    "Email",
					Type:    fibery.Text,
					SubType: fibery.Email,
				},
				"Taxable": {
					Name:    "Taxable",
					Type:    fibery.Text,
					SubType: fibery.Boolean,
				},
				"ResaleNum": {
					Name: "Resale ID",
					Type: fibery.Text,
				},
				"PrimaryPhone": {
					Name: "Phone",
					Type: fibery.Text,
					Format: map[string]any{
						"format": "phone",
					},
				},
				"AlternatePhone": {
					Name: "Alternate Phone",
					Type: fibery.Text,
					Format: map[string]any{
						"format": "phone",
					},
				},
				"Mobile": {
					Name: "Mobile",
					Type: fibery.Text,
					Format: map[string]any{
						"format": "phone",
					},
				},
				"Fax": {
					Name: "Fax",
					Type: fibery.Text,
					Format: map[string]any{
						"format": "phone",
					},
				},
				"Job": {
					Name:    "Job",
					Type:    fibery.Text,
					SubType: fibery.Boolean,
				},
				"BillWithParent": {
					Name:    "Bill With Parent",
					Type:    fibery.Text,
					SubType: fibery.Boolean,
				},
				"Notes": {
					Name:    "Notes",
					Type:    fibery.Text,
					SubType: fibery.MD,
				},
				"Website": {
					Name:    "Website",
					Type:    fibery.Text,
					SubType: fibery.URL,
				},
				"Balance": {
					Name: "Balance",
					Type: fibery.Number,
					Format: map[string]any{
						"format":               "Money",
						"currencyCode":         "USD",
						"hasThousandSeperator": true,
						"precision":            2,
					},
				},
				"BalanceWithJobs": {
					Name: "Balance With Jobs",
					Type: fibery.Number,
					Format: map[string]any{
						"format":               "Money",
						"currencyCode":         "USD",
						"hasThousandSeperator": true,
						"precision":            2,
					},
				},
				"ShippingLine1": {
					Name: "Shipping Line 1",
					Type: fibery.Text,
				},
				"ShippingLine2": {
					Name: "Shipping Line 2",
					Type: fibery.Text,
				},
				"ShippingLine3": {
					Name: "Shipping Line 3",
					Type: fibery.Text,
				},
				"ShippingLine4": {
					Name: "Shipping Line 4",
					Type: fibery.Text,
				},
				"ShippingLine5": {
					Name: "Shipping Line 5",
					Type: fibery.Text,
				},
				"ShippingCity": {
					Name: "Shipping City",
					Type: fibery.Text,
				},
				"ShippingState": {
					Name: "Shipping State",
					Type: fibery.Text,
				},
				"ShippingPostalCode": {
					Name: "Shipping Postal Code",
					Type: fibery.Text,
				},
				"ShippingCountry": {
					Name: "Shipping Country",
					Type: fibery.Text,
				},
				"ShippingLat": {
					Name: "Shipping Latitude",
					Type: fibery.Text,
				},
				"ShippingLong": {
					Name: "Shipping Longitude",
					Type: fibery.Text,
				},
				"BillingLine1": {
					Name: "Billing Line 1",
					Type: fibery.Text,
				},
				"BillingLine2": {
					Name: "Billing Line 2",
					Type: fibery.Text,
				},
				"BillingLine3": {
					Name: "Billing Line 3",
					Type: fibery.Text,
				},
				"BillingLine4": {
					Name: "Billing Line 4",
					Type: fibery.Text,
				},
				"BillingLine5": {
					Name: "Billing Line 5",
					Type: fibery.Text,
				},
				"BillingCity": {
					Name: "Billing City",
					Type: fibery.Text,
				},
				"BillingState": {
					Name: "Billing State",
					Type: fibery.Text,
				},
				"BillingPostalCode": {
					Name: "Billing Postal Code",
					Type: fibery.Text,
				},
				"BillingCountry": {
					Name: "Billing Country",
					Type: fibery.Text,
				},
				"BillingLat": {
					Name: "Billing Latitude",
					Type: fibery.Text,
				},
				"BillingLong": {
					Name: "Billing Longitude",
					Type: fibery.Text,
				},
				"TaxExemptionId": {
					Name: "Tax Exemption ID",
					Type: fibery.Text,
					Relation: &fibery.Relation{
						Cardinality:   fibery.MTO,
						Name:          "Tax Exemption",
						TargetName:    "Customers",
						TargetType:    "TaxExemption",
						TargetFieldID: "Id",
					},
				},
				"DefaultTaxCodeId": {
					Name: "Default Tax Code ID",
					Type: fibery.Text,
					Relation: &fibery.Relation{
						Cardinality:   fibery.MTO,
						Name:          "Default Tax Code",
						TargetName:    "Customers",
						TargetType:    "TaxCode",
						TargetFieldID: "Id",
					},
				},
				"CustomerTypeId": {
					Name: "Customer Type ID",
					Type: fibery.Text,
					Relation: &fibery.Relation{
						Cardinality:   fibery.MTO,
						Name:          "Customer Type",
						TargetName:    "Customers",
						TargetType:    "CustomerType",
						TargetFieldID: "Id",
					},
				},
				"SalesTermId": {
					Name: "Sales Term ID",
					Type: fibery.Text,
					Relation: &fibery.Relation{
						Cardinality:   fibery.MTO,
						Name:          "Sales Term",
						TargetName:    "Customers",
						TargetType:    "SalesTerm",
						TargetFieldID: "Id",
					},
				},
				"PaymentMethodId": {
					Name: "Payment Method ID",
					Type: fibery.Text,
					Relation: &fibery.Relation{
						Cardinality:   fibery.MTO,
						Name:          "Payment Method",
						TargetName:    "Customers",
						TargetType:    "PaymentMethod",
						TargetFieldID: "Id",
					},
				},
				"ParentId": {
					Name: "Parent ID",
					Type: fibery.Text,
					Relation: &fibery.Relation{
						Cardinality:   fibery.MTO,
						Name:          "Parent",
						TargetName:    "Jobs",
						TargetType:    "Customer",
						TargetFieldID: "Id",
					},
				},
			},
		},
		schemaGen: func(entity any) (map[string]any, error) {
			customer, ok := entity.(quickbooks.Customer)
			if !ok {
				return nil, fmt.Errorf("unable to convert entity to Customer")
			}

			job := false
			if customer.Job.Valid {
				job = customer.Job.Bool
			}

			data := map[string]any{
				"Id":                 customer.Id,
				"QBOId":              customer.Id,
				"DisplayName":        customer.DisplayName,
				"SyncToken":          customer.SyncToken,
				"__syncAction":       fibery.SET,
				"Active":             customer.Active,
				"Title":              customer.Title,
				"GivenName":          customer.GivenName,
				"MiddleName":         customer.MiddleName,
				"FamilyName":         customer.FamilyName,
				"Suffix":             customer.Suffix,
				"CompanyName":        customer.CompanyName,
				"PrimaryEmail":       customer.PrimaryEmailAddr,
				"Taxable":            customer.Taxable,
				"ResaleNum":          customer.ResaleNum,
				"PrimaryPhone":       customer.PrimaryPhone.FreeFormNumber,
				"AlternatePhone":     customer.AlternatePhone.FreeFormNumber,
				"Mobile":             customer.Mobile.FreeFormNumber,
				"Fax":                customer.Fax.FreeFormNumber,
				"Job":                job,
				"BillWithParent":     customer.BillWithParent,
				"Notes":              customer.Notes,
				"Website":            customer.WebAddr.URI,
				"Balance":            customer.Balance,
				"BalanceWithJobs":    customer.BalanceWithJobs,
				"ShippingLine1":      customer.ShipAddr.Line1,
				"ShippingLine2":      customer.ShipAddr.Line2,
				"ShippingLine3":      customer.ShipAddr.Line3,
				"ShippingLine4":      customer.ShipAddr.Line4,
				"ShippingLine5":      customer.ShipAddr.Line5,
				"ShippingCity":       customer.ShipAddr.City,
				"ShippingState":      customer.ShipAddr.CountrySubDivisionCode,
				"ShippingPostalCode": customer.ShipAddr.PostalCode,
				"ShippingCountry":    customer.ShipAddr.Country,
				"ShippingLat":        customer.ShipAddr.Lat,
				"ShippingLong":       customer.ShipAddr.Long,
				"BillingLine1":       customer.BillAddr.Line1,
				"BillingLine2":       customer.BillAddr.Line2,
				"BillingLine3":       customer.BillAddr.Line3,
				"BillingLine4":       customer.BillAddr.Line4,
				"BillingLine5":       customer.BillAddr.Line5,
				"BillingCity":        customer.BillAddr.City,
				"BillingState":       customer.BillAddr.CountrySubDivisionCode,
				"BillingPostalCode":  customer.BillAddr.PostalCode,
				"BillingCountry":     customer.BillAddr.Country,
				"BillingLat":         customer.BillAddr.Lat,
				"BillingLong":        customer.BillAddr.Long,
				"TaxExemptionId":     customer.TaxExemptionReasonId,
				"DefaultTaxCodeId":   customer.DefaultTaxCodeRef.Value,
				"CustomerTypeId":     customer.CustomerTypeRef.Value,
				"SalesTermId":        customer.SalesTermRef.Value,
				"PaymentMethodId":    customer.PaymentMethodRef.Value,
				"ParentId":           customer.ParentRef.Value,
			}
			return data, nil
		},
		query:          func(req Request) (Response, error) {},
		queryProcessor: func(entityArray any, schemaGen schemaGenFunc) ([]map[string]any, error) {},
	},
	cdcProcessor: func(cdc quickbooks.ChangeDataCapture, schemaGen schemaGenFunc) ([]map[string]any, error) {},
	whBatchProcessor: func(itemResponse quickbooks.BatchItemResponse, response *map[string][]map[string]any, cache *cache.Cache, realmId string, queryProcessor queryProcessorFunc, schemaGen schemaGenFunc, typeId string) error {
	},
}
