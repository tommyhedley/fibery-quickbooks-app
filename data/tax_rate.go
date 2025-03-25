package data

import (
	"github.com/tommyhedley/fibery-quickbooks-app/pkgs/fibery"
	"github.com/tommyhedley/quickbooks-go"
)

var TaxRate = QuickBooksType[quickbooks.TaxRate]{
	BaseType: fibery.BaseType{
		TypeId:   "TaxRate",
		TypeName: "Tax Rate",
		TypeSchema: map[string]fibery.Field{
			"id": {
				Name: "Id",
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
			"Description": {
				Name: "Description",
				Type: fibery.Text,
			},
			"SyncToken": {
				Name:     "Sync Token",
				Type:     fibery.Text,
				ReadOnly: true,
			},
			"RateValue": {
				Name: "Rate",
				Type: fibery.Number,
				Format: map[string]any{
					"format":    "Percent",
					"precision": 2,
				},
			},
			"Active": {
				Name:    "Active",
				Type:    fibery.Text,
				SubType: fibery.Boolean,
			},
		},
	},
	schemaGen: func(tr quickbooks.TaxRate) (map[string]any, error) {
		return map[string]any{
			"id":          tr.Id,
			"QBOId":       tr.Id,
			"Name":        tr.Name,
			"Description": tr.Description,
			"SyncToken":   tr.SyncToken,
			"RateValue":   tr.RateValue,
			"Active":      tr.Active,
		}, nil
	},
	pageQuery: func(req Request) ([]quickbooks.TaxRate, error) {
		params := quickbooks.RequestParameters{
			Ctx:     req.Ctx,
			RealmId: req.RealmId,
			Token:   req.Token,
		}

		items, err := req.Client.FindTaxRatesByPage(params, req.StartPosition, req.PageSize)
		if err != nil {
			return nil, err
		}

		return items, nil
	},
}

func init() {
	registerType(&TaxRate)
}
