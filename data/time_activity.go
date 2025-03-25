package data

import (
	"github.com/tommyhedley/fibery-quickbooks-app/pkgs/fibery"
	"github.com/tommyhedley/quickbooks-go"
)

var TimeActivity = QuickBooksWHType[quickbooks.TimeActivity]{
	QuickBooksType: QuickBooksType[quickbooks.TimeActivity]{
		BaseType: fibery.BaseType{
			TypeId:   "TimeActivity",
			TypeName: "Time Activity",
			TypeSchema: map[string]fibery.Field{
				"id": {
					Name: "ID",
					Type: fibery.Id,
				},
				"QBOId": {
					Name: "QBO ID",
					Type: fibery.Text,
				},
				"Description": {
					Name:    "Description",
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
				"ActivityType": {
					Name:     "Activity Type",
					Type:     fibery.Text,
					SubType:  fibery.SingleSelect,
					ReadOnly: true,
					Options: []map[string]any{
						{
							"name": "Employee",
						},
						{
							"name": "Vendor",
						},
					},
				},
				"TxnDate": {
					Name: "Invoice Date",
					Type: fibery.DateType,
				},
				"Hours": {
					Name: "Hours",
					Type: fibery.Number,
					Format: map[string]any{
						"format":               "Number",
						"currencyCode":         "h",
						"hasThousandSeperator": true,
						"precision":            0,
					},
				},
				"Minutes": {
					Name: "Minutes",
					Type: fibery.Number,
					Format: map[string]any{
						"format":               "Number",
						"currencyCode":         "m",
						"hasThousandSeperator": true,
						"precision":            0,
					},
				},
				"BreakHours": {
					Name: "Break Hours",
					Type: fibery.Number,
					Format: map[string]any{
						"format":               "Number",
						"currencyCode":         "h",
						"hasThousandSeperator": true,
						"precision":            0,
					},
				},
				"BreakMinutes": {
					Name: "Break Minutes",
					Type: fibery.Number,
					Format: map[string]any{
						"format":               "Number",
						"currencyCode":         "m",
						"hasThousandSeperator": true,
						"precision":            0,
					},
				},
				"StartTime": {
					Name: "Start Time",
					Type: fibery.DateType,
				},
				"EndTime": {
					Name: "End Time",
					Type: fibery.DateType,
				},
				"HourlyRate": {
					Name: "Rate",
					Type: fibery.Number,
					Format: map[string]any{
						"format":               "Money",
						"currencyCode":         "USD",
						"hasThousandSeperator": true,
						"precision":            2,
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
				"Taxable": {
					Name:    "Taxable",
					Type:    fibery.Text,
					SubType: fibery.Boolean,
				},
				"Billable": {
					Name:    "Billable",
					Type:    fibery.Text,
					SubType: fibery.Boolean,
				},
				"Billed": {
					Name:    "Billed",
					Type:    fibery.Text,
					SubType: fibery.Boolean,
				},
				"VendorId": {
					Name: "Vendor ID",
					Type: fibery.Text,
					Relation: &fibery.Relation{
						Cardinality:   fibery.MTO,
						Name:          "Vendor",
						TargetName:    "Time Activity",
						TargetType:    "Vendor",
						TargetFieldID: "id",
					},
				},
				"EmployeeId": {
					Name: "Employee ID",
					Type: fibery.Text,
					Relation: &fibery.Relation{
						Cardinality:   fibery.MTO,
						Name:          "Employee",
						TargetName:    "Time Activity",
						TargetType:    "Employee",
						TargetFieldID: "id",
					},
				},
				"CustomerId": {
					Name: "Customer ID",
					Type: fibery.Text,
					Relation: &fibery.Relation{
						Cardinality:   fibery.MTO,
						Name:          "Customer",
						TargetName:    "Time Activity",
						TargetType:    "Customer",
						TargetFieldID: "id",
					},
				},
				"ClassId": {
					Name: "Class ID",
					Type: fibery.Text,
					Relation: &fibery.Relation{
						Cardinality:   fibery.MTO,
						Name:          "Class",
						TargetName:    "Time Activity",
						TargetType:    "Class",
						TargetFieldID: "id",
					},
				},
				"ItemId": {
					Name: "Item ID",
					Type: fibery.Text,
					Relation: &fibery.Relation{
						Cardinality:   fibery.MTO,
						Name:          "Item",
						TargetName:    "Time Activity",
						TargetType:    "Item",
						TargetFieldID: "id",
					},
				},
			},
		},
		schemaGen: func(ta quickbooks.TimeActivity) (map[string]any, error) {
			var billable bool
			switch ta.BillableStatus {
			case quickbooks.BillableStatusType:
				billable = true
			case quickbooks.HasBeenBilledStatusType:
				billable = true
			default:
				billable = false
			}

			billed := false
			if ta.BillableStatus == quickbooks.HasBeenBilledStatusType {
				billed = true
			}

			var classId string
			if ta.ClassRef != nil {
				classId = ta.ClassRef.Value
			}

			var itemId string
			if ta.ItemRef != nil {
				itemId = ta.ItemRef.Value
			}

			return map[string]any{
				"id":           ta.Id,
				"QBOId":        ta.Id,
				"Description":  ta.Description,
				"SyncToken":    ta.SyncToken,
				"__syncAction": fibery.SET,
				"ActivityType": ta.NameOf,
				"TxnDate":      ta.TxnDate.Format(fibery.DateFormat),
				"Hours":        ta.Hours,
				"Minutes":      ta.Minutes,
				"BreakHours":   ta.BreakHours,
				"BreakMinutes": ta.BreakMinutes,
				"StartTime":    ta.StartTime.Format(fibery.DateFormat),
				"EndTime":      ta.EndTime.Format(fibery.DateFormat),
				"HourlyRate":   ta.HourlyRate,
				"CostRate":     ta.CostRate,
				"Taxable":      ta.Taxable,
				"Billable":     billable,
				"Billed":       billed,
				"VendorId":     ta.VendorRef.Value,
				"EmployeeId":   ta.EmployeeRef.Value,
				"CustomerId":   ta.CustomerRef.Value,
				"ClassId":      classId,
				"ItemId":       itemId,
			}, nil
		},
		pageQuery: func(req Request) ([]quickbooks.TimeActivity, error) {
			params := quickbooks.RequestParameters{
				Ctx:     req.Ctx,
				RealmId: req.RealmId,
				Token:   req.Token,
			}

			items, err := req.Client.FindTimeActivitiesByPage(params, req.StartPosition, req.PageSize)
			if err != nil {
				return nil, err
			}

			return items, nil
		},
	},
	entityId: func(ta quickbooks.TimeActivity) string {
		return ta.Id
	},
}

func init() {
	registerType(&TimeActivity)
}
