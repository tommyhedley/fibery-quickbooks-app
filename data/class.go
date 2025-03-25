package data

import (
	"github.com/tommyhedley/fibery-quickbooks-app/pkgs/fibery"
	"github.com/tommyhedley/quickbooks-go"
)

var Class = QuickBooksDualType[quickbooks.Class]{
	QuickBooksType: QuickBooksType[quickbooks.Class]{
		BaseType: fibery.BaseType{
			TypeId:   "Class",
			TypeName: "Class",
			TypeSchema: map[string]fibery.Field{
				"id": {
					Name: "ID",
					Type: fibery.Id,
				},
				"QBOId": {
					Name: "QBOId",
					Type: fibery.Text,
				},
				"Name": {
					Name: "Base Name",
					Type: fibery.Text,
				},
				"FullyQualifiedName": {
					Name:    "Full Name",
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
				"ParentClassId": {
					Name: "Parent Class ID",
					Type: fibery.Text,
					Relation: &fibery.Relation{
						Cardinality:   fibery.MTO,
						Name:          "Parent Class",
						TargetName:    "Sub-Classes",
						TargetType:    "Class",
						TargetFieldID: "id",
					},
				},
			},
		},
		schemaGen: func(c quickbooks.Class) (map[string]any, error) {
			return map[string]any{
				"id":                 c.Id,
				"QBOId":              c.Id,
				"Name":               c.Name,
				"FullyQualifiedName": c.FullyQualifiedName,
				"SyncToken":          c.SyncToken,
				"__syncAction":       fibery.SET,
				"Active":             c.Active,
				"ParentClassId":      c.ParentRef.Value,
			}, nil
		},
		pageQuery: func(req Request) ([]quickbooks.Class, error) {
			params := quickbooks.RequestParameters{
				Ctx:     req.Ctx,
				RealmId: req.RealmId,
				Token:   req.Token,
			}

			items, err := req.Client.FindClassesByPage(params, req.StartPosition, req.PageSize)
			if err != nil {
				return nil, err
			}

			return items, nil
		},
	},
	entityId: func(c quickbooks.Class) string {
		return c.Id
	},
	entityStatus: func(c quickbooks.Class) string {
		return c.Status
	},
}
