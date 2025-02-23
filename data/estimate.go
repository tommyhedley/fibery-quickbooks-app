package data

import (
	"encoding/json"
	"fmt"

	"github.com/tommyhedley/fibery-quickbooks-app/pkgs/fibery"
	"github.com/tommyhedley/quickbooks-go"
)

var Estimate = QuickBooksDualType{
	QuickBooksType: QuickBooksType{
		fiberyType: fiberyType{
			id:   "Estimate",
			name: "Estimate",
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
				"TxnStatus": {
					Name:     "Status",
					Type:     fibery.Text,
					SubType:  fibery.SingleSelect,
					ReadOnly: true,
					Options: []map[string]any{
						{
							"name":    "Pending",
							"default": true,
						},
						{
							"name": "Accepted",
						},
						{
							"name": "Closed",
						},
						{
							"name": "Rejected",
						},
					},
				},
				"AcceptedBy": {
					Name: "Accepted By",
					Type: fibery.Text,
				},
				"AcceptedDate": {
					Name: "Accepted Date",
					Type: fibery.DateType,
				},
				"ExpirationDate": {
					Name: "Expiration Date",
					Type: fibery.DateType,
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
				"ClassId": {
					Name: "Class ID",
					Type: fibery.Text,
					Relation: &fibery.Relation{
						Cardinality:   fibery.MTO,
						Name:          "Class",
						TargetName:    "Estimates",
						TargetType:    "Class",
						TargetFieldID: "Id",
					},
				},
				"TaxCodeId": {
					Name: "Tax Code ID",
					Type: fibery.Text,
					Relation: &fibery.Relation{
						Cardinality:   fibery.MTO,
						Name:          "Sales Tax Rate",
						TargetName:    "Estimates",
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
						TargetName:    "Estimates",
						TargetType:    "TaxExemption",
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
			estimate, ok := entity.(quickbooks.Estimate)
			if !ok {
				return nil, fmt.Errorf("unable to convert entity to estimate")
			}

			var discountType = map[bool]string{
				true:  "Percentage",
				false: "Amount",
			}
			var discountLine *quickbooks.Line
			var subtotalLine *quickbooks.Line
			for _, line := range estimate.Line {
				if line.DetailType == "DiscountLineDetail" {
					if discountLine != nil {
						return nil, fmt.Errorf("estimate %s has more than one discount line", estimate.Id)
					}
					discountLine = &line
				}
				if line.DetailType == "SubTotalLineDetail" {
					if subtotalLine != nil {
						return nil, fmt.Errorf("estimate %s has more than one subtotal line", estimate.Id)
					}
					subtotalLine = &line
				}
			}

			if subtotalLine == nil {
				return nil, fmt.Errorf("estimate %s has no subtotal lines", estimate.Id)
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
			if estimate.DeliveryInfo != nil && !estimate.DeliveryInfo.DeliveryTime.IsZero() {
				emailSendTime = estimate.DeliveryInfo.DeliveryTime.Format(fibery.DateFormat)
			}

			var name string
			if estimate.CustomerRef.Name == "" {
				name = estimate.DocNumber
			} else {
				name = estimate.DocNumber + " " + estimate.CustomerRef.Name
			}

			return map[string]any{
				"Id":                     estimate.Id,
				"QBOId":                  estimate.Id,
				"Name":                   name,
				"SyncToken":              estimate.SyncToken,
				"__syncAction":           fibery.SET,
				"ShippingLine1":          estimate.ShipAddr.Line1,
				"ShippingLine2":          estimate.ShipAddr.Line2,
				"ShippingLine3":          estimate.ShipAddr.Line3,
				"ShippingLine4":          estimate.ShipAddr.Line4,
				"ShippingLine5":          estimate.ShipAddr.Line5,
				"ShippingCity":           estimate.ShipAddr.City,
				"ShippingState":          estimate.ShipAddr.CountrySubDivisionCode,
				"ShippingPostalCode":     estimate.ShipAddr.PostalCode,
				"ShippingCountry":        estimate.ShipAddr.Country,
				"ShippingLat":            estimate.ShipAddr.Lat,
				"ShippingLong":           estimate.ShipAddr.Long,
				"ShippingFromLine1":      estimate.ShipFromAddr.Line1,
				"ShippingFromLine2":      estimate.ShipFromAddr.Line2,
				"ShippingFromLine3":      estimate.ShipFromAddr.Line3,
				"ShippingFromLine4":      estimate.ShipFromAddr.Line4,
				"ShippingFromLine5":      estimate.ShipFromAddr.Line5,
				"ShippingFromCity":       estimate.ShipFromAddr.City,
				"ShippingFromState":      estimate.ShipFromAddr.CountrySubDivisionCode,
				"ShippingFromPostalCode": estimate.ShipFromAddr.PostalCode,
				"ShippingFromCountry":    estimate.ShipFromAddr.Country,
				"ShippingFromLat":        estimate.ShipFromAddr.Lat,
				"ShippingFromLong":       estimate.ShipFromAddr.Long,
				"BillingLine1":           estimate.BillAddr.Line1,
				"BillingLine2":           estimate.BillAddr.Line2,
				"BillingLine3":           estimate.BillAddr.Line3,
				"BillingLine4":           estimate.BillAddr.Line4,
				"BillingLine5":           estimate.BillAddr.Line5,
				"BillingCity":            estimate.BillAddr.City,
				"BillingState":           estimate.BillAddr.CountrySubDivisionCode,
				"BillingPostalCode":      estimate.BillAddr.PostalCode,
				"BillingCountry":         estimate.BillAddr.Country,
				"BillingLat":             estimate.BillAddr.Lat,
				"BillingLong":            estimate.BillAddr.Long,
				"DocNumber":              estimate.DocNumber,
				"TxnStatus":              estimate.TxnStatus,
				"AcceptedBy":             estimate.AcceptedBy,
				"AcceptedDate":           estimate.AcceptedDate.Format(fibery.DateFormat),
				"ExpirationDate":         estimate.ExpirationDate.Format(fibery.DateFormat),
				"Email":                  estimate.BillEmail.Address,
				"EmailCC":                estimate.BillEmailCC.Address,
				"EmailBCC":               estimate.BillEmailBCC.Address,
				"EmailSendLater":         estimate.EmailStatus == "NeedToSend",
				"EmailSent":              estimate.EmailStatus == "EmailSent",
				"EmailSendTime":          emailSendTime,
				"TxnDate":                estimate.TxnDate.Format(fibery.DateFormat),
				"PrivateNote":            estimate.PrivateNote,
				"CustomerMemo":           estimate.CustomerMemo.Value,
				"DiscountPosition":       estimate.ApplyTaxAfterDiscount,
				"DiscountType":           discountTypeValue,
				"DiscountPercent":        discountPercent,
				"DiscountAmount":         discountAmount,
				"Tax":                    estimate.TxnTaxDetail.TotalTax,
				"SubtotalAmt":            subtotalAmount,
				"TotalAmt":               estimate.TotalAmt,
				"ClassId":                estimate.ClassRef.Value,
				"TaxCodeId":              estimate.TxnTaxDetail.TxnTaxCodeRef.Value,
				"TaxExemptionId":         estimate.TaxExemptionRef.Value,
				"CustomerId":             estimate.CustomerRef.Value,
			}, nil
		},
		query: func(req Request) (Response, error) {
			estimates, err := req.Client.FindEstimatesByPage(req.StartPosition, req.PageSize)
			if err != nil {
				return Response{}, fmt.Errorf("unable to find invoices: %w", err)
			}

			return Response{
				Data:     estimates,
				MoreData: len(estimates) >= req.PageSize,
			}, nil
		},
		queryProcessor: func(entityArray any, schemaGen schemaGenFunc) ([]map[string]any, error) {
		},
	},
	cdcProcessor: func(cdc quickbooks.ChangeDataCapture, schemaGen schemaGenFunc) ([]map[string]any, error) {
	},
	whBatchProcessor: func(itemResponse quickbooks.BatchItemResponse, response *map[string][]map[string]any, cache *cache.Cache, realmId string, queryProcessor queryProcessorFunc, schemaGen schemaGenFunc, typeId string) error {
	},
}

var EstimateLine = DependentDualType{
	dependentBaseType: dependentBaseType{
		fiberyType: fiberyType{
			id:   "EstimateLine",
			name: "Estimate Line",
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
				"EstimateId": {
					Name: "Estimate ID",
					Type: fibery.Text,
					Relation: &fibery.Relation{
						Cardinality:   fibery.MTO,
						Name:          "Estimate",
						TargetName:    "Lines",
						TargetType:    "Estimate",
						TargetFieldID: "Id",
					},
				},
			},
		},
		schemaGen: func(entity, source any) (map[string]any, error) {
			line, ok := entity.(quickbooks.Line)
			if !ok {
				return nil, fmt.Errorf("unable to convert entity to line")
			}

			estimate, ok := source.(quickbooks.Estimate)
			if !ok {
				return nil, fmt.Errorf("unable to convert source to Estimate")
			}

			if line.DetailType == quickbooks.GroupLine {
				return map[string]any{
					"Id":           fmt.Sprintf("%s:%s", estimate.Id, line.Id),
					"QBOId":        line.Id,
					"Description":  line.Description,
					"__syncAction": fibery.SET,
					"Qty":          line.GroupLineDetail.Quantity,
					"LineNum":      line.LineNum,
					"ItemId":       line.GroupLineDetail.GroupItemRef.Value,
					"estimateId":   estimate.Id,
				}, nil
			} else if line.DetailType == quickbooks.DescriptionLine || line.DetailType == quickbooks.SalesItemLine {
				return map[string]any{
					"Id":           fmt.Sprintf("%s:%s", estimate.Id, line.Id),
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
					"EstimateId":   estimate.Id,
				}, nil
			}
			return nil, nil
		},
		queryProcessor: func(sourceArray any, schemaGen depSchemaGenFunc) ([]map[string]any, error) {
			estimates, ok := sourceArray.([]quickbooks.Estimate)
			if !ok {
				return nil, fmt.Errorf("unable to convert sourceArray to Estimates")
			}
			items := []map[string]any{}
			for _, estimate := range estimates {
				for _, line := range estimate.Line {
					if line.DetailType == quickbooks.DescriptionLine || line.DetailType == quickbooks.SalesItemLine {
						item, err := schemaGen(line, estimate)
						if err != nil {
							return nil, fmt.Errorf("unable to transform Line data: %w", err)
						}
						items = append(items, item)
					}
					if line.DetailType == "GroupLineDetail" {
						for _, groupLine := range line.GroupLineDetail.Line {
							item, err := schemaGen(groupLine, estimate)
							if err != nil {
								return nil, fmt.Errorf("unable to transform Line data: %w", err)
							}
							item["Id"] = fmt.Sprintf("%s:%s:%s", estimate.Id, line.Id, groupLine.Id)
							item["GroupLineId"] = line.Id
							items = append(items, item)
						}
						item, err := schemaGen(line, estimate)
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
	source:       Estimate,
	sourceMapper: func(source any) (map[string]bool, error) {},
	typeMapper:   func(sourceArray any, sourceMapper sourceMapperFunc) (map[string]map[string]bool, error) {},
	cdcProcessor: func(cdc quickbooks.ChangeDataCapture, cacheEntry *IdCache, sourceMapper sourceMapperFunc, schemaGen depSchemaGenFunc) ([]map[string]any, error) {
	},
	whBatchProcessor: func(sourceArray any, cacheEntry *IdCache, sourceMapper sourceMapperFunc, schemaGen depSchemaGenFunc) ([]map[string]any, error) {
	},
}
