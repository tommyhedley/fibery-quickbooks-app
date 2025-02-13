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
		FiberyType: FiberyType{
			id:   "Invoice",
			name: "Invoice",
			schema: map[string]fibery.Field{
				"id": {
					Name: "id",
					Type: fibery.ID,
				},
				"qbo_id": {
					Name: "QBO ID",
					Type: fibery.Text,
				},
				"customer_id": {
					Name: "Customer ID",
					Type: fibery.Text,
					Relation: &fibery.Relation{
						Cardinality:   fibery.MTO,
						Name:          "Customer",
						TargetName:    "Invoices",
						TargetType:    "Customer",
						TargetFieldID: "id",
					},
				},
				"sync_token": {
					Name:     "Sync Token",
					Type:     fibery.Text,
					ReadOnly: true,
				},
				"shipping_address": {
					Name:    "Shipping Address",
					Type:    fibery.Text,
					SubType: fibery.MD,
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
				"billing_address": {
					Name:    "Billing Address",
					Type:    fibery.Text,
					SubType: fibery.MD,
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
					Name:    "Invoice",
					Type:    fibery.Text,
					SubType: fibery.Title,
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
				"date": {
					Name: "Date",
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
				"term_id": {
					Name: "Term ID",
					Type: fibery.Text,
					Relation: &fibery.Relation{
						Cardinality:   fibery.MTO,
						Name:          "Term",
						TargetName:    "Invoices",
						TargetType:    "Term",
						TargetFieldID: "id",
					},
				},
				"statement_memo": {
					Name: "Statement Message",
					Type: fibery.Text,
				},
				"customer_memo": {
					Name: "Invoice Message",
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
						Name:          "Tax Code",
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
				"created_qbo": {
					Name: "Creation Date (QBO)",
					Type: fibery.DateType,
				},
				"last_updated_qbo": {
					Name: "Last Updated (QBO)",
					Type: fibery.DateType,
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

			var billingAddress string
			if invoice.BillAddr.Line1 != "" {
				billingAddress += (invoice.BillAddr.Line1 + "  \n")
			}
			if invoice.BillAddr.Line2 != "" {
				billingAddress += (invoice.BillAddr.Line2 + "  \n")
			}
			if invoice.BillAddr.Line3 != "" {
				billingAddress += (invoice.BillAddr.Line3 + "  \n")
			}
			if invoice.BillAddr.Line4 != "" {
				billingAddress += (invoice.BillAddr.Line4 + "  \n")
			}
			if invoice.BillAddr.Line5 != "" {
				billingAddress += (invoice.BillAddr.Line5 + "  \n")
			}
			var shippingAddress string
			if invoice.ShipAddr.Line1 != "" {
				shippingAddress += (invoice.ShipAddr.Line1 + "  \n")
			}
			if invoice.ShipAddr.Line2 != "" {
				shippingAddress += (invoice.ShipAddr.Line2 + "  \n")
			}
			if invoice.ShipAddr.Line3 != "" {
				shippingAddress += (invoice.ShipAddr.Line3 + "  \n")
			}
			if invoice.ShipAddr.Line4 != "" {
				shippingAddress += (invoice.ShipAddr.Line4 + "  \n")
			}
			if invoice.ShipAddr.Line5 != "" {
				shippingAddress += (invoice.ShipAddr.Line5 + "  \n")
			}

			data = map[string]any{
				"id":                   invoice.Id,
				"qbo_id":               invoice.Id,
				"customer_id":          invoice.CustomerRef.Value,
				"sync_token":           invoice.SyncToken,
				"shipping_address":     shippingAddress,
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
				"billing_address":      billingAddress,
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
				"date":                 invoice.TxnDate.Format(fibery.DateFormat),
				"due_date":             invoice.DueDate.Format(fibery.DateFormat),
				"class_id":             invoice.ClassRef.Value,
				"print_status":         printStatus[invoice.PrintStatus],
				"term_id":              invoice.SalesTermRef.Value,
				"statement_memo":       invoice.PrivateNote,
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
					dependent := (*dependentPointer).(DepWHReceivable)
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
	DependentBaseType: DependentBaseType{
		FiberyType: FiberyType{
			id:   "Invoice_line",
			name: "Invoice Line",
			schema: map[string]fibery.Field{
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
					Relation: &fibery.Relation{
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
					Relation: &fibery.Relation{
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
					Relation: &fibery.Relation{
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
					Relation: &fibery.Relation{
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
					Relation: &fibery.Relation{
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
		schemaGen: func(entity any, source any) (map[string]any, error) {
			line, ok := entity.(quickbooks.Line)
			if !ok {
				return nil, fmt.Errorf("unable to convert entity to invoice line")
			}
			invoice, ok := source.(quickbooks.Invoice)
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
		queryProcessor: func(sourceArray any, schemaGen depSchemaGenFunc) ([]map[string]any, error) {
			invoices, ok := sourceArray.([]quickbooks.Invoice)
			if !ok {
				return nil, fmt.Errorf("unable to convert sourceArray to invoices")
			}
			items := []map[string]any{}
			for _, invoice := range invoices {
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
}

func init() {
	RegisterType(Invoice)
	RegisterType(InvoiceLine)
}
