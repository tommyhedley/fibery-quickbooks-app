package types

import (
	"github.com/tommyhedley/fibery-quickbooks-app/pkgs/fibery"
	"github.com/tommyhedley/fibery-quickbooks-app/pkgs/integration"
	"github.com/tommyhedley/quickbooks-go"
)

var employee = integration.NewDualType(
	"Employee",
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
		"qboId": {
			Params: fibery.Field{
				Name:     "QBO ID",
				Type:     fibery.Text,
				ReadOnly: true,
			},
			Convert: func(sd integration.StandardData[quickbooks.Employee]) (any, error) {
				return sd.Item.Id, nil
			},
		},
		"displayName": {
			Params: fibery.Field{
				Name:    "Display Name",
				Type:    fibery.Text,
				SubType: fibery.Title,
			},
			Convert: func(sd integration.StandardData[quickbooks.Employee]) (any, error) {
				return sd.Item.DisplayName, nil
			},
		},
		"syncToken": {
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
		"active": {
			Params: fibery.Field{
				Name:    "Active",
				Type:    fibery.Text,
				SubType: fibery.Boolean,
			},
			Convert: func(sd integration.StandardData[quickbooks.Employee]) (any, error) {
				return sd.Item.Active, nil
			},
		},
		"title": {
			Params: fibery.Field{
				Name: "Title",
				Type: fibery.Text,
			},
			Convert: func(sd integration.StandardData[quickbooks.Employee]) (any, error) {
				return sd.Item.Title, nil
			},
		},
		"givenName": {
			Params: fibery.Field{
				Name: "First Name",
				Type: fibery.Text,
			},
			Convert: func(sd integration.StandardData[quickbooks.Employee]) (any, error) {
				return sd.Item.GivenName, nil
			},
		},
		"middleName": {
			Params: fibery.Field{
				Name: "Middle Name",
				Type: fibery.Text,
			},
			Convert: func(sd integration.StandardData[quickbooks.Employee]) (any, error) {
				return sd.Item.MiddleName, nil
			},
		},
		"familyName": {
			Params: fibery.Field{
				Name: "Last Name",
				Type: fibery.Text,
			},
			Convert: func(sd integration.StandardData[quickbooks.Employee]) (any, error) {
				return sd.Item.FamilyName, nil
			},
		},
		"suffix": {
			Params: fibery.Field{
				Name: "Suffix",
				Type: fibery.Text,
			},
			Convert: func(sd integration.StandardData[quickbooks.Employee]) (any, error) {
				return sd.Item.Suffix, nil
			},
		},
		"primaryEmailAddr": {
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
		"billableTime": {
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
		"birthDate": {
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
		"primaryPhone": {
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
		"mobile": {
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
		"costRate": {
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
		"billRate": {
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
		"employeeNumber": {
			Params: fibery.Field{
				Name: "Employee ID",
				Type: fibery.Text,
			},
			Convert: func(sd integration.StandardData[quickbooks.Employee]) (any, error) {
				return sd.Item.EmployeeNumber, nil
			},
		},
		"addressLine1": {
			Params: fibery.Field{
				Name: "Address Line 1",
				Type: fibery.Text,
			},
			Convert: func(sd integration.StandardData[quickbooks.Employee]) (any, error) {
				return sd.Item.PrimaryAddr.Line1, nil
			},
		},
		"addressLine2": {
			Params: fibery.Field{
				Name: "Address Line 2",
				Type: fibery.Text,
			},
			Convert: func(sd integration.StandardData[quickbooks.Employee]) (any, error) {
				return sd.Item.PrimaryAddr.Line2, nil
			},
		},
		"addressLine3": {
			Params: fibery.Field{
				Name: "Address Line 3",
				Type: fibery.Text,
			},
			Convert: func(sd integration.StandardData[quickbooks.Employee]) (any, error) {
				return sd.Item.PrimaryAddr.Line3, nil
			},
		},
		"addressLine4": {
			Params: fibery.Field{
				Name: "Address Line 4",
				Type: fibery.Text,
			},
			Convert: func(sd integration.StandardData[quickbooks.Employee]) (any, error) {
				return sd.Item.PrimaryAddr.Line4, nil
			},
		},
		"addressLine5": {
			Params: fibery.Field{
				Name: "Address Line 5",
				Type: fibery.Text,
			},
			Convert: func(sd integration.StandardData[quickbooks.Employee]) (any, error) {
				return sd.Item.PrimaryAddr.Line5, nil
			},
		},
		"addressCity": {
			Params: fibery.Field{
				Name: "Address City",
				Type: fibery.Text,
			},
			Convert: func(sd integration.StandardData[quickbooks.Employee]) (any, error) {
				return sd.Item.PrimaryAddr.City, nil
			},
		},
		"addressState": {
			Params: fibery.Field{
				Name: "Address State",
				Type: fibery.Text,
			},
			Convert: func(sd integration.StandardData[quickbooks.Employee]) (any, error) {
				return sd.Item.PrimaryAddr.CountrySubDivisionCode, nil
			},
		},
		"addressPostalCode": {
			Params: fibery.Field{
				Name: "Address Postal Code",
				Type: fibery.Text,
			},
			Convert: func(sd integration.StandardData[quickbooks.Employee]) (any, error) {
				return sd.Item.PrimaryAddr.PostalCode, nil
			},
		},
		"addressCountry": {
			Params: fibery.Field{
				Name: "Address Country",
				Type: fibery.Text,
			},
			Convert: func(sd integration.StandardData[quickbooks.Employee]) (any, error) {
				return sd.Item.PrimaryAddr.Country, nil
			},
		},
		"addressLat": {
			Params: fibery.Field{
				Name: "Address Latitude",
				Type: fibery.Text,
			},
			Convert: func(sd integration.StandardData[quickbooks.Employee]) (any, error) {
				return sd.Item.PrimaryAddr.Lat, nil
			},
		},
		"addressLong": {
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
