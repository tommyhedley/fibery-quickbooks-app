package types

import (
	"github.com/tommyhedley/fibery-quickbooks-app/pkgs/fibery"
	"github.com/tommyhedley/fibery-quickbooks-app/pkgs/integration"
	"github.com/tommyhedley/quickbooks-go"
)

var customer = integration.NewDualType(
	"customer",
	"customer",
	"Customer",
	func(c quickbooks.Customer) string {
		return c.Id
	},
	func(c quickbooks.Customer) string {
		return c.Status
	},
	func(id string) quickbooks.Customer {
		return quickbooks.Customer{
			Id: id,
		}
	},
	func(bir quickbooks.BatchItemResponse) quickbooks.Customer {
		return bir.Customer
	},
	func(bqr quickbooks.BatchQueryResponse) []quickbooks.Customer {
		return bqr.Customer
	},
	func(cr quickbooks.CDCQueryResponse) []quickbooks.Customer {
		return cr.Customer
	},
	map[string]integration.FieldDef[quickbooks.Customer]{
		"QBOId": {
			Params: fibery.Field{
				Name: "QBO ID",
				Type: fibery.Text,
			},
			Convert: func(sd integration.StandardData[quickbooks.Customer]) (any, error) {
				return sd.Item.Id, nil
			},
		},
		"DisplayName": {
			Params: fibery.Field{
				Name:    "Display Name",
				Type:    fibery.Text,
				SubType: fibery.Title,
			},
			Convert: func(sd integration.StandardData[quickbooks.Customer]) (any, error) {
				return sd.Item.DisplayName, nil
			},
		},
		"SyncToken": {
			Params: fibery.Field{
				Name:     "Sync Token",
				Type:     fibery.Text,
				ReadOnly: true,
			},
			Convert: func(sd integration.StandardData[quickbooks.Customer]) (any, error) {
				return sd.Item.SyncToken, nil
			},
		},
		"__syncAction": {
			Params: fibery.Field{
				Type: fibery.Text,
				Name: "Sync Action",
			},
			Convert: func(sd integration.StandardData[quickbooks.Customer]) (any, error) {
				return fibery.SET, nil
			},
		},
		"Active": {
			Params: fibery.Field{
				Name:    "Active",
				Type:    fibery.Text,
				SubType: fibery.Boolean,
			},
			Convert: func(sd integration.StandardData[quickbooks.Customer]) (any, error) {
				return sd.Item.Active, nil
			},
		},
		"Title": {
			Params: fibery.Field{
				Name: "Title",
				Type: fibery.Text,
			},
			Convert: func(sd integration.StandardData[quickbooks.Customer]) (any, error) {
				return sd.Item.Title, nil
			},
		},
		"GivenName": {
			Params: fibery.Field{
				Name: "First Name",
				Type: fibery.Text,
			},
			Convert: func(sd integration.StandardData[quickbooks.Customer]) (any, error) {
				return sd.Item.GivenName, nil
			},
		},
		"MiddleName": {
			Params: fibery.Field{
				Name: "Middle Name",
				Type: fibery.Text,
			},
			Convert: func(sd integration.StandardData[quickbooks.Customer]) (any, error) {
				return sd.Item.MiddleName, nil
			},
		},
		"FamilyName": {
			Params: fibery.Field{
				Name: "Last Name",
				Type: fibery.Text,
			},
			Convert: func(sd integration.StandardData[quickbooks.Customer]) (any, error) {
				return sd.Item.FamilyName, nil
			},
		},
		"Suffix": {
			Params: fibery.Field{
				Name: "Suffix",
				Type: fibery.Text,
			},
			Convert: func(sd integration.StandardData[quickbooks.Customer]) (any, error) {
				return sd.Item.Suffix, nil
			},
		},
		"CompanyName": {
			Params: fibery.Field{
				Name: "Company Name",
				Type: fibery.Text,
			},
			Convert: func(sd integration.StandardData[quickbooks.Customer]) (any, error) {
				return sd.Item.CompanyName, nil
			},
		},
		"PrimaryEmail": {
			Params: fibery.Field{
				Name:    "Email",
				Type:    fibery.Text,
				SubType: fibery.Email,
			},
			Convert: func(sd integration.StandardData[quickbooks.Customer]) (any, error) {
				if sd.Item.PrimaryEmailAddr != nil {
					return sd.Item.PrimaryEmailAddr.Address, nil
				}
				return "", nil
			},
		},
		"Taxable": {
			Params: fibery.Field{
				Name:    "Taxable",
				Type:    fibery.Text,
				SubType: fibery.Boolean,
			},
			Convert: func(sd integration.StandardData[quickbooks.Customer]) (any, error) {
				return sd.Item.Taxable, nil
			},
		},
		"ResaleNum": {
			Params: fibery.Field{
				Name: "Resale ID",
				Type: fibery.Text,
			},
			Convert: func(sd integration.StandardData[quickbooks.Customer]) (any, error) {
				return sd.Item.ResaleNum, nil
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
			Convert: func(sd integration.StandardData[quickbooks.Customer]) (any, error) {
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
			Convert: func(sd integration.StandardData[quickbooks.Customer]) (any, error) {
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
			Convert: func(sd integration.StandardData[quickbooks.Customer]) (any, error) {
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
			Convert: func(sd integration.StandardData[quickbooks.Customer]) (any, error) {
				if sd.Item.Fax != nil {
					return sd.Item.Fax.FreeFormNumber, nil
				}
				return "", nil
			},
		},
		"Job": {
			Params: fibery.Field{
				Name:    "Job",
				Type:    fibery.Text,
				SubType: fibery.Boolean,
			},
			Convert: func(sd integration.StandardData[quickbooks.Customer]) (any, error) {
				if sd.Item.Job.Valid {
					return sd.Item.Job.Bool, nil
				}
				return false, nil
			},
		},
		"BillWithParent": {
			Params: fibery.Field{
				Name:    "Bill With Parent",
				Type:    fibery.Text,
				SubType: fibery.Boolean,
			},
			Convert: func(sd integration.StandardData[quickbooks.Customer]) (any, error) {
				return sd.Item.BillWithParent, nil
			},
		},
		"Notes": {
			Params: fibery.Field{
				Name:    "Notes",
				Type:    fibery.Text,
				SubType: fibery.MD,
			},
			Convert: func(sd integration.StandardData[quickbooks.Customer]) (any, error) {
				return sd.Item.Notes, nil
			},
		},
		"Website": {
			Params: fibery.Field{
				Name:    "Website",
				Type:    fibery.Text,
				SubType: fibery.URL,
			},
			Convert: func(sd integration.StandardData[quickbooks.Customer]) (any, error) {
				if sd.Item.WebAddr != nil {
					return sd.Item.WebAddr.URI, nil
				}
				return "", nil
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
			Convert: func(sd integration.StandardData[quickbooks.Customer]) (any, error) {
				return sd.Item.Balance, nil
			},
		},
		"BalanceWithJobs": {
			Params: fibery.Field{
				Name: "Balance With Jobs",
				Type: fibery.Number,
				Format: map[string]any{
					"format":               "Money",
					"currencyCode":         "USD",
					"hasThousandSeperator": true,
					"precision":            2,
				},
			},
			Convert: func(sd integration.StandardData[quickbooks.Customer]) (any, error) {
				return sd.Item.BalanceWithJobs, nil
			},
		},
		"ShippingLine1": {
			Params: fibery.Field{
				Name: "Shipping Line 1",
				Type: fibery.Text,
			},
			Convert: func(sd integration.StandardData[quickbooks.Customer]) (any, error) {
				if sd.Item.ShipAddr != nil {
					return sd.Item.ShipAddr.Line1, nil
				}
				return "", nil
			},
		},
		"ShippingLine2": {
			Params: fibery.Field{
				Name: "Shipping Line 2",
				Type: fibery.Text,
			},
			Convert: func(sd integration.StandardData[quickbooks.Customer]) (any, error) {
				if sd.Item.ShipAddr != nil {
					return sd.Item.ShipAddr.Line2, nil
				}
				return "", nil
			},
		},
		"ShippingLine3": {
			Params: fibery.Field{
				Name: "Shipping Line 3",
				Type: fibery.Text,
			},
			Convert: func(sd integration.StandardData[quickbooks.Customer]) (any, error) {
				if sd.Item.ShipAddr != nil {
					return sd.Item.ShipAddr.Line3, nil
				}
				return "", nil
			},
		},
		"ShippingLine4": {
			Params: fibery.Field{
				Name: "Shipping Line 4",
				Type: fibery.Text,
			},
			Convert: func(sd integration.StandardData[quickbooks.Customer]) (any, error) {
				if sd.Item.ShipAddr != nil {
					return sd.Item.ShipAddr.Line4, nil
				}
				return "", nil
			},
		},
		"ShippingLine5": {
			Params: fibery.Field{
				Name: "Shipping Line 5",
				Type: fibery.Text,
			},
			Convert: func(sd integration.StandardData[quickbooks.Customer]) (any, error) {
				if sd.Item.ShipAddr != nil {
					return sd.Item.ShipAddr.Line5, nil
				}
				return "", nil
			},
		},
		"ShippingCity": {
			Params: fibery.Field{
				Name: "Shipping City",
				Type: fibery.Text,
			},
			Convert: func(sd integration.StandardData[quickbooks.Customer]) (any, error) {
				if sd.Item.ShipAddr != nil {
					return sd.Item.ShipAddr.City, nil
				}
				return "", nil
			},
		},
		"ShippingState": {
			Params: fibery.Field{
				Name: "Shipping State",
				Type: fibery.Text,
			},
			Convert: func(sd integration.StandardData[quickbooks.Customer]) (any, error) {
				if sd.Item.ShipAddr != nil {
					return sd.Item.ShipAddr.CountrySubDivisionCode, nil
				}
				return "", nil
			},
		},
		"ShippingPostalCode": {
			Params: fibery.Field{
				Name: "Shipping Postal Code",
				Type: fibery.Text,
			},
			Convert: func(sd integration.StandardData[quickbooks.Customer]) (any, error) {
				if sd.Item.ShipAddr != nil {
					return sd.Item.ShipAddr.PostalCode, nil
				}
				return "", nil
			},
		},
		"ShippingCountry": {
			Params: fibery.Field{
				Name: "Shipping Country",
				Type: fibery.Text,
			},
			Convert: func(sd integration.StandardData[quickbooks.Customer]) (any, error) {
				if sd.Item.ShipAddr != nil {
					return sd.Item.ShipAddr.Country, nil
				}
				return "", nil
			},
		},
		"ShippingLat": {
			Params: fibery.Field{
				Name: "Shipping Latitude",
				Type: fibery.Text,
			},
			Convert: func(sd integration.StandardData[quickbooks.Customer]) (any, error) {
				if sd.Item.ShipAddr != nil {
					return sd.Item.ShipAddr.Lat, nil
				}
				return "", nil
			},
		},
		"ShippingLong": {
			Params: fibery.Field{
				Name: "Shipping Longitude",
				Type: fibery.Text,
			},
			Convert: func(sd integration.StandardData[quickbooks.Customer]) (any, error) {
				if sd.Item.ShipAddr != nil {
					return sd.Item.ShipAddr.Long, nil
				}
				return "", nil
			},
		},
		"BillingLine1": {
			Params: fibery.Field{
				Name: "Billing Line 1",
				Type: fibery.Text,
			},
			Convert: func(sd integration.StandardData[quickbooks.Customer]) (any, error) {
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
			Convert: func(sd integration.StandardData[quickbooks.Customer]) (any, error) {
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
			Convert: func(sd integration.StandardData[quickbooks.Customer]) (any, error) {
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
			Convert: func(sd integration.StandardData[quickbooks.Customer]) (any, error) {
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
			Convert: func(sd integration.StandardData[quickbooks.Customer]) (any, error) {
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
			Convert: func(sd integration.StandardData[quickbooks.Customer]) (any, error) {
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
			Convert: func(sd integration.StandardData[quickbooks.Customer]) (any, error) {
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
			Convert: func(sd integration.StandardData[quickbooks.Customer]) (any, error) {
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
			Convert: func(sd integration.StandardData[quickbooks.Customer]) (any, error) {
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
			Convert: func(sd integration.StandardData[quickbooks.Customer]) (any, error) {
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
			Convert: func(sd integration.StandardData[quickbooks.Customer]) (any, error) {
				if sd.Item.BillAddr != nil {
					return sd.Item.BillAddr.Long, nil
				}
				return "", nil
			},
		},
		"TaxExemptionId": {
			Params: fibery.Field{
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
			Convert: func(sd integration.StandardData[quickbooks.Customer]) (any, error) {
				return sd.Item.TaxExemptionReasonId, nil
			},
		},
		"DefaultTaxCodeId": {
			Params: fibery.Field{
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
			Convert: func(sd integration.StandardData[quickbooks.Customer]) (any, error) {
				if sd.Item.DefaultTaxCodeRef != nil {
					return sd.Item.DefaultTaxCodeRef.Value, nil
				}
				return "", nil
			},
		},
		"CustomerTypeId": {
			Params: fibery.Field{
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
			Convert: func(sd integration.StandardData[quickbooks.Customer]) (any, error) {
				if sd.Item.CustomerTypeRef != nil {
					return sd.Item.CustomerTypeRef.Value, nil
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
					TargetName:    "Customers",
					TargetType:    "SalesTerm",
					TargetFieldID: "id",
				},
			},
			Convert: func(sd integration.StandardData[quickbooks.Customer]) (any, error) {
				if sd.Item.SalesTermRef != nil {
					return sd.Item.SalesTermRef.Value, nil
				}
				return "", nil
			},
		},
		"PaymentMethodId": {
			Params: fibery.Field{
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
			Convert: func(sd integration.StandardData[quickbooks.Customer]) (any, error) {
				if sd.Item.PaymentMethodRef != nil {
					return sd.Item.PaymentMethodRef.Value, nil
				}
				return "", nil
			},
		},
		"ParentId": {
			Params: fibery.Field{
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
			Convert: func(sd integration.StandardData[quickbooks.Customer]) (any, error) {
				if sd.Item.ParentRef != nil {
					return sd.Item.ParentRef.Value, nil
				}
				return "", nil
			},
		},
	},
	nil,
)

func init() {
	integration.Types.Register(customer)
}

