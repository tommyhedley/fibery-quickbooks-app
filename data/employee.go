package data

import (
	"github.com/tommyhedley/fibery-quickbooks-app/pkgs/fibery"
	"github.com/tommyhedley/quickbooks-go"
)

var Employee = QuickBooksDualType[quickbooks.Employee]{
	QuickBooksType: QuickBooksType[quickbooks.Employee]{
		BaseType: fibery.BaseType{
			TypeId:   "Employee",
			TypeName: "Employee",
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
				"Mobile": {
					Name: "Mobile",
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
		schemaGen: func(e quickbooks.Employee) (map[string]any, error) {
			var email string
			if e.PrimaryEmailAddr != nil {
				email = e.PrimaryEmailAddr.Address
			}

			var primaryPhone string
			if e.PrimaryPhone != nil {
				primaryPhone = e.PrimaryPhone.FreeFormNumber
			}

			var mobile string
			if e.Mobile != nil {
				mobile = e.Mobile.FreeFormNumber
			}

			return map[string]any{
				"id":                e.Id,
				"QBOId":             e.Id,
				"DisplayName":       e.DisplayName,
				"SyncToken":         e.SyncToken,
				"__syncAction":      fibery.SET,
				"Active":            e.Active,
				"Title":             e.Title,
				"GivenName":         e.GivenName,
				"MiddleName":        e.MiddleName,
				"FamilyName":        e.FamilyName,
				"Suffix":            e.Suffix,
				"PrimaryEmailAddr":  email,
				"BillableTime":      e.BillableTime,
				"BirthDate":         e.BirthDate.Format(fibery.DateFormat),
				"PrimaryPhone":      primaryPhone,
				"Mobile":            mobile,
				"CostRate":          e.CostRate,
				"BillRate":          e.BillRate,
				"EmployeeNumber":    e.EmployeeNumber,
				"AddressLine1":      e.PrimaryAddr.Line1,
				"AddressLine2":      e.PrimaryAddr.Line2,
				"AddressLine3":      e.PrimaryAddr.Line3,
				"AddressLine4":      e.PrimaryAddr.Line4,
				"AddressLine5":      e.PrimaryAddr.Line1,
				"AddressCity":       e.PrimaryAddr.City,
				"AddressState":      e.PrimaryAddr.CountrySubDivisionCode,
				"AddressPostalCode": e.PrimaryAddr.PostalCode,
				"AddressCountry":    e.PrimaryAddr.Country,
				"AddressLat":        e.PrimaryAddr.Lat,
				"AddressLong":       e.PrimaryAddr.Long,
			}, nil
		},
		pageQuery: func(req Request) ([]quickbooks.Employee, error) {
			params := quickbooks.RequestParameters{
				Ctx:     req.Ctx,
				RealmId: req.RealmId,
				Token:   req.Token,
			}

			items, err := req.Client.FindEmployeesByPage(params, req.StartPosition, req.PageSize)
			if err != nil {
				return nil, err
			}

			return items, nil
		},
	},
	entityId: func(e quickbooks.Employee) string {
		return e.Id
	},
	entityStatus: func(e quickbooks.Employee) string {
		return e.Status
	},
}

func init() {
	registerType(&Employee)
}
