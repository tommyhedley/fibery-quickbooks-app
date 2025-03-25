package data

import (
	"github.com/tommyhedley/fibery-quickbooks-app/pkgs/fibery"
	"github.com/tommyhedley/quickbooks-go"
)

var CustomerType = QuickBooksCDCType[quickbooks.CustomerType]{
	QuickBooksType: QuickBooksType[quickbooks.CustomerType]{
		BaseType: fibery.BaseType{
			TypeId:   "CustomerType",
			TypeName: "Customer Type",
			TypeSchema: map[string]fibery.Field{
				"id": {
					Name: "ID",
					Type: fibery.Text,
				},
				"QBOId": {
					Name: "QBO ID",
					Type: fibery.Text,
				},
				"Name": {
					Name:    "Name",
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
			},
		},
		schemaGen: func(ct quickbooks.CustomerType) (map[string]any, error) {
			return map[string]any{
				"id":           ct.Id,
				"QBOId":        ct.Id,
				"Name":         ct.Name,
				"SyncToken":    ct.SyncToken,
				"__syncAction": fibery.SET,
				"Active":       ct.Active,
			}, nil
		},
		pageQuery: func(req Request) ([]quickbooks.CustomerType, error) {
			params := quickbooks.RequestParameters{
				Ctx:     req.Ctx,
				RealmId: req.RealmId,
				Token:   req.Token,
			}

			items, err := req.Client.FindCustomerTypesByPage(params, req.StartPosition, req.PageSize)
			if err != nil {
				return nil, err
			}

			return items, nil

		},
	},
	entityId: func(ct quickbooks.CustomerType) string {
		return ct.Id
	},
	entityStatus: func(ct quickbooks.CustomerType) string {
		return ct.Status
	},
}

func init() {
	registerType(&CustomerType)
}
