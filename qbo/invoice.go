// Copyright (c) 2018, Randy Westlund. All rights reserved.
// This code is under the BSD-2-Clause license.

package qbo

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"

	"github.com/tommyhedley/fibery/fibery-qbo-integration/sync"
)

var InvoiceType = sync.DataType{
	ID:   "invoice",
	Name: "Invoice",
	Schema: map[string]sync.Field{
		"id": {
			Name: "id",
			Type: sync.ID,
		},
		"customer_id": {
			Name: "Customer ID",
			Type: sync.Text,
			Relation: &sync.Relation{
				Cardinality:   sync.MTO,
				Name:          "Customer",
				TargetName:    "Invoices",
				TargetType:    "customer",
				TargetFieldID: "id",
			},
		},
		"sync_token": {
			Name:     "Sync Token",
			Type:     sync.Text,
			ReadOnly: true,
		},
		"shipping_address": {
			Name:    "Shipping Address",
			Type:    sync.Text,
			SubType: sync.MD,
		},
		"billing_address": {
			Name:    "Billing Address",
			Type:    sync.Text,
			SubType: sync.MD,
		},
		"invoice_num": {
			Name: "Invoice",
			Type: sync.Text,
		},
		"email": {
			Name:    "To",
			Type:    sync.Text,
			SubType: sync.Email,
		},
		"email_cc": {
			Name:    "CC",
			Type:    sync.Text,
			SubType: sync.Email,
		},
		"email_bcc": {
			Name:    "BCC",
			Type:    sync.Text,
			SubType: sync.Email,
		},
		"email_status": {
			Name:     "Email Status",
			Type:     sync.Text,
			SubType:  sync.SingleSelect,
			ReadOnly: true,
			Options: []map[string]any{
				{
					"name": "Not Sent",
				},
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
			Type: sync.Date,
		},
		"date": {
			Name: "Date",
			Type: sync.Date,
		},
		"due_date": {
			Name: "Due Date",
			Type: sync.Date,
		},
		"class_id": {
			Name: "Class ID",
			Type: sync.Text,
			Relation: &sync.Relation{
				Cardinality:   sync.MTO,
				Name:          "Class",
				TargetName:    "Invoices",
				TargetType:    "class",
				TargetFieldID: "id",
			},
		},
		"print_status": {
			Name:     "Print Status",
			Type:     sync.Text,
			SubType:  sync.SingleSelect,
			ReadOnly: true,
			Options: []map[string]any{
				{
					"name": "Not Set",
				},
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
			Type: sync.Text,
			Relation: &sync.Relation{
				Cardinality:   sync.MTO,
				Name:          "Term",
				TargetName:    "Invoices",
				TargetType:    "term",
				TargetFieldID: "id",
			},
		},
		"statement_memo": {
			Name:    "Statement Message",
			Type:    sync.Text,
			SubType: sync.Title,
		},
		"customer_memo": {
			Name: "Invoice Message",
			Type: sync.Text,
		},
		"allow_ach": {
			Name:    "ACH Payments",
			Type:    sync.Text,
			SubType: sync.Boolean,
		},
		"allow_cc": {
			Name:    "Credit Card Payments",
			Type:    sync.Text,
			SubType: sync.Boolean,
		},
		"tax_code_id": {
			Name: "Tax Code ID",
			Type: sync.Text,
			Relation: &sync.Relation{
				Cardinality:   sync.MTO,
				Name:          "Tax Code",
				TargetName:    "Invoices",
				TargetType:    "tax_code",
				TargetFieldID: "id",
			},
		},
		"tax_position": {
			Name:     "Apply Tax",
			Type:     sync.Text,
			SubType:  sync.SingleSelect,
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
			Type: sync.Text,
			Relation: &sync.Relation{
				Cardinality:   sync.MTO,
				Name:          "Tax Exemption",
				TargetName:    "Invoices",
				TargetType:    "tax_exemption",
				TargetFieldID: "id",
			},
		},
		"deposit_account_id": {
			Name: "Deposit Account ID",
			Type: sync.Text,
			Relation: &sync.Relation{
				Cardinality:   sync.MTO,
				Name:          "Deposit Account",
				TargetName:    "Invoice Deposits",
				TargetType:    "account",
				TargetFieldID: "id",
			},
		},
		"deposit_field": {
			Name: "Deposit Amount",
			Type: sync.Number,
			Format: map[string]any{
				"format":               "Money",
				"currencyCode":         "USD",
				"hasThousandSeperator": true,
				"precision":            2,
			},
		},
		"discount_type": {
			Name:     "Discount Type",
			Type:     sync.Text,
			SubType:  sync.SingleSelect,
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
			Type: sync.Number,
			Format: map[string]any{
				"format":    "Percent",
				"precision": 2,
			},
		},
		"discount_amount": {
			Name: "Discount Amount",
			Type: sync.Number,
			Format: map[string]any{
				"format":               "Money",
				"currencyCode":         "USD",
				"hasThousandSeperator": true,
				"precision":            2,
			},
		},
		"tax": {
			Name: "Tax",
			Type: sync.Number,
			Format: map[string]any{
				"format":               "Money",
				"currencyCode":         "USD",
				"hasThousandSeperator": true,
				"precision":            2,
			},
		},
		"total": {
			Name: "Total",
			Type: sync.Number,
			Format: map[string]any{
				"format":               "Money",
				"currencyCode":         "USD",
				"hasThousandSeperator": true,
				"precision":            2,
			},
		},
		"balance": {
			Name: "Balance",
			Type: sync.Number,
			Format: map[string]any{
				"format":               "Money",
				"currencyCode":         "USD",
				"hasThousandSeperator": true,
				"precision":            2,
			},
		},
		"link": {
			Name:    "URL",
			Type:    sync.Text,
			SubType: sync.URL,
		},
		"created_qbo": {
			Name: "Creation Date (QBO)",
			Type: sync.Date,
		},
		"last_updated_qbo": {
			Name: "Last Updated (QBO)",
			Type: sync.Date,
		},
	},
	DataRequest: getInvoiceSubtype("invoice"),
}

var InvoiceLineType = sync.DataType{
	ID:   "invoice_line",
	Name: "Invoice Line",
	Schema: map[string]sync.Field{
		"id": {
			Name: "id",
			Type: sync.ID,
		},
		"description": {
			Name:    "Description",
			Type:    sync.Text,
			SubType: sync.Title,
		},
		"type": {
			Name:    "Type",
			Type:    sync.Text,
			SubType: sync.SingleSelect,
			Options: []map[string]any{
				{
					"name": "Sales Item Line",
				},
				{
					"name": "Group Line",
				},
				{
					"name": "Description Line",
				},
			},
		},
		"quantity": {
			Name: "Quantity",
			Type: sync.Number,
			Format: map[string]any{
				"format":               "Number",
				"hasThousandSeparator": true,
				"precision":            2,
			},
		},
		"unit_price": {
			Name: "Unit Price",
			Type: sync.Number,
			Format: map[string]any{
				"format":               "Money",
				"currencyCode":         "USD",
				"hasThousandSeperator": true,
				"precision":            2,
			},
		},
		"amount": {
			Name: "Amount",
			Type: sync.Number,
			Format: map[string]any{
				"format":               "Money",
				"currencyCode":         "USD",
				"hasThousandSeperator": true,
				"precision":            2,
			},
		},
		"line_num": {
			Name:    "Line",
			Type:    sync.Number,
			SubType: sync.Integer,
		},
		"group_line_id": {
			Name: "Group Line ID",
			Type: sync.Text,
			Relation: &sync.Relation{
				Cardinality:   sync.MTO,
				Name:          "Group",
				TargetName:    "Lines",
				TargetType:    "invoice_line",
				TargetFieldID: "id",
			},
		},
		"item_id": {
			Name: "Item",
			Type: sync.Text,
			Relation: &sync.Relation{
				Cardinality:   sync.MTO,
				Name:          "Item",
				TargetName:    "Invoice Lines",
				TargetType:    "item",
				TargetFieldID: "id",
			},
		},
		"class_id": {
			Name: "Class ID",
			Type: sync.Text,
			Relation: &sync.Relation{
				Cardinality:   sync.MTO,
				Name:          "Class",
				TargetName:    "Expense Account Line(s)",
				TargetType:    "class",
				TargetFieldID: "id",
			},
		},
		"tax_code_id": {
			Name: "Tax Code ID",
			Type: sync.Text,
			Relation: &sync.Relation{
				Cardinality:   sync.MTO,
				Name:          "Tax Code",
				TargetName:    "Invoice Lines",
				TargetType:    "tax_code",
				TargetFieldID: "id",
			},
		},
		"markup_percent": {
			Name: "Markup",
			Type: sync.Number,
			Format: map[string]any{
				"format":    "Percent",
				"precision": 2,
			},
		},
		"service_date": {
			Name:    "Date",
			Type:    sync.Date,
			SubType: sync.Day,
		},
	},
	DataRequest: getInvoiceSubtype("invoice_line"),
}

func getInvoiceSubtype(subTypeID string) sync.DataRequest {
	return func(req sync.RequestParameters) ([]map[string]any, bool, error) {
		return getInvoices(subTypeID, req)
	}
}

func getInvoices(subTypeID string, req sync.RequestParameters) ([]map[string]any, bool, error) {
	if subTypeID == "invoice" {

		return nil, false, nil
	}
	if subTypeID == "invoice_line" {
		return nil, false, nil
	}
	return nil, false, fmt.Errorf("invalid subtype: %s", subTypeID)
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
	Balance                      json.Number   `json:",omitempty"`
	TxnSource                    string        `json:",omitempty"`
	AllowOnlineCreditCardPayment bool          `json:",omitempty"`
	AllowOnlineACHPayment        bool          `json:",omitempty"`
	Deposit                      json.Number   `json:",omitempty"`
	DepositToAccountRef          ReferenceType `json:",omitempty"`
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
	ItemRef   ReferenceType `json:",omitempty"`
	ClassRef  ReferenceType `json:",omitempty"`
	UnitPrice json.Number   `json:",omitempty"`
	// MarkupInfo
	Qty             float32       `json:",omitempty"`
	ItemAccountRef  ReferenceType `json:",omitempty"`
	TaxCodeRef      ReferenceType `json:",omitempty"`
	ServiceDate     Date          `json:",omitempty"`
	TaxInclusiveAmt json.Number   `json:",omitempty"`
	DiscountRate    json.Number   `json:",omitempty"`
	DiscountAmt     json.Number   `json:",omitempty"`
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

	for i := 0; i < resp.QueryResponse.TotalCount; i += queryPageSize {
		query := "SELECT * FROM Invoice ORDERBY Id STARTPOSITION " + strconv.Itoa(i+1) + " MAXRESULTS " + strconv.Itoa(queryPageSize)

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
