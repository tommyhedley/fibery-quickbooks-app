package data

import (
	"fmt"

	"github.com/tommyhedley/fibery-quickbooks-app/pkgs/fibery"
	"github.com/tommyhedley/quickbooks-go"
)

var TaxCode = QuickBooksType[quickbooks.TaxCode]{
	BaseType: fibery.BaseType{
		TypeId:   "TaxCode",
		TypeName: "Tax Code",
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
			"Description": {
				Name: "Description",
				Type: fibery.Text,
			},
			"SyncToken": {
				Name: "Sync Token",
				Type: fibery.Text,
			},
			"TaxGroup": {
				Name:    "Tax Group",
				Type:    fibery.Text,
				SubType: fibery.Boolean,
			},
			"Taxable": {
				Name:    "Taxable",
				Type:    fibery.Text,
				SubType: fibery.Boolean,
			},
			"Active": {
				Name:    "Active",
				Type:    fibery.Text,
				SubType: fibery.Boolean,
			},
			"Hidden": {
				Name:    "Hidden",
				Type:    fibery.Text,
				SubType: fibery.Boolean,
			},
			"TaxCodeType": {
				Name:     "Type",
				Type:     fibery.Text,
				SubType:  fibery.SingleSelect,
				ReadOnly: true,
				Options: []map[string]any{
					{
						"name": "User Defined",
					},
					{
						"name": "System Generated",
					},
				},
			},
		},
	},
	schemaGen: func(tc quickbooks.TaxCode) (map[string]any, error) {
		var taxCodeType string
		switch tc.TaxCodeConfigType {
		case "SYSTEM_GENERATED":
			taxCodeType = "System Generated"
		case "USER_DEFINED":
			taxCodeType = "User Defined"
		}
		return map[string]any{
			"id":          tc.Id,
			"QBOId":       tc.Id,
			"Name":        tc.Name,
			"Description": tc.Description,
			"SyncToken":   tc.SyncToken,
			"TaxGroup":    tc.TaxGroup,
			"Taxable":     tc.Taxable,
			"Active":      tc.Active,
			"Hidden":      tc.Hidden,
			"TaxCodeType": taxCodeType,
		}, nil
	},
	pageQuery: func(req Request) ([]quickbooks.TaxCode, error) {
		params := quickbooks.RequestParameters{
			Ctx:     req.Ctx,
			RealmId: req.RealmId,
			Token:   req.Token,
		}

		items, err := req.Client.FindTaxCodesByPage(params, req.StartPosition, req.PageSize)
		if err != nil {
			return nil, err
		}

		return items, nil
	},
}

var TaxCodeLine = DependentDataType[quickbooks.TaxCode]{
	dependentBaseType: dependentBaseType[quickbooks.TaxCode]{
		BaseType: fibery.BaseType{
			TypeId:   "TaxCodeLine",
			TypeName: "Tax Code Line",
			TypeSchema: map[string]fibery.Field{
				"id": {
					Name: "ID",
					Type: fibery.Text,
				},
				"Name": {
					Name:    "Name",
					Type:    fibery.Text,
					SubType: fibery.Title,
				},
				"TaxRateId": {
					Name: "Tax Rate ID",
					Type: fibery.Text,
					Relation: &fibery.Relation{
						Cardinality:   fibery.MTO,
						Name:          "Tax Rate",
						TargetName:    "Tax Code Lines",
						TargetType:    "TaxRate",
						TargetFieldID: "id",
					},
				},
				"TaxType": {
					Name:     "Tax Type",
					Type:     fibery.Text,
					SubType:  fibery.SingleSelect,
					ReadOnly: true,
					Options: []map[string]any{
						{
							"name": "Tax On Amount",
						},
						{
							"name": "Tax On Amount Plus Tax",
						},
						{
							"name": "Tax On Tax",
						},
					},
				},
				"TaxOrder": {
					Name: "Tax Order",
					Type: fibery.Number,
				},
				"TaxCodeIdPurchase": {
					Name: "Purchase Tax Code ID",
					Type: fibery.Text,
					Relation: &fibery.Relation{
						Cardinality:   fibery.MTO,
						Name:          "Purchase Tax Code",
						TargetName:    "Purchase Tax Rates",
						TargetType:    "TaxCode",
						TargetFieldID: "id",
					},
				},
				"TaxCodeIdSale": {
					Name: "Sale Tax Code ID",
					Type: fibery.Text,
					Relation: &fibery.Relation{
						Cardinality:   fibery.MTO,
						Name:          "Sales Tax Code",
						TargetName:    "Sales Tax Rates",
						TargetType:    "TaxCode",
						TargetFieldID: "id",
					},
				},
			},
		},
		schemaGen: func(tc quickbooks.TaxCode) ([]map[string]any, error) {
			items := []map[string]any{}
			for _, ptRate := range tc.PurchaseTaxRateList.TaxRateDetail {
				var taxType string
				switch ptRate.TaxTypeApplicable {
				case "TaxOnAmount":
					taxType = "Tax On Amount"
				case "TaxOnAmountPlusTax":
					taxType = "Tax On Amount Plus Tax"
				case "TaxOnTax":
					taxType = "Tax On Tax"
				}
				item := map[string]any{
					"id":                fmt.Sprintf("pt:%s", ptRate.TaxOrder.String()),
					"Name":              ptRate.TaxRateRef.Name,
					"TaxRateId":         ptRate.TaxRateRef.Value,
					"TaxType":           taxType,
					"TaxOrder":          ptRate.TaxOrder,
					"TaxCodeIdPurchase": tc.Id,
				}
				items = append(items, item)
			}
			for _, stRate := range tc.SalesTaxRateList.TaxRateDetail {
				var taxType string
				switch stRate.TaxTypeApplicable {
				case "TaxOnAmount":
					taxType = "Tax On Amount"
				case "TaxOnAmountPlusTax":
					taxType = "Tax On Amount Plus Tax"
				case "TaxOnTax":
					taxType = "Tax On Tax"
				}
				item := map[string]any{
					"id":            fmt.Sprintf("st:%s", stRate.TaxOrder.String()),
					"Name":          stRate.TaxRateRef.Name,
					"TaxRateId":     stRate.TaxRateRef.Value,
					"TaxType":       taxType,
					"TaxOrder":      stRate.TaxOrder,
					"TaxCodeIdSale": tc.Id,
				}
				items = append(items, item)
			}
			return items, nil
		},
	},
	sourceType: &TaxCode,
}

func init() {
	registerType(&TaxCode)
	registerType(&TaxCodeLine)
}
