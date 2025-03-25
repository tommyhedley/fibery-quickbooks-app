package data

import (
	"fmt"

	"github.com/tommyhedley/fibery-quickbooks-app/pkgs/fibery"
	"github.com/tommyhedley/quickbooks-go"
)

var BillPayment = QuickBooksDualType[quickbooks.BillPayment]{
	QuickBooksType: QuickBooksType[quickbooks.BillPayment]{
		BaseType: fibery.BaseType{
			TypeId:   "BillPayment",
			TypeName: "Bill Payment",
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
				"DocNumber": {
					Name: "Reference Number",
					Type: fibery.Text,
				},
				"TxnDate": {
					Name:    "Payment Date",
					Type:    fibery.DateType,
					SubType: fibery.Day,
				},
				"PrivateNote": {
					Name:    "Memo",
					Type:    fibery.Text,
					SubType: fibery.MD,
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
				"PayType": {
					Name:     "Payment Type",
					Type:     fibery.Text,
					SubType:  fibery.SingleSelect,
					ReadOnly: true,
					Options: []map[string]any{
						{
							"name": "Check",
						},
						{
							"name": "Credit Card",
						},
					},
				},
				"VendorId": {
					Name: "Vendor ID",
					Type: fibery.Text,
					Relation: &fibery.Relation{
						Cardinality:   fibery.MTO,
						Name:          "Vendor",
						TargetName:    "Bill Payments",
						TargetType:    "Vendor",
						TargetFieldID: "id",
					},
				},
				"PaymentAccountId": {
					Name: "Payment Account ID",
					Type: fibery.Text,
					Relation: &fibery.Relation{
						Cardinality:   fibery.MTO,
						Name:          "Payment Account",
						TargetName:    "Bill Payments",
						TargetType:    "Account",
						TargetFieldID: "id",
					},
				},
			},
		},
		schemaGen: func(bp quickbooks.BillPayment) (map[string]any, error) {
			var paymentAccountId string
			if bp.APAccountRef != nil {
				paymentAccountId = bp.APAccountRef.Value
			}

			var payType string
			switch bp.PayType {
			case quickbooks.CreditCardPaymentType:
				payType = "Credit Card"
				paymentAccountId = bp.CreditCardPayment.CCAccountRef.Value
			case quickbooks.CheckPaymentType:
				payType = "Check"
				paymentAccountId = bp.CheckPayment.BankAccountRef.Value
			}

			return map[string]any{
				"id":               bp.Id,
				"QBOId":            bp.Id,
				"Name":             bp.VendorRef.Name,
				"SyncToken":        bp.SyncToken,
				"__syncAction":     fibery.SET,
				"DocNumber":        bp.DocNumber,
				"TxnDate":          bp.TxnDate.Format(fibery.DateFormat),
				"PrivateNote":      bp.PrivateNote,
				"TotalAmt":         bp.TotalAmt,
				"PayType":          payType,
				"VendorId":         bp.VendorRef.Value,
				"PaymentAccountId": paymentAccountId,
			}, nil
		},
		pageQuery: func(req Request) ([]quickbooks.BillPayment, error) {
			params := quickbooks.RequestParameters{
				Ctx:     req.Ctx,
				RealmId: req.RealmId,
				Token:   req.Token,
			}

			items, err := req.Client.FindBillPaymentsByPage(params, req.StartPosition, req.PageSize)
			if err != nil {
				return nil, err
			}

			return items, nil
		},
	},
	entityId: func(bp quickbooks.BillPayment) string {
		return bp.Id
	},
	entityStatus: func(bp quickbooks.BillPayment) string {
		return bp.Status
	},
}

var BillPaymentLine = DependentDualType[quickbooks.BillPayment]{
	dependentBaseType: dependentBaseType[quickbooks.BillPayment]{
		BaseType: fibery.BaseType{
			TypeId:   "BillPaymentLine",
			TypeName: "Bill Payment Line",
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
				"BillId": {
					Name: "Bill ID",
					Type: fibery.Text,
					Relation: &fibery.Relation{
						Cardinality:   fibery.MTO,
						Name:          "Bill",
						TargetName:    "Bill Payment Lines",
						TargetType:    "Bill",
						TargetFieldID: "id",
					},
				},
				"VendorCreditId": {
					Name: "Vendor Credit ID",
					Type: fibery.Text,
					Relation: &fibery.Relation{
						Cardinality:   fibery.MTO,
						Name:          "Vendor Credit",
						TargetName:    "Bill Payment Lines",
						TargetType:    "VendorCredit",
						TargetFieldID: "id",
					},
				},
				"DepositId": {
					Name: "Deposit ID",
					Type: fibery.Text,
					Relation: &fibery.Relation{
						Cardinality:   fibery.MTO,
						Name:          "Deposit",
						TargetName:    "Bill Payment Lines",
						TargetType:    "Deposit",
						TargetFieldID: "id",
					},
				},
			},
		},
		schemaGen: func(bp quickbooks.BillPayment) ([]map[string]any, error) {
			items := []map[string]any{}
			for _, line := range bp.Line {
				var description string
				var billId string
				var vendorCreditId string
				switch line.LinkedTxn[0].TxnType {
				case "Bill":
					description = "Bill Payment"
					billId = line.LinkedTxn[0].TxnId
				case "VendorCredit":
					description = "Vendor Credit"
					vendorCreditId = line.LinkedTxn[0].TxnId
				}

				item := map[string]any{
					"id":             fmt.Sprintf("%s:%s", bp.Id, line.Id),
					"QBOId":          line.Id,
					"Description":    description,
					"__syncAction":   fibery.SET,
					"Amount":         line.Amount,
					"BillId":         billId,
					"VendorCreditId": vendorCreditId,
				}

				items = append(items, item)
			}
			return items, nil
		},
	},
	sourceType: &BillPayment,
	sourceId: func(bp quickbooks.BillPayment) string {
		return bp.Id
	},
	sourceStatus: func(bp quickbooks.BillPayment) string {
		return bp.Status
	},
	sourceMapper: func(bp quickbooks.BillPayment) map[string]struct{} {
		sourceMap := map[string]struct{}{}
		for _, line := range bp.Line {
			sourceMap[fmt.Sprintf("%s:%s", bp.Id, line.Id)] = struct{}{}
		}
		return sourceMap
	},
}

func init() {
	registerType(&BillPayment)
	registerType(&BillPaymentLine)
}
