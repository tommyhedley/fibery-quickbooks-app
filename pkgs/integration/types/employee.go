package types

import (
	"github.com/tommyhedley/fibery-quickbooks-app/pkgs/fibery"
	"github.com/tommyhedley/fibery-quickbooks-app/pkgs/integration"
	"github.com/tommyhedley/quickbooks-go"
)

var employee = integration.NewDualType(
	"employee",
	"employee",
	"Employee",
	func(e quickbooks.Employee) string {
		return e.Id
	},
	func(e quickbooks.Employee) string {
		return e.Status
	},
	func(id string) quickbooks.Employee {
		return quickbooks.Employee{
			Id: id,
		}
	},
	func(bir quickbooks.BatchItemResponse) quickbooks.Employee {
		return bir.Employee
	},
	func(bqr quickbooks.BatchQueryResponse) []quickbooks.Employee {
		return bqr.Employee
	},
	func(cr quickbooks.CDCQueryResponse) []quickbooks.Employee {
		return cr.Employee
	},
	map[string]integration.FieldDef[quickbooks.Employee]{
		"QBOId": {
			Params: fibery.Field{
				Name: "QBO ID",
				Type: fibery.Text,
			},
			Convert: func(sd integration.StandardData[quickbooks.Employee]) (any, error) {
				return sd.Item.Id, nil
			},
		},
		"DisplayName": {
			Params: fibery.Field{
				Name:    "Display Name",
				Type:    fibery.Text,
				SubType: fibery.Title,
			},
			Convert: func(sd integration.StandardData[quickbooks.Employee]) (any, error) {
				return sd.Item.DisplayName, nil
			},
		},
		"SyncToken": {
			Params: fibery.Field{
				Name:     "Sync Token",
				Type:     fibery.Text,
				ReadOnly: true,
			},
			Convert: func(sd integration.StandardData[quickbooks.Employee]) (any, error) {
				return sd.Item.SyncToken, nil
			},
		},
		"__syncAction": {
			Params: fibery.Field{
				Type: fibery.Text,
				Name: "Sync Action",
			},
			Convert: func(sd integration.StandardData[quickbooks.Employee]) (any, error) {
				return fibery.SET, nil
			},
		},
		"Active": {
			Params: fibery.Field{
				Name:    "Active",
				Type:    fibery.Text,
				SubType: fibery.Boolean,
			},
			Convert: func(sd integration.StandardData[quickbooks.Employee]) (any, error) {
				return sd.Item.Active, nil
			},
		},
		"Title": {
			Params: fibery.Field{
				Name: "Title",
				Type: fibery.Text,
			},
			Convert: func(sd integration.StandardData[quickbooks.Employee]) (any, error) {
				return sd.Item.Title, nil
			},
		},
		"GivenName": {
			Params: fibery.Field{
				Name: "First Name",
				Type: fibery.Text,
			},
			Convert: func(sd integration.StandardData[quickbooks.Employee]) (any, error) {
				return sd.Item.GivenName, nil
			},
		},
		"MiddleName": {
			Params: fibery.Field{
				Name: "Middle Name",
				Type: fibery.Text,
			},
			Convert: func(sd integration.StandardData[quickbooks.Employee]) (any, error) {
				return sd.Item.MiddleName, nil
			},
		},
		"FamilyName": {
			Params: fibery.Field{
				Name: "Last Name",
				Type: fibery.Text,
			},
			Convert: func(sd integration.StandardData[quickbooks.Employee]) (any, error) {
				return sd.Item.FamilyName, nil
			},
		},
		"Suffix": {
			Params: fibery.Field{
				Name: "Suffix",
				Type: fibery.Text,
			},
			Convert: func(sd integration.StandardData[quickbooks.Employee]) (any, error) {
				return sd.Item.Suffix, nil
			},
		},
		"PrimaryEmailAddr": {
			Params: fibery.Field{
				Name:    "Email",
				Type:    fibery.Text,
				SubType: fibery.Email,
			},
			Convert: func(sd integration.StandardData[quickbooks.Employee]) (any, error) {
				if sd.Item.PrimaryEmailAddr != nil {
					return sd.Item.PrimaryEmailAddr.Address, nil
				}
				return "", nil
			},
		},
		"BillableTime": {
			Params: fibery.Field{
				Name:        "Billable",
				Type:        fibery.Text,
				SubType:     fibery.Boolean,
				Description: "Is the entity enabled for use in QuickBooks?",
			},
			Convert: func(sd integration.StandardData[quickbooks.Employee]) (any, error) {
				return sd.Item.BillableTime, nil
			},
		},
		"BirthDate": {
			Params: fibery.Field{
				Name:    "Date of Birth",
				Type:    fibery.DateType,
				SubType: fibery.Day,
			},
			Convert: func(sd integration.StandardData[quickbooks.Employee]) (any, error) {
				if sd.Item.BirthDate != nil {
					return sd.Item.BirthDate.Format(fibery.DateFormat), nil
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
			Convert: func(sd integration.StandardData[quickbooks.Employee]) (any, error) {
				if sd.Item.PrimaryPhone != nil {
					return sd.Item.PrimaryPhone.FreeFormNumber, nil
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
			Convert: func(sd integration.StandardData[quickbooks.Employee]) (any, error) {
				if sd.Item.Mobile != nil {
					return sd.Item.Mobile.FreeFormNumber, nil
				}
				return "", nil
			},
		},
		"CostRate": {
			Params: fibery.Field{
				Name: "Cost Rate",
				Type: fibery.Number,
				Format: map[string]any{
					"format":               "Money",
					"currencyCode":         "USD",
					"hasThousandSeperator": true,
					"precision":            2,
				},
			},
			Convert: func(sd integration.StandardData[quickbooks.Employee]) (any, error) {
				return sd.Item.CostRate, nil
			},
		},
		"BillRate": {
			Params: fibery.Field{
				Name: "Bill Rate",
				Type: fibery.Number,
				Format: map[string]any{
					"format":               "Money",
					"currencyCode":         "USD",
					"hasThousandSeperator": true,
					"precision":            2,
				},
			},
			Convert: func(sd integration.StandardData[quickbooks.Employee]) (any, error) {
				return sd.Item.BillRate, nil
			},
		},
		"EmployeeNumber": {
			Params: fibery.Field{
				Name: "Employee ID",
				Type: fibery.Text,
			},
			Convert: func(sd integration.StandardData[quickbooks.Employee]) (any, error) {
				return sd.Item.EmployeeNumber, nil
			},
		},
		"AddressLine1": {
			Params: fibery.Field{
				Name: "Address Line 1",
				Type: fibery.Text,
			},
			Convert: func(sd integration.StandardData[quickbooks.Employee]) (any, error) {
				return sd.Item.PrimaryAddr.Line1, nil
			},
		},
		"AddressLine2": {
			Params: fibery.Field{
				Name: "Address Line 2",
				Type: fibery.Text,
			},
			Convert: func(sd integration.StandardData[quickbooks.Employee]) (any, error) {
				return sd.Item.PrimaryAddr.Line2, nil
			},
		},
		"AddressLine3": {
			Params: fibery.Field{
				Name: "Address Line 3",
				Type: fibery.Text,
			},
			Convert: func(sd integration.StandardData[quickbooks.Employee]) (any, error) {
				return sd.Item.PrimaryAddr.Line3, nil
			},
		},
		"AddressLine4": {
			Params: fibery.Field{
				Name: "Address Line 4",
				Type: fibery.Text,
			},
			Convert: func(sd integration.StandardData[quickbooks.Employee]) (any, error) {
				return sd.Item.PrimaryAddr.Line4, nil
			},
		},
		"AddressLine5": {
			Params: fibery.Field{
				Name: "Address Line 5",
				Type: fibery.Text,
			},
			Convert: func(sd integration.StandardData[quickbooks.Employee]) (any, error) {
				return sd.Item.PrimaryAddr.Line5, nil
			},
		},
		"AddressCity": {
			Params: fibery.Field{
				Name: "Address City",
				Type: fibery.Text,
			},
			Convert: func(sd integration.StandardData[quickbooks.Employee]) (any, error) {
				return sd.Item.PrimaryAddr.City, nil
			},
		},
		"AddressState": {
			Params: fibery.Field{
				Name: "Address State",
				Type: fibery.Text,
			},
			Convert: func(sd integration.StandardData[quickbooks.Employee]) (any, error) {
				return sd.Item.PrimaryAddr.CountrySubDivisionCode, nil
			},
		},
		"AddressPostalCode": {
			Params: fibery.Field{
				Name: "Address Postal Code",
				Type: fibery.Text,
			},
			Convert: func(sd integration.StandardData[quickbooks.Employee]) (any, error) {
				return sd.Item.PrimaryAddr.PostalCode, nil
			},
		},
		"AddressCountry": {
			Params: fibery.Field{
				Name: "Address Country",
				Type: fibery.Text,
			},
			Convert: func(sd integration.StandardData[quickbooks.Employee]) (any, error) {
				return sd.Item.PrimaryAddr.Country, nil
			},
		},
		"AddressLat": {
			Params: fibery.Field{
				Name: "Address Latitude",
				Type: fibery.Text,
			},
			Convert: func(sd integration.StandardData[quickbooks.Employee]) (any, error) {
				return sd.Item.PrimaryAddr.Lat, nil
			},
		},
		"AddressLong": {
			Params: fibery.Field{
				Name: "Address Longitude",
				Type: fibery.Text,
			},
			Convert: func(sd integration.StandardData[quickbooks.Employee]) (any, error) {
				return sd.Item.PrimaryAddr.Long, nil
			},
		},
	},
	nil,
)

func init() {
	integration.Types.Register(employee)
}

