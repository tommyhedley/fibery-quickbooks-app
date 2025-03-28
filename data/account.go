package data

import (
	"github.com/tommyhedley/fibery-quickbooks-app/pkgs/fibery"
	"github.com/tommyhedley/quickbooks-go"
)

var Account = QuickBooksDualType[quickbooks.Account]{
	QuickBooksType: QuickBooksType[quickbooks.Account]{
		BaseType: fibery.BaseType{
			TypeId:   "Account",
			TypeName: "Account",
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
		},
		schemaGen: func(a quickbooks.Account) (map[string]any, error) {
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
		pageQuery: func(req Request) ([]quickbooks.Account, error) {
			params := quickbooks.RequestParameters{
				Ctx:     req.Ctx,
				RealmId: req.RealmId,
				Token:   req.Token,
			}

			items, err := req.Client.FindAccountsByPage(params, req.StartPosition, req.PageSize)
			if err != nil {
				return nil, err
			}

			return items, nil
		},
	},
	entityId: func(a quickbooks.Account) string {
		return a.Id
	},
	entityStatus: func(a quickbooks.Account) string {
		return a.Status
	},
}

func init() {
	registerType(&Account)
}
