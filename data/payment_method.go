package data

import (
	"github.com/tommyhedley/fibery-quickbooks-app/pkgs/fibery"
)
import "github.com/tommyhedley/quickbooks-go"

var PaymentMethod = QuickBooksDualType[quickbooks.PaymentMethod]{
	QuickBooksType: QuickBooksType[quickbooks.PaymentMethod]{
		BaseType: fibery.BaseType{
			TypeId:   "PaymentMethod",
			TypeName: "Payment Method",
			TypeSchema: map[string]fibery.Field{
				"id": {
					Name: "ID",
					Type: fibery.Id,
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
				"Type": {
					Name:     "Type",
					Type:     fibery.Text,
					SubType:  fibery.SingleSelect,
					ReadOnly: true,
					Options: []map[string]any{
						{
							"name": "Credit Card",
						},
						{
							"name": "Non-Credit Card",
						},
					},
				},
			},
		},
		schemaGen: func(pm quickbooks.PaymentMethod) (map[string]any, error) {
			var paymentType string
			switch pm.Type {
			case "CREDIT_CARD":
				paymentType = "Credit Card"
			case "NON_CREDIT_CARD":
				paymentType = "Non-Credit Card"
			}
			return map[string]any{
				"id":          pm.Id,
				"QBOId":       pm.Id,
				"Name":        pm.Name,
				"SyncToken":   pm.SyncToken,
				"__syncToken": fibery.SET,
				"Active":      pm.Active,
				"Type":        paymentType,
			}, nil
		},
		pageQuery: func(req Request) ([]quickbooks.PaymentMethod, error) {
			params := quickbooks.RequestParameters{
				Ctx:     req.Ctx,
				RealmId: req.RealmId,
				Token:   req.Token,
			}

			items, err := req.Client.FindPaymentMethodsByPage(params, req.StartPosition, req.PageSize)
			if err != nil {
				return nil, err
			}

			return items, nil
		},
	},
	entityId: func(pm quickbooks.PaymentMethod) string {
		return pm.Id
	},
	entityStatus: func(pm quickbooks.PaymentMethod) string {
		return pm.Status
	},
}

func init() {
	registerType(&PaymentMethod)
}
