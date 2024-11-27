package investment_analysis
import (
  "testing";
  // "math";
  // "log";
)

type TestROI struct {
    // DealInformation
    purchasePrice                       int
    closingAndRenovations               int
    goinginCapRate                      float64
    initialRevenue                      float64
    initialExpenses                     float64
    initialCapitalReserves              float64
    projectedRevenueGrowth              float64
    projectedExpensesGrowth             float64
    projectedCapitalReservesGrowth      float64
    // LoanSizer
    maxLTV                              float64
    minDSCR                             float64
    amortization                        int
    term                                int
    interestRate                        float64
    ioPeriod                            int
    requestedLoanAmount                 int
    loanOriginationFees                 float64
    // TaxAssumptions
    lanBuildingValue                    float64
    fixDepreciationTimeLine             int
    incomeTaxRate                       float64
    capitalGainsTaxRate                 float64
    depreciationRecaptureTaxRate        float64
    // SaleTerms
    exitCapRate                         float64
    costOfSale                          float64
    saleYear                            int
}


func TestROIMethods(t *testing.T) {
    var testCases = []struct {
        input  TestROI
        want  []map[string]interface{}
    }{
      {
          TestROI{
              // DealInformation
              purchasePrice: 6500000,
              closingAndRenovations: -225000,
              goinginCapRate: 0.0596,
              initialRevenue: 687500,
              initialExpenses: -300000,
              // initialCapitalReserves: -7500,
              initialCapitalReserves: 7500,
              projectedRevenueGrowth: 0.0350,
              projectedExpensesGrowth: 0.0250,
              projectedCapitalReservesGrowth: 0.0250,
              // LoanSizer
              maxLTV: 0.70,
              minDSCR: 1.25,
              amortization: 30,
              term: 10,
              interestRate: 0.045,
              ioPeriod: 2,
              requestedLoanAmount: 6500000,
              loanOriginationFees: 0.01,
              // TaxAssumptions
              lanBuildingValue: 0.30,
              fixDepreciationTimeLine: 27,
              incomeTaxRate: 0.25,
              capitalGainsTaxRate: 0.15,
              depreciationRecaptureTaxRate: 0.25,
              // SaleTerms
              exitCapRate: 0.0650,
              costOfSale: 0.0250,
              saleYear: 10,
          },
          []map[string]interface{}{
              {
                  "net_cash_flow": -2220500.0,
              },
              {
                  "cash_on_cash_return": 0.0505,
                  "cashflow_after_debt_service": 115668.48,
                  "depreciation_expense": -168518.52,
                  "expense": -300000.0,
                  "implied_income_tax": 0.0308,
                  "income_tax": -3557.87,
                  "interest_payment": -204750.0,
                  "net_cash_flow": 112110.61,
                  "noi": 387500.0,
                  "principal_payment":0.0,
                  "reserve": 7500.0,
                  "revenue":687500.0,
                  "year":1,
              },
              {
                  "cash_on_cash_return": 0.0738,
                  "cashflow_after_debt_service": 186725.27,
                  "depreciation_expense": -168518.52,
                  "expense": -331143.87,
                  "implied_income_tax": 0.1223,
                  "income_tax": -22843.26,
                  "interest_payment": -197886.64,
                  "net_cash_flow": 163882.01,
                  "noi": 457778.19,
                  "principal_payment": -81444.88,
                  "reserve": 8278.6,
                  "revenue": 788922.06,
                  "year": 5,
              },
              {
                  "cash_on_cash_return": 0.1074,
                  "cashflow_after_debt_service": 292368.0,
                  "depreciation_expense": -168518.52,
                  "expense": -374658.89,
                  "implied_income_tax": 0.1847,
                  "income_tax": -53994.54,
                  "interest_payment": -177836.38,
                  "net_cash_flow": 238373.46,
                  "noi": 562333.04,
                  "principal_payment": -101495.14,
                  "reserve": 9366.48,
                  "revenue": 936991.93,
                  "year": 10,
              },
          },
      },
    }

    for _, test := range testCases {
      roi, err := NewReturnOfInvestment(
          // DealInformation
          test.input.purchasePrice,
          test.input.closingAndRenovations,
          test.input.goinginCapRate,
          test.input.initialRevenue,
          test.input.initialExpenses,
          test.input.initialCapitalReserves,
          test.input.projectedRevenueGrowth,
          test.input.projectedExpensesGrowth,
          test.input.projectedCapitalReservesGrowth,
          // LoanSizer
          test.input.maxLTV,
          test.input.minDSCR,
          test.input.amortization,
          test.input.term,
          test.input.interestRate,
          test.input.ioPeriod,
          test.input.requestedLoanAmount,
          test.input.loanOriginationFees,
          // TaxAssumptions
          test.input.lanBuildingValue,
          test.input.fixDepreciationTimeLine,
          test.input.incomeTaxRate,
          test.input.capitalGainsTaxRate,
          test.input.depreciationRecaptureTaxRate,
          // SaleTerms
          test.input.exitCapRate,
          test.input.costOfSale,
          test.input.saleYear,
      )
      if err != nil {
          t.Errorf("ReturnOfInvestment internal error: %v", err)
          return
      }
      netCashFlowProjection, err := roi.NetCashFlowProjection()
      if err != nil {
          t.Errorf("NetCashFlowProjection internal error: %v", err)
          return
      }
      // Because in the test we picked 4 values, year 0, year 1, year 5 and
      // year 10, an slice with those elements from the netCashFlowProjection
      // will be created. Then we will compare both of the lists.
      // XXX: This only works becase this is a giant singular test. Change
      // the approach if necessary.
      var sliceForComparisson []map[string]interface{}

      sliceForComparisson = append(sliceForComparisson, netCashFlowProjection[0])
      sliceForComparisson = append(sliceForComparisson, netCashFlowProjection[1])
      sliceForComparisson = append(sliceForComparisson, netCashFlowProjection[5])
      sliceForComparisson = append(sliceForComparisson, netCashFlowProjection[10])

      // t.Errorf("roi: %v", netCashFlowProjection)
      // return

      for i := 0; i < len(test.want); i++  {
          for key,value_test := range test.want[i] {
              value_want, ok := sliceForComparisson[i][key]
              if !ok || value_test != value_want {
                  t.Errorf("there is a difference in the element: %v, between the values in the key: %v, value_test: %v, value_want: %v", i, key, value_test, value_want)
              }
          }
      }
      // log.Printf("roi: %+v\n", net_cash_flow_projection)
      // t.Errorf("got: %v", net_cash_flow_projection)
      // ppmt, ipmt, err := roi.loanMetrics.PaymentDistribution()
      // if err != nil {
      //     t.Errorf("NetCashFlowProjection error: %v", err)
      //     return
      // }
      // t.Errorf("LOOK HERE ppmt: %v, ipmt: %v", ppmt, ipmt)
    }
}
