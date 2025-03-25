package data

import (
	"github.com/tommyhedley/fibery-quickbooks-app/pkgs/fibery"
	"github.com/tommyhedley/quickbooks-go"
)

var Term = QuickBooksDualType[quickbooks.Term]{
	QuickBooksType: QuickBooksType[quickbooks.Term]{
		BaseType: fibery.BaseType{
			TypeId:   "Term",
			TypeName: "Term",
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
				"Active": {
					Name:    "Active",
					Type:    fibery.Text,
					SubType: fibery.Boolean,
				},
			},
		},
		schemaGen: func(t quickbooks.Term) (map[string]any, error) {
			return map[string]any{
				"id":        t.Id,
				"QBOId":     t.Id,
				"Name":      t.Name,
				"SyncToken": t.SyncToken,
				"Active":    t.Active,
			}, nil
		},
		pageQuery: func(req Request) ([]quickbooks.Term, error) {
			params := quickbooks.RequestParameters{
				Ctx:     req.Ctx,
				RealmId: req.RealmId,
				Token:   req.Token,
			}

			items, err := req.Client.FindTermsByPage(params, req.StartPosition, req.PageSize)
			if err != nil {
				return nil, err
			}

			return items, nil
		},
	},
}
