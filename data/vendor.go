package data

import (
	"github.com/tommyhedley/fibery-quickbooks-app/pkgs/fibery"
	"github.com/tommyhedley/quickbooks-go"
)

var Vendor = QuickBooksDualType[quickbooks.Vendor]{
	QuickBooksType: QuickBooksType[quickbooks.Vendor]{
		BaseType: fibery.BaseType{
			TypeId:   "Vendor",
			TypeName: "Vendor",
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
				"1099": {
					Name:        "1099",
					Type:        fibery.Text,
					SubType:     fibery.Boolean,
					Description: "Is the Vendor a 1099 contractor?",
				},
				"CostRate": {
					Name:        "Cost Rate",
					Type:        fibery.Number,
					Description: "Default cost rate of the Vendor",
					Format: map[string]any{
						"format":               "Money",
						"currencyCode":         "USD",
						"hasThousandSeperator": true,
						"precision":            2,
					},
				},
				"BillRate": {
					Name:        "Bill Rate",
					Type:        fibery.Number,
					Description: "Default billing rate of the Vendor",
					Format: map[string]any{
						"format":               "Money",
						"currencyCode":         "USD",
						"hasThousandSeperator": true,
						"precision":            2,
					},
				},
				"Website": {
					Name:    "Website",
					Type:    fibery.Text,
					SubType: fibery.URL,
				},
				"AccountNumber": {
					Name:        "Account Number",
					Type:        fibery.Text,
					Description: "Name or number of the account associated with this vendor",
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
				"BillingAddress": {
					Name:    "Billing Address",
					Type:    fibery.Text,
					SubType: fibery.MD,
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
			},
		},
		schemaGen: func(v quickbooks.Vendor) (map[string]any, error) {
			var email string
			if v.PrimaryEmailAddr != nil {
				email = v.PrimaryEmailAddr.Address
			}

			var primaryPhone string
			if v.PrimaryPhone != nil {
				primaryPhone = v.PrimaryPhone.FreeFormNumber
			}

			var alternatePhone string
			if v.AlternatePhone != nil {
				alternatePhone = v.AlternatePhone.FreeFormNumber
			}

			var mobile string
			if v.Mobile != nil {
				mobile = v.Mobile.FreeFormNumber
			}

			var fax string
			if v.Fax != nil {
				fax = v.Fax.FreeFormNumber
			}

			var website string
			if v.WebAddr != nil {
				website = v.WebAddr.URI
			}

			var billAddr quickbooks.PhysicalAddress
			if v.BillAddr != nil {
				billAddr = *v.BillAddr
			}

			var termId string
			if v.TermRef != nil {
				termId = v.TermRef.Value
			}

			return map[string]any{
				"id":                v.Id,
				"QBOId":             v.Id,
				"DisplayName":       v.DisplayName,
				"SyncToken":         v.SyncToken,
				"__syncAction":      fibery.SET,
				"Active":            v.Active,
				"Title":             v.Title,
				"GivenName":         v.GivenName,
				"MiddleName":        v.MiddleName,
				"FamilyName":        v.FamilyName,
				"Suffix":            v.Suffix,
				"CompanyName":       v.CompanyName,
				"PrimaryEmail":      email,
				"SalesTermId":       termId,
				"PrimaryPhone":      primaryPhone,
				"AlternatePhone":    alternatePhone,
				"Mobile":            mobile,
				"Fax":               fax,
				"1099":              v.Vendor1099,
				"CostRate":          v.CostRate,
				"BillRate":          v.BillRate,
				"Website":           website,
				"Balance":           v.Balance,
				"BillingLine1":      billAddr.Line1,
				"BillingLine2":      billAddr.Line2,
				"BillingLine3":      billAddr.Line3,
				"BillingLine4":      billAddr.Line4,
				"BillingLine5":      billAddr.Line5,
				"BillingCity":       billAddr.City,
				"BillingState":      billAddr.CountrySubDivisionCode,
				"BillingPostalCode": billAddr.PostalCode,
				"BillingCountry":    billAddr.Country,
				"BillingLat":        billAddr.Lat,
				"BillingLong":       billAddr.Long,
			}, nil
		},
		pageQuery: func(req Request) ([]quickbooks.Vendor, error) {
			params := quickbooks.RequestParameters{
				Ctx:     req.Ctx,
				RealmId: req.RealmId,
				Token:   req.Token,
			}

			items, err := req.Client.FindVendorsByPage(params, req.StartPosition, req.PageSize)
			if err != nil {
				return nil, err
			}

			return items, nil
		},
	},
	entityId: func(v quickbooks.Vendor) string {
		return v.Id
	},
	entityStatus: func(v quickbooks.Vendor) string {
		return v.Status
	},
}
