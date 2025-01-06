package data

import (
	"fmt"

	"github.com/tommyhedley/fibery/fibery-qbo-integration/pkgs/fibery"
	"github.com/tommyhedley/fibery/fibery-qbo-integration/pkgs/qbo"
)

var Invoice = QuickbooksDualType{
	QuickbooksType: QuickbooksType{
		FiberyType: FiberyType{
			Id:     "invoice",
			Name:   "Invoice",
			Schema: map[string]fibery.Field{},
		},
		SchemaTransformer: func(entity any) (map[string]any, error) {
			return map[string]any{}, nil
		},
		DataQuery: func(req Request) (Response, error) {
			return Response{}, nil
		},
	},
}

var InvoiceLine = DependentDataType{
	FiberyType: FiberyType{
		Id:   "Invoice_line",
		Name: "Invoice Line",
		Schema: map[string]fibery.Field{
			"id": {
				Name: "id",
				Type: fibery.ID,
			},
			"qbo_id": {
				Name: "QBO ID",
				Type: fibery.Text,
			},
			"invoice_sync_token": {
				Name:     "Invoice Sync Token",
				Type:     fibery.Text,
				ReadOnly: true,
			},
			"invoice_id": {
				Name: "Invoice ID",
				Type: fibery.Text,
				Relation: fibery.Relation{
					Cardinality:   fibery.MTO,
					Name:          "Invoice",
					TargetName:    "Invoice Lines",
					TargetType:    "Invoice",
					TargetFieldID: "id",
				},
			},
			"description": {
				Name:    "Description",
				Type:    fibery.Text,
				SubType: fibery.Title,
			},
			"line_type": {
				Name:    "Line Type",
				Type:    fibery.Text,
				SubType: fibery.SingleSelect,
				Options: []map[string]any{
					{
						"name": "Sales Item",
					},
					{
						"name": "Group",
					},
					{
						"name": "Description",
					},
					{
						"name": "Group",
					},
				},
			},
			"quantity": {
				Name: "Quantity",
				Type: fibery.Number,
				Format: map[string]any{
					"format":               "Number",
					"hasThousandSeparator": true,
					"precision":            2,
				},
			},
			"unit_price": {
				Name: "Unit Price",
				Type: fibery.Number,
				Format: map[string]any{
					"format":               "Money",
					"currencyCode":         "USD",
					"hasThousandSeperator": true,
					"precision":            2,
				},
			},
			"amount": {
				Name: "Amount",
				Type: fibery.Number,
				Format: map[string]any{
					"format":               "Money",
					"currencyCode":         "USD",
					"hasThousandSeperator": true,
					"precision":            2,
				},
			},
			"line_num": {
				Name:    "Line",
				Type:    fibery.Number,
				SubType: fibery.Integer,
			},
			"group_line_id": {
				Name: "Group Line ID",
				Type: fibery.Text,
				Relation: fibery.Relation{
					Cardinality:   fibery.MTO,
					Name:          "Group",
					TargetName:    "Lines",
					TargetType:    "Invoice_line",
					TargetFieldID: "id",
				},
			},
			"item_id": {
				Name: "Item",
				Type: fibery.Text,
				Relation: fibery.Relation{
					Cardinality:   fibery.MTO,
					Name:          "Item",
					TargetName:    "Invoice Lines",
					TargetType:    "Item",
					TargetFieldID: "id",
				},
			},
			"class_id": {
				Name: "Class ID",
				Type: fibery.Text,
				Relation: fibery.Relation{
					Cardinality:   fibery.MTO,
					Name:          "Class",
					TargetName:    "Expense Account Line(s)",
					TargetType:    "Class",
					TargetFieldID: "id",
				},
			},
			"tax_code_id": {
				Name: "Tax Code ID",
				Type: fibery.Text,
				Relation: fibery.Relation{
					Cardinality:   fibery.MTO,
					Name:          "Tax Code",
					TargetName:    "Invoice Lines",
					TargetType:    "TaxCode",
					TargetFieldID: "id",
				},
			},
			"markup_percent": {
				Name: "Markup",
				Type: fibery.Number,
				Format: map[string]any{
					"format":    "Percent",
					"precision": 2,
				},
			},
			"service_date": {
				Name:    "Date",
				Type:    fibery.DateType,
				SubType: fibery.Day,
			},
			"__syncAction": {
				Type: fibery.Text,
				Name: "Sync Action",
			}},
	},
	Source: Invoice,
	SchemaTransformer: func(entity any, source any) (map[string]any, error) {
		line, ok := entity.(qbo.Line)
		if !ok {
			return nil, fmt.Errorf("unable to convert entity to invoice line")
		}
		invoice, ok := source.(qbo.Invoice)
		if !ok {
			return nil, fmt.Errorf("unable to convert source to invoice")
		}
		var lineTypes = map[string]string{
			"SalesItemLineDetail": "Sales Item Line",
			"GroupLineDetail":     "Group Line",
			"DescriptionOnly":     "Description Line",
		}

		if line.DetailType == "GroupLineDetail" {
			return map[string]any{
				"id":                 fmt.Sprintf("%s:%s", invoice.Id, line.Id),
				"qbo_id":             line.Id,
				"invoice_sync_token": invoice.SyncToken,
				"invoice_id":         invoice.Id,
				"description":        line.Description,
				"line_type":          lineTypes[line.DetailType],
				"quantity":           line.GroupLineDetail.Quantity,
				"line_num":           line.LineNum,
				"item_id":            line.GroupLineDetail.GroupItemRef.Value,
				"__syncAction":       fibery.SET,
			}, nil
		} else if line.DetailType == "DescriptionOnly" || line.DetailType == "SalesItemLineDetail" {
			return map[string]any{
				"id":                 fmt.Sprintf("%s:%s", invoice.Id, line.Id),
				"qbo_id":             line.Id,
				"invoice_sync_token": invoice.SyncToken,
				"invoice_id":         invoice.Id,
				"description":        line.Description,
				"type":               lineTypes[line.DetailType],
				"quantity":           line.SalesItemLineDetail.Qty,
				"unit_price":         line.SalesItemLineDetail.UnitPrice,
				"amount":             line.Amount,
				"line_num":           line.LineNum,
				"item_id":            line.SalesItemLineDetail.ItemRef.Value,
				"class_id":           line.SalesItemLineDetail.ClassRef.Value,
				"tax_code_id":        line.SalesItemLineDetail.TaxCodeRef.Value,
				"markup_percent":     line.SalesItemLineDetail.MarkupInfo.Percent,
				"service_date":       line.SalesItemLineDetail.ServiceDate.Format(fibery.DateFormat),
				"__syncAction":       fibery.SET,
			}, nil
		}
		return nil, nil
	},
	SourceMapper: func(source any) (map[string]bool, error) {
		invoice, ok := source.(qbo.Invoice)
		if !ok {
			return nil, fmt.Errorf("unable to convert source to invoice")
		}
		sourceMap := map[string]bool{}
		for _, line := range invoice.Line {
			if line.DetailType == "GroupLineDetail" {
				for _, groupLine := range line.GroupLineDetail.Line {
					sourceMap[fmt.Sprintf("%s:%s:%s", invoice.Id, line.Id, groupLine.Id)] = true
				}
				sourceMap[fmt.Sprintf("%s:%s", invoice.Id, line.Id)] = true
			}
			if line.DetailType == "DescriptionOnly" || line.DetailType == "SalesItemLineDetail" {
				sourceMap[fmt.Sprintf("%s:%s", invoice.Id, line.Id)] = true
			}
		}
		return sourceMap, nil
	},
	TypeMapper: func(sourceArray any, sourceMapper SourceMapperFunc) (map[string]map[string]bool, error) {
		invoices, ok := sourceArray.([]qbo.Invoice)
		if !ok {
			return nil, fmt.Errorf("unable to convert sourceArray to invoices")
		}
		idMap := map[string]map[string]bool{}
		for _, invoice := range invoices {
			sourceMap, err := sourceMapper(invoice)
			if err != nil {
				return nil, fmt.Errorf("unable to map source: %w", err)
			}
			idMap[invoice.Id] = sourceMap
		}
		return idMap, nil
	},
	QueryProcessor: func(sourceArray any, schemaTransformer DepSchemaTransformerFunc) ([]map[string]any, error) {
		invoices, ok := sourceArray.([]qbo.Invoice)
		if !ok {
			return nil, fmt.Errorf("unable to convert sourceArray to invoices")
		}
		items := []map[string]any{}
		for _, invoice := range invoices {
			for _, line := range invoice.Line {
				if line.DetailType == "DescriptionOnly" || line.DetailType == "SalesItemLineDetail" {
					item, err := schemaTransformer(line, invoice)
					if err != nil {
						return nil, fmt.Errorf("unable to invoice line transform data: %w", err)
					}
					items = append(items, item)
				}
				if line.DetailType == "GroupLineDetail" {
					for _, groupLine := range line.GroupLineDetail.Line {
						item, err := schemaTransformer(groupLine, invoice)
						if err != nil {
							return nil, fmt.Errorf("unable to invoice line transform data: %w", err)
						}
						item["id"] = fmt.Sprintf("%s:%s:%s", invoice.Id, line.Id, groupLine.Id)
						item["group_line_id"] = line.Id
						items = append(items, item)
					}
					item, err := schemaTransformer(line, invoice)
					if err != nil {
						return nil, fmt.Errorf("unable to invoice line transform data: %w", err)
					}
					items = append(items, item)
				}
			}
		}
		return items, nil
	},
	ChangeDataCaptureProcessor: func(cdc qbo.ChangeDataCapture, cacheEntry *IdCache, sourceMapper SourceMapperFunc, schemaTransformer DepSchemaTransformerFunc) ([]map[string]any, error) {
		items := []map[string]any{}
		cacheEntry.Mu.Lock()
		defer cacheEntry.Mu.Unlock()
		for _, cdcResponse := range cdc.CDCResponse {
			for _, queryResponse := range cdcResponse.QueryResponse {
				for _, invoice := range queryResponse.Invoice {
					// map lines in cdc response
					cdcItemIds, err := sourceMapper(invoice)
					if err != nil {
						return nil, fmt.Errorf("unable to map source: %w", err)
					}

					// handle lines on deleted invoices
					if invoice.Status == "Deleted" {
						cachedIds := cacheEntry.Entries[invoice.Id]
						fmt.Printf("cachedIds from invoice %s: %s\n", invoice.Id, FormatJSON(cachedIds))
						for cachedId := range cachedIds {
							items = append(items, map[string]any{
								"id":           cachedId,
								"__syncAction": fibery.REMOVE,
							})
						}
						delete(cacheEntry.Entries, invoice.Id)
						if _, ok := cacheEntry.Entries[invoice.Id]; !ok {
							fmt.Printf("cache entry for invoice %s deleted\n", invoice.Id)
						}
						continue
					}

					fmt.Printf("items after delete: %s\n", FormatJSON(items))

					// transform line data on added or updated invoices
					for _, line := range invoice.Line {
						if line.DetailType == "DescriptionOnly" || line.DetailType == "SalesItemLineDetail" {
							item, err := schemaTransformer(line, invoice)
							if err != nil {
								return nil, fmt.Errorf("unable to invoice line transform data: %w", err)
							}
							items = append(items, item)
						}
						if line.DetailType == "GroupLineDetail" {
							for _, groupLine := range line.GroupLineDetail.Line {
								item, err := schemaTransformer(groupLine, invoice)
								if err != nil {
									return nil, fmt.Errorf("unable to invoice line transform data: %w", err)
								}
								item["id"] = fmt.Sprintf("%s:%s:%s", invoice.Id, line.Id, groupLine.Id)
								item["group_line_id"] = line.Id
								items = append(items, item)
							}
							item, err := schemaTransformer(line, invoice)
							if err != nil {
								return nil, fmt.Errorf("unable to invoice line transform data: %w", err)
							}
							items = append(items, item)
						}
					}

					fmt.Printf("items after transform: %s\n", FormatJSON(items))

					// check for lines in cache but not in cdc response
					if _, ok := cacheEntry.Entries[invoice.Id]; ok {
						cachedIds := cacheEntry.Entries[invoice.Id]
						fmt.Printf("cachedLines: %s\n", FormatJSON(cachedIds))
						for cachedId := range cachedIds {
							if !cdcItemIds[cachedId] {
								items = append(items, map[string]any{
									"id":           cachedId,
									"__syncAction": fibery.REMOVE,
								})
							}
						}
					}

					fmt.Printf("items after remove: %s\n", FormatJSON(items))

					// update cache with new line ids
					cacheEntry.Entries[invoice.Id] = cdcItemIds

					fmt.Printf("cache after transform: %s\n", FormatJSON(cacheEntry.Entries[invoice.Id]))
				}
			}
		}
		return items, nil
	},
}

func init() {
	fmt.Println(Invoice.GetName())
	TestFiberyType(Invoice)
}
