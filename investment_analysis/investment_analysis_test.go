package investment_analysis
import (
  "testing";
  "log";
  ls "github.com/jacobitosuperstar/go-cre-loan-calculations/loan_sizer";
)


func TestROI(t *testing.T) {
    var testCases = []struct {
        input  ReturnOfInvestment
        // want  []map[string]interface{}
    }{
      {
        ReturnOfInvestment{
          taxMetrics: TaxAssumptions{
            LanBuildingValue: 0.3,
            FixDepreciationTimeLine: 27,
            IncomeTaxRate: 0.25,
            CapitalGainsTaxRate: 0.15,
            DepreciationRecaptureTaxRate: 0.25,
          },
          dealMetrics: DealInformation{
            PurchasePrice: 6500000,
            ClosingAndRenovations: 225000,
            GoingInCapRate: 0.0596,
            InitRevenue: 687500,
            InitOperatingExpenses: 300000,
            InitCapitalReserves: 7500,
            ProjRevenueGrowth: 0.0350,
            ProjOperatingExpensesGrowth: 0.0250,
            ProjCapitalReservesGrowth: 0.0250,
          },
          loanMetrics: ls.LoanSizer{
            MaxLTV: 0.70,
            MinDSCR: 1.25,
            Amortization: 30,
            Term: 10,
            IOPeriod: 2,
            Rate: 0.0450,
            PropertyValue: 10000000,
            NOI: 0.0,
            RequestedLoanAmount: 10000000,
            LoanOriginationFees: 0.01,
          },
          saleMetrics: SaleTerms{
            ExitCapRate: 0.0650,
            CostOfSale: 0.0250,
            SaleYear: 10,
          },
        },
      },
    }

    for _, test := range testCases {
      net_cash_flow_projection, err := test.input.NetCashFlowProjection()
      log.Printf("roi: %+v\n", net_cash_flow_projection)
      if err != nil {
          t.Errorf("NetCashFlowProjection got: %v, error: %v", net_cash_flow_projection, err)
          return
      }
      // log.Printf("roi: %+v\n", net_cash_flow_projection)
      t.Errorf("got: %v", net_cash_flow_projection)
      ppmt, ipmt, err := test.input.loanMetrics.PaymentDistribution()
      if err != nil {
          t.Errorf("NetCashFlowProjection error: %v", err)
          return
      }
      t.Errorf("LOOK HERE ppmt: %v, ipmt: %v", ppmt, ipmt)
    }
}
