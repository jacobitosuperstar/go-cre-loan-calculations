// The main idea here is do as the package name says investment analysis in
// our real state buys

// TODO: how is the projected return of our investment?
// [X] TaxAssumptions
// [X] DealMetrics
// [X] LoanTerms
// [X] SaleMetrics
// [ ] ReturnOfInvestment
// [ ] Target ROI
// [ ] Objective Search

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

// NewTaxAssumptions returns a new TaxAssumptions struct with the passed in
// values. If there is an error within the values the function returns a
// default struct and an error.
func NewTaxAssumptions (
    lanBuildingValue                float64,
    fixDepreciationTimeLine         int,
    incomeTaxRate                   float64,
    capitalGainsTaxRate             float64,
    depreciationRecaptureTaxRate    float64,
) (
    TaxAssumptions,
    error,
) {
    // Data Validation
    if lanBuildingValue > 1 {
        return TaxAssumptions{}, fmt.Errorf("lanBuildingValue cannot be greater than 1")
    }
    if fixDepreciationTimeLine < 0 {
        return TaxAssumptions{}, fmt.Errorf("fixDepreciationTimeLine cannot be lower than 0")
    }
    if incomeTaxRate > 1 {
        return TaxAssumptions{}, fmt.Errorf("incomeTaxRate cannot be greater than 1")
    }
    if capitalGainsTaxRate > 1 {
        return TaxAssumptions{}, fmt.Errorf("capitalGainsTaxRate cannot be greater than 1")
    }
    if depreciationRecaptureTaxRate > 1 {
        return TaxAssumptions{}, fmt.Errorf("depreciationRecaptureTaxRate cannot be greater than 1")
    }
    // Struct Creation
    taxAssumptions := TaxAssumptions{
        LanBuildingValue: lanBuildingValue,
        FixDepreciationTimeLine: fixDepreciationTimeLine,
        IncomeTaxRate: incomeTaxRate,
        CapitalGainsTaxRate: capitalGainsTaxRate,
        DepreciationRecaptureTaxRate: depreciationRecaptureTaxRate,
    }
    return taxAssumptions, nil
}

// DealInformation is a struct that has all the information regarding the
// buying of the commercial property.
type DealInformation struct {
    PurchasePrice                   int
    ClosingAndRenovations           int
    GoingInCapRate                  float64
    InitRevenue                     float64
    InitOperatingExpenses           float64
    InitCapitalReserves             float64
    ProjRevenueGrowth               float64
    ProjOperatingExpensesGrowth     float64
    ProjCapitalReservesGrowth       float64
}

// NewDealInformation returns a DealInformation struct with the passed in
// values. If there is an error within the values the function returns a
// default struct and an error.
func NewDealInformation (
    purchasePrice               int,
    closingAndRenovations       int,
    goingInCapRate              float64,
    initialRevenue              float64,
    initialExpenses             float64,
    initialCapitalReserves      float64,
) (
    DealInformation,
    error,
) {
    // Data Validation
    if purchasePrice < 0 {
        return DealInformation{}, fmt.Errorf("purchasePrice of the property cannot be lower than 0.")
    }
    if closingAndRenovations > 0 {
        return DealInformation{}, fmt.Errorf("closingAndRenovations of the property cannot be greater than 0.")
    }
    if initialRevenue < 0 {
        return DealInformation{}, fmt.Errorf("initialRevenue cannot be lower than 0.")
    }
    if initialExpenses > 0 {
        return DealInformation{}, fmt.Errorf("initialExpenses cannot be lower than 0.")
    }
    if initialCapitalReserves < 0 {
        return DealInformation{}, fmt.Errorf("initialCapitalReserves cannot be lower than 0.")
    }
    if purchasePrice == 0 {
        initialNOI := initialRevenue + initialExpenses
        purchasePrice = int(math.Floor(initialNOI/goingInCapRate))
    }

    // Struct Creation
    dealInformation := DealInformation{
        PurchasePrice: purchasePrice,
        ClosingAndRenovations: closingAndRenovations,
        GoingInCapRate: goingInCapRate,
        InitRevenue: initialRevenue,
        InitOperatingExpenses: initialExpenses,
        InitCapitalReserves: initialCapitalReserves,
    }
    return dealInformation, nil
}


// SaleTerms is a struc that has all the sale information regarding the sale of
// the sale of the property.
type SaleTerms struct {
    ExitCapRate     float64
    CostOfSale      float64
    SaleYear        int
}

// ProjectedSalePrice returns the projected sale price of real state.
func (st SaleTerms) ProjectedSalePrice (noi float64) float64 {
    projected_sale_price := utils.Round2(noi / st.ExitCapRate)
    projected_sale_price = projected_sale_price - (projected_sale_price * st.CostOfSale)
    return utils.Round2(projected_sale_price)
}


// NewSaleTerms creates a new SaleTerms struct. If the data is valid, a new
// SaleTerms struct is returned, if not an initialized struct is returned with
// the errors.
func NewSaleTerms(
    exitCapRate     float64,
    costOfSale      float64,
    saleYear        int,
) (
    saleTerms SaleTerms,
    err error,
){
    // Data Validation
    if costOfSale <= 0 {
        return saleTerms, fmt.Errorf("costOfSale cannot be lower than 0")
    }
    if saleYear <= 0 {
        return saleTerms, fmt.Errorf("saleYear cannot be lower than 0")
    }
    // Struct Creation
    saleTerms = SaleTerms{
        ExitCapRate: exitCapRate,
        CostOfSale: costOfSale,
        SaleYear: saleYear,
    }
    return saleTerms, nil
}

// ROI of the totallity of the deal.
type ReturnOfInvestment struct {
    taxMetrics      TaxAssumptions
    dealMetrics     DealInformation
    loanMetrics     ls.LoanSizer
    saleMetrics     SaleTerms
}

// Constructor

// NewReturnOfInvestment constructs the ReturnOfInvestment struct with the
// composite information of the deal numbers.
func NewReturnOfInvestment(
    // DealInformation
    purchasePrice                       int,
    closingAndRenovations               int,
    goinginCapRate                      float64,
    initialRevenue                      float64,
    initialExpenses                     float64,
    initialCapitalReserves              float64,
    projectedRevenueGrowth              float64,
    projectedExpensesGrowth             float64,
    projectedCapitalReservesGrowth      float64,
    // LoanSizer
    maxLTV                              float64,
    minDSCR                             float64,
    amortization                        int,
    term                                int,
    interestRate                        float64,
    ioPeriod                            int,
    requestedLoanAmount                 int,
    loanOriginationFees                 float64,
    // TaxAssumptions
    lanBuildingValue                    float64,
    fixDepreciationTimeLine             int,
    incomeTaxRate                       float64,
    capitalGainsTaxRate                 float64,
    depreciationRecaptureTaxRate        float64,
    // SaleTerms
    exitCapRate                         float64,
    costOfSale                          float64,
    saleYear                            int,
) (
    ReturnOfInvestment,
    error,
) {
    dealMetrics, err := NewDealInformation(
        purchasePrice,
        closingAndRenovations,
        goinginCapRate,
        initialRevenue,
        initialExpenses,
        initialCapitalReserves,
    )
    if err != nil {
        return ReturnOfInvestment{}, fmt.Errorf("DealMetrics Internal error: %v", err)
    }

    initialNOI := dealMetrics.InitRevenue + dealMetrics.InitOperatingExpenses
    loanSizer, err := ls.NewLoanSizer(
        maxLTV,
        minDSCR,
        amortization,
        term,
        ioPeriod,
        interestRate,
        dealMetrics.PurchasePrice,
        initialNOI,
        requestedLoanAmount,
        loanOriginationFees,
    )
    if err != nil {
        return ReturnOfInvestment{}, fmt.Errorf("LoanSizer Internal error: %v", err)
    }
    taxAssumptions, err := NewTaxAssumptions(
        lanBuildingValue,
        fixDepreciationTimeLine,
        incomeTaxRate,
        capitalGainsTaxRate,
        depreciationRecaptureTaxRate,
    )
    if err != nil {
        return ReturnOfInvestment{}, fmt.Errorf("TaxAssumptions Internal error: %v", err)
    }

    saleTerms, err := NewSaleTerms(
        exitCapRate,
        costOfSale,
        saleYear,
    )
    if err != nil {
        return ReturnOfInvestment{}, fmt.Errorf("SaleTerms Internal error: %v", err)
    }

    roi := ReturnOfInvestment{
        taxMetrics: taxAssumptions,
        dealMetrics: dealMetrics,
        loanMetrics: loanSizer,
        saleMetrics: saleTerms,
    }
    return roi, nil
}

// Calculation methods

// Internal
// ...

// External

// AdquisitionCost returns the AdquisitionCost of Deal
func (roi ReturnOfInvestment) AdquisitionCost () (float64, error)  {
    mla, err := roi.loanMetrics.MaximumLoanAmount()
    if err != nil {
        return 0.0, fmt.Errorf("MaximumLoanAmount internal error: %v", err)
    }
    adquisitionCost := - float64(roi.dealMetrics.PurchasePrice) +
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
        current_noi := utils.Round2(revenue + expense)
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

// Target ROI of the investment
type TargetReturnOfInvestment struct {
    IRR                         float64
    EquityMultiple              float64
    AverageCashOnCashReturn     float64
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
