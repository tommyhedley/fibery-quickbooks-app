package data

import (
	"github.com/tommyhedley/fibery-quickbooks-app/pkgs/fibery"
	"github.com/tommyhedley/quickbooks-go"
)

var Estimate = QuickBooksDualType{
	QuickBooksType: QuickBooksType{
		fiberyType: fiberyType{
			id:   "Estimate",
			name: "Estimate",
			schema: map[string]fibery.Field{
				"id": {
					Name: "id",
					Type: fibery.ID,
				},
				"qbo_id": {
					Name: "QBO ID",
					Type: fibery.Text,
				},
				"name": {
					Name: "Name",
					Type: fibery.Text,
				},
				"customer_id": {
					Name: "Customer ID",
					Type: fibery.Text,
					Relation: &fibery.Relation{
						Cardinality:   fibery.MTO,
						Name:          "Customer",
						TargetName:    "Estimates",
						TargetType:    "Customer",
						TargetFieldID: "id",
					},
				},
				"sync_token": {
					Name:     "Sync Token",
					Type:     fibery.Text,
					ReadOnly: true,
				},
				"shipping_from_line_1": {
					Name: "Sale Line 1",
					Type: fibery.Text,
				},
				"shipping_from_line_2": {
					Name: "Sale Line 2",
					Type: fibery.Text,
				},
				"shipping_from_line_3": {
					Name: "Sale Line 3",
					Type: fibery.Text,
				},
				"shipping_from_line_4": {
					Name: "Sale Line 4",
					Type: fibery.Text,
				},
				"shipping_from_line_5": {
					Name: "Sale Line 5",
					Type: fibery.Text,
				},
				"shipping_from_city": {
					Name: "Sale City",
					Type: fibery.Text,
				},
				"shipping_from_state": {
					Name: "Sale State",
					Type: fibery.Text,
				},
				"shipping_from_postal_code": {
					Name: "Sale Postal Code",
					Type: fibery.Text,
				},
				"shipping_from_country": {
					Name: "Sale Country",
					Type: fibery.Text,
				},
				"shipping_from_lat": {
					Name: "Sale Latitude",
					Type: fibery.Text,
				},
				"shipping_from_long": {
					Name: "Sale Longitude",
					Type: fibery.Text,
				},
				"shipping_line_1": {
					Name: "Shipping Line 1",
					Type: fibery.Text,
				},
				"shipping_line_2": {
					Name: "Shipping Line 2",
					Type: fibery.Text,
				},
				"shipping_line_3": {
					Name: "Shipping Line 3",
					Type: fibery.Text,
				},
				"shipping_line_4": {
					Name: "Shipping Line 4",
					Type: fibery.Text,
				},
				"shipping_line_5": {
					Name: "Shipping Line 5",
					Type: fibery.Text,
				},
				"shipping_city": {
					Name: "Shipping City",
					Type: fibery.Text,
				},
				"shipping_state": {
					Name: "Shipping State",
					Type: fibery.Text,
				},
				"shipping_postal_code": {
					Name: "Shipping Postal Code",
					Type: fibery.Text,
				},
				"shipping_country": {
					Name: "Shipping Country",
					Type: fibery.Text,
				},
				"shipping_lat": {
					Name: "Shipping Latitude",
					Type: fibery.Text,
				},
				"shipping_long": {
					Name: "Shipping Longitude",
					Type: fibery.Text,
				},
				"billing_line_1": {
					Name: "Billing Line 1",
					Type: fibery.Text,
				},
				"billing_line_2": {
					Name: "Billing Line 2",
					Type: fibery.Text,
				},
				"billing_line_3": {
					Name: "Billing Line 3",
					Type: fibery.Text,
				},
				"billing_line_4": {
					Name: "Billing Line 4",
					Type: fibery.Text,
				},
				"billing_line_5": {
					Name: "Billing Line 5",
					Type: fibery.Text,
				},
				"billing_city": {
					Name: "Billing City",
					Type: fibery.Text,
				},
				"billing_state": {
					Name: "Billing State",
					Type: fibery.Text,
				},
				"billing_postal_code": {
					Name: "Billing Postal Code",
					Type: fibery.Text,
				},
				"billing_country": {
					Name: "Billing Country",
					Type: fibery.Text,
				},
				"billing_lat": {
					Name: "Billing Latitude",
					Type: fibery.Text,
				},
				"billing_long": {
					Name: "Billing Longitude",
					Type: fibery.Text,
				},
				"invoice_num": {
					Name: "Number",
					Type: fibery.Text,
				},
				"email": {
					Name:    "To",
					Type:    fibery.Text,
					SubType: fibery.Email,
				},
				"email_cc": {
					Name:    "CC",
					Type:    fibery.Text,
					SubType: fibery.Email,
				},
				"email_bcc": {
					Name:    "BCC",
					Type:    fibery.Text,
					SubType: fibery.Email,
				},
				"email_status": {
					Name:     "Email Status",
					Type:     fibery.Text,
					SubType:  fibery.SingleSelect,
					ReadOnly: true,
					Options: []map[string]any{
						{
							"name": "Need To Send",
						},
						{
							"name": "Sent",
						},
					},
				},
				"email_send_time": {
					Name: "Send Time",
					Type: fibery.DateType,
				},
				"txn_date": {
					Name: "Invoice Date",
					Type: fibery.DateType,
				},
				"due_date": {
					Name: "Due Date",
					Type: fibery.DateType,
				},
				"class_id": {
					Name: "Class ID",
					Type: fibery.Text,
					Relation: &fibery.Relation{
						Cardinality:   fibery.MTO,
						Name:          "Class",
						TargetName:    "Invoices",
						TargetType:    "Class",
						TargetFieldID: "id",
					},
				},
				"print_status": {
					Name:     "Print Status",
					Type:     fibery.Text,
					SubType:  fibery.SingleSelect,
					ReadOnly: true,
					Options: []map[string]any{
						{
							"name": "Need To Print",
						},
						{
							"name": "Print Complete",
						},
					},
				},
				"sales_term_id": {
					Name: "Sales Term ID",
					Type: fibery.Text,
					Relation: &fibery.Relation{
						Cardinality:   fibery.MTO,
						Name:          "Terms",
						TargetName:    "Invoices",
						TargetType:    "Term",
						TargetFieldID: "id",
					},
				},
				"private_note": {
					Name:    "Message on Statement",
					Type:    fibery.Text,
					SubType: fibery.MD,
				},
				"customer_memo": {
					Name: "Message on Invoice",
					Type: fibery.Text,
				},
				"allow_ach": {
					Name:    "ACH Payments",
					Type:    fibery.Text,
					SubType: fibery.Boolean,
				},
				"allow_cc": {
					Name:    "Credit Card Payments",
					Type:    fibery.Text,
					SubType: fibery.Boolean,
				},
				"tax_code_id": {
					Name: "Tax Code ID",
					Type: fibery.Text,
					Relation: &fibery.Relation{
						Cardinality:   fibery.MTO,
						Name:          "Sales Tax Rate",
						TargetName:    "Invoices",
						TargetType:    "TaxCode",
						TargetFieldID: "id",
					},
				},
				"tax_position": {
					Name:     "Apply Tax",
					Type:     fibery.Text,
					SubType:  fibery.SingleSelect,
					ReadOnly: true,
					Options: []map[string]any{
						{
							"name":    "Before Discount",
							"default": true,
						},
						{
							"name": "After Discount",
						},
					},
				},
				"tax_exemption_id": {
					Name: "Tax Exemption ID",
					Type: fibery.Text,
					Relation: &fibery.Relation{
						Cardinality:   fibery.MTO,
						Name:          "Tax Exemption",
						TargetName:    "Invoices",
						TargetType:    "TaxExemption",
						TargetFieldID: "id",
					},
				},
				"deposit_account_id": {
					Name: "Deposit Account ID",
					Type: fibery.Text,
					Relation: &fibery.Relation{
						Cardinality:   fibery.MTO,
						Name:          "Deposit Account",
						TargetName:    "Invoice Deposits",
						TargetType:    "Account",
						TargetFieldID: "id",
					},
				},
				"deposit_field": {
					Name: "Deposit Amount",
					Type: fibery.Number,
					Format: map[string]any{
						"format":               "Money",
						"currencyCode":         "USD",
						"hasThousandSeperator": true,
						"precision":            2,
					},
				},
				"discount_type": {
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
				"discount_percent": {
					Name: "Discount Percent",
					Type: fibery.Number,
					Format: map[string]any{
						"format":    "Percent",
						"precision": 2,
					},
				},
				"discount_amount": {
					Name: "Discount Amount",
					Type: fibery.Number,
					Format: map[string]any{
						"format":               "Money",
						"currencyCode":         "USD",
						"hasThousandSeperator": true,
						"precision":            2,
					},
				},
				"tax": {
					Name: "Tax",
					Type: fibery.Number,
					Format: map[string]any{
						"format":               "Money",
						"currencyCode":         "USD",
						"hasThousandSeperator": true,
						"precision":            2,
					},
				},
				"subtotal": {
					Name: "Subtotal",
					Type: fibery.Number,
					Format: map[string]any{
						"format":               "Money",
						"currencyCode":         "USD",
						"hasThousandSeperator": true,
						"precision":            2,
					},
				},
				"total": {
					Name: "Total",
					Type: fibery.Number,
					Format: map[string]any{
						"format":               "Money",
						"currencyCode":         "USD",
						"hasThousandSeperator": true,
						"precision":            2,
					},
				},
				"balance": {
					Name: "Balance",
					Type: fibery.Number,
					Format: map[string]any{
						"format":               "Money",
						"currencyCode":         "USD",
						"hasThousandSeperator": true,
						"precision":            2,
					},
				},
				"__syncAction": {
					Type: fibery.Text,
					Name: "Sync Action",
				},
			},
		},
		schemaGen: func(entity any) (map[string]any, error) {
			invoice, ok := entity.(quickbooks.Invoice)
			if !ok {
				return nil, fmt.Errorf("unable to convert entity to invoice")
			}
			var emailStatus = map[string]string{
				"NotSet":     "",
				"NeedToSend": "Need To Send",
				"EmailSent":  "Sent",
			}

			var printStatus = map[string]string{
				"NotSet":        "",
				"NeedToPrint":   "Need To Print",
				"PrintComplete": "Print Complete",
			}

			var taxPosition = map[bool]string{
				true:  "After Discount",
				false: "Before Discount",
			}

			var discountType = map[bool]string{
				true:  "Percentage",
				false: "Amount",
			}
			var data map[string]any
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

			var discountTypeValue string
			var discountPercent float32
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

			data = map[string]any{
				"id":                   invoice.Id,
				"qbo_id":               invoice.Id,
				"name":                 name,
				"customer_id":          invoice.CustomerRef.Value,
				"sync_token":           invoice.SyncToken,
				"shipping_line_1":      invoice.ShipAddr.Line1,
				"shipping_line_2":      invoice.ShipAddr.Line2,
				"shipping_line_3":      invoice.ShipAddr.Line3,
				"shipping_line_4":      invoice.ShipAddr.Line4,
				"shipping_line_5":      invoice.ShipAddr.Line5,
				"shipping_city":        invoice.ShipAddr.City,
				"shipping_state":       invoice.ShipAddr.CountrySubDivisionCode,
				"shipping_postal_code": invoice.ShipAddr.PostalCode,
				"shipping_country":     invoice.ShipAddr.Country,
				"shipping_lat":         invoice.ShipAddr.Lat,
				"shipping_long":        invoice.ShipAddr.Long,
				"billing_line_1":       invoice.BillAddr.Line1,
				"billing_line_2":       invoice.BillAddr.Line2,
				"billing_line_3":       invoice.BillAddr.Line3,
				"billing_line_4":       invoice.BillAddr.Line4,
				"billing_line_5":       invoice.BillAddr.Line5,
				"billing_city":         invoice.BillAddr.City,
				"billing_state":        invoice.BillAddr.CountrySubDivisionCode,
				"billing_postal_code":  invoice.BillAddr.PostalCode,
				"billing_country":      invoice.BillAddr.Country,
				"billing_lat":          invoice.BillAddr.Lat,
				"billing_long":         invoice.BillAddr.Long,
				"invoice_num":          invoice.DocNumber,
				"email":                invoice.BillEmail.Address,
				"email_cc":             invoice.BillEmailCC.Address,
				"email_bcc":            invoice.BillEmailBCC.Address,
				"email_status":         emailStatus[invoice.EmailStatus],
				"email_send_time":      emailSendTime,
				"txn_date":             invoice.TxnDate.Format(fibery.DateFormat),
				"due_date":             invoice.DueDate.Format(fibery.DateFormat),
				"class_id":             invoice.ClassRef.Value,
				"print_status":         printStatus[invoice.PrintStatus],
				"term_id":              invoice.SalesTermRef.Value,
				"private_note":         invoice.PrivateNote,
				"customer_memo":        invoice.CustomerMemo.Value,
				"allow_ach":            invoice.AllowOnlineACHPayment,
				"allow_cc":             invoice.AllowOnlineCreditCardPayment,
				"tax_code_id":          invoice.TxnTaxDetail.TxnTaxCodeRef.Value,
				"tax_position":         taxPosition[invoice.ApplyTaxAfterDiscount],
				"tax_exemption_id":     invoice.TaxExemptionRef.Value,
				"deposit_account_id":   invoice.DepositToAccountRef.Value,
				"deposit_field":        invoice.Deposit,
				"discount_type":        discountTypeValue,
				"discount_percent":     discountPercent,
				"discount_amount":      discountAmount,
				"tax":                  invoice.TxnTaxDetail.TotalTax,
				"subtotal":             subtotalAmount,
				"total":                invoice.TotalAmt,
				"balance":              invoice.Balance,
				"created_qbo":          invoice.MetaData.CreateTime.Format(fibery.DateFormat),
				"last_updated_qbo":     invoice.MetaData.LastUpdatedTime.Format(fibery.DateFormat),
				"__syncAction":         fibery.SET,
			}
			return data, nil
		},
		query: func(req Request) (Response, error) {
			invoices, err := req.Client.FindInvoicesByPage(req.StartPosition)
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
					dependent := (*dependentPointer).(DepWHQueryable)
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
