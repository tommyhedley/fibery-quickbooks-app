package data

import (
	"fmt"

	"github.com/tommyhedley/fibery-quickbooks-app/pkgs/fibery"
	"github.com/tommyhedley/quickbooks-go"
)

var Employee = QuickBooksDualType{
	QuickBooksType: QuickBooksType{
		fiberyType: fiberyType{
			id:   "Employee",
			name: "Employee",
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
				"PrimaryEmailAddr": {
					Name:    "Email",
					Type:    fibery.Text,
					SubType: fibery.Email,
				},
				"BillableTime": {
					Name:        "Billable",
					Type:        fibery.Text,
					SubType:     fibery.Boolean,
					Description: "Is the entity enabled for use in QuickBooks?",
				},
				"BirthDate": {
					Name:    "Date of Birth",
					Type:    fibery.DateType,
					SubType: fibery.Day,
				},
				"PrimaryPhone": {
					Name: "Phone",
					Type: fibery.Text,
					Format: map[string]any{
						"format": "phone",
					},
				},
				"CostRate": {
					Name: "Cost Rate",
					Type: fibery.Number,
					Format: map[string]any{
						"format":               "Money",
						"currencyCode":         "USD",
						"hasThousandSeperator": true,
						"precision":            2,
					},
				},
				"BillRate": {
					Name: "Bill Rate",
					Type: fibery.Number,
					Format: map[string]any{
						"format":               "Money",
						"currencyCode":         "USD",
						"hasThousandSeperator": true,
						"precision":            2,
					},
				},
				"EmployeeNumber": {
					Name: "Employee ID",
					Type: fibery.Text,
				},
				"AddressLine1": {
					Name: "Address Line 1",
					Type: fibery.Text,
				},
				"AddressLine2": {
					Name: "Address Line 2",
					Type: fibery.Text,
				},
				"AddressLine3": {
					Name: "Address Line 3",
					Type: fibery.Text,
				},
				"AddressLine4": {
					Name: "Address Line 4",
					Type: fibery.Text,
				},
				"AddressLine5": {
					Name: "Address Line 5",
					Type: fibery.Text,
				},
				"AddressCity": {
					Name: "Address City",
					Type: fibery.Text,
				},
				"AddressState": {
					Name: "Address State",
					Type: fibery.Text,
				},
				"AddressPostalCode": {
					Name: "Address Postal Code",
					Type: fibery.Text,
				},
				"AddressCountry": {
					Name: "Address Country",
					Type: fibery.Text,
				},
				"AddressLat": {
					Name: "Address Latitude",
					Type: fibery.Text,
				},
				"AddressLong": {
					Name: "Address Longitude",
					Type: fibery.Text,
				},
			},
		},
		schemaGen: func(entity any) (map[string]any, error) {
			employee, ok := entity.(quickbooks.Employee)
			if !ok {
				return nil, fmt.Errorf("unable to convert entity to Employee")
			}

			return map[string]any{
				"Id":                employee.Id,
				"QBOId":             employee.Id,
				"DisplayName":       employee.DisplayName,
				"SyncToken":         employee.SyncToken,
				"__syncAction":      fibery.SET,
				"Active":            employee.Active,
				"Title":             employee.Title,
				"GivenName":         employee.GivenName,
				"MiddleName":        employee.MiddleName,
				"FamilyName":        employee.FamilyName,
				"Suffix":            employee.Suffix,
				"PrimaryEmailAddr":  employee.PrimaryEmailAddr,
				"BillableTime":      employee.BillableTime,
				"BirthDate":         employee.BirthDate.Format(fibery.DateFormat),
				"PrimaryPhone":      employee.PrimaryPhone.FreeFormNumber,
				"CostRate":          employee.CostRate,
				"BillRate":          employee.BillRate,
				"EmployeeNumber":    employee.EmployeeNumber,
				"AddressLine1":      employee.PrimaryAddr.Line1,
				"AddressLine2":      employee.PrimaryAddr.Line2,
				"AddressLine3":      employee.PrimaryAddr.Line3,
				"AddressLine4":      employee.PrimaryAddr.Line4,
				"AddressLine5":      employee.PrimaryAddr.Line1,
				"AddressCity":       employee.PrimaryAddr.City,
				"AddressState":      employee.PrimaryAddr.CountrySubDivisionCode,
				"AddressPostalCode": employee.PrimaryAddr.PostalCode,
				"AddressCountry":    employee.PrimaryAddr.Country,
				"AddressLat":        employee.PrimaryAddr.Lat,
				"AddressLong":       employee.PrimaryAddr.Long,
			}, nil
		},
		query:          func(req Request) (Response, error) {},
		queryProcessor: func(entityArray any, schemaGen schemaGenFunc) ([]map[string]any, error) {},
	},
	cdcProcessor: func(cdc quickbooks.ChangeDataCapture, schemaGen schemaGenFunc) ([]map[string]any, error) {},
	whBatchProcessor: func(itemResponse quickbooks.BatchItemResponse, response *map[string][]map[string]any, cache *cache.Cache, realmId string, queryProcessor queryProcessorFunc, schemaGen schemaGenFunc, typeId string) error {
	},
}
