package data

import (
	"encoding/json"
	"fmt"

	"github.com/patrickmn/go-cache"
	"github.com/tommyhedley/fibery-quickbooks-app/pkgs/fibery"
	"github.com/tommyhedley/quickbooks-go"
)

var Invoice = QuickBooksDualType{
	QuickBooksType: QuickBooksType{
		fiberyType: fiberyType{
			id:   "Invoice",
			name: "Invoice",
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
				"ShippingFromLine1": {
					Name: "Sale Line 1",
					Type: fibery.Text,
				},
				"ShippingFromLine2": {
					Name: "Sale Line 2",
					Type: fibery.Text,
				},
				"ShippingFromLine3": {
					Name: "Sale Line 3",
					Type: fibery.Text,
				},
				"ShippingFromLine4": {
					Name: "Sale Line 4",
					Type: fibery.Text,
				},
				"ShippingFromLine5": {
					Name: "Sale Line 5",
					Type: fibery.Text,
				},
				"ShippingFromCity": {
					Name: "Sale City",
					Type: fibery.Text,
				},
				"ShippingFromState": {
					Name: "Sale State",
					Type: fibery.Text,
				},
				"ShippingFromPostalCode": {
					Name: "Sale Postal Code",
					Type: fibery.Text,
				},
				"ShippingFromCountry": {
					Name: "Sale Country",
					Type: fibery.Text,
				},
				"shippingFromLat": {
					Name: "Sale Latitude",
					Type: fibery.Text,
				},
				"ShippingFromLong": {
					Name: "Sale Longitude",
					Type: fibery.Text,
				},
				"ShippingLine1": {
					Name: "Shipping Line 1",
					Type: fibery.Text,
				},
				"ShippingLine2": {
					Name: "Shipping Line 2",
					Type: fibery.Text,
				},
				"ShippingLine3": {
					Name: "Shipping Line 3",
					Type: fibery.Text,
				},
				"ShippingLine4": {
					Name: "Shipping Line 4",
					Type: fibery.Text,
				},
				"ShippingLine5": {
					Name: "Shipping Line 5",
					Type: fibery.Text,
				},
				"ShippingCity": {
					Name: "Shipping City",
					Type: fibery.Text,
				},
				"ShippingState": {
					Name: "Shipping State",
					Type: fibery.Text,
				},
				"ShippingPostalCode": {
					Name: "Shipping Postal Code",
					Type: fibery.Text,
				},
				"ShippingCountry": {
					Name: "Shipping Country",
					Type: fibery.Text,
				},
				"ShippingLat": {
					Name: "Shipping Latitude",
					Type: fibery.Text,
				},
				"ShippingLong": {
					Name: "Shipping Longitude",
					Type: fibery.Text,
				},
				"BillingLine1": {
					Name: "Billing Line 1",
					Type: fibery.Text,
				},
				"BillingLine2": {
					Name: "Billing Line 2",
					Type: fibery.Text,
				},
				"BillingLine3": {
					Name: "Billing Line 3",
					Type: fibery.Text,
				},
				"BillingLine4": {
					Name: "Billing Line 4",
					Type: fibery.Text,
				},
				"BillingLine5": {
					Name: "Billing Line 5",
					Type: fibery.Text,
				},
				"BillingCity": {
					Name: "Billing City",
					Type: fibery.Text,
				},
				"BillingState": {
					Name: "Billing State",
					Type: fibery.Text,
				},
				"BillingPostalCode": {
					Name: "Billing Postal Code",
					Type: fibery.Text,
				},
				"BillingCountry": {
					Name: "Billing Country",
					Type: fibery.Text,
				},
				"BillingLat": {
					Name: "Billing Latitude",
					Type: fibery.Text,
				},
				"BillingLong": {
					Name: "Billing Longitude",
					Type: fibery.Text,
				},
				"DocNumber": {
					Name: "Number",
					Type: fibery.Text,
				},
				"Email": {
					Name:    "To",
					Type:    fibery.Text,
					SubType: fibery.Email,
				},
				"EmailCC": {
					Name:    "CC",
					Type:    fibery.Text,
					SubType: fibery.Email,
				},
				"EmailBCC": {
					Name:    "BCC",
					Type:    fibery.Text,
					SubType: fibery.Email,
				},
				"EmailSendLater": {
					Name:    "Send Later",
					Type:    fibery.Text,
					SubType: fibery.Boolean,
				},
				"EmailSent": {
					Name:    "Sent",
					Type:    fibery.Text,
					SubType: fibery.Boolean,
				},
				"EmailSendTime": {
					Name: "Send Time",
					Type: fibery.DateType,
				},
				"TxnDate": {
					Name: "Date",
					Type: fibery.DateType,
				},
				"DueDate": {
					Name: "Due Date",
					Type: fibery.DateType,
				},
				"PrivateNote": {
					Name:    "Message on Statement",
					Type:    fibery.Text,
					SubType: fibery.MD,
				},
				"CustomerMemo": {
					Name:    "Message on Invoice",
					Type:    fibery.Text,
					SubType: fibery.MD,
				},
				"DiscountPosition": {
					Name:        "Discount Before Tax",
					Type:        fibery.Text,
					SubType:     fibery.Boolean,
					Description: "Should the discount be applied before or after the sales tax calculation? Default is false as tax should generally be calculated first before a discount is given.",
				},
				"DepositField": {
					Name: "Deposit Amount",
					Type: fibery.Number,
					Format: map[string]any{
						"format":               "Money",
						"currencyCode":         "USD",
						"hasThousandSeperator": true,
						"precision":            2,
					},
				},
				"DiscountType": {
					Name:     "Discount Type",
					Type:     fibery.Text,
					SubType:  fibery.SingleSelect,
					ReadOnly: true,
					Options: []map[string]any{
						{
							"name": "Percent",
						},
						{
							"name": "Amount",
						},
					},
				},
				"DiscountPercent": {
					Name: "Discount Percent",
					Type: fibery.Number,
					Format: map[string]any{
						"format":    "Percent",
						"precision": 2,
					},
				},
				"DiscountAmount": {
					Name: "Discount Amount",
					Type: fibery.Number,
					Format: map[string]any{
						"format":               "Money",
						"currencyCode":         "USD",
						"hasThousandSeperator": true,
						"precision":            2,
					},
				},
				"Tax": {
					Name: "Tax",
					Type: fibery.Number,
					Format: map[string]any{
						"format":               "Money",
						"currencyCode":         "USD",
						"hasThousandSeperator": true,
						"precision":            2,
					},
				},
				"SubtotalAmt": {
					Name: "Subtotal",
					Type: fibery.Number,
					Format: map[string]any{
						"format":               "Money",
						"currencyCode":         "USD",
						"hasThousandSeperator": true,
						"precision":            2,
					},
				},
				"TotalAmt": {
					Name: "Total",
					Type: fibery.Number,
					Format: map[string]any{
						"format":               "Money",
						"currencyCode":         "USD",
						"hasThousandSeperator": true,
						"precision":            2,
					},
				},
				"Balance": {
					Name: "Balance",
					Type: fibery.Number,
					Format: map[string]any{
						"format":               "Money",
						"currencyCode":         "USD",
						"hasThousandSeperator": true,
						"precision":            2,
					},
				},
				"AllowACH": {
					Name:    "ACH Payments",
					Type:    fibery.Text,
					SubType: fibery.Boolean,
				},
				"AllowCC": {
					Name:    "Credit Card Payments",
					Type:    fibery.Text,
					SubType: fibery.Boolean,
				},
				"ClassId": {
					Name: "Class ID",
					Type: fibery.Text,
					Relation: &fibery.Relation{
						Cardinality:   fibery.MTO,
						Name:          "Class",
						TargetName:    "Invoices",
						TargetType:    "Class",
						TargetFieldID: "Id",
					},
				},
				"SalesTermId": {
					Name: "Sales Term ID",
					Type: fibery.Text,
					Relation: &fibery.Relation{
						Cardinality:   fibery.MTO,
						Name:          "Terms",
						TargetName:    "Invoices",
						TargetType:    "Term",
						TargetFieldID: "Id",
					},
				},
				"TaxCodeId": {
					Name: "Tax Code ID",
					Type: fibery.Text,
					Relation: &fibery.Relation{
						Cardinality:   fibery.MTO,
						Name:          "Sales Tax Rate",
						TargetName:    "Invoices",
						TargetType:    "TaxCode",
						TargetFieldID: "Id",
					},
				},
				"TaxExemptionId": {
					Name: "Tax Exemption ID",
					Type: fibery.Text,
					Relation: &fibery.Relation{
						Cardinality:   fibery.MTO,
						Name:          "Tax Exemption",
						TargetName:    "Invoices",
						TargetType:    "TaxExemption",
						TargetFieldID: "Id",
					},
				},
				"DepositAccountId": {
					Name: "Deposit Account ID",
					Type: fibery.Text,
					Relation: &fibery.Relation{
						Cardinality:   fibery.MTO,
						Name:          "Deposit Account",
						TargetName:    "Invoice Deposits",
						TargetType:    "Account",
						TargetFieldID: "Id",
					},
				},
				"CustomerId": {
					Name: "Customer ID",
					Type: fibery.Text,
					Relation: &fibery.Relation{
						Cardinality:   fibery.MTO,
						Name:          "Customer",
						TargetName:    "Estimates",
						TargetType:    "Customer",
						TargetFieldID: "Id",
					},
				},
			},
		},
		schemaGen: func(entity any) (map[string]any, error) {
			invoice, ok := entity.(quickbooks.Invoice)
			if !ok {
				return nil, fmt.Errorf("unable to convert entity to invoice")
			}

			var discountType = map[bool]string{
				true:  "Percentage",
				false: "Amount",
			}
			var discountLine *quickbooks.Line
			var subtotalLine *quickbooks.Line
			for _, line := range invoice.Line {
				if line.DetailType == "DiscountLineDetail" {
					if discountLine != nil {
						return nil, fmt.Errorf("invoice %s has more than one discount line", invoice.Id)
					}
					discountLine = &line
				}
				if line.DetailType == "SubTotalLineDetail" {
					if subtotalLine != nil {
						return nil, fmt.Errorf("invoice %s has more than one subtotal line", invoice.Id)
					}
					subtotalLine = &line
				}
			}

			if subtotalLine == nil {
				return nil, fmt.Errorf("estimate %s has no subtotal lines", invoice.Id)
			}

			var discountTypeValue string
			var discountPercent json.Number
			var discountAmount json.Number

			if discountLine != nil {
				discountTypeValue = discountType[discountLine.DiscountLineDetail.PercentBased]
				discountPercent = discountLine.DiscountLineDetail.DiscountPercent
				discountAmount = discountLine.Amount
			}

			var subtotalAmount json.Number
			if subtotalLine != nil {
				subtotalAmount = subtotalLine.Amount
			}

			var emailSendTime string
			if invoice.DeliveryInfo != nil && !invoice.DeliveryInfo.DeliveryTime.IsZero() {
				emailSendTime = invoice.DeliveryInfo.DeliveryTime.Format(fibery.DateFormat)
			}

			var name string
			if invoice.CustomerRef.Name == "" {
				name = invoice.DocNumber
			} else {
				name = invoice.DocNumber + " " + invoice.CustomerRef.Name
			}

			return map[string]any{
				"Id":                     invoice.Id,
				"QBOId":                  invoice.Id,
				"Name":                   name,
				"SyncToken":              invoice.SyncToken,
				"__syncAction":           fibery.SET,
				"ShippingLine1":          invoice.ShipAddr.Line1,
				"ShippingLine2":          invoice.ShipAddr.Line2,
				"ShippingLine3":          invoice.ShipAddr.Line3,
				"ShippingLine4":          invoice.ShipAddr.Line4,
				"ShippingLine5":          invoice.ShipAddr.Line5,
				"ShippingCity":           invoice.ShipAddr.City,
				"ShippingState":          invoice.ShipAddr.CountrySubDivisionCode,
				"ShippingPostalCode":     invoice.ShipAddr.PostalCode,
				"ShippingCountry":        invoice.ShipAddr.Country,
				"ShippingLat":            invoice.ShipAddr.Lat,
				"ShippingLong":           invoice.ShipAddr.Long,
				"ShippingFromLine1":      invoice.ShipFromAddr.Line1,
				"ShippingFromLine2":      invoice.ShipFromAddr.Line2,
				"ShippingFromLine3":      invoice.ShipFromAddr.Line3,
				"ShippingFromLine4":      invoice.ShipFromAddr.Line4,
				"ShippingFromLine5":      invoice.ShipFromAddr.Line5,
				"ShippingFromCity":       invoice.ShipFromAddr.City,
				"ShippingFromState":      invoice.ShipFromAddr.CountrySubDivisionCode,
				"ShippingFromPostalCode": invoice.ShipFromAddr.PostalCode,
				"ShippingFromCountry":    invoice.ShipFromAddr.Country,
				"ShippingFromLat":        invoice.ShipFromAddr.Lat,
				"ShippingFromLong":       invoice.ShipFromAddr.Long,
				"BillingLine1":           invoice.BillAddr.Line1,
				"BillingLine2":           invoice.BillAddr.Line2,
				"BillingLine3":           invoice.BillAddr.Line3,
				"BillingLine4":           invoice.BillAddr.Line4,
				"BillingLine5":           invoice.BillAddr.Line5,
				"BillingCity":            invoice.BillAddr.City,
				"BillingState":           invoice.BillAddr.CountrySubDivisionCode,
				"BillingPostalCode":      invoice.BillAddr.PostalCode,
				"BillingCountry":         invoice.BillAddr.Country,
				"BillingLat":             invoice.BillAddr.Lat,
				"BillingLong":            invoice.BillAddr.Long,
				"DocNumber":              invoice.DocNumber,
				"Email":                  invoice.BillEmail.Address,
				"EmailCC":                invoice.BillEmailCC.Address,
				"EmailBCC":               invoice.BillEmailBCC.Address,
				"EmailSendLater":         invoice.EmailStatus == "NeedToSend",
				"EmailSent":              invoice.EmailStatus == "EmailSent",
				"EmailSendTime":          emailSendTime,
				"TxnDate":                invoice.TxnDate.Format(fibery.DateFormat),
				"DueDate":                invoice.DueDate.Format(fibery.DateFormat),
				"PrivateNote":            invoice.PrivateNote,
				"CustomerMemo":           invoice.CustomerMemo.Value,
				"DiscountPosition":       invoice.ApplyTaxAfterDiscount,
				"DepositField":           invoice.Deposit,
				"DiscountType":           discountTypeValue,
				"DiscountPercent":        discountPercent,
				"DiscountAmount":         discountAmount,
				"Tax":                    invoice.TxnTaxDetail.TotalTax,
				"SubtotalAmt":            subtotalAmount,
				"TotalAmt":               invoice.TotalAmt,
				"Balance":                invoice.Balance,
				"AllowACH":               invoice.AllowOnlineACHPayment,
				"AllowCC":                invoice.AllowOnlineCreditCardPayment,
				"ClassId":                invoice.ClassRef.Value,
				"TaxCodeId":              invoice.TxnTaxDetail.TxnTaxCodeRef.Value,
				"TaxExemptionId":         invoice.TaxExemptionRef.Value,
				"DepositAccountId":       invoice.DepartmentRef.Value,
				"CustomerId":             invoice.CustomerRef.Value,
			}, nil
		},
		query: func(req Request) (Response, error) {
			invoices, err := req.Client.FindInvoicesByPage(req.StartPosition, req.PageSize)
			if err != nil {
				return Response{}, fmt.Errorf("unable to find invoices: %w", err)
			}

			return Response{
				Data:     invoices,
				MoreData: len(invoices) >= quickbooks.QueryPageSize,
			}, nil
		},
		queryProcessor: func(entityArray any, schemaGen schemaGenFunc) ([]map[string]any, error) {
			invoices, ok := entityArray.([]quickbooks.Invoice)
			if !ok {
				return nil, fmt.Errorf("unable to convert entityArray to invoices")
			}
			items := []map[string]any{}
			for _, invoice := range invoices {
				item, err := schemaGen(invoice)
				if err != nil {
					return nil, fmt.Errorf("unable to transform data: %w", err)
				}
				items = append(items, item)
			}
			return items, nil
		},
	},
	cdcProcessor: func(cdc quickbooks.ChangeDataCapture, schemaGen schemaGenFunc) ([]map[string]any, error) {
		items := []map[string]any{}
		for _, cdcResponse := range cdc.CDCResponse {
			for _, queryResponse := range cdcResponse.QueryResponse {
				for _, cdcInvoice := range queryResponse.Invoice {
					if cdcInvoice.Status == "Deleted" {
						items = append(items, map[string]any{
							"id":           cdcInvoice.Id,
							"__syncAction": fibery.REMOVE,
						})
					} else {
						item, err := schemaGen(cdcInvoice.Invoice)
						if err != nil {
							return nil, fmt.Errorf("unable to transform data: %w", err)
						}
						items = append(items, item)
					}
				}
			}
		}
		return items, nil
	},
	whBatchProcessor: func(itemResponse quickbooks.BatchItemResponse, response *map[string][]map[string]any, cache *cache.Cache, realmId string, queryProcessor queryProcessorFunc, schemaGen schemaGenFunc, typeId string) error {
		if len(itemResponse.Fault.Faults) > 0 {
			return fmt.Errorf("batch request failed: %v", itemResponse.Fault.Faults)
		}
		if invoices := itemResponse.QueryResponse.Invoice; invoices != nil {
			invoiceData, err := queryProcessor(invoices, schemaGen)
			if err != nil {
				return fmt.Errorf("unable to process invoice query data: %w", err)
			}
			(*response)[typeId] = append((*response)[typeId], invoiceData...)
			if dependents, ok := SourceDependents[typeId]; ok {
				for _, dependentPointer := range dependents {
					// change type assertion to WHType
					dependent := (*dependentPointer).(DependentWHQueryable)
					cacheKey := fmt.Sprintf("%s:%s", realmId, dependent.GetId())
					if cacheEntry, found := cache.Get(cacheKey); found {
						cacheEntry, ok := cacheEntry.(*IdCache)
						if !ok {
							return fmt.Errorf("unable to convert cache entry to IdCache")
						}
						dependentData, err := dependent.ProcessWHBatch(invoices, cacheEntry)
						if err != nil {
							return fmt.Errorf("unable to process dependent %s query data: %w", dependent.GetId(), err)
						}
						(*response)[dependent.GetId()] = append((*response)[dependent.GetId()], dependentData...)
					}
				}
			}
		}
		return nil
	},
}

var InvoiceLine = DependentDualType{
	dependentBaseType: dependentBaseType{
		fiberyType: fiberyType{
			id:   "InvoiceLine",
			name: "Invoice Line",
			schema: map[string]fibery.Field{
				"Id": {
					Name: "ID",
					Type: fibery.ID,
				},
				"QBOId": {
					Name: "QBO ID",
					Type: fibery.Text,
				},
				"Description": {
					Name:    "Description",
					Type:    fibery.Text,
					SubType: fibery.Title,
				},
				"__syncAction": {
					Type: fibery.Text,
					Name: "Sync Action",
				},
				"LineNum": {
					Name: "Line",
					Type: fibery.Number,
				},
				"Taxed": {
					Name:    "Taxed",
					Type:    fibery.Text,
					SubType: fibery.Boolean,
				},
				"ServiceDate": {
					Name:    "Service Date",
					Type:    fibery.DateType,
					SubType: fibery.Day,
				},
				"Qty": {
					Name: "Quantity",
					Type: fibery.Number,
					Format: map[string]any{
						"format":               "Number",
						"hasThousandSeparator": true,
						"precision":            2,
					},
				},
				"UnitPrice": {
					Name: "Unit Price",
					Type: fibery.Number,
					Format: map[string]any{
						"format":               "Money",
						"currencyCode":         "USD",
						"hasThousandSeperator": true,
						"precision":            2,
					},
				},
				"Amount": {
					Name: "Amount",
					Type: fibery.Number,
					Format: map[string]any{
						"format":               "Money",
						"currencyCode":         "USD",
						"hasThousandSeperator": true,
						"precision":            2,
					},
				},
				"GroupLineId": {
					Name: "Group Line ID",
					Type: fibery.Text,
					Relation: &fibery.Relation{
						Cardinality:   fibery.MTO,
						Name:          "Group",
						TargetName:    "Lines",
						TargetType:    "EstimateLine",
						TargetFieldID: "Id",
					},
				},
				"ItemId": {
					Name: "Item",
					Type: fibery.Text,
					Relation: &fibery.Relation{
						Cardinality:   fibery.MTO,
						Name:          "Item",
						TargetName:    "Estimate Lines",
						TargetType:    "Item",
						TargetFieldID: "Id",
					},
				},
				"ClassId": {
					Name: "Class ID",
					Type: fibery.Text,
					Relation: &fibery.Relation{
						Cardinality:   fibery.MTO,
						Name:          "Class",
						TargetName:    "Expense Lines",
						TargetType:    "Class",
						TargetFieldID: "Id",
					},
				},
				"InvoiceId": {
					Name: "Estimate ID",
					Type: fibery.Text,
					Relation: &fibery.Relation{
						Cardinality:   fibery.MTO,
						Name:          "Invoice",
						TargetName:    "Lines",
						TargetType:    "Invoice",
						TargetFieldID: "Id",
					},
				},
			},
		},
		schemaGen: func(entity any, source any) (map[string]any, error) {
			line, ok := entity.(quickbooks.Line)
			if !ok {
				return nil, fmt.Errorf("unable to convert entity to Line")
			}

			invoice, ok := source.(quickbooks.Invoice)
			if !ok {
				return nil, fmt.Errorf("unable to convert source to Invoice")
			}

			if line.DetailType == quickbooks.GroupLine {
				return map[string]any{
					"Id":           fmt.Sprintf("%s:%s", invoice.Id, line.Id),
					"QBOId":        line.Id,
					"Description":  line.Description,
					"__syncAction": fibery.SET,
					"Qty":          line.GroupLineDetail.Quantity,
					"LineNum":      line.LineNum,
					"ItemId":       line.GroupLineDetail.GroupItemRef.Value,
					"InvoiceId":    invoice.Id,
				}, nil
			} else if line.DetailType == quickbooks.DescriptionLine || line.DetailType == quickbooks.SalesItemLine {
				return map[string]any{
					"Id":           fmt.Sprintf("%s:%s", invoice.Id, line.Id),
					"QBOId":        line.Id,
					"Description":  line.Description,
					"__syncAction": fibery.SET,
					"LineNum":      line.LineNum,
					"Taxed":        line.SalesItemLineDetail.TaxCodeRef.Value == "TAX",
					"ServiceDate":  line.SalesItemLineDetail.ServiceDate.Format(fibery.DateFormat),
					"Qty":          line.GroupLineDetail.Quantity,
					"UnitPrice":    line.SalesItemLineDetail.UnitPrice,
					"Amount":       line.Amount,
					"ItemId":       line.GroupLineDetail.GroupItemRef.Value,
					"ClassId":      line.SalesItemLineDetail.ClassRef.Value,
					"InvoiceId":    invoice.Id,
				}, nil
			}
			return nil, nil
		},
		queryProcessor: func(sourceArray any, schemaGen depSchemaGenFunc) ([]map[string]any, error) {
			invoices, ok := sourceArray.([]quickbooks.Invoice)
			if !ok {
				return nil, fmt.Errorf("unable to convert sourceArray to Invoices")
			}
			items := []map[string]any{}
			for _, invoice := range invoices {
				for _, line := range invoice.Line {
					if line.DetailType == quickbooks.DescriptionLine || line.DetailType == quickbooks.SalesItemLine {
						item, err := schemaGen(line, invoice)
						if err != nil {
							return nil, fmt.Errorf("unable to transform Line data: %w", err)
						}
						items = append(items, item)
					}
					if line.DetailType == quickbooks.GroupLine {
						for _, groupLine := range line.GroupLineDetail.Line {
							item, err := schemaGen(groupLine, invoice)
							if err != nil {
								return nil, fmt.Errorf("unable to transform Line data: %w", err)
							}
							item["Id"] = fmt.Sprintf("%s:%s:%s", invoice.Id, line.Id, groupLine.Id)
							item["GroupLineId"] = line.Id
							items = append(items, item)
						}
						item, err := schemaGen(line, invoice)
						if err != nil {
							return nil, fmt.Errorf("unable to transform Line data: %w", err)
						}
						items = append(items, item)
					}
				}
			}
			return items, nil
		},
	},
	source: Invoice,
	sourceMapper: func(source any) (map[string]bool, error) {
		invoice, ok := source.(quickbooks.Invoice)
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
	typeMapper: func(sourceArray any, sourceMapper sourceMapperFunc) (map[string]map[string]bool, error) {
		invoices, ok := sourceArray.([]quickbooks.Invoice)
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
	cdcProcessor: func(cdc quickbooks.ChangeDataCapture, cacheEntry *IdCache, sourceMapper sourceMapperFunc, schemaGen depSchemaGenFunc) ([]map[string]any, error) {
		items := []map[string]any{}
		cacheEntry.Mu.Lock()
		defer cacheEntry.Mu.Unlock()
		for _, cdcResponse := range cdc.CDCResponse {
			for _, queryResponse := range cdcResponse.QueryResponse {
				for _, cdcInvoice := range queryResponse.Invoice {
					// map lines in cdc response
					cdcItemIds, err := sourceMapper(cdcInvoice.Invoice)
					if err != nil {
						return nil, fmt.Errorf("unable to map source: %w", err)
					}

					// handle lines on deleted invoices
					if cdcInvoice.Status == "Deleted" {
						cachedIds := cacheEntry.Entries[cdcInvoice.Id]
						for cachedId := range cachedIds {
							items = append(items, map[string]any{
								"id":           cachedId,
								"__syncAction": fibery.REMOVE,
							})
						}
						delete(cacheEntry.Entries, cdcInvoice.Id)
						continue
					}

					// transform line data on added or updated invoices
					for _, line := range cdcInvoice.Line {
						if line.DetailType == "DescriptionOnly" || line.DetailType == "SalesItemLineDetail" {
							item, err := schemaGen(line, cdcInvoice.Invoice)
							if err != nil {
								return nil, fmt.Errorf("unable to invoice line transform data: %w", err)
							}
							items = append(items, item)
						}
						if line.DetailType == "GroupLineDetail" {
							for _, groupLine := range line.GroupLineDetail.Line {
								item, err := schemaGen(groupLine, cdcInvoice.Invoice)
								if err != nil {
									return nil, fmt.Errorf("unable to invoice line transform data: %w", err)
								}
								item["id"] = fmt.Sprintf("%s:%s:%s", cdcInvoice.Id, line.Id, groupLine.Id)
								item["group_line_id"] = line.Id
								items = append(items, item)
							}
							item, err := schemaGen(line, cdcInvoice.Invoice)
							if err != nil {
								return nil, fmt.Errorf("unable to invoice line transform data: %w", err)
							}
							items = append(items, item)
						}
					}

					// check for lines in cache but not in cdc response
					if _, ok := cacheEntry.Entries[cdcInvoice.Id]; ok {
						cachedIds := cacheEntry.Entries[cdcInvoice.Id]
						for cachedId := range cachedIds {
							if !cdcItemIds[cachedId] {
								items = append(items, map[string]any{
									"id":           cachedId,
									"__syncAction": fibery.REMOVE,
								})
							}
						}
					}

					// update cache with new line ids
					cacheEntry.Entries[cdcInvoice.Id] = cdcItemIds
				}
			}
		}
		return items, nil
	},
	whBatchProcessor: func(sourceArray any, cacheEntry *IdCache, sourceMapper sourceMapperFunc, schemaGen depSchemaGenFunc) ([]map[string]any, error) {
		invoices, ok := sourceArray.([]quickbooks.Invoice)
		if !ok {
			return nil, fmt.Errorf("unable to convert sourceArray to []qbo.Invoice")
		}
		items := []map[string]any{}
		cacheEntry.Mu.Lock()
		defer cacheEntry.Mu.Unlock()
		for _, invoice := range invoices {
			sourceItemIds, err := sourceMapper(invoice)
			if err != nil {
				return nil, fmt.Errorf("unable to map source: %w", err)
			}

			for _, line := range invoice.Line {
				if line.DetailType == "DescriptionOnly" || line.DetailType == "SalesItemLineDetail" {
					item, err := schemaGen(line, invoice)
					if err != nil {
						return nil, fmt.Errorf("unable to invoice line transform data: %w", err)
					}
					items = append(items, item)
				}
				if line.DetailType == "GroupLineDetail" {
					for _, groupLine := range line.GroupLineDetail.Line {
						item, err := schemaGen(groupLine, invoice)
						if err != nil {
							return nil, fmt.Errorf("unable to invoice line transform data: %w", err)
						}
						item["id"] = fmt.Sprintf("%s:%s:%s", invoice.Id, line.Id, groupLine.Id)
						item["group_line_id"] = line.Id
						items = append(items, item)
					}
					item, err := schemaGen(line, invoice)
					if err != nil {
						return nil, fmt.Errorf("unable to invoice line transform data: %w", err)
					}
					items = append(items, item)
				}
			}

			// check for lines in cache but not in cdc response
			if _, ok := cacheEntry.Entries[invoice.Id]; ok {
				cachedIds := cacheEntry.Entries[invoice.Id]
				for cachedId := range cachedIds {
					if !sourceItemIds[cachedId] {
						items = append(items, map[string]any{
							"id":           cachedId,
							"__syncAction": fibery.REMOVE,
						})
					}
				}
			}

			// update cache with new line ids
			cacheEntry.Entries[invoice.Id] = sourceItemIds
		}
		return items, nil
	},
}

func init() {
	RegisterType(Invoice)
	RegisterType(InvoiceLine)
}
