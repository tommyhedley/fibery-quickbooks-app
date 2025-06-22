package types

import (
	"github.com/tommyhedley/fibery-quickbooks-app/pkgs/app"
	"github.com/tommyhedley/fibery-quickbooks-app/pkgs/fibery"
	"github.com/tommyhedley/quickbooks-go"
)

var customer = app.NewDualType(
	"Customer",
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
	map[string]app.FieldDef[quickbooks.Customer]{
		"qboId": {
			Params: fibery.Field{
				Name:     "QBO ID",
				Type:     fibery.Text,
				ReadOnly: true,
				Ignore:   true,
			},
			Convert: func(sd app.StandardData[quickbooks.Customer]) (any, error) {
				return sd.Item.Id, nil
			},
		},
		"displayName": {
			Params: fibery.Field{
				Name:    "Display Name",
				Type:    fibery.Text,
				SubType: fibery.Title,
			},
			Convert: func(sd app.StandardData[quickbooks.Customer]) (any, error) {
				return sd.Item.DisplayName, nil
			},
		},
		"syncToken": {
			Params: fibery.Field{
				Name:     "Sync Token",
				Type:     fibery.Text,
				ReadOnly: true,
				Ignore:   true,
			},
			Convert: func(sd app.StandardData[quickbooks.Customer]) (any, error) {
				return sd.Item.SyncToken, nil
			},
		},
		"__syncAction": {
			Params: fibery.Field{
				Type: fibery.Text,
				Name: "Sync Action",
			},
			Convert: func(sd app.StandardData[quickbooks.Customer]) (any, error) {
				return fibery.SET, nil
			},
		},
		"active": {
			Params: fibery.Field{
				Name:    "Active",
				Type:    fibery.Text,
				SubType: fibery.Boolean,
			},
			Convert: func(sd app.StandardData[quickbooks.Customer]) (any, error) {
				return sd.Item.Active, nil
			},
		},
		"title": {
			Params: fibery.Field{
				Name: "Title",
				Type: fibery.Text,
			},
			Convert: func(sd app.StandardData[quickbooks.Customer]) (any, error) {
				return sd.Item.Title, nil
			},
		},
		"givenName": {
			Params: fibery.Field{
				Name: "First Name",
				Type: fibery.Text,
			},
			Convert: func(sd app.StandardData[quickbooks.Customer]) (any, error) {
				return sd.Item.GivenName, nil
			},
		},
		"middleName": {
			Params: fibery.Field{
				Name: "Middle Name",
				Type: fibery.Text,
			},
			Convert: func(sd app.StandardData[quickbooks.Customer]) (any, error) {
				return sd.Item.MiddleName, nil
			},
		},
		"familyName": {
			Params: fibery.Field{
				Name: "Last Name",
				Type: fibery.Text,
			},
			Convert: func(sd app.StandardData[quickbooks.Customer]) (any, error) {
				return sd.Item.FamilyName, nil
			},
		},
		"suffix": {
			Params: fibery.Field{
				Name: "Suffix",
				Type: fibery.Text,
			},
			Convert: func(sd app.StandardData[quickbooks.Customer]) (any, error) {
				return sd.Item.Suffix, nil
			},
		},
		"companyName": {
			Params: fibery.Field{
				Name: "Company Name",
				Type: fibery.Text,
			},
			Convert: func(sd app.StandardData[quickbooks.Customer]) (any, error) {
				return sd.Item.CompanyName, nil
			},
		},
		"primaryEmail": {
			Params: fibery.Field{
				Name:    "Email",
				Type:    fibery.Text,
				SubType: fibery.Email,
			},
			Convert: func(sd app.StandardData[quickbooks.Customer]) (any, error) {
				if sd.Item.PrimaryEmailAddr != nil {
					return sd.Item.PrimaryEmailAddr.Address, nil
				}
				return "", nil
			},
		},
		"taxable": {
			Params: fibery.Field{
				Name:    "Taxable",
				Type:    fibery.Text,
				SubType: fibery.Boolean,
			},
			Convert: func(sd app.StandardData[quickbooks.Customer]) (any, error) {
				return sd.Item.Taxable, nil
			},
		},
		"resaleNum": {
			Params: fibery.Field{
				Name: "Resale ID",
				Type: fibery.Text,
			},
			Convert: func(sd app.StandardData[quickbooks.Customer]) (any, error) {
				return sd.Item.ResaleNum, nil
			},
		},
		"primaryPhone": {
			Params: fibery.Field{
				Name: "Phone",
				Type: fibery.Text,
				Format: map[string]any{
					"format": "phone",
				},
			},
			Convert: func(sd app.StandardData[quickbooks.Customer]) (any, error) {
				if sd.Item.PrimaryPhone != nil {
					return sd.Item.PrimaryPhone.FreeFormNumber, nil
				}
				return "", nil
			},
		},
		"alternatePhone": {
			Params: fibery.Field{
				Name: "Alternate Phone",
				Type: fibery.Text,
				Format: map[string]any{
					"format": "phone",
				},
			},
			Convert: func(sd app.StandardData[quickbooks.Customer]) (any, error) {
				if sd.Item.AlternatePhone != nil {
					return sd.Item.AlternatePhone.FreeFormNumber, nil
				}
				return "", nil
			},
		},
		"mobile": {
			Params: fibery.Field{
				Name: "Mobile",
				Type: fibery.Text,
				Format: map[string]any{
					"format": "phone",
				},
			},
			Convert: func(sd app.StandardData[quickbooks.Customer]) (any, error) {
				if sd.Item.Mobile != nil {
					return sd.Item.Mobile.FreeFormNumber, nil
				}
				return "", nil
			},
		},
		"fax": {
			Params: fibery.Field{
				Name: "Fax",
				Type: fibery.Text,
				Format: map[string]any{
					"format": "phone",
				},
			},
			Convert: func(sd app.StandardData[quickbooks.Customer]) (any, error) {
				if sd.Item.Fax != nil {
					return sd.Item.Fax.FreeFormNumber, nil
				}
				return "", nil
			},
		},
		"job": {
			Params: fibery.Field{
				Name:    "Job",
				Type:    fibery.Text,
				SubType: fibery.Boolean,
			},
			Convert: func(sd app.StandardData[quickbooks.Customer]) (any, error) {
				if sd.Item.Job.Valid {
					return sd.Item.Job.Bool, nil
				}
				return false, nil
			},
		},
		"billWithParent": {
			Params: fibery.Field{
				Name:    "Bill With Parent",
				Type:    fibery.Text,
				SubType: fibery.Boolean,
			},
			Convert: func(sd app.StandardData[quickbooks.Customer]) (any, error) {
				return sd.Item.BillWithParent, nil
			},
		},
		"notes": {
			Params: fibery.Field{
				Name:    "Notes",
				Type:    fibery.Text,
				SubType: fibery.MD,
			},
			Convert: func(sd app.StandardData[quickbooks.Customer]) (any, error) {
				return sd.Item.Notes, nil
			},
		},
		"website": {
			Params: fibery.Field{
				Name:    "Website",
				Type:    fibery.Text,
				SubType: fibery.URL,
			},
			Convert: func(sd app.StandardData[quickbooks.Customer]) (any, error) {
				if sd.Item.WebAddr != nil {
					return sd.Item.WebAddr.URI, nil
				}
				return "", nil
			},
		},
		"balance": {
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
			Convert: func(sd app.StandardData[quickbooks.Customer]) (any, error) {
				return sd.Item.Balance, nil
			},
		},
		"balanceWithJobs": {
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
			Convert: func(sd app.StandardData[quickbooks.Customer]) (any, error) {
				return sd.Item.BalanceWithJobs, nil
			},
		},
		"shippingLine1": {
			Params: fibery.Field{
				Name: "Shipping Line 1",
				Type: fibery.Text,
			},
			Convert: func(sd app.StandardData[quickbooks.Customer]) (any, error) {
				if sd.Item.ShipAddr != nil {
					return sd.Item.ShipAddr.Line1, nil
				}
				return "", nil
			},
		},
		"shippingLine2": {
			Params: fibery.Field{
				Name: "Shipping Line 2",
				Type: fibery.Text,
			},
			Convert: func(sd app.StandardData[quickbooks.Customer]) (any, error) {
				if sd.Item.ShipAddr != nil {
					return sd.Item.ShipAddr.Line2, nil
				}
				return "", nil
			},
		},
		"shippingLine3": {
			Params: fibery.Field{
				Name: "Shipping Line 3",
				Type: fibery.Text,
			},
			Convert: func(sd app.StandardData[quickbooks.Customer]) (any, error) {
				if sd.Item.ShipAddr != nil {
					return sd.Item.ShipAddr.Line3, nil
				}
				return "", nil
			},
		},
		"shippingLine4": {
			Params: fibery.Field{
				Name: "Shipping Line 4",
				Type: fibery.Text,
			},
			Convert: func(sd app.StandardData[quickbooks.Customer]) (any, error) {
				if sd.Item.ShipAddr != nil {
					return sd.Item.ShipAddr.Line4, nil
				}
				return "", nil
			},
		},
		"shippingLine5": {
			Params: fibery.Field{
				Name: "Shipping Line 5",
				Type: fibery.Text,
			},
			Convert: func(sd app.StandardData[quickbooks.Customer]) (any, error) {
				if sd.Item.ShipAddr != nil {
					return sd.Item.ShipAddr.Line5, nil
				}
				return "", nil
			},
		},
		"shippingCity": {
			Params: fibery.Field{
				Name: "Shipping City",
				Type: fibery.Text,
			},
			Convert: func(sd app.StandardData[quickbooks.Customer]) (any, error) {
				if sd.Item.ShipAddr != nil {
					return sd.Item.ShipAddr.City, nil
				}
				return "", nil
			},
		},
		"shippingState": {
			Params: fibery.Field{
				Name: "Shipping State",
				Type: fibery.Text,
			},
			Convert: func(sd app.StandardData[quickbooks.Customer]) (any, error) {
				if sd.Item.ShipAddr != nil {
					return sd.Item.ShipAddr.CountrySubDivisionCode, nil
				}
				return "", nil
			},
		},
		"shippingPostalCode": {
			Params: fibery.Field{
				Name: "Shipping Postal Code",
				Type: fibery.Text,
			},
			Convert: func(sd app.StandardData[quickbooks.Customer]) (any, error) {
				if sd.Item.ShipAddr != nil {
					return sd.Item.ShipAddr.PostalCode, nil
				}
				return "", nil
			},
		},
		"shippingCountry": {
			Params: fibery.Field{
				Name: "Shipping Country",
				Type: fibery.Text,
			},
			Convert: func(sd app.StandardData[quickbooks.Customer]) (any, error) {
				if sd.Item.ShipAddr != nil {
					return sd.Item.ShipAddr.Country, nil
				}
				return "", nil
			},
		},
		"shippingLat": {
			Params: fibery.Field{
				Name: "Shipping Latitude",
				Type: fibery.Text,
			},
			Convert: func(sd app.StandardData[quickbooks.Customer]) (any, error) {
				if sd.Item.ShipAddr != nil {
					return sd.Item.ShipAddr.Lat, nil
				}
				return "", nil
			},
		},
		"shippingLong": {
			Params: fibery.Field{
				Name: "Shipping Longitude",
				Type: fibery.Text,
			},
			Convert: func(sd app.StandardData[quickbooks.Customer]) (any, error) {
				if sd.Item.ShipAddr != nil {
					return sd.Item.ShipAddr.Long, nil
				}
				return "", nil
			},
		},
		"billingLine1": {
			Params: fibery.Field{
				Name: "Billing Line 1",
				Type: fibery.Text,
			},
			Convert: func(sd app.StandardData[quickbooks.Customer]) (any, error) {
				if sd.Item.BillAddr != nil {
					return sd.Item.BillAddr.Line1, nil
				}
				return "", nil
			},
		},
		"billingLine2": {
			Params: fibery.Field{
				Name: "Billing Line 2",
				Type: fibery.Text,
			},
			Convert: func(sd app.StandardData[quickbooks.Customer]) (any, error) {
				if sd.Item.BillAddr != nil {
					return sd.Item.BillAddr.Line2, nil
				}
				return "", nil
			},
		},
		"billingLine3": {
			Params: fibery.Field{
				Name: "Billing Line 3",
				Type: fibery.Text,
			},
			Convert: func(sd app.StandardData[quickbooks.Customer]) (any, error) {
				if sd.Item.BillAddr != nil {
					return sd.Item.BillAddr.Line3, nil
				}
				return "", nil
			},
		},
		"billingLine4": {
			Params: fibery.Field{
				Name: "Billing Line 4",
				Type: fibery.Text,
			},
			Convert: func(sd app.StandardData[quickbooks.Customer]) (any, error) {
				if sd.Item.BillAddr != nil {
					return sd.Item.BillAddr.Line4, nil
				}
				return "", nil
			},
		},
		"billingLine5": {
			Params: fibery.Field{
				Name: "Billing Line 5",
				Type: fibery.Text,
			},
			Convert: func(sd app.StandardData[quickbooks.Customer]) (any, error) {
				if sd.Item.BillAddr != nil {
					return sd.Item.BillAddr.Line5, nil
				}
				return "", nil
			},
		},
		"billingCity": {
			Params: fibery.Field{
				Name: "Billing City",
				Type: fibery.Text,
			},
			Convert: func(sd app.StandardData[quickbooks.Customer]) (any, error) {
				if sd.Item.BillAddr != nil {
					return sd.Item.BillAddr.City, nil
				}
				return "", nil
			},
		},
		"billingState": {
			Params: fibery.Field{
				Name: "Billing State",
				Type: fibery.Text,
			},
			Convert: func(sd app.StandardData[quickbooks.Customer]) (any, error) {
				if sd.Item.BillAddr != nil {
					return sd.Item.BillAddr.CountrySubDivisionCode, nil
				}
				return "", nil
			},
		},
		"billingPostalCode": {
			Params: fibery.Field{
				Name: "Billing Postal Code",
				Type: fibery.Text,
			},
			Convert: func(sd app.StandardData[quickbooks.Customer]) (any, error) {
				if sd.Item.BillAddr != nil {
					return sd.Item.BillAddr.PostalCode, nil
				}
				return "", nil
			},
		},
		"billingCountry": {
			Params: fibery.Field{
				Name: "Billing Country",
				Type: fibery.Text,
			},
			Convert: func(sd app.StandardData[quickbooks.Customer]) (any, error) {
				if sd.Item.BillAddr != nil {
					return sd.Item.BillAddr.Country, nil
				}
				return "", nil
			},
		},
		"billingLat": {
			Params: fibery.Field{
				Name: "Billing Latitude",
				Type: fibery.Text,
			},
			Convert: func(sd app.StandardData[quickbooks.Customer]) (any, error) {
				if sd.Item.BillAddr != nil {
					return sd.Item.BillAddr.Lat, nil
				}
				return "", nil
			},
		},
		"billingLong": {
			Params: fibery.Field{
				Name: "Billing Longitude",
				Type: fibery.Text,
			},
			Convert: func(sd app.StandardData[quickbooks.Customer]) (any, error) {
				if sd.Item.BillAddr != nil {
					return sd.Item.BillAddr.Long, nil
				}
				return "", nil
			},
		},
		"taxExemptionId": {
			Params: fibery.Field{
				Name: "Tax Exemption ID",
				Type: fibery.Text,
				Relation: &fibery.Relation{
					Cardinality:   fibery.MTO,
					Name:          "Tax Exemption",
					TargetName:    "Customers",
					TargetType:    "taxExemption",
					TargetFieldID: "id",
				},
			},
			Convert: func(sd app.StandardData[quickbooks.Customer]) (any, error) {
				return sd.Item.TaxExemptionReasonId, nil
			},
		},
		"defaultTaxCodeId": {
			Params: fibery.Field{
				Name: "Default Tax Code ID",
				Type: fibery.Text,
				Relation: &fibery.Relation{
					Cardinality:   fibery.MTO,
					Name:          "Default Tax Code",
					TargetName:    "Customers",
					TargetType:    "taxCode",
					TargetFieldID: "id",
				},
			},
			Convert: func(sd app.StandardData[quickbooks.Customer]) (any, error) {
				if sd.Item.DefaultTaxCodeRef != nil {
					return sd.Item.DefaultTaxCodeRef.Value, nil
				}
				return "", nil
			},
		},
		"customerTypeId": {
			Params: fibery.Field{
				Name: "Customer Type ID",
				Type: fibery.Text,
				Relation: &fibery.Relation{
					Cardinality:   fibery.MTO,
					Name:          "Customer Type",
					TargetName:    "Customers",
					TargetType:    "customerType",
					TargetFieldID: "id",
				},
			},
			Convert: func(sd app.StandardData[quickbooks.Customer]) (any, error) {
				if sd.Item.CustomerTypeRef != nil {
					return sd.Item.CustomerTypeRef.Value, nil
				}
				return "", nil
			},
		},
		"salesTermId": {
			Params: fibery.Field{
				Name: "Sales Term ID",
				Type: fibery.Text,
				Relation: &fibery.Relation{
					Cardinality:   fibery.MTO,
					Name:          "Sales Term",
					TargetName:    "Customers",
					TargetType:    "salesTerm",
					TargetFieldID: "id",
				},
			},
			Convert: func(sd app.StandardData[quickbooks.Customer]) (any, error) {
				if sd.Item.SalesTermRef != nil {
					return sd.Item.SalesTermRef.Value, nil
				}
				return "", nil
			},
		},
		"paymentMethodId": {
			Params: fibery.Field{
				Name: "Payment Method ID",
				Type: fibery.Text,
				Relation: &fibery.Relation{
					Cardinality:   fibery.MTO,
					Name:          "Payment Method",
					TargetName:    "Customers",
					TargetType:    "paymentMethod",
					TargetFieldID: "id",
				},
			},
			Convert: func(sd app.StandardData[quickbooks.Customer]) (any, error) {
				if sd.Item.PaymentMethodRef != nil {
					return sd.Item.PaymentMethodRef.Value, nil
				}
				return "", nil
			},
		},
		"parentId": {
			Params: fibery.Field{
				Name: "Parent ID",
				Type: fibery.Text,
				Relation: &fibery.Relation{
					Cardinality:   fibery.MTO,
					Name:          "Parent",
					TargetName:    "Jobs",
					TargetType:    "customer",
					TargetFieldID: "id",
				},
			},
			Convert: func(sd app.StandardData[quickbooks.Customer]) (any, error) {
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
	app.Types.Register(customer)
}
