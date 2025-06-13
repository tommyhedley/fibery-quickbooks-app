package types

import (
	"github.com/tommyhedley/fibery-quickbooks-app/pkgs/fibery"
	"github.com/tommyhedley/fibery-quickbooks-app/pkgs/integration"
	"github.com/tommyhedley/quickbooks-go"
)

var vendor = integration.NewDualType(
	"vendor",
	"vendor",
	"Vendor",
	func(v quickbooks.Vendor) string {
		return v.Id
	},
	func(v quickbooks.Vendor) string {
		return v.Status
	},
	func(id string) quickbooks.Vendor {
		return quickbooks.Vendor{
			Id: id,
		}
	},
	func(bir quickbooks.BatchItemResponse) quickbooks.Vendor {
		return bir.Vendor
	},
	func(bqr quickbooks.BatchQueryResponse) []quickbooks.Vendor {
		return bqr.Vendor
	},
	func(cr quickbooks.CDCQueryResponse) []quickbooks.Vendor {
		return cr.Vendor
	},
	map[string]integration.FieldDef[quickbooks.Vendor]{
		"QBOId": {
			Params: fibery.Field{
				Name: "QBO ID",
				Type: fibery.Text,
			},
			Convert: func(sd integration.StandardData[quickbooks.Vendor]) (any, error) {
				return sd.Item.Id, nil
			},
		},
		"DisplayName": {
			Params: fibery.Field{
				Name:    "Display Name",
				Type:    fibery.Text,
				SubType: fibery.Title,
			},
			Convert: func(sd integration.StandardData[quickbooks.Vendor]) (any, error) {
				return sd.Item.DisplayName, nil
			},
		},
		"SyncToken": {
			Params: fibery.Field{
				Name:     "Sync Token",
				Type:     fibery.Text,
				ReadOnly: true,
			},
			Convert: func(sd integration.StandardData[quickbooks.Vendor]) (any, error) {
				return sd.Item.SyncToken, nil
			},
		},
		"__syncAction": {
			Params: fibery.Field{
				Type: fibery.Text,
				Name: "Sync Action",
			},
			Convert: func(sd integration.StandardData[quickbooks.Vendor]) (any, error) {
				return fibery.SET, nil
			},
		},
		"Active": {
			Params: fibery.Field{
				Name:    "Active",
				Type:    fibery.Text,
				SubType: fibery.Boolean,
			},
			Convert: func(sd integration.StandardData[quickbooks.Vendor]) (any, error) {
				return sd.Item.Active, nil
			},
		},
		"Title": {
			Params: fibery.Field{
				Name: "Title",
				Type: fibery.Text,
			},
			Convert: func(sd integration.StandardData[quickbooks.Vendor]) (any, error) {
				return sd.Item.Title, nil
			},
		},
		"GivenName": {
			Params: fibery.Field{
				Name: "First Name",
				Type: fibery.Text,
			},
			Convert: func(sd integration.StandardData[quickbooks.Vendor]) (any, error) {
				return sd.Item.GivenName, nil
			},
		},
		"MiddleName": {
			Params: fibery.Field{
				Name: "Middle Name",
				Type: fibery.Text,
			},
			Convert: func(sd integration.StandardData[quickbooks.Vendor]) (any, error) {
				return sd.Item.MiddleName, nil
			},
		},
		"FamilyName": {
			Params: fibery.Field{
				Name: "Last Name",
				Type: fibery.Text,
			},
			Convert: func(sd integration.StandardData[quickbooks.Vendor]) (any, error) {
				return sd.Item.FamilyName, nil
			},
		},
		"Suffix": {
			Params: fibery.Field{
				Name: "Suffix",
				Type: fibery.Text,
			},
			Convert: func(sd integration.StandardData[quickbooks.Vendor]) (any, error) {
				return sd.Item.Suffix, nil
			},
		},
		"CompanyName": {
			Params: fibery.Field{
				Name: "Company Name",
				Type: fibery.Text,
			},
			Convert: func(sd integration.StandardData[quickbooks.Vendor]) (any, error) {
				return sd.Item.CompanyName, nil
			},
		},
		"PrimaryEmail": {
			Params: fibery.Field{
				Name:    "Email",
				Type:    fibery.Text,
				SubType: fibery.Email,
			},
			Convert: func(sd integration.StandardData[quickbooks.Vendor]) (any, error) {
				if sd.Item.PrimaryEmailAddr != nil {
					return sd.Item.PrimaryEmailAddr.Address, nil
				}
				return "", nil
			},
		},
		"SalesTermId": {
			Params: fibery.Field{
				Name: "Sales Term ID",
				Type: fibery.Text,
				Relation: &fibery.Relation{
					Cardinality:   fibery.MTO,
					Name:          "Sales Term",
					TargetName:    "Vendors",
					TargetType:    "SalesTerm",
					TargetFieldID: "id",
				},
			},
			Convert: func(sd integration.StandardData[quickbooks.Vendor]) (any, error) {
				if sd.Item.TermRef != nil {
					return sd.Item.TermRef.Value, nil
				}
				return "", nil
			},
		},
		"PrimaryPhone": {
			Params: fibery.Field{
				Name: "Phone",
				Type: fibery.Text,
				Format: map[string]any{
					"format": "phone",
				},
			},
			Convert: func(sd integration.StandardData[quickbooks.Vendor]) (any, error) {
				if sd.Item.PrimaryPhone != nil {
					return sd.Item.PrimaryPhone.FreeFormNumber, nil
				}
				return "", nil
			},
		},
		"AlternatePhone": {
			Params: fibery.Field{
				Name: "Alternate Phone",
				Type: fibery.Text,
				Format: map[string]any{
					"format": "phone",
				},
			},
			Convert: func(sd integration.StandardData[quickbooks.Vendor]) (any, error) {
				if sd.Item.AlternatePhone != nil {
					return sd.Item.AlternatePhone.FreeFormNumber, nil
				}
				return "", nil
			},
		},
		"Mobile": {
			Params: fibery.Field{
				Name: "Mobile",
				Type: fibery.Text,
				Format: map[string]any{
					"format": "phone",
				},
			},
			Convert: func(sd integration.StandardData[quickbooks.Vendor]) (any, error) {
				if sd.Item.Mobile != nil {
					return sd.Item.Mobile.FreeFormNumber, nil
				}
				return "", nil
			},
		},
		"Fax": {
			Params: fibery.Field{
				Name: "Fax",
				Type: fibery.Text,
				Format: map[string]any{
					"format": "phone",
				},
			},
			Convert: func(sd integration.StandardData[quickbooks.Vendor]) (any, error) {
				if sd.Item.Fax != nil {
					return sd.Item.Fax.FreeFormNumber, nil
				}
				return "", nil
			},
		},
		"1099": {
			Params: fibery.Field{
				Name:        "1099",
				Type:        fibery.Text,
				SubType:     fibery.Boolean,
				Description: "Is the Vendor a 1099 contractor?",
			},
			Convert: func(sd integration.StandardData[quickbooks.Vendor]) (any, error) {
				return sd.Item.Vendor1099, nil
			},
		},
		"CostRate": {
			Params: fibery.Field{
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
			Convert: func(sd integration.StandardData[quickbooks.Vendor]) (any, error) {
				return sd.Item.CostRate, nil
			},
		},
		"BillRate": {
			Params: fibery.Field{
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
			Convert: func(sd integration.StandardData[quickbooks.Vendor]) (any, error) {
				return sd.Item.BillRate, nil
			},
		},
		"Website": {
			Params: fibery.Field{
				Name:    "Website",
				Type:    fibery.Text,
				SubType: fibery.URL,
			},
			Convert: func(sd integration.StandardData[quickbooks.Vendor]) (any, error) {
				if sd.Item.WebAddr != nil {
					return sd.Item.WebAddr.URI, nil
				}
				return "", nil
			},
		},
		"AccountNumber": {
			Params: fibery.Field{
				Name:        "Account Number",
				Type:        fibery.Text,
				Description: "Name or number of the account associated with this vendor",
			},
			Convert: func(sd integration.StandardData[quickbooks.Vendor]) (any, error) {
				return sd.Item.AcctNum, nil
			},
		},
		"Balance": {
			Params: fibery.Field{
				Name: "Balance",
				Type: fibery.Number,
				Format: map[string]any{
					"format":               "Money",
					"currencyCode":         "USD",
					"hasThousandSeperator": true,
					"precision":            2,
				},
			},
			Convert: func(sd integration.StandardData[quickbooks.Vendor]) (any, error) {
				return sd.Item.Balance, nil
			},
		},
		"BillingLine1": {
			Params: fibery.Field{
				Name: "Billing Line 1",
				Type: fibery.Text,
			},
			Convert: func(sd integration.StandardData[quickbooks.Vendor]) (any, error) {
				if sd.Item.BillAddr != nil {
					return sd.Item.BillAddr.Line1, nil
				}
				return "", nil
			},
		},
		"BillingLine2": {
			Params: fibery.Field{
				Name: "Billing Line 2",
				Type: fibery.Text,
			},
			Convert: func(sd integration.StandardData[quickbooks.Vendor]) (any, error) {
				if sd.Item.BillAddr != nil {
					return sd.Item.BillAddr.Line2, nil
				}
				return "", nil
			},
		},
		"BillingLine3": {
			Params: fibery.Field{
				Name: "Billing Line 3",
				Type: fibery.Text,
			},
			Convert: func(sd integration.StandardData[quickbooks.Vendor]) (any, error) {
				if sd.Item.BillAddr != nil {
					return sd.Item.BillAddr.Line3, nil
				}
				return "", nil
			},
		},
		"BillingLine4": {
			Params: fibery.Field{
				Name: "Billing Line 4",
				Type: fibery.Text,
			},
			Convert: func(sd integration.StandardData[quickbooks.Vendor]) (any, error) {
				if sd.Item.BillAddr != nil {
					return sd.Item.BillAddr.Line4, nil
				}
				return "", nil
			},
		},
		"BillingLine5": {
			Params: fibery.Field{
				Name: "Billing Line 5",
				Type: fibery.Text,
			},
			Convert: func(sd integration.StandardData[quickbooks.Vendor]) (any, error) {
				if sd.Item.BillAddr != nil {
					return sd.Item.BillAddr.Line5, nil
				}
				return "", nil
			},
		},
		"BillingCity": {
			Params: fibery.Field{
				Name: "Billing City",
				Type: fibery.Text,
			},
			Convert: func(sd integration.StandardData[quickbooks.Vendor]) (any, error) {
				if sd.Item.BillAddr != nil {
					return sd.Item.BillAddr.City, nil
				}
				return "", nil
			},
		},
		"BillingState": {
			Params: fibery.Field{
				Name: "Billing State",
				Type: fibery.Text,
			},
			Convert: func(sd integration.StandardData[quickbooks.Vendor]) (any, error) {
				if sd.Item.BillAddr != nil {
					return sd.Item.BillAddr.CountrySubDivisionCode, nil
				}
				return "", nil
			},
		},
		"BillingPostalCode": {
			Params: fibery.Field{
				Name: "Billing Postal Code",
				Type: fibery.Text,
			},
			Convert: func(sd integration.StandardData[quickbooks.Vendor]) (any, error) {
				if sd.Item.BillAddr != nil {
					return sd.Item.BillAddr.PostalCode, nil
				}
				return "", nil
			},
		},
		"BillingCountry": {
			Params: fibery.Field{
				Name: "Billing Country",
				Type: fibery.Text,
			},
			Convert: func(sd integration.StandardData[quickbooks.Vendor]) (any, error) {
				if sd.Item.BillAddr != nil {
					return sd.Item.BillAddr.Country, nil
				}
				return "", nil
			},
		},
		"BillingLat": {
			Params: fibery.Field{
				Name: "Billing Latitude",
				Type: fibery.Text,
			},
			Convert: func(sd integration.StandardData[quickbooks.Vendor]) (any, error) {
				if sd.Item.BillAddr != nil {
					return sd.Item.BillAddr.Lat, nil
				}
				return "", nil
			},
		},
		"BillingLong": {
			Params: fibery.Field{
				Name: "Billing Longitude",
				Type: fibery.Text,
			},
			Convert: func(sd integration.StandardData[quickbooks.Vendor]) (any, error) {
				if sd.Item.BillAddr != nil {
					return sd.Item.BillAddr.Long, nil
				}
				return "", nil
			},
		},
	},
	nil,
)

func init() {
	integration.Types.Register(vendor)
}

