package investment_analysis
import (
  "testing";
  "log";
  ls "github.com/jacobitosuperstar/go-cre-loan-calculations/loan_sizer";
)


func TestROI(t *testing.T) {
    var testCases = []struct {
        input  map[string]interface{}
        // want  []map[string]interface{}
    }{
      {
          map[string]interface{} {
              "purchase_price": 0,
              "closing_an_renovations": 0,
              "going_in_cap_rate": 0,
              "initial_revenue": 0,
              "initial_expenses": 0,
              "initial_capital_reserves": 0,
              "projected_revenue_growth": 0,
              "projected_expenses_growth": 0,
              "projected_capital_reserves_growth": 0,
              // LoanSizer
              "max_ltv": 0,
              "min_dscr": 0,
              "amortization": 0,
              "term": 0,
              "interest_rate": 0,
              "io_period": 0,
              "requested_loan_amount": 0,
              "loan_origination_fees": 0,
              // TaxAssumptions
              "lan_building_value": 0,
              "fix_depreciation_time_line": 0,
              "income_tax_rate": 0,
              "capitalGains_tax_rate": 0,
              "depreciation_recapture_tax_rate": 0,
              // SaleTerms
              "exit_cap_rate": 0,
              "cost_of_sale": 0,
              "sale_year": 0,
          },
      },
    }

    for _, test := range testCases {
      roi, err := NewReturnOfInvestment()
      if err != nil {
          t.Errorf("ReturnOfInvestment internal error: %v", err)
          return
      }
      log.Printf("roi: %v",roi)
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
