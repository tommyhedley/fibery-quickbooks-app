package data

import (
	"fmt"

	"github.com/tommyhedley/fibery-quickbooks-app/pkgs/fibery"
)

var TaxExemption = QuickBooksType{
	fiberyType: fiberyType{
		id:   "TaxExemption",
		name: "Tax Exemption",
		schema: map[string]fibery.Field{
			"id": {
				Name: "id",
				Type: fibery.ID,
			},
			"name": {
				Name: "Name",
				Type: fibery.Text,
			},
		},
	},
	schemaGen: func(entity any) (map[string]any, error) {
		qualifiedEntity, ok := entity.(map[string]any)
		if !ok {
			return nil, fmt.Errorf("entity is not in the format map[string]any")
		}
		return qualifiedEntity, nil
	},
	query: func(req Request) (Response, error) {
		taxExemptions := []map[string]any{
			{
				"id":   "1",
				"name": "Federal government",
			},
			{
				"id":   "2",
				"name": "State government",
			},
			{
				"id":   "3",
				"name": "Local government",
			},
			{
				"id":   "4",
				"name": "Tribal government",
			},
			{
				"id":   "5",
				"name": "Charitable organization",
			},
			{
				"id":   "6",
				"name": "Religious organization",
			},
			{
				"id":   "7",
				"name": "Educational organization",
			},
			{
				"id":   "8",
				"name": "Hospital",
			},
			{
				"id":   "9",
				"name": "Resale",
			},
			{
				"id":   "10",
				"name": "Direct pay permit",
			},
			{
				"id":   "11",
				"name": "Multiple points of use",
			},
			{
				"id":   "12",
				"name": "Direct mail",
			},
			{
				"id":   "13",
				"name": "Agricultural production",
			},
			{
				"id":   "14",
				"name": "Industrial production / manufacturing",
			},
			{
				"id":   "15",
				"name": "Foreign diplomat",
			},
		}
		return Response{Data: taxExemptions, MoreData: false}, nil
	},
	queryProcessor: func(entityArray any, schemaGen schemaGenFunc) ([]map[string]any, error) {
		entities, ok := entityArray.([]map[string]any)
		if !ok {
			return nil, fmt.Errorf("entityArray is not in the format []map[string]any")
		}
		return entities, nil
	},
}

func init() {
	RegisterType(TaxExemption)
}
