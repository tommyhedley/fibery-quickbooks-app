package main

import (
	"context"
	"fmt"

	"github.com/tommyhedley/fibery-quickbooks-app/pkgs/fibery"
	"github.com/tommyhedley/quickbooks-go"
)

type TypeRegistry map[string]Type

func NewTypeRegistry() TypeRegistry {
	tr := TypeRegistry{}
	BuildTypes(tr)
	return tr
}

func (tr TypeRegistry) Register(t Type) {
	tr[t.Id()] = t
}

func (tr TypeRegistry) GetType(id string) (Type, bool) {
	if typ, exists := tr[id]; exists {
		return typ, true
	}
	return nil, false
}

func BuildTypes(tr TypeRegistry) {
	account := NewQuickBooksDualType(
		"Account",
		"Account",
		map[string]fibery.Field{
			"id": {
				Name: "ID",
				Type: fibery.Id,
			},
			"QBOId": {
				Name: "QBO ID",
				Type: fibery.Text,
			},
			"Name": {
				Name: "Base Name",
				Type: fibery.Text,
			},
			"FullyQualifiedName": {
				Name:    "Full Name",
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
			"Active": {
				Name:    "Active",
				Type:    fibery.Text,
				SubType: fibery.Boolean,
			},
			"Description": {
				Name:    "Description",
				Type:    fibery.Text,
				SubType: fibery.MD,
			},
			"AcctNum": {
				Name: "Account Number",
				Type: fibery.Text,
			},
			"CurrentBalance": {
				Name: "Balance",
				Type: fibery.Number,
				Format: map[string]any{
					"format":               "Money",
					"currencyCode":         "USD",
					"hasThousandSeperator": true,
					"precision":            2,
				},
			},
			"CurrentBalanceWithSubAccounts": {
				Name: "Balance With Sub-Accounts",
				Type: fibery.Number,
				Format: map[string]any{
					"format":               "Money",
					"currencyCode":         "USD",
					"hasThousandSeperator": true,
					"precision":            2,
				},
			},
			"Classification": {
				Name:     "Classification",
				Type:     fibery.Text,
				SubType:  fibery.SingleSelect,
				ReadOnly: true,
				Options: []map[string]any{
					{
						"name": "Asset",
					},
					{
						"name": "Equity",
					},
					{
						"name": "Expense",
					},
					{
						"name": "Liability",
					},
					{
						"name": "Revenue",
					},
				},
			},
			"AccountType": {
				Name:     "Account Type",
				Type:     fibery.Text,
				SubType:  fibery.SingleSelect,
				ReadOnly: true,
				Options: []map[string]any{
					{
						"name": "Bank",
					},
					{
						"name": "Other Current Asset",
					},
					{
						"name": "Fixed Asset",
					},
					{
						"name": "Other Asset",
					},
					{
						"name": "Accounts Receivable",
					},
					{
						"name": "Equity",
					},
					{
						"name": "Expense",
					},
					{
						"name": "Other Expense",
					},
					{
						"name": "Cost of Goods Sold",
					},
					{
						"name": "Accounts Payable",
					},
					{
						"name": "Credit Card",
					},
					{
						"name": "Long Term Liability",
					},
					{
						"name": "Other Current Liability",
					},
					{
						"name": "Income",
					},
					{
						"name": "Other Income",
					},
				},
			},
			"AccountSubType": {
				Name:     "Account Sub-Type",
				Type:     fibery.Text,
				SubType:  fibery.SingleSelect,
				ReadOnly: true,
				Options: []map[string]any{
					{
						"name": "CashOnHand",
					},
					{
						"name": "Checking",
					},
					{
						"name": "MoneyMarket",
					},
					{
						"name": "RentsHeldInTrust",
					},
					{
						"name": "Savings",
					},
					{
						"name": "TrustAccounts",
					},
					{
						"name": "CashAndCashEquivalents",
					},
					{
						"name": "OtherEarMarkedBankAccounts",
					},
					{
						"name": "AllowanceForBadDebts",
					},
					{
						"name": "DevelopmentCosts",
					},
					{
						"name": "EmployeeCashAdvances",
					},
					{
						"name": "OtherCurrentAssets",
					},
					{
						"name": "Inventory",
					},
					{
						"name": "Investment_MortgageRealEstateLoans",
					},
					{
						"name": "Investment_Other",
					},
					{
						"name": "Investment_TaxExemptSecurities",
					},
					{
						"name": "Investment_USGovernmentObligations",
					},
					{
						"name": "LoansToOfficers",
					},
					{
						"name": "LoansToOthers",
					},
					{
						"name": "LoansToStockholders",
					},
					{
						"name": "PrepaidExpenses",
					},
					{
						"name": "Retainage",
					},
					{
						"name": "UndepositedFunds",
					},
					{
						"name": "AssetsAvailableForSale",
					},
					{
						"name": "BalWithGovtAuthorities",
					},
					{
						"name": "CalledUpShareCapitalNotPaid",
					},
					{
						"name": "ExpenditureAuthorisationsAndLettersOfCredit",
					},
					{
						"name": "GlobalTaxDeferred",
					},
					{
						"name": "GlobalTaxRefund",
					},
					{
						"name": "InternalTransfers",
					},
					{
						"name": "OtherConsumables",
					},
					{
						"name": "ProvisionsCurrentAssets",
					},
					{
						"name": "ShortTermInvestmentsInRelatedParties",
					},
					{
						"name": "ShortTermLoansAndAdvancesToRelatedParties",
					},
					{
						"name": "TradeAndOtherReceivables",
					},
					{
						"name": "AccumulatedDepletion",
					},
					{
						"name": "AccumulatedDepreciation",
					},
					{
						"name": "DepletableAssets",
					},
					{
						"name": "FixedAssetComputers",
					},
					{
						"name": "FixedAssetCopiers",
					},
					{
						"name": "FixedAssetFurniture",
					},
					{
						"name": "FixedAssetPhone",
					},
					{
						"name": "FixedAssetPhotoVideo",
					},
					{
						"name": "FixedAssetSoftware",
					},
					{
						"name": "FixedAssetOtherToolsEquipment",
					},
					{
						"name": "FurnitureAndFixtures",
					},
					{
						"name": "Land",
					},
					{
						"name": "LeaseholdImprovements",
					},
					{
						"name": "OtherFixedAssets",
					},
					{
						"name": "AccumulatedAmortization",
					},
					{
						"name": "Buildings",
					},
					{
						"name": "IntangibleAssets",
					},
					{
						"name": "MachineryAndEquipment",
					},
					{
						"name": "Vehicles",
					},
					{
						"name": "AssetsInCourseOfConstruction",
					},
					{
						"name": "CapitalWip",
					},
					{
						"name": "CumulativeDepreciationOnIntangibleAssets",
					},
					{
						"name": "IntangibleAssetsUnderDevelopment",
					},
					{
						"name": "LandAsset",
					},
					{
						"name": "NonCurrentAssets",
					},
					{
						"name": "ParticipatingInterests",
					},
					{
						"name": "ProvisionsFixedAssets",
					},
					{
						"name": "LeaseBuyout",
					},
					{
						"name": "OtherLongTermAssets",
					},
					{
						"name": "SecurityDeposits",
					},
					{
						"name": "AccumulatedAmortizationOfOtherAssets",
					},
					{
						"name": "Goodwill",
					},
					{
						"name": "Licenses",
					},
					{
						"name": "OrganizationalCosts",
					},
					{
						"name": "AssetsHeldForSale",
					},
					{
						"name": "AvailableForSaleFinancialAssets",
					},
					{
						"name": "DeferredTax",
					},
					{
						"name": "Investments",
					},
					{
						"name": "LongTermInvestments",
					},
					{
						"name": "LongTermLoansAndAdvancesToRelatedParties",
					},
					{
						"name": "OtherIntangibleAssets",
					},
					{
						"name": "OtherLongTermInvestments",
					},
					{
						"name": "OtherLongTermLoansAndAdvances",
					},
					{
						"name": "PrepaymentsAndAccruedIncome",
					},
					{
						"name": "ProvisionsNonCurrentAssets",
					},
					{
						"name": "OpeningBalanceEquity",
					},
					{
						"name": "PartnersEquity",
					},
					{
						"name": "RetainedEarnings",
					},
					{
						"name": "AccumulatedAdjustment",
					},
					{
						"name": "OwnersEquity",
					},
					{
						"name": "PaidInCapitalOrSurplus",
					},
					{
						"name": "PartnerContributions",
					},
					{
						"name": "PartnerDistributions",
					},
					{
						"name": "PreferredStock",
					},
					{
						"name": "CommonStock",
					},
					{
						"name": "TreasuryStock",
					},
					{
						"name": "EstimatedTaxes",
					},
					{
						"name": "Healthcare",
					},
					{
						"name": "PersonalIncome",
					},
					{
						"name": "PersonalExpense",
					},
					{
						"name": "AccumulatedOtherComprehensiveIncome",
					},
					{
						"name": "CalledUpShareCapital",
					},
					{
						"name": "CapitalReserves",
					},
					{
						"name": "DividendDisbursed",
					},
					{
						"name": "EquityInEarningsOfSubsiduaries",
					},
					{
						"name": "InvestmentGrants",
					},
					{
						"name": "MoneyReceivedAgainstShareWarrants",
					},
					{
						"name": "OtherFreeReserves",
					},
					{
						"name": "ShareApplicationMoneyPendingAllotment",
					},
					{
						"name": "ShareCapital",
					},
					{
						"name": "Funds",
					},
					{
						"name": "AdvertisingPromotional",
					},
					{
						"name": "BadDebts",
					},
					{
						"name": "BankCharges",
					},
					{
						"name": "CharitableContributions",
					},
					{
						"name": "CommissionsAndFees",
					},
					{
						"name": "Entertainment",
					},
					{
						"name": "EntertainmentMeals",
					},
					{
						"name": "EquipmentRental",
					},
					{
						"name": "FinanceCosts",
					},
					{
						"name": "GlobalTaxExpense",
					},
					{
						"name": "Insurance",
					},
					{
						"name": "InterestPaid",
					},
					{
						"name": "LegalProfessionalFees",
					},
					{
						"name": "OfficeExpenses",
					},
					{
						"name": "OfficeGeneralAdministrativeExpenses",
					},
					{
						"name": "OtherBusinessExpenses",
					},
					{
						"name": "OtherMiscellaneousServiceCost",
					},
					{
						"name": "PromotionalMeals",
					},
					{
						"name": "RentOrLeaseOfBuildings",
					},
					{
						"name": "RepairMaintenance",
					},
					{
						"name": "ShippingFreightDelivery",
					},
					{
						"name": "SuppliesMaterials",
					},
					{
						"name": "Travel",
					},
					{
						"name": "TravelMeals",
					},
					{
						"name": "Utilities",
					},
					{
						"name": "Auto",
					},
					{
						"name": "CostOfLabor",
					},
					{
						"name": "DuesSubscriptions",
					},
					{
						"name": "PayrollExpenses",
					},
					{
						"name": "TaxesPaid",
					},
					{
						"name": "UnappliedCashBillPaymentExpense",
					},
					{
						"name": "Utilities",
					},
					{
						"name": "AmortizationExpense",
					},
					{
						"name": "AppropriationsToDepreciation",
					},
					{
						"name": "BorrowingCost",
					},
					{
						"name": "CommissionsAndFees",
					},
					{
						"name": "DistributionCosts",
					},
					{
						"name": "ExternalServices",
					},
					{
						"name": "ExtraordinaryCharges",
					},
					{
						"name": "IncomeTaxExpense",
					},
					{
						"name": "LossOnDiscontinuedOperationsNetOfTax",
					},
					{
						"name": "ManagementCompensation",
					},
					{
						"name": "OtherCurrentOperatingCharges",
					},
					{
						"name": "OtherExternalServices",
					},
					{
						"name": "OtherRentalCosts",
					},
					{
						"name": "OtherSellingExpenses",
					},
					{
						"name": "ProjectStudiesSurveysAssessments",
					},
					{
						"name": "PurchasesRebates",
					},
					{
						"name": "ShippingAndDeliveryExpense",
					},
					{
						"name": "StaffCosts",
					},
					{
						"name": "Sundry",
					},
					{
						"name": "TravelExpensesGeneralAndAdminExpenses",
					},
					{
						"name": "TravelExpensesSellingExpense",
					},
					{
						"name": "Depreciation",
					},
					{
						"name": "ExchangeGainOrLoss",
					},
					{
						"name": "OtherMiscellaneousExpense",
					},
					{
						"name": "PenaltiesSettlements",
					},
					{
						"name": "Amortization",
					},
					{
						"name": "GasAndFuel",
					},
					{
						"name": "HomeOffice",
					},
					{
						"name": "HomeOwnerRentalInsurance",
					},
					{
						"name": "OtherHomeOfficeExpenses",
					},
					{
						"name": "MortgageInterest",
					},
					{
						"name": "RentAndLease",
					},
					{
						"name": "RepairsAndMaintenance",
					},
					{
						"name": "ParkingAndTolls",
					},
					{
						"name": "Vehicle",
					},
					{
						"name": "VehicleInsurance",
					},
					{
						"name": "VehicleLease",
					},
					{
						"name": "VehicleLoanInterest",
					},
					{
						"name": "VehicleLoan",
					},
					{
						"name": "VehicleRegistration",
					},
					{
						"name": "VehicleRepairs",
					},
					{
						"name": "OtherVehicleExpenses",
					},
					{
						"name": "Utilities",
					},
					{
						"name": "WashAndRoadServices",
					},
					{
						"name": "DeferredTaxExpense",
					},
					{
						"name": "Depletion",
					},
					{
						"name": "ExceptionalItems",
					},
					{
						"name": "ExtraordinaryItems",
					},
					{
						"name": "IncomeTaxOtherExpense",
					},
					{
						"name": "MatCredit",
					},
					{
						"name": "PriorPeriodItems",
					},
					{
						"name": "TaxRoundoffGainOrLoss",
					},
					{
						"name": "EquipmentRentalCos",
					},
					{
						"name": "OtherCostsOfServiceCos",
					},
					{
						"name": "ShippingFreightDeliveryCos",
					},
					{
						"name": "SuppliesMaterialsCogs",
					},
					{
						"name": "CostOfLaborCos",
					},
					{
						"name": "CostOfSales",
					},
					{
						"name": "FreightAndDeliveryCost",
					},
					{
						"name": "Accounts Payable",
					},
					{
						"name": "OutstandingDuesMicroSmallEnterprise",
					},
					{
						"name": "OutstandingDuesOtherThanMicroSmallEnterprise",
					},
					{
						"name": "Credit Card",
					},
					{
						"name": "Long Term Liability",
					},
					{
						"name": "NotesPayable",
					},
					{
						"name": "OtherLongTermLiabilities",
					},
					{
						"name": "ShareholderNotesPayable",
					},
					{
						"name": "AccrualsAndDeferredIncome",
					},
					{
						"name": "AccruedLongLermLiabilities",
					},
					{
						"name": "AccruedVacationPayable",
					},
					{
						"name": "BankLoans",
					},
					{
						"name": "DebtsRelatedToParticipatingInterests",
					},
					{
						"name": "DeferredTaxLiabilities",
					},
					{
						"name": "GovernmentAndOtherPublicAuthorities",
					},
					{
						"name": "GroupAndAssociates",
					},
					{
						"name": "LiabilitiesRelatedToAssetsHeldForSale",
					},
					{
						"name": "LongTermBorrowings",
					},
					{
						"name": "LongTermDebit",
					},
					{
						"name": "LongTermEmployeeBenefitObligations",
					},
					{
						"name": "ObligationsUnderFinanceLeases",
					},
					{
						"name": "OtherLongTermProvisions",
					},
					{
						"name": "ProvisionForLiabilities",
					},
					{
						"name": "ProvisionsNonCurrentLiabilities",
					},
					{
						"name": "StaffAndRelatedLongTermLiabilityAccounts",
					},
					{
						"name": "DirectDepositPayable",
					},
					{
						"name": "LineOfCredit",
					},
					{
						"name": "LoanPayable",
					},
					{
						"name": "GlobalTaxPayable",
					},
					{
						"name": "GlobalTaxSuspense",
					},
					{
						"name": "OtherCurrentLiabilities",
					},
					{
						"name": "PayrollClearing",
					},
					{
						"name": "PayrollTaxPayable",
					},
					{
						"name": "PrepaidExpensesPayable",
					},
					{
						"name": "RentsInTrustLiability",
					},
					{
						"name": "TrustAccountsLiabilities",
					},
					{
						"name": "FederalIncomeTaxPayable",
					},
					{
						"name": "InsurancePayable",
					},
					{
						"name": "SalesTaxPayable",
					},
					{
						"name": "StateLocalIncomeTaxPayable",
					},
					{
						"name": "AccruedLiabilities",
					},
					{
						"name": "CurrentLiabilities",
					},
					{
						"name": "CurrentPortionEmployeeBenefitsObligations",
					},
					{
						"name": "CurrentPortionOfObligationsUnderFinanceLeases",
					},
					{
						"name": "CurrentTaxLiability",
					},
					{
						"name": "DividendsPayable",
					},
					{
						"name": "DutiesAndTaxes",
					},
					{
						"name": "InterestPayables",
					},
					{
						"name": "ProvisionForWarrantyObligations",
					},
					{
						"name": "ProvisionsCurrentLiabilities",
					},
					{
						"name": "ShortTermBorrowings",
					},
					{
						"name": "SocialSecurityAgencies",
					},
					{
						"name": "StaffAndRelatedLiabilityAccounts",
					},
					{
						"name": "SundryDebtorsAndCreditors",
					},
					{
						"name": "TradeAndOtherPayables",
					},
					{
						"name": "NonProfitIncome",
					},
					{
						"name": "OtherPrimaryIncome",
					},
					{
						"name": "SalesOfProductIncome",
					},
					{
						"name": "ServiceFeeIncome",
					},
					{
						"name": "DiscountsRefundsGiven",
					},
					{
						"name": "UnappliedCashPaymentIncome",
					},
					{
						"name": "CashReceiptIncome",
					},
					{
						"name": "OperatingGrants",
					},
					{
						"name": "OtherCurrentOperatingIncome",
					},
					{
						"name": "OwnWorkCapitalized",
					},
					{
						"name": "RevenueGeneral",
					},
					{
						"name": "SalesRetail",
					},
					{
						"name": "SalesWholesale",
					},
					{
						"name": "SavingsByTaxScheme",
					},
					{
						"name": "DividendIncome",
					},
					{
						"name": "InterestEarned",
					},
					{
						"name": "OtherInvestmentIncome",
					},
					{
						"name": "OtherMiscellaneousIncome",
					},
					{
						"name": "TaxExemptInterest",
					},
					{
						"name": "GainLossOnSaleOfFixedAssets",
					},
					{
						"name": "GainLossOnSaleOfInvestments",
					},
					{
						"name": "LossOnDisposalOfAssets",
					},
					{
						"name": "OtherOperatingIncome",
					},
					{
						"name": "UnrealisedLossOnSecuritiesNetOfTax",
					},
				},
			},
			"ParentAccountId": {
				Name: "Parent Account ID",
				Type: fibery.Text,
				Relation: &fibery.Relation{
					Cardinality:   fibery.MTO,
					Name:          "Parent Account",
					TargetName:    "Sub-Accounts",
					TargetType:    "Account",
					TargetFieldID: "id",
				},
			},
		},
		func(a quickbooks.Account) (map[string]any, error) {
			var parentAccountId string
			if a.ParentRef != nil {
				parentAccountId = a.ParentRef.Value
			}

			return map[string]any{
				"id":                            a.Id,
				"QBOId":                         a.Id,
				"Name":                          a.Name,
				"FullyQualifiedName":            a.FullyQualifiedName,
				"SyncToken":                     a.SyncToken,
				"__syncAction":                  fibery.SET,
				"Active":                        a.Active,
				"Description":                   a.Description,
				"AcctNum":                       a.AcctNum,
				"CurrentBalance":                a.CurrentBalance,
				"CurrentBalanceWithSubAccounts": a.CurrentBalanceWithSubAccounts,
				"Classification":                a.Classification,
				"AccountType":                   a.AccountType,
				"AccountSubType":                a.AccountSubType,
				"ParentAccountId":               parentAccountId,
			}, nil
		},
		func(client *quickbooks.Client, ctx context.Context, realmId string, token *quickbooks.BearerToken, startPosition, pageSize int) ([]quickbooks.Account, error) {
			params := quickbooks.RequestParameters{
				Ctx:     ctx,
				RealmId: realmId,
				Token:   token,
			}

			items, err := client.FindAccountsByPage(params, startPosition, pageSize)
			if err != nil {
				return nil, err
			}

			return items, nil
		},
		func(a quickbooks.Account) string {
			return a.Id
		},
		func(a quickbooks.Account) string {
			return a.Status
		},
	)

	tr.Register(account)

	bill := NewQuickBooksDualType(
		"Bill",
		"Bill",
		map[string]fibery.Field{
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
				Name: "Bill Number",
				Type: fibery.Text,
			},
			"TxnDate": {
				Name:    "Bill Date",
				Type:    fibery.DateType,
				SubType: fibery.Day,
			},
			"DueDate": {
				Name:    "Due Date",
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
			"VendorId": {
				Name: "Vendor ID",
				Type: fibery.Text,
				Relation: &fibery.Relation{
					Cardinality:   fibery.MTO,
					Name:          "Vendor",
					TargetName:    "Bills",
					TargetType:    "Vendor",
					TargetFieldID: "id",
				},
			},
			"APAccountId": {
				Name: "AP Account ID",
				Type: fibery.Text,
				Relation: &fibery.Relation{
					Cardinality:   fibery.MTO,
					Name:          "AP Account",
					TargetName:    "Bills",
					TargetType:    "Account",
					TargetFieldID: "id",
				},
			},
			"SalesTermId": {
				Name: "Sales Term ID",
				Type: fibery.Text,
				Relation: &fibery.Relation{
					Cardinality:   fibery.MTO,
					Name:          "Terms",
					TargetName:    "Bills",
					TargetType:    "SalesTerm",
					TargetFieldID: "id",
				},
			},
		},
		func(b quickbooks.Bill) (map[string]any, error) {
			var apAccountId string
			if b.APAccountRef != nil {
				apAccountId = b.APAccountRef.Value
			}

			var salesTermId string
			if b.SalesTermRef != nil {
				salesTermId = b.SalesTermRef.Value
			}

			return map[string]any{
				"id":           b.Id,
				"QBOId":        b.Id,
				"Name":         b.PrivateNote,
				"SyncToken":    b.SyncToken,
				"__syncAction": fibery.SET,
				"DocNumber":    b.DocNumber,
				"TxnDate":      b.TxnDate.Format(fibery.DateFormat),
				"DueDate":      b.DueDate.Format(fibery.DateFormat),
				"PrivateNote":  b.PrivateNote,
				"TotalAmt":     b.TotalAmt,
				"Balance":      b.Balance,
				"VendorId":     b.VendorRef.Value,
				"APAccountId":  apAccountId,
				"SalesTermId":  salesTermId,
			}, nil
		},
		func(client *quickbooks.Client, ctx context.Context, realmId string, token *quickbooks.BearerToken, startPosition, pageSize int) ([]quickbooks.Bill, error) {
			params := quickbooks.RequestParameters{
				Ctx:     ctx,
				RealmId: realmId,
				Token:   token,
			}

			items, err := client.FindBillsByPage(params, startPosition, pageSize)
			if err != nil {
				return nil, err
			}

			return items, nil
		},
		func(b quickbooks.Bill) string {
			return b.Id
		},
		func(b quickbooks.Bill) string {
			return b.Status
		},
	)

	tr.Register(bill)

	billItemLine := NewDependentDualType(
		"BillItemLine",
		"Bill Item Line",
		map[string]fibery.Field{
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
			"Tax": {
				Name:    "Tax",
				Type:    fibery.Text,
				SubType: fibery.Boolean,
			},
			"Billable": {
				Name:    "Billable",
				Type:    fibery.Text,
				SubType: fibery.Boolean,
			},
			"Billed": {
				Name:    "Billed",
				Type:    fibery.Text,
				SubType: fibery.Boolean,
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
			"MarkupPercent": {
				Name: "Markup",
				Type: fibery.Number,
				Format: map[string]any{
					"format":    "Percent",
					"precision": 2,
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
			"BillId": {
				Name: "Bill ID",
				Type: fibery.Text,
				Relation: &fibery.Relation{
					Cardinality:   fibery.MTO,
					Name:          "Bill",
					TargetName:    "Item Lines",
					TargetType:    "Bill",
					TargetFieldID: "id",
				},
			},
			"ItemId": {
				Name: "Item ID",
				Type: fibery.Text,
				Relation: &fibery.Relation{
					Cardinality:   fibery.MTO,
					Name:          "Item",
					TargetName:    "Bill Item Lines",
					TargetType:    "Item",
					TargetFieldID: "id",
				},
			},
			"CustomerId": {
				Name: "Customer ID",
				Type: fibery.Text,
				Relation: &fibery.Relation{
					Cardinality:   fibery.MTO,
					Name:          "Customer",
					TargetName:    "Bill Item Lines",
					TargetType:    "Customer",
					TargetFieldID: "id",
				},
			},
			"ClassId": {
				Name: "Class ID",
				Type: fibery.Text,
				Relation: &fibery.Relation{
					Cardinality:   fibery.MTO,
					Name:          "Class",
					TargetName:    "Bill Item Lines",
					TargetType:    "Class",
					TargetFieldID: "id",
				},
			},
			"MarkupAccountId": {
				Name: "Markup Account ID",
				Type: fibery.Text,
				Relation: &fibery.Relation{
					Cardinality:   fibery.MTO,
					Name:          "Markup Income Account",
					TargetName:    "Bill Item Line Markup",
					TargetType:    "Account",
					TargetFieldID: "id",
				},
			},
		},
		func(b quickbooks.Bill) ([]map[string]any, error) {
			items := []map[string]any{}
			for _, line := range b.Line {
				if line.DetailType == quickbooks.ItemExpenseLine {
					tax := false
					if line.ItemBasedExpenseLineDetail.TaxCodeRef.Value == "TAX" {
						tax = true
					}

					var billable bool
					switch line.ItemBasedExpenseLineDetail.BillableStatus {
					case quickbooks.BillableStatusType:
						billable = true
					case quickbooks.HasBeenBilledStatusType:
						billable = true
					default:
						billable = false
					}

					billed := false
					if line.ItemBasedExpenseLineDetail.BillableStatus == quickbooks.HasBeenBilledStatusType {
						billed = true
					}

					item := map[string]any{
						"id":              fmt.Sprintf("%s:i:%s", b.Id, line.Id),
						"QBOId":           line.Id,
						"Description":     line.Description,
						"__syncAction":    fibery.SET,
						"LineNum":         line.LineNum,
						"Tax":             tax,
						"Billable":        billable,
						"Billed":          billed,
						"Qty":             line.ItemBasedExpenseLineDetail.Qty,
						"UnitPrice":       line.ItemBasedExpenseLineDetail.UnitPrice,
						"MarkupPercent":   line.ItemBasedExpenseLineDetail.MarkupInfo.Percent,
						"Amount":          line.Amount,
						"BillId":          b.Id,
						"ItemId":          line.ItemBasedExpenseLineDetail.ItemRef.Value,
						"CustomerId":      line.AccountBasedExpenseLineDetail.CustomerRef.Value,
						"ClassId":         line.ItemBasedExpenseLineDetail.ClassRef.Value,
						"MarkupAccountId": line.ItemBasedExpenseLineDetail.MarkupInfo.MarkUpIncomeAccountRef.Value,
					}
					items = append(items, item)
				}
			}
			return items, nil
		},
		bill,
		func(b quickbooks.Bill) string {
			return b.Id
		},
		func(b quickbooks.Bill) string {
			return b.Status
		},
		func(b quickbooks.Bill) map[string]struct{} {
			sourceMap := map[string]struct{}{}
			for _, line := range b.Line {
				if line.DetailType == quickbooks.ItemExpenseLine {
					sourceMap[line.Id] = struct{}{}
				}
			}
			return sourceMap
		},
	)

	tr.Register(billItemLine)

	billAccountLine := NewDependentDualType(
		"BillAccountLine",
		"Bill Account Line",
		map[string]fibery.Field{
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
			"Tax": {
				Name:    "Tax",
				Type:    fibery.Text,
				SubType: fibery.Boolean,
			},
			"Billable": {
				Name:    "Billable",
				Type:    fibery.Text,
				SubType: fibery.Boolean,
			},
			"Billed": {
				Name:    "Billed",
				Type:    fibery.Text,
				SubType: fibery.Boolean,
			},
			"MarkupPercent": {
				Name: "Markup",
				Type: fibery.Number,
				Format: map[string]any{
					"format":    "Percent",
					"precision": 2,
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
			"BillId": {
				Name: "Bill ID",
				Type: fibery.Text,
				Relation: &fibery.Relation{
					Cardinality:   fibery.MTO,
					Name:          "Bill",
					TargetName:    "Account Lines",
					TargetType:    "Bill",
					TargetFieldID: "id",
				},
			},
			"AccountId": {
				Name: "Account ID",
				Type: fibery.Text,
				Relation: &fibery.Relation{
					Cardinality:   fibery.MTO,
					Name:          "Category",
					TargetName:    "Bill Account Lines",
					TargetType:    "Account",
					TargetFieldID: "id",
				},
			},
			"CustomerId": {
				Name: "Customer ID",
				Type: fibery.Text,
				Relation: &fibery.Relation{
					Cardinality:   fibery.MTO,
					Name:          "Customer",
					TargetName:    "Bill Account Lines",
					TargetType:    "Customer",
					TargetFieldID: "id",
				},
			},
			"ClassId": {
				Name: "Class ID",
				Type: fibery.Text,
				Relation: &fibery.Relation{
					Cardinality:   fibery.MTO,
					Name:          "Class",
					TargetName:    "Bill Account Lines",
					TargetType:    "Class",
					TargetFieldID: "id",
				},
			},
			"MarkupAccountId": {
				Name: "Markup Account ID",
				Type: fibery.Text,
				Relation: &fibery.Relation{
					Cardinality:   fibery.MTO,
					Name:          "Markup Income Account",
					TargetName:    "Bill Account Line Markup",
					TargetType:    "Account",
					TargetFieldID: "id",
				},
			},
		},
		func(b quickbooks.Bill) ([]map[string]any, error) {
			items := []map[string]any{}
			for _, line := range b.Line {
				if line.DetailType == quickbooks.AccountExpenseLine {
					tax := false
					if line.AccountBasedExpenseLineDetail.TaxCodeRef.Value == "TAX" {
						tax = true
					}

					var billable bool
					switch line.AccountBasedExpenseLineDetail.BillableStatus {
					case quickbooks.BillableStatusType:
						billable = true
					case quickbooks.HasBeenBilledStatusType:
						billable = true
					default:
						billable = false
					}

					billed := false
					if line.AccountBasedExpenseLineDetail.BillableStatus == quickbooks.HasBeenBilledStatusType {
						billed = true
					}

					item := map[string]any{
						"id":              fmt.Sprintf("%s:a:%s", b.Id, line.Id),
						"QBOId":           line.Id,
						"Description":     line.Description,
						"__syncAction":    fibery.SET,
						"LineNum":         line.LineNum,
						"Tax":             tax,
						"Billable":        billable,
						"Billed":          billed,
						"MarkupPercent":   line.AccountBasedExpenseLineDetail.MarkupInfo.Percent,
						"Amount":          line.Amount,
						"BillId":          b.Id,
						"AccountId":       line.AccountBasedExpenseLineDetail.AccountRef.Value,
						"CustomerId":      line.AccountBasedExpenseLineDetail.CustomerRef.Value,
						"ClassId":         line.AccountBasedExpenseLineDetail.ClassRef.Value,
						"MarkupAccountId": line.AccountBasedExpenseLineDetail.MarkupInfo.MarkUpIncomeAccountRef.Value,
					}

					items = append(items, item)
				}
			}
			return items, nil
		},
		bill,
		func(b quickbooks.Bill) string {
			return b.Id
		},
		func(b quickbooks.Bill) string {
			return b.Status
		},
		func(b quickbooks.Bill) map[string]struct{} {
			sourceMap := map[string]struct{}{}
			for _, line := range b.Line {
				if line.DetailType == quickbooks.AccountExpenseLine {
					sourceMap[line.Id] = struct{}{}
				}
			}
			return sourceMap
		},
	)

	tr.Register(billAccountLine)

	billPayment := NewQuickBooksDualType(
		"BillPayment",
		"Bill Payment",
		map[string]fibery.Field{
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
		func(bp quickbooks.BillPayment) (map[string]any, error) {
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
		func(client *quickbooks.Client, ctx context.Context, realmId string, token *quickbooks.BearerToken, startPosition, pageSize int) ([]quickbooks.BillPayment, error) {
			params := quickbooks.RequestParameters{
				Ctx:     ctx,
				RealmId: realmId,
				Token:   token,
			}

			items, err := client.FindBillPaymentsByPage(params, startPosition, pageSize)
			if err != nil {
				return nil, err
			}

			return items, nil
		},
		func(bp quickbooks.BillPayment) string {
			return bp.Id
		},
		func(bp quickbooks.BillPayment) string {
			return bp.Status
		},
	)

	tr.Register(billPayment)

	billPaymentLine := NewDependentDualType(
		"BillPaymentLine",
		"Bill Payment Line",
		map[string]fibery.Field{
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
		func(bp quickbooks.BillPayment) ([]map[string]any, error) {
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
		billPayment,
		func(bp quickbooks.BillPayment) string {
			return bp.Id
		},
		func(bp quickbooks.BillPayment) string {
			return bp.Status
		},
		func(bp quickbooks.BillPayment) map[string]struct{} {
			sourceMap := map[string]struct{}{}
			for _, line := range bp.Line {
				sourceMap[fmt.Sprintf("%s:%s", bp.Id, line.Id)] = struct{}{}
			}
			return sourceMap
		},
	)

	tr.Register(billPaymentLine)

	class := NewQuickBooksDualType(
		"Class",
		"Class",
		map[string]fibery.Field{
			"id": {
				Name: "ID",
				Type: fibery.Id,
			},
			"QBOId": {
				Name: "QBOId",
				Type: fibery.Text,
			},
			"Name": {
				Name: "Base Name",
				Type: fibery.Text,
			},
			"FullyQualifiedName": {
				Name:    "Full Name",
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
			"Active": {
				Name:    "Active",
				Type:    fibery.Text,
				SubType: fibery.Boolean,
			},
			"ParentClassId": {
				Name: "Parent Class ID",
				Type: fibery.Text,
				Relation: &fibery.Relation{
					Cardinality:   fibery.MTO,
					Name:          "Parent Class",
					TargetName:    "Sub-Classes",
					TargetType:    "Class",
					TargetFieldID: "id",
				},
			},
		},
		func(c quickbooks.Class) (map[string]any, error) {
			return map[string]any{
				"id":                 c.Id,
				"QBOId":              c.Id,
				"Name":               c.Name,
				"FullyQualifiedName": c.FullyQualifiedName,
				"SyncToken":          c.SyncToken,
				"__syncAction":       fibery.SET,
				"Active":             c.Active,
				"ParentClassId":      c.ParentRef.Value,
			}, nil
		},
		func(client *quickbooks.Client, ctx context.Context, realmId string, token *quickbooks.BearerToken, startPosition, pageSize int) ([]quickbooks.Class, error) {
			params := quickbooks.RequestParameters{
				Ctx:     ctx,
				RealmId: realmId,
				Token:   token,
			}

			items, err := client.FindClassesByPage(params, startPosition, pageSize)
			if err != nil {
				return nil, err
			}

			return items, nil
		},
		func(c quickbooks.Class) string {
			return c.Id
		},
		func(c quickbooks.Class) string {
			return c.Status
		},
	)

	tr.Register(class)
}
