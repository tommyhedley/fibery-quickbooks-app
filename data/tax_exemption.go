package data

import (
	"github.com/tommyhedley/fibery-quickbooks-app/pkgs/fibery"
)

type TaxExemptionEntity struct {
	Id   string
	Name string
}

var TaxExemption = QuickBooksType[TaxExemptionEntity]{
	BaseType: fibery.BaseType{
		TypeId:   "TaxExemption",
		TypeName: "Tax Exemption",
		TypeSchema: map[string]fibery.Field{
			"id": {
				Name: "ID",
				Type: fibery.Id,
			},
			"Name": {
				Name:    "Name",
				Type:    fibery.Text,
				SubType: fibery.Title,
			},
		},
	},
	schemaGen: func(te TaxExemptionEntity) (map[string]any, error) {
		return map[string]any{
			"id":   te.Id,
			"Name": te.Name,
		}, nil
	},
	pageQuery: func(req Request) ([]TaxExemptionEntity, error) {
		return []TaxExemptionEntity{
			{
				Id:   "1",
				Name: "Federal government",
			},
			{
				Id:   "2",
				Name: "State government",
			},
			{
				Id:   "3",
				Name: "Local government",
			},
			{
				Id:   "4",
				Name: "Tribal government",
			},
			{
				Id:   "5",
				Name: "Charitable organization",
			},
			{
				Id:   "6",
				Name: "Religious organization",
			},
			{
				Id:   "7",
				Name: "Educational organization",
			},
			{
				Id:   "8",
				Name: "Hospital",
			},
			{
				Id:   "9",
				Name: "Resale",
			},
			{
				Id:   "10",
				Name: "Direct pay permit",
			},
			{
				Id:   "11",
				Name: "Multiple points of use",
			},
			{
				Id:   "12",
				Name: "Direct mail",
			},
			{
				Id:   "13",
				Name: "Agricultural production",
			},
			{
				Id:   "14",
				Name: "Industrial production / manufacturing",
			},
			{
				Id:   "15",
				Name: "Foreign diplomat",
			},
		}, nil
	},
}

func init() {
	registerType(&TaxExemption)
}
