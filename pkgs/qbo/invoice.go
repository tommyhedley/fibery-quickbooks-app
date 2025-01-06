// Copyright (c) 2018, Randy Westlund. All rights reserved.
// This code is under the BSD-2-Clause license.

package qbo

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"sync"
)

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
	Domain                       string        `json:"domain"`
	Status                       string        `json:"status"`
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

// FindInvoicesByPage gets a page of invoices from the QuickBooks account at the current max results size.
func (c *Client) FindInvoicesByPage(StartPosition int) ([]Invoice, error) {
	var resp struct {
		QueryResponse struct {
			Invoices      []Invoice `json:"Invoice"`
			MaxResults    int
			StartPosition int
			TotalCount    int
		}
	}

	query := "SELECT * FROM Invoice ORDERBY Id STARTPOSITION " + strconv.Itoa(StartPosition) + " MAXRESULTS " + strconv.Itoa(QueryPageSize)

	if err := c.query(query, &resp); err != nil {
		return nil, err
	}

	if resp.QueryResponse.Invoices == nil {
		return nil, errors.New("no invoices could be found")
	}

	return resp.QueryResponse.Invoices, nil
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

// Invoice FiberyType Implementation

func (Invoice) ID() string {
	return "Invoice"
}

func (Invoice) Name() string {
	return "Invoice"
}

func (Invoice) Schema() map[string]Field {
	return map[string]Field{
		"id": {
			Name: "id",
			Type: ID,
		},
		"qbo_id": {
			Name: "QBO ID",
			Type: Text,
		},
		"customer_id": {
			Name: "Customer ID",
			Type: Text,
			Relation: &Relation{
				Cardinality:   MTO,
				Name:          "Customer",
				TargetName:    "Invoices",
				TargetType:    "Customer",
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
			Name:    "Invoice",
			Type:    Text,
			SubType: Title,
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
				TargetType:    "Class",
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
				TargetType:    "Term",
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
				TargetType:    "TaxCode",
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
				TargetType:    "TaxExemption",
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
				TargetType:    "Account",
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
	}
}

func (I Invoice) Dependents() map[string][]DependentDataType {
	children := make(map[string][]DependentDataType)
	lineCount := len(I.Line)
	if lineCount > 0 {
		lines := make([]DependentDataType, 0, lineCount)
		for i := range I.Line {
			lines = append(lines, I.Line[i])
		}
		children[lines[0].ID()] = lines
	}
	return children
}

func (Invoice) getFullData(req *DataRequest) (DataResponse, error) {
	client, err := NewClient(req.RealmID, req.Token)
	if err != nil {
		return DataResponse{}, fmt.Errorf("unable to create new client: %w", err)
	}

	invoices, err := client.FindInvoicesByPage(req.StartPosition)
	if err != nil {
		return DataResponse{}, fmt.Errorf("unable to find invoices: %w", err)
	}

	return DataResponse{
		Data: invoices,
		More: len(invoices) >= QueryPageSize,
	}, nil
}

func (I Invoice) TransformItem() (map[string]any, error) {
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
	var discountLine *InvoiceLine
	var subtotalLine *InvoiceLine
	for _, line := range I.Line {
		if line.DetailType == "DiscountLineDetail" {
			if discountLine != nil {
				return nil, fmt.Errorf("invoice %s has more than one discount line", I.Id)
			}
			discountLine = &line
		}
		if line.DetailType == "SubTotalLineDetail" {
			if subtotalLine != nil {
				return nil, fmt.Errorf("invoice %s has more than one subtotal line", I.Id)
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
	if I.DeliveryInfo != nil && !I.DeliveryInfo.DeliveryTime.IsZero() {
		emailSendTime = I.DeliveryInfo.DeliveryTime.Format(fiberyDateFormat)
	}

	data = map[string]any{
		"id":                   I.Id,
		"qbo_id":               I.Id,
		"customer_id":          I.CustomerRef.Value,
		"sync_token":           I.SyncToken,
		"shipping_line_1":      I.ShipAddr.Line1,
		"shipping_line_2":      I.ShipAddr.Line2,
		"shipping_line_3":      I.ShipAddr.Line3,
		"shipping_line_4":      I.ShipAddr.Line4,
		"shipping_line_5":      I.ShipAddr.Line5,
		"shipping_city":        I.ShipAddr.City,
		"shipping_state":       I.ShipAddr.CountrySubDivisionCode,
		"shipping_postal_code": I.ShipAddr.PostalCode,
		"shipping_country":     I.ShipAddr.Country,
		"shipping_lat":         I.ShipAddr.Lat,
		"shipping_long":        I.ShipAddr.Long,
		"billing_line_1":       I.BillAddr.Line1,
		"billing_line_2":       I.BillAddr.Line2,
		"billing_line_3":       I.BillAddr.Line3,
		"billing_line_4":       I.BillAddr.Line4,
		"billing_line_5":       I.BillAddr.Line5,
		"billing_city":         I.BillAddr.City,
		"billing_state":        I.BillAddr.CountrySubDivisionCode,
		"billing_postal_code":  I.BillAddr.PostalCode,
		"billing_country":      I.BillAddr.Country,
		"billing_lat":          I.BillAddr.Lat,
		"billing_long":         I.BillAddr.Long,
		"invoice_num":          I.DocNumber,
		"email":                I.BillEmail.Address,
		"email_cc":             I.BillEmailCC.Address,
		"email_bcc":            I.BillEmailBCC.Address,
		"email_status":         emailStatus[I.EmailStatus],
		"email_send_time":      emailSendTime,
		"date":                 I.TxnDate.Format(fiberyDateFormat),
		"due_date":             I.DueDate.Format(fiberyDateFormat),
		"class_id":             I.ClassRef.Value,
		"print_status":         printStatus[I.PrintStatus],
		"term_id":              I.SalesTermRef.Value,
		"statement_memo":       I.PrivateNote,
		"customer_memo":        I.CustomerMemo.Value,
		"allow_ach":            I.AllowOnlineACHPayment,
		"allow_cc":             I.AllowOnlineCreditCardPayment,
		"tax_code_id":          I.TxnTaxDetail.TxnTaxCodeRef.Value,
		"tax_position":         taxPosition[I.ApplyTaxAfterDiscount],
		"tax_exemption_id":     I.TaxExemptionRef.Value,
		"deposit_account_id":   I.DepositToAccountRef.Value,
		"deposit_field":        I.Deposit,
		"discount_type":        discountTypeValue,
		"discount_percent":     discountPercent,
		"discount_amount":      discountAmount,
		"tax":                  I.TxnTaxDetail.TotalTax,
		"subtotal":             subtotalAmount,
		"total":                I.TotalAmt,
		"balance":              I.Balance,
		"created_qbo":          I.MetaData.CreateTime.Format(fiberyDateFormat),
		"last_updated_qbo":     I.MetaData.LastUpdatedTime.Format(fiberyDateFormat),
		"__syncAction":         SET,
	}
	return data, nil
}

func (I Invoice) transformFullData(data DataResponse) ([]map[string]any, error) {
	invoices := data.Data.([]Invoice)
	items := []map[string]any{}
	for _, invoice := range invoices {
		item, err := invoice.TransformItem()
		if err != nil {
			return nil, fmt.Errorf("unable to invoice transform data: %w", err)
		}
		items = append(items, item)
	}
	return items, nil
}

func (I Invoice) transformChangeDataCapture(data ChangeDataCapture) ([]map[string]any, error) {
	items := []map[string]any{}
	for _, cdcResponse := range data.CDCResponse {
		for _, queryResponse := range cdcResponse.CDCQueryResponse {
			if queryResponse.QueryResponse.Items != nil {
				for _, item := range queryResponse.QueryResponse.Items {
					invoice, ok := item.(Invoice)
					if !ok {
						break
					}
					if invoice.Status == "Deleted" {
						items = append(items, map[string]any{
							"id":           invoice.Id,
							"__syncAction": REMOVE,
						})
					} else {
						item, err := invoice.TransformItem()
						if err != nil {
							return nil, fmt.Errorf("unable to invoice transform data: %w", err)
						}
						items = append(items, item)
					}
				}
			}
		}
	}
	return items, nil
}

func (I Invoice) GetData(req *DataRequest) (DataHandlerResponse, error) {
	if req.LastSynced == "" {
		groupKey := fmt.Sprintf("%s:%s:%d", req.OperationID, "invoice", req.StartPosition)
		res, err, _ := req.Group.Do(groupKey, func() (interface{}, error) {
			return I.getFullData(req)
		})

		if err != nil {
			return DataHandlerResponse{}, fmt.Errorf("unable to query invoices from qbo: %w", err)
		}

		items, err := I.transformFullData(res.(DataResponse))
		if err != nil {
			return DataHandlerResponse{}, fmt.Errorf("unable to transform data: %w", err)
		}

		return DataHandlerResponse{
			Items: items,
			Pagination: Pagination{
				HasNext: res.(DataResponse).More,
				NextPageConfig: NextPageConfig{
					StartPosition: req.StartPosition + QueryPageSize,
				},
			},
			SynchronizationType: FullSync,
		}, nil
	} else {
		groupKey := req.OperationID
		res, err, _ := req.Group.Do(groupKey, func() (interface{}, error) {
			return getChangeDataCapture(req)
		})

		if err != nil {
			return DataHandlerResponse{}, fmt.Errorf("unable to cdc query from qbo: %w", err)
		}

		items, err := I.transformChangeDataCapture(res.(ChangeDataCapture))
		if err != nil {
			return DataHandlerResponse{}, fmt.Errorf("unable to transform data: %w", err)
		}

		return DataHandlerResponse{
			Items: items,
			Pagination: Pagination{
				HasNext: false,
				NextPageConfig: NextPageConfig{
					StartPosition: 0,
				},
			},
			SynchronizationType: DeltaSync,
		}, nil
	}
}

func (I Invoice) TransformItemAndDependents() (map[string][]map[string]any, error) {
	res := make(map[string][]map[string]any)
	invoiceData, err := I.TransformItem()
	if err != nil {
		return nil, fmt.Errorf("unable to transform invoice item: %w", err)
	}
	res[I.ID()] = []map[string]any{invoiceData}
	dependents := I.Dependents()
	for typ, deps := range dependents {
		for _, dep := range deps {
			depData, err := dep.transformItem(I)
			if err != nil {
				return nil, fmt.Errorf("unable to transform %s item: %w", typ, err)
			}
			res[typ] = append(res[typ], depData)
		}
	}
	return res, nil
}

// InvoiceLine FiberyType Implementation

func (InvoiceLine) ID() string {
	return "Invoice_line"
}

func (InvoiceLine) Name() string {
	return "Invoice Line"
}

func (InvoiceLine) Schema() map[string]Field {
	return map[string]Field{
		"id": {
			Name: "id",
			Type: ID,
		},
		"qbo_id": {
			Name: "QBO ID",
			Type: Text,
		},
		"invoice_sync_token": {
			Name:     "Invoice Sync Token",
			Type:     Text,
			ReadOnly: true,
		},
		"invoice_id": {
			Name: "Invoice ID",
			Type: Text,
			Relation: &Relation{
				Cardinality:   MTO,
				Name:          "Invoice",
				TargetName:    "Invoice Lines",
				TargetType:    "Invoice",
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
				TargetType:    "Invoice_line",
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
				TargetType:    "Item",
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
				TargetType:    "Class",
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
				TargetType:    "TaxCode",
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
	}
}

func (InvoiceLine) Parent() string {
	return "Invoice"
}

func (InvoiceLine) getFullData(req *DataRequest) (DataResponse, error) {
	var i Invoice
	res, err := i.getFullData(req)
	if err != nil {
		return DataResponse{}, fmt.Errorf("unable to get %s data: %w", req.RequestedType, err)
	}
	return res, nil
}

func (Il InvoiceLine) transformItem(parent any) (map[string]any, error) {
	I := parent.(Invoice)
	var lineTypes = map[string]string{
		"SalesItemLineDetail": "Sales Item Line",
		"GroupLineDetail":     "Group Line",
		"DescriptionOnly":     "Description Line",
	}

	if Il.DetailType == "GroupLineDetail" {
		return map[string]any{
			"id":                 fmt.Sprintf("%s:%s", I.Id, Il.Id),
			"qbo_id":             Il.Id,
			"invoice_sync_token": I.SyncToken,
			"invoice_id":         I.Id,
			"description":        Il.Description,
			"line_type":          lineTypes[Il.DetailType],
			"quantity":           Il.GroupLineDetail.Quantity,
			"line_num":           Il.LineNum,
			"item_id":            Il.GroupLineDetail.GroupItemRef.Value,
			"__syncAction":       SET,
		}, nil
	} else if Il.DetailType == "DescriptionOnly" || Il.DetailType == "SalesItemLineDetail" {
		return map[string]any{
			"id":                 fmt.Sprintf("%s:%s", I.Id, Il.Id),
			"qbo_id":             Il.Id,
			"invoice_sync_token": I.SyncToken,
			"invoice_id":         I.Id,
			"description":        Il.Description,
			"type":               lineTypes[Il.DetailType],
			"quantity":           Il.SalesItemLineDetail.Qty,
			"unit_price":         Il.SalesItemLineDetail.UnitPrice,
			"amount":             Il.Amount,
			"line_num":           Il.LineNum,
			"item_id":            Il.SalesItemLineDetail.ItemRef.Value,
			"class_id":           Il.SalesItemLineDetail.ClassRef.Value,
			"tax_code_id":        Il.SalesItemLineDetail.TaxCodeRef.Value,
			"markup_percent":     Il.SalesItemLineDetail.MarkupInfo.Percent,
			"service_date":       Il.SalesItemLineDetail.ServiceDate.Format(fiberyDateFormat),
			"__syncAction":       SET,
		}, nil
	}
	return nil, nil
}

func (Il InvoiceLine) transformFullData(data DataResponse) ([]map[string]any, error) {
	invoices := data.Data.([]Invoice)
	items := []map[string]any{}
	for _, invoice := range invoices {
		for _, line := range invoice.Line {
			if line.DetailType == "DescriptionOnly" || line.DetailType == "SalesItemLineDetail" {
				item, err := line.transformItem(invoice)
				if err != nil {
					return nil, fmt.Errorf("unable to invoice line transform data: %w", err)
				}
				items = append(items, item)
			}
			if line.DetailType == "GroupLineDetail" {
				for _, groupLine := range line.GroupLineDetail.Line {
					item, err := groupLine.transformItem(invoice)
					if err != nil {
						return nil, fmt.Errorf("unable to invoice line transform data: %w", err)
					}
					item["id"] = fmt.Sprintf("%s:%s:%s", invoice.Id, line.Id, groupLine.Id)
					item["group_line_id"] = line.Id
					items = append(items, item)
				}
				item, err := line.transformItem(invoice)
				if err != nil {
					return nil, fmt.Errorf("unable to invoice line transform data: %w", err)
				}
				items = append(items, item)
			}
		}
	}
	return items, nil
}

func (Il InvoiceLine) transformChangeDataCapture(data ChangeDataCapture, idCache *DependentDataIDCache) ([]map[string]any, error) {
	items := []map[string]any{}
	for _, cdcResponse := range data.CDCResponse {
		for _, queryResponse := range cdcResponse.CDCQueryResponse {
			if queryResponse.QueryResponse.Items != nil {
				for _, item := range queryResponse.QueryResponse.Items {
					invoice, ok := item.(Invoice)
					if !ok {
						break
					}

					// map lines in cdc response
					newLineIDs := map[string]bool{}
					for _, line := range invoice.Line {
						if line.DetailType == "GroupLineDetail" {
							for _, groupLine := range line.GroupLineDetail.Line {
								newLineIDs[fmt.Sprintf("%s:%s:%s", invoice.Id, line.Id, groupLine.Id)] = true
							}
							newLineIDs[fmt.Sprintf("%s:%s", invoice.Id, line.Id)] = true
						}
						if line.DetailType == "DescriptionOnly" || line.DetailType == "SalesItemLineDetail" {
							newLineIDs[fmt.Sprintf("%s:%s", invoice.Id, line.Id)] = true
						}
					}

					fmt.Printf("newLineIDs: %s\n", FormatJSON(newLineIDs))

					// handle lines on deleted invoices
					if invoice.Status == "Deleted" {
						cachedLines := idCache.IDs[invoice.Id]
						fmt.Printf("cachedLines: %s\n", FormatJSON(cachedLines))
						for lineID := range cachedLines {
							items = append(items, map[string]any{
								"id":           lineID,
								"__syncAction": REMOVE,
							})
						}
						delete(idCache.IDs, invoice.Id)
						if _, ok := idCache.IDs[invoice.Id]; !ok {
							fmt.Printf("cache entry for invoice %s deleted\n", invoice.Id)
						}
						continue
					}

					fmt.Printf("items after delete: %s\n", FormatJSON(items))

					// transform line data on added or updated invoices
					for _, line := range invoice.Line {
						if line.DetailType == "DescriptionOnly" || line.DetailType == "SalesItemLineDetail" {
							item, err := line.transformItem(invoice)
							if err != nil {
								return nil, fmt.Errorf("unable to invoice line transform data: %w", err)
							}
							items = append(items, item)
						}
						if line.DetailType == "GroupLineDetail" {
							for _, groupLine := range line.GroupLineDetail.Line {
								item, err := groupLine.transformItem(invoice)
								if err != nil {
									return nil, fmt.Errorf("unable to invoice line transform data: %w", err)
								}
								item["id"] = fmt.Sprintf("%s:%s:%s", invoice.Id, line.Id, groupLine.Id)
								item["group_line_id"] = line.Id
								items = append(items, item)
							}
							item, err := line.transformItem(invoice)
							if err != nil {
								return nil, fmt.Errorf("unable to invoice line transform data: %w", err)
							}
							items = append(items, item)
						}
					}

					fmt.Printf("items after transform: %s\n", FormatJSON(items))

					// check for lines in cache but not in cdc response
					if _, ok := idCache.IDs[invoice.Id]; ok {
						cachedLines := idCache.IDs[invoice.Id]
						fmt.Printf("cachedLines: %s\n", FormatJSON(cachedLines))
						for cachedLineID := range cachedLines {
							if !newLineIDs[cachedLineID] {
								items = append(items, map[string]any{
									"id":           cachedLineID,
									"__syncAction": REMOVE,
								})
							}
						}
					}

					fmt.Printf("items after remove: %s\n", FormatJSON(items))

					// update cache with new line ids
					idCache.IDs[invoice.Id] = newLineIDs

					fmt.Printf("cache after transform: %s\n", FormatJSON(idCache.IDs[invoice.Id]))
				}
			}
		}
	}
	return items, nil
}

func (Il InvoiceLine) GetData(req *DataRequest) (DataHandlerResponse, error) {
	var cacheMutex sync.Mutex
	cacheKey := fmt.Sprintf("%s:%s", req.RealmID, "invoice")
	cacheMutex.Lock()
	cacheInterface, cacheExists := req.Cache.Get(cacheKey)
	cacheMutex.Unlock()

	if req.LastSynced == "" || !cacheExists {
		groupKey := fmt.Sprintf("%s:%s:%d", req.OperationID, "invoice", req.StartPosition)
		res, err, _ := req.Group.Do(groupKey, func() (interface{}, error) {
			return Il.getFullData(req)
		})

		if err != nil {
			return DataHandlerResponse{}, fmt.Errorf("unable to query invoices from qbo: %w", err)
		}

		IDMap := map[string]map[string]bool{}

		for _, invoice := range res.(DataResponse).Data.([]Invoice) {
			IDMap[invoice.Id] = map[string]bool{}
			for _, line := range invoice.Line {
				if line.DetailType == "GroupLineDetail" {
					for _, groupLine := range line.GroupLineDetail.Line {
						IDMap[invoice.Id][fmt.Sprintf("%s:%s:%s", invoice.Id, line.Id, groupLine.Id)] = true
					}
					IDMap[invoice.Id][fmt.Sprintf("%s:%s", invoice.Id, line.Id)] = true
				}
				if line.DetailType == "DescriptionOnly" || line.DetailType == "SalesItemLineDetail" {
					IDMap[invoice.Id][fmt.Sprintf("%s:%s", invoice.Id, line.Id)] = true
				}
			}
		}

		cacheMutex.Lock()
		if !cacheExists {
			IDEntry := DependentDataIDCache{
				OperationID: req.OperationID,
				IDs:         IDMap,
			}
			err = req.Cache.Add(cacheKey, IDEntry, IDCacheLifetime)
			if err != nil {
				cacheMutex.Unlock()
				return DataHandlerResponse{}, fmt.Errorf("unable to add cache entry: %w", err)
			}
		} else {
			existingIDCache := cacheInterface.(DependentDataIDCache)
			if existingIDCache.OperationID == req.OperationID {
				for invoiceID, linesMap := range IDMap {
					if _, ok := existingIDCache.IDs[invoiceID]; !ok {
						existingIDCache.IDs[invoiceID] = map[string]bool{}
					}
					for lineID := range linesMap {
						existingIDCache.IDs[invoiceID][lineID] = true
					}
				}
				req.Cache.Set(cacheKey, existingIDCache, IDCacheLifetime)
			} else {
				IDEntry := DependentDataIDCache{
					OperationID: req.OperationID,
					IDs:         IDMap,
				}
				req.Cache.Set(cacheKey, IDEntry, IDCacheLifetime)
			}
		}
		cacheMutex.Unlock()

		items, err := Il.transformFullData(res.(DataResponse))
		if err != nil {
			return DataHandlerResponse{}, fmt.Errorf("unable to transform data: %w", err)
		}

		return DataHandlerResponse{
			Items: items,
			Pagination: Pagination{
				HasNext: res.(DataResponse).More,
				NextPageConfig: NextPageConfig{
					StartPosition: req.StartPosition + QueryPageSize,
				},
			},
			SynchronizationType: FullSync,
		}, nil
	} else {
		groupKey := req.OperationID
		existingIDCache := cacheInterface.(DependentDataIDCache)
		res, err, _ := req.Group.Do(groupKey, func() (interface{}, error) {
			return getChangeDataCapture(req)
		})

		if err != nil {
			return DataHandlerResponse{}, fmt.Errorf("unable to cdc query from qbo: %w", err)
		}

		cacheMutex.Lock()
		items, err := Il.transformChangeDataCapture(res.(ChangeDataCapture), &existingIDCache)
		cacheMutex.Unlock()
		if err != nil {
			return DataHandlerResponse{}, fmt.Errorf("unable to transform data: %w", err)
		}

		cacheMutex.Lock()
		IDEntry := DependentDataIDCache{
			OperationID: req.OperationID,
			IDs:         existingIDCache.IDs,
		}
		req.Cache.Set(cacheKey, IDEntry, IDCacheLifetime)
		cacheMutex.Unlock()

		return DataHandlerResponse{
			Items: items,
			Pagination: Pagination{
				HasNext: false,
				NextPageConfig: NextPageConfig{
					StartPosition: 0,
				},
			},
			SynchronizationType: DeltaSync,
		}, nil
	}
}

func init() {
	RegisterType(Invoice{})
	RegisterType(InvoiceLine{})
	TestDependentDataType(InvoiceLine{})
	TestQuickbooksDataType(Invoice{})
	TestParentDataType(Invoice{})
}
