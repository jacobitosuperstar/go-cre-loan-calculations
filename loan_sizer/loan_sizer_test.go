// Testing the Loan Sizer functionality
// TODO: Testing needs
// [X] maximum loan amount
// [X] loan io period payment
// [X] loan payment
// [X] balloon payment
// [X] payments distribution

package loan_sizer
import "testing"


type PaymentDistributionResults struct {
  ipmt [10]float64
  ppmt [10]float64
}

type TestResults struct {
  MaximumLoanAmount float64
  IOLoanPayment float64
  LoanPayment float64
  BalloonPayment float64
  PaymentDistribution PaymentDistributionResults
}

func TestLoanSizerFunctionality(t *testing.T){
    var testCases = []struct {
        name string
        ls LoanSizer
        want TestResults
    }{
        {
            name: "Zero Values",
            ls: LoanSizer{
                MaxLTV: 0.70,
                MinDSCR: 1.40,
                Amortization: 30,
                Term: 10,
                IOPeriod: 3,
                Rate: 0.0045,
                PropertyValue: 0,
                NOI: 0,
                RequestedLoanAmount: 0,
                LoanOriginationFees: 0.001,
            },
            want: TestResults{
              MaximumLoanAmount: 0.0,
              IOLoanPayment: 0.0,
              LoanPayment: 0.0,
              BalloonPayment: 0.0,
              PaymentDistribution: PaymentDistributionResults{
                ipmt: [10]float64{0,0,0,0,0,0,0,0,0,0},
                ppmt: [10]float64{0,0,0,0,0,0,0,0,0,0},
              },
            },
        },
        {
            name: "Valid Inputs, LTV as MLA",
            ls: LoanSizer{
                MaxLTV: 0.70,
                MinDSCR: 1.40,
                Amortization: 30,
                Term: 10,
                IOPeriod: 3,
                Rate: 0.0045,
                PropertyValue: 1000,
                NOI: 500,
                RequestedLoanAmount: 900,
                LoanOriginationFees: 0.001,
            },
            want: TestResults{
              MaximumLoanAmount: 700,
              IOLoanPayment: -3.15,
              LoanPayment: -25,
              BalloonPayment: 544.95,
              PaymentDistribution: PaymentDistributionResults{
                ipmt: [10]float64{-3.15, -3.15, -3.15, -3.15, -3.05, -2.95, -2.85, -2.75, -2.65, -2.55},
                ppmt: [10]float64{0, 0, 0, -21.85, -21.95, -22.05, -22.15, -22.25, -22.35, -22.45},
              },
            },
        },
        {
            name: "Valid Inputs, MinDSCR as MLA",
            ls: LoanSizer{
                MaxLTV: 0.70,
                MinDSCR: 1.40,
                Amortization: 30,
                Term: 10,
                IOPeriod: 3,
                Rate: 0.0045,
                PropertyValue: 1000,
                NOI: 250,
                RequestedLoanAmount: 900,
                LoanOriginationFees: 0.001,
            },
            want: TestResults{
              MaximumLoanAmount: 700,
              IOLoanPayment: -3.15,
              LoanPayment: -25,
              BalloonPayment: 544.95,
              PaymentDistribution: PaymentDistributionResults{
                ipmt: [10]float64{-3.15, -3.15, -3.15, -3.15, -3.05, -2.95, -2.85, -2.75, -2.65, -2.55},
                ppmt: [10]float64{0,0,0, -21.85, -21.95, -22.05, -22.15, -22.25, -22.35, -22.45},
              },
            },
        },
        {
            name: "Valid Inputs, RequestedLoanAmount as MLA",
            ls: LoanSizer{
                MaxLTV: 0.70,
                MinDSCR: 1.40,
                Amortization: 30,
                Term: 10,
                IOPeriod: 3,
                Rate: 0.0045,
                PropertyValue: 1000,
                NOI: 400,
                RequestedLoanAmount: 400,
                LoanOriginationFees: 0.001,
            },
            want: TestResults{
              MaximumLoanAmount: 400,
              IOLoanPayment: -1.8,
              LoanPayment: -14.28,
              BalloonPayment: 311.45,
              PaymentDistribution: PaymentDistributionResults{
                ipmt: [10]float64{-1.8,-1.8,-1.8,-1.8,-1.74,-1.69,-1.63,-1.57,-1.52,-1.46},
                ppmt: [10]float64{0,0,0,-12.48,-12.54,-12.59,-12.65,-12.71,-12.76,-12.82},
              },
            },
        },
        {
            name: "Valid Inputs, Term and Amortization equal",
            ls: LoanSizer{
                MaxLTV: 0.70,
                MinDSCR: 1.40,
                Amortization: 10,
                Term: 10,
                IOPeriod: 0,
                Rate: 0.0045,
                PropertyValue: 1000,
                NOI: 250,
                RequestedLoanAmount: 900,
                LoanOriginationFees: 0.001,
            },
            want: TestResults{
              MaximumLoanAmount: 700,
              IOLoanPayment: -3.15,
              LoanPayment: -71.74,
              BalloonPayment: 0.0,
              PaymentDistribution: PaymentDistributionResults{
                ipmt: [10]float64{-3.15, -2.84, -2.53, -2.22, -1.91, -1.59, -1.28, -0.96, -0.64, -0.32},
                ppmt: [10]float64{-68.59, -68.9, -69.21, -69.52, -69.83, -70.15, -70.46, -70.78, -71.1, -71.46},
              },
            },
        },
    }

    for _, test := range testCases {
        // Maximum Loan Amount tests
        t.Run(test.name, func(t *testing.T) {
            got, err := test.ls.MaximumLoanAmount()
            if err != nil {
                t.Errorf(
                  "MaximumLoanAmount error: %v",
                  err,
                )
            }
            if got != test.want.MaximumLoanAmount {
                t.Errorf(
                  "MaximumLoanAmount got: %g, wanted: %g",
                  got,
                  test.want.MaximumLoanAmount,
                )
            }
        })
        // IO Loan Payment tests
        t.Run(test.name, func(t *testing.T) {
            got, err := test.ls.IOLoanPayment()
            if err != nil {
                t.Errorf(
                  "IOLoanPayment error: %v",
                  err,
                )
            }
            if got != test.want.IOLoanPayment {
                t.Errorf(
                  "IOLoanPayment got: %g, wanted: %g",
                  got,
                  test.want.IOLoanPayment,
                )
            }
        })
        // Loan Payment tests
        t.Run(test.name, func(t *testing.T) {
            got, err := test.ls.LoanPayment()
            if err != nil {
                t.Errorf(
                  "LoanPayment error: %v",
                  err,
                )
            }
            if got != test.want.LoanPayment {
                t.Errorf(
                  "LoanPayment got: %g, wanted: %g",
                  got,
                  test.want.LoanPayment,
                )
            }
        })
        // BalloonPayment tests
        t.Run(test.name, func(t *testing.T) {
            got, err := test.ls.BalloonPayment()
            if err != nil {
                t.Errorf(
                  "BalloonPayment error: %v",
                  err,
                )
            }
            if got != test.want.BalloonPayment {
                t.Errorf(
                  "BalloonPayment got: %g, wanted: %g",
                  got,
                  test.want.BalloonPayment,
                )
            }
        })
        // PaymentDistribution tests
        t.Run(test.name, func(t *testing.T) {
            ppmt, ipmt, err := test.ls.PaymentDistribution()

            if err != nil {
                t.Errorf("PaymentDistribution internal error: %v", err)
                return
            }

            if len(ppmt) != len(test.want.PaymentDistribution.ppmt) {
                t.Errorf(
                  "PaymentDistribution ppmt got: %g, wanted: %g",
                  ppmt,
                  test.want.PaymentDistribution.ppmt,
                )
                return
            }

            for i := range ppmt {
                if ppmt[i] != test.want.PaymentDistribution.ppmt[i] {
                    t.Errorf(
                      "PaymentDistribution ppmt got: %g, wanted: %g",
                      ppmt[i],
                      test.want.PaymentDistribution.ppmt[i],
                    )
                }
            }

            if len(ipmt) != len(test.want.PaymentDistribution.ipmt) {
                t.Errorf(
                  "PaymentDistribution ipmt got: %g, wanted: %g",
                  ipmt,
                  test.want.PaymentDistribution.ipmt,
                )
                return
            }

            for i := range ppmt {
                if ipmt[i] != test.want.PaymentDistribution.ipmt[i] {
                    t.Errorf(
                      "PaymentDistribution ipmt got: %g, wanted: %g",
                      ipmt[i],
                      test.want.PaymentDistribution.ipmt[i],
                    )
                }
            }
        })
    }
}
