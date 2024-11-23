// The main idea here is do as the package name says investment analysis in
// our real state buys

// TODO: how is the projected return of our investment?
// [ ] TaxAssumptions
// [ ] DealMetrics
// [ ] LoanTerms
// [ ] SaleMetrics
// [ ] ReturnOfInvestment

package investment_analysis

import (
    "fmt";
    "math";
    utils "github.com/jacobitosuperstar/go-cre-loan-calculations/internal/utils";
    ls "github.com/jacobitosuperstar/go-cre-loan-calculations/loan_sizer";
)

// TaxAssumptions is a struct that has all the taxes information regarding the
// deal.
type TaxAssumptions struct {
    LanBuildingValue                float64
    FixDepreciationTimeLine         int
    IncomeTaxRate                   float64
    CapitalGainsTaxRate             float64
    DepreciationRecaptureTaxRate    float64
}

// DealInformation is a struct that has all the information regarding the
// buying of the commercial property.
type DealInformation struct {
    PurchasePrice               int
    ClosingAndRenovations       int
    GoingInCapRate              float64
    InitRevenue                 float64
    InitOperatingExpenses       float64
    InitCapitalReserves         float64
    ProjRevenueGrowth           float64
    ProjOperatingExpensesGrowth float64
    ProjCapitalReservesGrowth   float64
}

// SaleTerms is a struc that has all the sale information regarding the sale of
// the sale of the property.
type SaleTerms struct {
    ExitCapRate         float64
    CostOfSale          float64
    SaleYear            int
}

// ProjectedSalePrice returns the projected sale price of real state.
func (st SaleTerms) ProjectedSalePrice (noi float64) float64 {
    projected_sale_price := utils.Round2(noi / st.ExitCapRate)
    projected_sale_price = projected_sale_price - (projected_sale_price * st.CostOfSale)
    return utils.Round2(projected_sale_price)
}

// ROI of the totallity of the deal.
type ReturnOfInvestment struct {
    taxMetrics              TaxAssumptions
    dealMetrics             DealInformation
    loanMetrics             ls.LoanSizer
    saleMetrics             SaleTerms
}

// AdquisitionCost returns the AdquisitionCost of Deal
func (roi ReturnOfInvestment) AdquisitionCost () (float64, error)  {
    mla, err := roi.loanMetrics.MaximumLoanAmount()
    if err != nil {
        return 0.0, fmt.Errorf("MaximumLoanAmount internal error: %v", err)
    }
    adquisitionCost := - float64(roi.dealMetrics.PurchasePrice) -
    float64(roi.dealMetrics.ClosingAndRenovations) -
    (roi.loanMetrics.LoanOriginationFees * mla) +
    mla
    return utils.Round2(adquisitionCost), nil
}

// CashOnCashReturn returns the made money in reference to the money invested
// to adquire the property.
func (roi ReturnOfInvestment) CashOnCashReturn (net_cash_flow float64)  (float64, error) {
    adq_cost, err := roi.AdquisitionCost()
    if err != nil {
        return 0.0, fmt.Errorf("AdquisitionCost internal error: %v", err)
    }
    return utils.Round4(math.Abs(net_cash_flow/adq_cost)), nil
}

// NetCashFlowProjection returns the net cash flow projection of the Deal.
// Negative values are payments that need to be done, Positive values are money
// given.
func (roi ReturnOfInvestment) NetCashFlowProjection () ([]map[string]interface{}, error) {
    var net_cash_flow_projection []map[string]interface{}

    adquisition_cost, err := roi.AdquisitionCost()
    if err != nil {
        return net_cash_flow_projection, fmt.Errorf("AdquisitionCost internal error: %v", err)
    }

    net_cash_flow_projection = append(
        net_cash_flow_projection,
        map[string]interface{} {
            "net_cash_flow": adquisition_cost,
        },
    )

    revenue := roi.dealMetrics.InitRevenue
    expense := roi.dealMetrics.InitOperatingExpenses
    reserve := roi.dealMetrics.InitCapitalReserves

    // getting the building value
    purchase_price := roi.dealMetrics.PurchasePrice
    building_value := utils.Round2(float64(purchase_price) * (1.0 - roi.taxMetrics.LanBuildingValue))

    // depreciation of the building
    building_depreciation := utils.Round2(- building_value/float64(roi.taxMetrics.FixDepreciationTimeLine))

    // payment distribution of the loan
    ppmt, ipmt, err := roi.loanMetrics.PaymentDistribution()
    if err != nil {
        return net_cash_flow_projection, fmt.Errorf("PaymentDistribution internal error: %v", err)
    }
    // payments of the loan
    current_pmt, err := roi.loanMetrics.LoanPayment()
    if err != nil {
        return net_cash_flow_projection, fmt.Errorf("LoanPayment internal error: %v", err)
    }

    // BalloonPayment at the end of the term
    balloonpayment, err := roi.loanMetrics.BalloonPayment()
    if err != nil {
        return net_cash_flow_projection, fmt.Errorf("BalloonPayment internal error: %v", err)
    }

    // Iterating over the term and appending the values to the
    // NetCashFlowProjection slice.
    for i := 0; i < roi.loanMetrics.Term; i++ {
        // this year NOI
        current_noi := utils.Round2(revenue - expense)
        // this year interest and principal payments
        current_ppmt := ppmt[i]
        current_ipmt := ipmt[i]
        // cashflow after debt service
        cfads := utils.Round2(current_noi + reserve + current_pmt)
        // depreciation expense
        depreciation_expense := 0.0
        if i < roi.taxMetrics.FixDepreciationTimeLine {
            depreciation_expense = building_depreciation
        }
        // income tax
        income_tax := utils.Round2(- (current_noi + current_ipmt + depreciation_expense) * roi.taxMetrics.IncomeTaxRate)
        implied_income_tax := utils.Round4(math.Abs(income_tax/cfads))
        // net cashflow
        ncf := utils.Round2(cfads + income_tax)
        // cash on cash return
        cocr, err := roi.CashOnCashReturn(ncf)
        if err != nil {
            return net_cash_flow_projection, fmt.Errorf("CashOnCashReturn internal error: %v", err)
        }

        net_cash_flow_projection = append(
            net_cash_flow_projection,
            map[string]interface{} {
                "year": i + 1,
                "revenue": revenue,
                "expense": expense,
                "noi": current_noi,
                "reserve": reserve,
                "principal_payment": current_ppmt,
                "interest_payment": current_ipmt,
                "cashflow_after_debt_service": cfads,
                "depreciation_expense": depreciation_expense,
                "income_tax": income_tax,
                "implied_income_tax": implied_income_tax,
                "net_cash_flow": ncf,
                "cash_on_cash_return": cocr,
            },
        )
        revenue = utils.Round2(revenue + revenue * roi.dealMetrics.ProjRevenueGrowth)
        expense = utils.Round2(expense + expense * roi.dealMetrics.ProjOperatingExpensesGrowth)
        reserve = utils.Round2(reserve + reserve * roi.dealMetrics.ProjCapitalReservesGrowth)
    }
    after_term_noi := utils.Round2(revenue + expense)
    // Adding the cashflow after the sell of the property
    // sale with the projected NOI
    projected_sale_price := roi.saleMetrics.ProjectedSalePrice(after_term_noi)
    // capital gains tax
    cg := projected_sale_price -
        float64(roi.dealMetrics.PurchasePrice) -
        float64(roi.dealMetrics.ClosingAndRenovations)
    cg = utils.Round2(cg)
    cgt := cg * roi.taxMetrics.CapitalGainsTaxRate
    // Depreciation Recapture tax
    drt := building_depreciation *
        float64(roi.saleMetrics.SaleYear) *
        roi.taxMetrics.DepreciationRecaptureTaxRate
    drt = utils.Round2(drt)
    // Sale calculations
    sale := net_cash_flow_projection[roi.loanMetrics.Term - 1]
    sale_net_cash_flow := sale["net_cash_flow"].(float64)
    sale_net_cash_flow = sale_net_cash_flow +
        projected_sale_price +
        drt +
        cgt +
        balloonpayment
    sale_net_cash_flow = utils.Round2(sale_net_cash_flow)
    sale["net_cash_flow"] = sale_net_cash_flow
    sale["sale_price"] = projected_sale_price
    sale["depreciation_recapture_tax"] = drt
    sale["capital_gains_tax"] = cgt
    // Setting the value
    return net_cash_flow_projection, nil
}

func (roi *ReturnOfInvestment) SetIRR () error {
    return nil
}

func (roi *ReturnOfInvestment) SetEquityMultiple () error {
    return nil
}

func (roi *ReturnOfInvestment) SetAverageCashOnCashReturn () error {
    return nil
}

// InitReturnOfInvestment sets the calculated terms in the ReturnOfInvestment
// struct.
func InitReturnOfInvestment(roi ReturnOfInvestment) ReturnOfInvestment {
    return roi
}

// Target ROI of the investment
type TargetReturnOfInvestment struct {
    IRR                     float64     `json:"internal_rate_of_return"`
    EquityMultiple          float64     `json:"equity_multiple"`
    AverageCashOnCashReturn float64     `json:"average_cash_on_cash_return"`
}

//
// func InitTargetReturnOfInvestment (roi ReturnOfInvestment)  (ReturnOfInvestment, error) {
//     roi.SetAdquisitionCost()
//     err := roi.SetNetCashFlowProjection()
//     if err != nil {
//         return roi, err
//     }
//     return roi, nil
// }
