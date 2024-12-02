// Copyright (c) 2018, Randy Westlund. All rights reserved.
// This code is under the BSD-2-Clause license.

package qbo

import (
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"strconv"

	"github.com/patrickmn/go-cache"
)

var InvoiceType = DataType{
	ID:   "invoice",
	Name: "Invoice",
	Schema: map[string]Field{
		"id": {
			Name: "id",
			Type: ID,
		},
		"qbo_id": {
			Name: "QBO ID",
			Type: Text,
		},
		"name": {
			Name: "Name",
			Type: Text,
		},
		"customer_id": {
			Name: "Customer ID",
			Type: Text,
			Relation: &Relation{
				Cardinality:   MTO,
				Name:          "Customer",
				TargetName:    "Invoices",
				TargetType:    "customer",
				TargetFieldID: "id",
			},
		},
		"sync_token": {
			Name:     "Sync Token",
			Type:     Text,
			ReadOnly: true,
		},
		"shipping_line_1": {
			Name: "Shipping Line 1",
			Type: Text,
		},
		"shipping_line_2": {
			Name: "Shipping Line 2",
			Type: Text,
		},
		"shipping_line_3": {
			Name: "Shipping Line 3",
			Type: Text,
		},
		"shipping_line_4": {
			Name: "Shipping Line 4",
			Type: Text,
		},
		"shipping_line_5": {
			Name: "Shipping Line 5",
			Type: Text,
		},
		"shipping_city": {
			Name: "Shipping City",
			Type: Text,
		},
		"shipping_state": {
			Name: "Shipping State",
			Type: Text,
		},
		"shipping_postal_code": {
			Name: "Shipping Postal Code",
			Type: Text,
		},
		"shipping_country": {
			Name: "Shipping Country",
			Type: Text,
		},
		"shipping_lat": {
			Name: "Shipping Latitude",
			Type: Text,
		},
		"shipping_long": {
			Name: "Shipping Longitude",
			Type: Text,
		},
		"billing_line_1": {
			Name: "Billing Line 1",
			Type: Text,
		},
		"billing_line_2": {
			Name: "Billing Line 2",
			Type: Text,
		},
		"billing_line_3": {
			Name: "Billing Line 3",
			Type: Text,
		},
		"billing_line_4": {
			Name: "Billing Line 4",
			Type: Text,
		},
		"billing_line_5": {
			Name: "Billing Line 5",
			Type: Text,
		},
		"billing_city": {
			Name: "Billing City",
			Type: Text,
		},
		"billing_state": {
			Name: "Billing State",
			Type: Text,
		},
		"billing_postal_code": {
			Name: "Billing Postal Code",
			Type: Text,
		},
		"billing_country": {
			Name: "Billing Country",
			Type: Text,
		},
		"billing_lat": {
			Name: "Billing Latitude",
			Type: Text,
		},
		"billing_long": {
			Name: "Billing Longitude",
			Type: Text,
		},
		"invoice_num": {
			Name: "Invoice",
			Type: Text,
		},
		"email": {
			Name:    "To",
			Type:    Text,
			SubType: Email,
		},
		"email_cc": {
			Name:    "CC",
			Type:    Text,
			SubType: Email,
		},
		"email_bcc": {
			Name:    "BCC",
			Type:    Text,
			SubType: Email,
		},
		"email_status": {
			Name:     "Email Status",
			Type:     Text,
			SubType:  SingleSelect,
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
			Type: DateType,
		},
		"date": {
			Name: "Date",
			Type: DateType,
		},
		"due_date": {
			Name: "Due Date",
			Type: DateType,
		},
		"class_id": {
			Name: "Class ID",
			Type: Text,
			Relation: &Relation{
				Cardinality:   MTO,
				Name:          "Class",
				TargetName:    "Invoices",
				TargetType:    "class",
				TargetFieldID: "id",
			},
		},
		"print_status": {
			Name:     "Print Status",
			Type:     Text,
			SubType:  SingleSelect,
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
			Type: Text,
			Relation: &Relation{
				Cardinality:   MTO,
				Name:          "Term",
				TargetName:    "Invoices",
				TargetType:    "term",
				TargetFieldID: "id",
			},
		},
		"statement_memo": {
			Name: "Statement Message",
			Type: Text,
		},
		"customer_memo": {
			Name: "Invoice Message",
			Type: Text,
		},
		"allow_ach": {
			Name:    "ACH Payments",
			Type:    Text,
			SubType: Boolean,
		},
		"allow_cc": {
			Name:    "Credit Card Payments",
			Type:    Text,
			SubType: Boolean,
		},
		"tax_code_id": {
			Name: "Tax Code ID",
			Type: Text,
			Relation: &Relation{
				Cardinality:   MTO,
				Name:          "Tax Code",
				TargetName:    "Invoices",
				TargetType:    "tax_code",
				TargetFieldID: "id",
			},
		},
		"tax_position": {
			Name:     "Apply Tax",
			Type:     Text,
			SubType:  SingleSelect,
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
			Type: Text,
			Relation: &Relation{
				Cardinality:   MTO,
				Name:          "Tax Exemption",
				TargetName:    "Invoices",
				TargetType:    "tax_exemption",
				TargetFieldID: "id",
			},
		},
		"deposit_account_id": {
			Name: "Deposit Account ID",
			Type: Text,
			Relation: &Relation{
				Cardinality:   MTO,
				Name:          "Deposit Account",
				TargetName:    "Invoice Deposits",
				TargetType:    "account",
				TargetFieldID: "id",
			},
		},
		"deposit_field": {
			Name: "Deposit Amount",
			Type: Number,
			Format: map[string]any{
				"format":               "Money",
				"currencyCode":         "USD",
				"hasThousandSeperator": true,
				"precision":            2,
			},
		},
		"discount_type": {
			Name:     "Discount Type",
			Type:     Text,
			SubType:  SingleSelect,
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
			Type: Number,
			Format: map[string]any{
				"format":    "Percent",
				"precision": 2,
			},
		},
		"discount_amount": {
			Name: "Discount Amount",
			Type: Number,
			Format: map[string]any{
				"format":               "Money",
				"currencyCode":         "USD",
				"hasThousandSeperator": true,
				"precision":            2,
			},
		},
		"tax": {
			Name: "Tax",
			Type: Number,
			Format: map[string]any{
				"format":               "Money",
				"currencyCode":         "USD",
				"hasThousandSeperator": true,
				"precision":            2,
			},
		},
		"subtotal": {
			Name: "Subtotal",
			Type: Number,
			Format: map[string]any{
				"format":               "Money",
				"currencyCode":         "USD",
				"hasThousandSeperator": true,
				"precision":            2,
			},
		},
		"total": {
			Name: "Total",
			Type: Number,
			Format: map[string]any{
				"format":               "Money",
				"currencyCode":         "USD",
				"hasThousandSeperator": true,
				"precision":            2,
			},
		},
		"balance": {
			Name: "Balance",
			Type: Number,
			Format: map[string]any{
				"format":               "Money",
				"currencyCode":         "USD",
				"hasThousandSeperator": true,
				"precision":            2,
			},
		},
		"created_qbo": {
			Name: "Creation Date (QBO)",
			Type: DateType,
		},
		"last_updated_qbo": {
			Name: "Last Updated (QBO)",
			Type: DateType,
		},
		"__syncAction": {
			Type: Text,
			Name: "Sync Action",
		},
	},
	DataRequest: func(req RequestParameters) ([]map[string]any, bool, error) {
		return getInvoiceData("invoice", req)
	},
}

var InvoiceLineType = DataType{
	ID:   "invoice_line",
	Name: "Invoice Line",
	Schema: map[string]Field{
		"id": {
			Name: "id",
			Type: ID,
		},
		"qbo_id": {
			Name: "QBO ID",
			Type: Text,
		},
		"invoice_id": {
			Name: "Invoice ID",
			Type: Text,
			Relation: &Relation{
				Cardinality:   MTO,
				Name:          "Invoice",
				TargetName:    "Invoice Lines",
				TargetType:    "invoice",
				TargetFieldID: "id",
			},
		},
		"description": {
			Name:    "Description",
			Type:    Text,
			SubType: Title,
		},
		"line_type": {
			Name:    "Line Type",
			Type:    Text,
			SubType: SingleSelect,
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
			Type: Number,
			Format: map[string]any{
				"format":               "Number",
				"hasThousandSeparator": true,
				"precision":            2,
			},
		},
		"unit_price": {
			Name: "Unit Price",
			Type: Number,
			Format: map[string]any{
				"format":               "Money",
				"currencyCode":         "USD",
				"hasThousandSeperator": true,
				"precision":            2,
			},
		},
		"amount": {
			Name: "Amount",
			Type: Number,
			Format: map[string]any{
				"format":               "Money",
				"currencyCode":         "USD",
				"hasThousandSeperator": true,
				"precision":            2,
			},
		},
		"line_num": {
			Name:    "Line",
			Type:    Number,
			SubType: Integer,
		},
		"group_line_id": {
			Name: "Group Line ID",
			Type: Text,
			Relation: &Relation{
				Cardinality:   MTO,
				Name:          "Group",
				TargetName:    "Lines",
				TargetType:    "invoice_line",
				TargetFieldID: "id",
			},
		},
		"item_id": {
			Name: "Item",
			Type: Text,
			Relation: &Relation{
				Cardinality:   MTO,
				Name:          "Item",
				TargetName:    "Invoice Lines",
				TargetType:    "item",
				TargetFieldID: "id",
			},
		},
		"class_id": {
			Name: "Class ID",
			Type: Text,
			Relation: &Relation{
				Cardinality:   MTO,
				Name:          "Class",
				TargetName:    "Expense Account Line(s)",
				TargetType:    "class",
				TargetFieldID: "id",
			},
		},
		"tax_code_id": {
			Name: "Tax Code ID",
			Type: Text,
			Relation: &Relation{
				Cardinality:   MTO,
				Name:          "Tax Code",
				TargetName:    "Invoice Lines",
				TargetType:    "tax_code",
				TargetFieldID: "id",
			},
		},
		"markup_percent": {
			Name: "Markup",
			Type: Number,
			Format: map[string]any{
				"format":    "Percent",
				"precision": 2,
			},
		},
		"service_date": {
			Name:    "Date",
			Type:    DateType,
			SubType: Day,
		},
		"__syncAction": {
			Type: Text,
			Name: "Sync Action",
		},
	},
	DataRequest: func(req RequestParameters) ([]map[string]any, bool, error) {
		return getInvoiceData("invoice_line", req)
	},
}

func getInvoiceData(subtypeID string, req RequestParameters) (data []map[string]any, morePages bool, err error) {
	// need to reconfigure so that singleflight is used to limit invoice requests, not for invoice conversion
	convertInvoiceData := func(subTypeID string, sync SyncType, invoices []Invoice) ([]map[string]any, error) {
		switch subTypeID {
		case "invoice":
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
			var data []map[string]any
			for _, invoice := range invoices {
				var discountLine *Line
				var subtotalLine *Line
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
					emailSendTime = invoice.DeliveryInfo.DeliveryTime.Format(fiberyDateFormat)
				}

				data = append(data, map[string]any{
					"id":                   invoice.Id,
					"qbo_id":               invoice.Id,
					"name":                 invoice.DocNumber + " - " + invoice.CustomerRef.Name,
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
					"date":                 invoice.TxnDate.Format(fiberyDateFormat),
					"due_date":             invoice.DueDate.Format(fiberyDateFormat),
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
					"created_qbo":          invoice.MetaData.CreateTime.Format(fiberyDateFormat),
					"last_updated_qbo":     invoice.MetaData.LastUpdatedTime.Format(fiberyDateFormat),
					"__syncAction":         sync,
				})
			}
			return data, nil
		case "invoice_line":
			var lineTypes = map[string]string{
				"SalesItemLineDetail": "Sales Item Line",
				"GroupLineDetail":     "Group Line",
				"DescriptionOnly":     "Description Line",
			}

			var data []map[string]any
			for _, invoice := range invoices {
				for _, line := range invoice.Line {
					if line.DetailType == "GroupLineDetail" {
						data = append(data, map[string]any{
							"id":           fmt.Sprintf("%s:%s", invoice.Id, line.Id),
							"qbo_id":       line.Id,
							"invoice_id":   invoice.Id,
							"description":  line.Description,
							"line_type":    lineTypes[line.DetailType],
							"quantity":     line.GroupLineDetail.Quantity,
							"line_num":     line.LineNum,
							"item_id":      line.GroupLineDetail.GroupItemRef.Value,
							"__syncAction": sync,
						})
						for _, groupLine := range line.GroupLineDetail.Line {
							data = append(data, map[string]any{
								"id":             fmt.Sprintf("%s:%s:%s", invoice.Id, line.Id, groupLine.Id),
								"qbo_id":         groupLine.Id,
								"invoice_id":     invoice.Id,
								"description":    groupLine.Description,
								"line_type":      lineTypes[groupLine.DetailType],
								"quantity":       groupLine.SalesItemLineDetail.Qty,
								"unit_price":     groupLine.SalesItemLineDetail.UnitPrice,
								"amount":         groupLine.Amount,
								"line_num":       groupLine.LineNum,
								"group_line_id":  line.Id,
								"item_id":        groupLine.SalesItemLineDetail.ItemRef.Value,
								"class_id":       groupLine.SalesItemLineDetail.ClassRef.Value,
								"tax_code_id":    groupLine.SalesItemLineDetail.TaxCodeRef.Value,
								"markup_percent": groupLine.SalesItemLineDetail.MarkupInfo.Percent,
								"service_date":   groupLine.SalesItemLineDetail.ServiceDate.Format(fiberyDateFormat),
								"__syncAction":   sync,
							})
						}
					} else if line.DetailType == "DescriptionOnly" || line.DetailType == "SalesItemLineDetail" {
						data = append(data, map[string]any{
							"id":             fmt.Sprintf("%s:%s", invoice.Id, line.Id),
							"qbo_id":         line.Id,
							"invoice_id":     invoice.Id,
							"description":    line.Description,
							"type":           lineTypes[line.DetailType],
							"quantity":       line.SalesItemLineDetail.Qty,
							"unit_price":     line.SalesItemLineDetail.UnitPrice,
							"amount":         line.Amount,
							"line_num":       line.LineNum,
							"item_id":        line.SalesItemLineDetail.ItemRef.Value,
							"class_id":       line.SalesItemLineDetail.ClassRef.Value,
							"tax_code_id":    line.SalesItemLineDetail.TaxCodeRef.Value,
							"markup_percent": line.SalesItemLineDetail.MarkupInfo.Percent,
							"service_date":   line.SalesItemLineDetail.ServiceDate.Format(fiberyDateFormat),
							"__syncAction":   sync,
						})
					}
				}
			}
			return data, nil
		default:
			return nil, fmt.Errorf("invalid subtype: %s", subTypeID)
		}
	}

	groupKey := fmt.Sprintf("%s:%s", req.OperationID, "invoice")
	cacheKey := fmt.Sprintf("%s:%s:%d", req.OperationID, "invoice", req.StartPosition)

	sync := fullSync
	if req.LastSynced != "" {
		sync = deltaSync
	}

	if cacheEntryInterface, exists := req.Cache.Get(cacheKey); exists {
		slog.Info(fmt.Sprintf("Cache hit for %s", cacheKey))
		cacheEntry := cacheEntryInterface.(*CacheEntry[Invoice])

		cacheEntry.mu.Lock()
		defer cacheEntry.mu.Unlock()

		invoices := cacheEntry.Data
		more := cacheEntry.More

		cacheEntry.ProcessedTypes[subtypeID] = true
		req.Cache.Set(cacheKey, cacheEntry, cache.DefaultExpiration)

		// If all subtypes have been processed, remove the cache entry
		if allSubtypesProcessed(cacheEntry.ProcessedTypes) {
			req.Cache.Delete(cacheKey)
			slog.Info(fmt.Sprintf("Deleted cache entry for %s", cacheKey))
		}

		data, err := convertInvoiceData(subtypeID, sync, invoices)
		if err != nil {
			return nil, false, err
		}

		return data, more, nil
	}

	type result struct {
		data []Invoice
		more bool
	}

	res, err, _ := req.Group.Do(groupKey, func() (interface{}, error) {
		if cacheEntryInterface, exists := req.Cache.Get(cacheKey); exists {
			slog.Info(fmt.Sprintf("Cache hit for %s", cacheKey))
			cacheEntry := cacheEntryInterface.(*CacheEntry[Invoice])

			cacheEntry.mu.Lock()
			defer cacheEntry.mu.Unlock()

			invoices := cacheEntry.Data
			more := cacheEntry.More

			cacheEntry.ProcessedTypes[subtypeID] = true

			req.Cache.Set(cacheKey, cacheEntry, cache.DefaultExpiration)

			if allSubtypesProcessed(cacheEntry.ProcessedTypes) {
				req.Cache.Delete(cacheKey)
				slog.Info(fmt.Sprintf("Deleted cache entry for %s", cacheKey))
			}

			return result{invoices, more}, nil
		}

		client, err := NewClient(req.RealmID, req.Token)
		if err != nil {
			return nil, fmt.Errorf("unable to create new client: %w", err)
		}

		var query string

		if sync == fullSync {
			query = fmt.Sprintf("SELECT * FROM Invoice STARTPOSITION %d MAXRESULTS %d", req.StartPosition, QueryPageSize)
		} else {
			query = fmt.Sprintf("SELECT * FROM Invoice WHERE MetaData.LastUpdatedTime >= '%s' STARTPOSITION %d MAXRESULTS %d", req.LastSynced, req.StartPosition, QueryPageSize)
		}

		var resp struct {
			QueryResponse struct {
				Invoices      []Invoice `json:"Invoice"`
				StartPosition int
				MaxResults    int
			}
		}

		if err := client.query(query, &resp); err != nil {
			return nil, err
		}

		invoices := resp.QueryResponse.Invoices
		more := len(invoices) == QueryPageSize
		processedTypes := map[string]bool{
			"invoice":      false,
			"invoice_line": false,
		}

		entry := &CacheEntry[Invoice]{
			Data:           invoices,
			ProcessedTypes: processedTypes,
			More:           more,
		}

		entry.ProcessedTypes[subtypeID] = true

		req.Cache.Set(cacheKey, entry, cache.DefaultExpiration)
		slog.Info(fmt.Sprintf("Created cache entry for %s", cacheKey))

		return result{invoices, more}, nil
	})
	if err != nil {
		return nil, false, err
	}

	data, err = convertInvoiceData(subtypeID, sync, res.(result).data)
	if err != nil {
		return nil, false, err
	}

	return data, res.(result).more, nil
}

func init() {
	InvoiceType.Register()
	InvoiceLineType.Register()
}

// Invoice represents a QuickBooks Invoice object.
type Invoice struct {
	Id            string        `json:"Id,omitempty"`
	SyncToken     string        `json:",omitempty"`
	MetaData      MetaData      `json:",omitempty"`
	CustomField   []CustomField `json:",omitempty"`
	DocNumber     string        `json:",omitempty"`
	TxnDate       Date          `json:",omitempty"`
	DepartmentRef ReferenceType `json:",omitempty"`
	PrivateNote   string        `json:",omitempty"`
	LinkedTxn     []LinkedTxn   `json:"LinkedTxn"`
	Line          []Line
	TxnTaxDetail  TxnTaxDetail `json:",omitempty"`
	CustomerRef   ReferenceType
	CustomerMemo  MemoRef         `json:",omitempty"`
	BillAddr      PhysicalAddress `json:",omitempty"`
	ShipAddr      PhysicalAddress `json:",omitempty"`
	ClassRef      ReferenceType   `json:",omitempty"`
	SalesTermRef  ReferenceType   `json:",omitempty"`
	DueDate       Date            `json:",omitempty"`
	// GlobalTaxCalculation
	ShipMethodRef                ReferenceType `json:",omitempty"`
	ShipDate                     Date          `json:",omitempty"`
	TrackingNum                  string        `json:",omitempty"`
	TotalAmt                     json.Number   `json:",omitempty"`
	CurrencyRef                  ReferenceType `json:",omitempty"`
	ExchangeRate                 json.Number   `json:",omitempty"`
	HomeAmtTotal                 json.Number   `json:",omitempty"`
	HomeBalance                  json.Number   `json:",omitempty"`
	ApplyTaxAfterDiscount        bool          `json:",omitempty"`
	PrintStatus                  string        `json:",omitempty"`
	EmailStatus                  string        `json:",omitempty"`
	BillEmail                    EmailAddress  `json:",omitempty"`
	BillEmailCC                  EmailAddress  `json:"BillEmailCc,omitempty"`
	BillEmailBCC                 EmailAddress  `json:"BillEmailBcc,omitempty"`
	DeliveryInfo                 *DeliveryInfo `json:",omitempty"`
	TaxExemptionRef              ReferenceType `json:",omitempty"`
	Balance                      json.Number   `json:",omitempty"`
	TxnSource                    string        `json:",omitempty"`
	AllowOnlineCreditCardPayment bool          `json:",omitempty"`
	AllowOnlineACHPayment        bool          `json:",omitempty"`
	Deposit                      json.Number   `json:",omitempty"`
	DepositToAccountRef          ReferenceType `json:",omitempty"`
}

type MarkupInfo struct {
	PriceLevelRef          ReferenceType `json:",omitempty"`
	Percent                json.Number   `json:",omitempty"`
	MarkUpIncomeAccountRef ReferenceType `json:",omitempty"`
}

type DeliveryInfo struct {
	DeliveryType string
	DeliveryTime Date
}

type LinkedTxn struct {
	TxnID   string `json:"TxnId"`
	TxnType string `json:"TxnType"`
}

type TxnTaxDetail struct {
	TxnTaxCodeRef ReferenceType `json:",omitempty"`
	TotalTax      json.Number   `json:",omitempty"`
	TaxLine       []Line        `json:",omitempty"`
}

type AccountBasedExpenseLineDetail struct {
	AccountRef ReferenceType
	TaxAmount  json.Number `json:",omitempty"`
	// TaxInclusiveAmt json.Number              `json:",omitempty"`
	// ClassRef        ReferenceType `json:",omitempty"`
	// TaxCodeRef      ReferenceType `json:",omitempty"`
	// MarkupInfo MarkupInfo `json:",omitempty"`
	// BillableStatus BillableStatusEnum       `json:",omitempty"`
	// CustomerRef    ReferenceType `json:",omitempty"`
}

type Line struct {
	Id                            string `json:",omitempty"`
	LineNum                       int    `json:",omitempty"`
	Description                   string `json:",omitempty"`
	Amount                        json.Number
	DetailType                    string
	AccountBasedExpenseLineDetail AccountBasedExpenseLineDetail `json:",omitempty"`
	SalesItemLineDetail           SalesItemLineDetail           `json:",omitempty"`
	GroupLineDetail               GroupLineDetail               `json:",omitempty"`
	DiscountLineDetail            DiscountLineDetail            `json:",omitempty"`
	TaxLineDetail                 TaxLineDetail                 `json:",omitempty"`
}

// TaxLineDetail ...
type TaxLineDetail struct {
	PercentBased     bool        `json:",omitempty"`
	NetAmountTaxable json.Number `json:",omitempty"`
	// TaxInclusiveAmount json.Number `json:",omitempty"`
	// OverrideDeltaAmount
	TaxPercent json.Number `json:",omitempty"`
	TaxRateRef ReferenceType
}

// SalesItemLineDetail ...
type SalesItemLineDetail struct {
	ItemRef         ReferenceType `json:",omitempty"`
	ClassRef        ReferenceType `json:",omitempty"`
	UnitPrice       json.Number   `json:",omitempty"`
	MarkupInfo      MarkupInfo    `json:",omitempty"`
	Qty             float32       `json:",omitempty"`
	ItemAccountRef  ReferenceType `json:",omitempty"`
	TaxCodeRef      ReferenceType `json:",omitempty"`
	ServiceDate     Date          `json:",omitempty"`
	TaxInclusiveAmt json.Number   `json:",omitempty"`
	DiscountRate    json.Number   `json:",omitempty"`
	DiscountAmt     json.Number   `json:",omitempty"`
}

// GroupLineDetail ...
type GroupLineDetail struct {
	Quantity     float32       `json:",omitempty"`
	GroupItemRef ReferenceType `json:",omitempty"`
	Line         []Line        `json:",omitempty"`
}

// DiscountLineDetail ...
type DiscountLineDetail struct {
	PercentBased    bool
	DiscountPercent float32 `json:",omitempty"`
}

// CreateInvoice creates the given Invoice on the QuickBooks server, returning
// the resulting Invoice object.
func (c *Client) CreateInvoice(invoice *Invoice) (*Invoice, error) {
	var resp struct {
		Invoice Invoice
		Time    Date
	}

	if err := c.post("invoice", invoice, &resp, nil); err != nil {
		return nil, err
	}

	return &resp.Invoice, nil
}

// DeleteInvoice deletes the invoice
//
// If the invoice was already deleted, QuickBooks returns 400 :(
// The response looks like this:
// {"Fault":{"Error":[{"Message":"Object Not Found","Detail":"Object Not Found : Something you're trying to use has been made inactive. Check the fields with accounts, invoices, items, vendors or employees.","code":"610","element":""}],"type":"ValidationFault"},"time":"2018-03-20T20:15:59.571-07:00"}
//
// This is slightly horrifying and not documented in their API. When this
// happens we just return success; the goal of deleting it has been
// accomplished, just not by us.
func (c *Client) DeleteInvoice(invoice *Invoice) error {
	if invoice.Id == "" || invoice.SyncToken == "" {
		return errors.New("missing id/sync token")
	}

	return c.post("invoice", invoice, nil, map[string]string{"operation": "delete"})
}

// FindInvoices gets the full list of Invoices in the QuickBooks account.
func (c *Client) FindInvoices() ([]Invoice, error) {
	var resp struct {
		QueryResponse struct {
			Invoices      []Invoice `json:"Invoice"`
			MaxResults    int
			StartPosition int
			TotalCount    int
		}
	}

	if err := c.query("SELECT COUNT(*) FROM Invoice", &resp); err != nil {
		return nil, err
	}

	if resp.QueryResponse.TotalCount == 0 {
		return nil, errors.New("no invoices could be found")
	}

	invoices := make([]Invoice, 0, resp.QueryResponse.TotalCount)

	for i := 0; i < resp.QueryResponse.TotalCount; i += QueryPageSize {
		query := "SELECT * FROM Invoice ORDERBY Id STARTPOSITION " + strconv.Itoa(i+1) + " MAXRESULTS " + strconv.Itoa(QueryPageSize)

		if err := c.query(query, &resp); err != nil {
			return nil, err
		}

		if resp.QueryResponse.Invoices == nil {
			return nil, errors.New("no invoices could be found")
		}

		invoices = append(invoices, resp.QueryResponse.Invoices...)
	}

	return invoices, nil
}

// FindInvoiceById finds the invoice by the given id
func (c *Client) FindInvoiceById(id string) (*Invoice, error) {
	var resp struct {
		Invoice Invoice
		Time    Date
	}

	if err := c.get("invoice/"+id, &resp, nil); err != nil {
		return nil, err
	}

	return &resp.Invoice, nil
}

// QueryInvoices accepts an SQL query and returns all invoices found using it
func (c *Client) QueryInvoices(query string) ([]Invoice, error) {
	var resp struct {
		QueryResponse struct {
			Invoices      []Invoice `json:"Invoice"`
			StartPosition int
			MaxResults    int
		}
	}

	if err := c.query(query, &resp); err != nil {
		return nil, err
	}

	if resp.QueryResponse.Invoices == nil {
		return nil, errors.New("could not find any invoices")
	}

	return resp.QueryResponse.Invoices, nil
}

// SendInvoice sends the invoice to the Invoice.BillEmail if emailAddress is left empty
func (c *Client) SendInvoice(invoiceId string, emailAddress string) error {
	queryParameters := make(map[string]string)

	if emailAddress != "" {
		queryParameters["sendTo"] = emailAddress
	}

	return c.post("invoice/"+invoiceId+"/send", nil, nil, queryParameters)
}

// UpdateInvoice updates the invoice
func (c *Client) UpdateInvoice(invoice *Invoice) (*Invoice, error) {
	if invoice.Id == "" {
		return nil, errors.New("missing invoice id")
	}

	existingInvoice, err := c.FindInvoiceById(invoice.Id)
	if err != nil {
		return nil, err
	}

	invoice.SyncToken = existingInvoice.SyncToken

	payload := struct {
		*Invoice
		Sparse bool `json:"sparse"`
	}{
		Invoice: invoice,
		Sparse:  true,
	}

	var invoiceData struct {
		Invoice Invoice
		Time    Date
	}

	if err = c.post("invoice", payload, &invoiceData, nil); err != nil {
		return nil, err
	}

	return &invoiceData.Invoice, err
}

func (c *Client) VoidInvoice(invoice Invoice) error {
	if invoice.Id == "" {
		return errors.New("missing invoice id")
	}

	existingInvoice, err := c.FindInvoiceById(invoice.Id)
	if err != nil {
		return err
	}

	invoice.SyncToken = existingInvoice.SyncToken

	return c.post("invoice", invoice, nil, map[string]string{"operation": "void"})
}
