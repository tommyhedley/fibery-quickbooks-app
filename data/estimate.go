package data

import (
	"encoding/json"
	"fmt"

	"github.com/tommyhedley/fibery-quickbooks-app/pkgs/fibery"
	"github.com/tommyhedley/quickbooks-go"
)

var Estimate = QuickBooksDualType[quickbooks.Estimate]{
	QuickBooksType: QuickBooksType[quickbooks.Estimate]{
		BaseType: fibery.BaseType{
			TypeId:   "Estimate",
			TypeName: "Estimate",
			TypeSchema: map[string]fibery.Field{
				"id": {
					Name: "ID",
					Type: fibery.Id,
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
						TargetFieldID: "id",
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
						TargetFieldID: "id",
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
						TargetFieldID: "id",
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
						TargetFieldID: "id",
					},
				},
			},
		},
		schemaGen: func(e quickbooks.Estimate) (map[string]any, error) {
			var discountType = map[bool]string{
				true:  "Percentage",
				false: "Amount",
			}
			var discountLine *quickbooks.Line
			var subtotalLine *quickbooks.Line
			for _, line := range e.Line {
				if line.DetailType == "DiscountLineDetail" {
					if discountLine != nil {
						return nil, fmt.Errorf("estimate %s has more than one discount line", e.Id)
					}
					discountLine = &line
				}
				if line.DetailType == "SubTotalLineDetail" {
					if subtotalLine != nil {
						return nil, fmt.Errorf("estimate %s has more than one subtotal line", e.Id)
					}
					subtotalLine = &line
				}
			}

			if subtotalLine == nil {
				return nil, fmt.Errorf("estimate %s has no subtotal lines", e.Id)
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
			if e.DeliveryInfo != nil && !e.DeliveryInfo.DeliveryTime.IsZero() {
				emailSendTime = e.DeliveryInfo.DeliveryTime.Format(fibery.DateFormat)
			}

			var name string
			if e.CustomerRef.Name == "" {
				name = e.DocNumber
			} else {
				name = e.DocNumber + " - " + e.CustomerRef.Name
			}

			var shipAddr quickbooks.PhysicalAddress
			if e.ShipAddr != nil {
				shipAddr = *e.ShipAddr
			}

			var billAddr quickbooks.PhysicalAddress
			if e.BillAddr != nil {
				billAddr = *e.BillAddr
			}

			var acceptedDate string
			if e.AcceptedDate != nil {
				acceptedDate = e.AcceptedDate.Format(fibery.DateFormat)
			}

			var expirationDate string
			if e.ExpirationDate != nil {
				expirationDate = e.ExpirationDate.Format(fibery.DateFormat)
			}

			var billEmailCC string
			if e.BillEmailCC != nil {
				billEmailCC = e.BillEmailCC.Address
			}

			var billEmailBCC string
			if e.BillEmailBCC != nil {
				billEmailBCC = e.BillEmailCC.Address
			}

			var txnDate string
			if e.TxnDate != nil {
				txnDate = e.TxnDate.Format(fibery.DateFormat)
			}

			var totalTax json.Number
			var taxCodeId string
			if e.TxnTaxDetail != nil {
				totalTax = e.TxnTaxDetail.TotalTax
				taxCodeId = e.TxnTaxDetail.TxnTaxCodeRef.Value
			}

			var classId string
			if e.ClassRef != nil {
				classId = e.ClassRef.Value
			}

			var taxExemptionId string
			if e.TaxExemptionRef != nil {
				taxExemptionId = e.TaxExemptionRef.Value
			}

			return map[string]any{
				"id":                     e.Id,
				"QBOId":                  e.Id,
				"Name":                   name,
				"SyncToken":              e.SyncToken,
				"__syncAction":           fibery.SET,
				"ShippingLine1":          shipAddr.Line1,
				"ShippingLine2":          shipAddr.Line2,
				"ShippingLine3":          shipAddr.Line3,
				"ShippingLine4":          shipAddr.Line4,
				"ShippingLine5":          shipAddr.Line5,
				"ShippingCity":           shipAddr.City,
				"ShippingState":          shipAddr.CountrySubDivisionCode,
				"ShippingPostalCode":     shipAddr.PostalCode,
				"ShippingCountry":        shipAddr.Country,
				"ShippingLat":            shipAddr.Lat,
				"ShippingLong":           shipAddr.Long,
				"ShippingFromLine1":      e.ShipFromAddr.Line1,
				"ShippingFromLine2":      e.ShipFromAddr.Line2,
				"ShippingFromLine3":      e.ShipFromAddr.Line3,
				"ShippingFromLine4":      e.ShipFromAddr.Line4,
				"ShippingFromLine5":      e.ShipFromAddr.Line5,
				"ShippingFromCity":       e.ShipFromAddr.City,
				"ShippingFromState":      e.ShipFromAddr.CountrySubDivisionCode,
				"ShippingFromPostalCode": e.ShipFromAddr.PostalCode,
				"ShippingFromCountry":    e.ShipFromAddr.Country,
				"ShippingFromLat":        e.ShipFromAddr.Lat,
				"ShippingFromLong":       e.ShipFromAddr.Long,
				"BillingLine1":           billAddr.Line1,
				"BillingLine2":           billAddr.Line2,
				"BillingLine3":           billAddr.Line3,
				"BillingLine4":           billAddr.Line4,
				"BillingLine5":           billAddr.Line5,
				"BillingCity":            billAddr.City,
				"BillingState":           billAddr.CountrySubDivisionCode,
				"BillingPostalCode":      billAddr.PostalCode,
				"BillingCountry":         billAddr.Country,
				"BillingLat":             billAddr.Lat,
				"BillingLong":            billAddr.Long,
				"DocNumber":              e.DocNumber,
				"TxnStatus":              e.TxnStatus,
				"AcceptedBy":             e.AcceptedBy,
				"AcceptedDate":           acceptedDate,
				"ExpirationDate":         expirationDate,
				"Email":                  e.BillEmail.Address,
				"EmailCC":                billEmailCC,
				"EmailBCC":               billEmailBCC,
				"EmailSendLater":         e.EmailStatus == "NeedToSend",
				"EmailSent":              e.EmailStatus == "EmailSent",
				"EmailSendTime":          emailSendTime,
				"TxnDate":                txnDate,
				"PrivateNote":            e.PrivateNote,
				"CustomerMemo":           e.CustomerMemo.Value,
				"DiscountPosition":       e.ApplyTaxAfterDiscount,
				"DiscountType":           discountTypeValue,
				"DiscountPercent":        discountPercent,
				"DiscountAmount":         discountAmount,
				"Tax":                    totalTax,
				"SubtotalAmt":            subtotalAmount,
				"TotalAmt":               e.TotalAmt,
				"ClassId":                classId,
				"TaxCodeId":              taxCodeId,
				"TaxExemptionId":         taxExemptionId,
				"CustomerId":             e.CustomerRef.Value,
			}, nil
		},
		pageQuery: func(req Request) ([]quickbooks.Estimate, error) {
			params := quickbooks.RequestParameters{
				Ctx:     req.Ctx,
				RealmId: req.RealmId,
				Token:   req.Token,
			}

			items, err := req.Client.FindEstimatesByPage(params, req.StartPosition, req.PageSize)
			if err != nil {
				return nil, err
			}

			return items, nil
		},
	},
	entityId: func(e quickbooks.Estimate) string {
		return e.Id
	},
	entityStatus: func(e quickbooks.Estimate) string {
		return e.Status
	},
}

var EstimateLine = DependentDualType[quickbooks.Estimate]{
	dependentBaseType: dependentBaseType[quickbooks.Estimate]{
		BaseType: fibery.BaseType{
			TypeId:   "EstimateLine",
			TypeName: "Estimate Line",
			TypeSchema: map[string]fibery.Field{
				"id": {
					Name: "ID",
					Type: fibery.Id,
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
						TargetFieldID: "id",
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
						TargetFieldID: "id",
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
						TargetFieldID: "id",
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
						TargetFieldID: "id",
					},
				},
			},
		},
		schemaGen: func(e quickbooks.Estimate) ([]map[string]any, error) {
			items := []map[string]any{}
			for _, line := range e.Line {
				if line.DetailType == quickbooks.DescriptionLine || line.DetailType == quickbooks.SalesItemLine {
					item := map[string]any{
						"id":           fmt.Sprintf("%s:%s", e.Id, line.Id),
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
						"EstimateId":   e.Id,
					}
					items = append(items, item)
				}
				if line.DetailType == quickbooks.GroupLine {
					for _, groupLine := range line.GroupLineDetail.Line {
						item := map[string]any{
							"id":           fmt.Sprintf("%s:%s:%s", e.Id, line.Id, groupLine.Id),
							"GroupLineId":  line.Id,
							"QBOId":        line.Id,
							"Description":  line.Description,
							"__syncAction": fibery.SET,
							"Qty":          line.GroupLineDetail.Quantity,
							"LineNum":      line.LineNum,
							"ItemId":       line.GroupLineDetail.GroupItemRef.Value,
							"EstimateId":   e.Id,
						}
						items = append(items, item)
					}
					item := map[string]any{
						"id":           fmt.Sprintf("%s:%s", e.Id, line.Id),
						"QBOId":        line.Id,
						"Description":  line.Description,
						"__syncAction": fibery.SET,
						"Qty":          line.GroupLineDetail.Quantity,
						"LineNum":      line.LineNum,
						"ItemId":       line.GroupLineDetail.GroupItemRef.Value,
						"EstimateId":   e.Id,
					}
					items = append(items, item)
				}
			}
			return items, nil
		},
	},
	sourceType: &Estimate,
	sourceId: func(e quickbooks.Estimate) string {
		return e.Id
	},
	sourceStatus: func(e quickbooks.Estimate) string {
		return e.Status
	},
	sourceMapper: func(e quickbooks.Estimate) map[string]struct{} {
		sourceMap := map[string]struct{}{}
		for _, line := range e.Line {
			if line.DetailType == quickbooks.GroupLine {
				for _, groupLine := range line.GroupLineDetail.Line {
					sourceMap[fmt.Sprintf("%s:%s:%s", e.Id, line.Id, groupLine.Id)] = struct{}{}
				}
				sourceMap[fmt.Sprintf("%s:%s", e.Id, line.Id)] = struct{}{}
			}
			if line.DetailType == quickbooks.DescriptionLine || line.DetailType == quickbooks.SalesItemLine {
				sourceMap[fmt.Sprintf("%s:%s", e.Id, line.Id)] = struct{}{}
			}
		}
		return sourceMap
	},
}

func init() {
	registerType(&Estimate)
	registerType(&EstimateLine)
}
