package data

import (
	"fmt"
	"github.com/tommyhedley/fibery-quickbooks-app/pkgs/fibery"
	"github.com/tommyhedley/quickbooks-go"
)

var Entity = QuickBooksDualType{
	QuickBooksType: QuickBooksType{
		fiberyType: fiberyType{
			id:   "Entity",
			name: "Entity",
			schema: map[string]fibery.Field{
				"Id": {
					Name: "ID",
					Type: fibery.ID,
				},
				"QBOId": {
					Name: "QBO ID",
					Type: fibery.Text,
				},
				"Name": {
					Name: "Name",
					Type: fibery.Text,
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
				"CustomerId": {
					Name: "Customer ID",
					Type: fibery.Text,
					Relation: &fibery.Relation{
						Cardinality:   fibery.OTO,
						Name:          "Customer",
						TargetName:    "Entity",
						TargetType:    "Customer",
						TargetFieldID: "Id",
					},
				},
				"VendorId": {
					Name: "Vendor ID",
					Type: fibery.Text,
					Relation: &fibery.Relation{
						Cardinality:   fibery.OTO,
						Name:          "Vendor",
						TargetName:    "Entity",
						TargetType:    "Vendor",
						TargetFieldID: "Id",
					},
				},
				"EmployeeId": {
					Name: "Sales Term ID",
					Type: fibery.Text,
					Relation: &fibery.Relation{
						Cardinality:   fibery.MTO,
						Name:          "Employee",
						TargetName:    "Entity",
						TargetType:    "Employee",
						TargetFieldID: "Id",
					},
				},
			},
		},
		schemaGen: func(entity any) (map[string]any, error) {
			switch dataType := entity.(type) {
			case quickbooks.Customer:
				return map[string]any{
					"Id":           fmt.Sprintf("c:%s", dataType.Id),
					"QBOId":        dataType.Id,
					"Name":         dataType.DisplayName,
					"SyncToken":    dataType.SyncToken,
					"__syncAction": fibery.SET,
					"CustomerId":   dataType.Id,
				}, nil
			case quickbooks.Vendor:
				return map[string]any{
					"Id":           fmt.Sprintf("v:%s", dataType.Id),
					"QBOId":        dataType.Id,
					"Name":         dataType.DisplayName,
					"SyncToken":    dataType.SyncToken,
					"__syncAction": fibery.SET,
					"CustomerId":   dataType.Id,
				}, nil
			case quickbooks.Employee:
				return map[string]any{
					"Id":           fmt.Sprintf("e:%s", dataType.Id),
					"QBOId":        dataType.Id,
					"Name":         dataType.DisplayName,
					"SyncToken":    dataType.SyncToken,
					"__syncAction": fibery.SET,
					"CustomerId":   dataType.Id,
				}, nil
			default:
				return nil, fmt.Errorf("enitity was not one of: Customer, Vendor, Employee")
			}
		},
		query:          func(req Request) (Response, error) {},
		queryProcessor: func(entityArray any, schemaGen schemaGenFunc) ([]map[string]any, error) {},
	},
	cdcProcessor: func(cdc quickbooks.ChangeDataCapture, schemaGen schemaGenFunc) ([]map[string]any, error) {},
	whBatchProcessor: func(itemResponse quickbooks.BatchItemResponse, response *map[string][]map[string]any, cache *cache.Cache, realmId string, queryProcessor queryProcessorFunc, schemaGen schemaGenFunc, typeId string) error {
	},
}
