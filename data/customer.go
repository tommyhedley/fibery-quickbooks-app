package data

import (
	"github.com/tommyhedley/fibery-quickbooks-app/pkgs/fibery"
	"github.com/tommyhedley/quickbooks-go"
)

var Customer = QuickBooksDualType[quickbooks.Customer]{
	QuickBooksType: QuickBooksType[quickbooks.Customer]{
		BaseType: fibery.BaseType{
			TypeId:   "Customer",
			TypeName: "Customer",
			TypeSchema: map[string]fibery.Field{
				"id": {
					Name: "ID",
					Type: fibery.Id,
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
						TargetFieldID: "id",
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
						TargetFieldID: "id",
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
						TargetFieldID: "id",
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
						TargetFieldID: "id",
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
						TargetFieldID: "id",
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
						TargetFieldID: "id",
					},
				},
			},
		},
		schemaGen: func(c quickbooks.Customer) (map[string]any, error) {
			var email string
			if c.PrimaryEmailAddr != nil {
				email = c.PrimaryEmailAddr.Address
			}

			var primaryPhone string
			if c.PrimaryPhone != nil {
				primaryPhone = c.PrimaryPhone.FreeFormNumber
			}

			var alternatePhone string
			if c.AlternatePhone != nil {
				alternatePhone = c.AlternatePhone.FreeFormNumber
			}

			var mobile string
			if c.Mobile != nil {
				mobile = c.Mobile.FreeFormNumber
			}

			var fax string
			if c.Fax != nil {
				fax = c.Fax.FreeFormNumber
			}

			var website string
			if c.WebAddr != nil {
				website = c.WebAddr.URI
			}

			var shipAddr quickbooks.PhysicalAddress
			if c.ShipAddr != nil {
				shipAddr = *c.ShipAddr
			}

			var billAddr quickbooks.PhysicalAddress
			if c.BillAddr != nil {
				billAddr = *c.BillAddr
			}

			job := false
			if c.Job.Valid {
				job = c.Job.Bool
			}

			var defaultTaxCodeId string
			if c.DefaultTaxCodeRef != nil {
				defaultTaxCodeId = c.DefaultTaxCodeRef.Value
			}

			var customerTypeId string
			if c.CustomerTypeRef != nil {
				customerTypeId = c.CustomerTypeRef.Value
			}

			var salesTermId string
			if c.SalesTermRef != nil {
				salesTermId = c.SalesTermRef.Value
			}

			var paymentMethodId string
			if c.PaymentMethodRef != nil {
				paymentMethodId = c.PaymentMethodRef.Value
			}

			var parentId string
			if c.ParentRef != nil {
				parentId = c.ParentRef.Value
			}

			return map[string]any{
				"id":                 c.Id,
				"QBOId":              c.Id,
				"DisplayName":        c.DisplayName,
				"SyncToken":          c.SyncToken,
				"__syncAction":       fibery.SET,
				"Active":             c.Active,
				"Title":              c.Title,
				"GivenName":          c.GivenName,
				"MiddleName":         c.MiddleName,
				"FamilyName":         c.FamilyName,
				"Suffix":             c.Suffix,
				"CompanyName":        c.CompanyName,
				"PrimaryEmail":       email,
				"Taxable":            c.Taxable,
				"ResaleNum":          c.ResaleNum,
				"PrimaryPhone":       primaryPhone,
				"AlternatePhone":     alternatePhone,
				"Mobile":             mobile,
				"Fax":                fax,
				"Job":                job,
				"BillWithParent":     c.BillWithParent,
				"Notes":              c.Notes,
				"Website":            website,
				"Balance":            c.Balance,
				"BalanceWithJobs":    c.BalanceWithJobs,
				"ShippingLine1":      shipAddr.Line1,
				"ShippingLine2":      shipAddr.Line2,
				"ShippingLine3":      shipAddr.Line3,
				"ShippingLine4":      shipAddr.Line4,
				"ShippingLine5":      shipAddr.Line5,
				"ShippingCity":       shipAddr.City,
				"ShippingState":      shipAddr.CountrySubDivisionCode,
				"ShippingPostalCode": shipAddr.PostalCode,
				"ShippingCountry":    shipAddr.Country,
				"ShippingLat":        shipAddr.Lat,
				"ShippingLong":       shipAddr.Long,
				"BillingLine1":       billAddr.Line1,
				"BillingLine2":       billAddr.Line2,
				"BillingLine3":       billAddr.Line3,
				"BillingLine4":       billAddr.Line4,
				"BillingLine5":       billAddr.Line5,
				"BillingCity":        billAddr.City,
				"BillingState":       billAddr.CountrySubDivisionCode,
				"BillingPostalCode":  billAddr.PostalCode,
				"BillingCountry":     billAddr.Country,
				"BillingLat":         billAddr.Lat,
				"BillingLong":        billAddr.Long,
				"TaxExemptionId":     c.TaxExemptionReasonId,
				"DefaultTaxCodeId":   defaultTaxCodeId,
				"CustomerTypeId":     customerTypeId,
				"SalesTermId":        salesTermId,
				"PaymentMethodId":    paymentMethodId,
				"ParentId":           parentId,
			}, nil
		},
		pageQuery: func(req Request) ([]quickbooks.Customer, error) {
			params := quickbooks.RequestParameters{
				Ctx:     req.Ctx,
				RealmId: req.RealmId,
				Token:   req.Token,
			}

			items, err := req.Client.FindCustomersByPage(params, req.StartPosition, req.PageSize)
			if err != nil {
				return nil, err
			}

			return items, nil
		},
	},
	entityId: func(c quickbooks.Customer) string {
		return c.Id
	},
	entityStatus: func(c quickbooks.Customer) string {
		return c.Status
	},
}

func init() {
	registerType(&Customer)
}
