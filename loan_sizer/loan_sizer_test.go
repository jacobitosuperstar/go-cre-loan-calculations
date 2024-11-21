// Testing the Loan Sizer functionality
// TODO: Testing needs
// [X] maximum loan amount
// [X] loan io period payment
// [X] loan payment
// [X] balloon payment
// [ ] payments distribution

package loan_sizer
import "testing"

func TestMaximumLoanAmount(t *testing.T){
    var testCases = []struct {
        name string
        ls LoanSizer
        want float64
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
            want: 0.00,
        },
        {
            name: "Valid Inputs, LTV",
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
            want: 700,
        },
        {
            name: "Valid Inputs, MinDSCR",
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
            want: 700,
        },
        {
            name: "Valid Inputs, RequestedLoanAmount",
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
            want: 400,
        },
    }

    for _, test := range testCases {
        t.Run(test.name, func(t *testing.T) {
            got, err := test.ls.MaximumLoanAmount()
            if err != nil {
                t.Errorf("error: %v", err)
            }
            if got != test.want {
                t.Errorf("got: %g, wanted: %g", got, test.want)
            }
        })
    }
}


func TestIOPayment(t *testing.T){
    var testCases = []struct {
        name string
        ls LoanSizer
        want float64
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
            want: 0,
        },
        {
            name: "Valid Inputs",
            ls: LoanSizer{
                MaxLTV: 0.70,
                MinDSCR: 1.40,
                Amortization: 30,
                Term: 10,
                IOPeriod: 0,
                Rate: 0.0045,
                PropertyValue: 1000,
                NOI: 250,
                RequestedLoanAmount: 900,
                LoanOriginationFees: 0.001,
            },
            want: -3.15,
        },
    }

    for _, test := range testCases {
        t.Run(test.name, func(t *testing.T) {
            got, err := test.ls.IOLoanPayment()
            if err != nil {
                t.Errorf("error: %v", err)
            }
            if got != test.want {
                t.Errorf("got: %g, wanted: %g", got, test.want)
            }
        })
    }
}


func TestPayment(t *testing.T){
    var testCases = []struct {
        name string
        ls LoanSizer
        want float64
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
            want: 0,
        },
        {
            name: "Valid Inputs",
            ls: LoanSizer{
                MaxLTV: 0.70,
                MinDSCR: 1.40,
                Amortization: 30,
                Term: 10,
                IOPeriod: 0,
                Rate: 0.0045,
                PropertyValue: 1000,
                NOI: 250,
                RequestedLoanAmount: 900,
                LoanOriginationFees: 0.001,
            },
            want: -25,
        },
    }

    for _, test := range testCases {
        t.Run(test.name, func(t *testing.T) {
            got, err := test.ls.LoanPayment()
            if err != nil {
                t.Errorf("error: %v", err)
            }
            if got != test.want {
                t.Errorf("got: %g, wanted: %g", got, test.want)
            }
        })
    }
}


func TestBalloonPayment(t *testing.T){
    var testCases = []struct {
        name string
        ls LoanSizer
        want float64
    }{
        {
            name: "Zero Values",
            ls: LoanSizer{
                MaxLTV: 0.70,
                MinDSCR: 1.40,
                Amortization: 30,
                Term: 10,
                IOPeriod: 0,
                Rate: 0.0045,
                PropertyValue: 0,
                NOI: 0,
                RequestedLoanAmount: 0,
                LoanOriginationFees: 0.001,
            },
            want: 0,
        },
        {
            name: "Valid Inputs, Term and Amortization equal",
            ls: LoanSizer{
                MaxLTV: 0.70,
                MinDSCR: 1.40,
                Amortization: 30,
                Term: 30,
                IOPeriod: 0,
                Rate: 0.0045,
                PropertyValue: 1000,
                NOI: 250,
                RequestedLoanAmount: 900,
                LoanOriginationFees: 0.001,
            },
            want: 0,
        },
        {
            name: "Valid Inputs, Term and Amortization different",
            ls: LoanSizer{
                MaxLTV: 0.70,
                MinDSCR: 1.40,
                Amortization: 30,
                Term: 10,
                IOPeriod: 0,
                Rate: 0.0045,
                PropertyValue: 1000,
                NOI: 250,
                RequestedLoanAmount: 900,
                LoanOriginationFees: 0.001,
            },
            want: 477,
        },
    }

    for _, test := range testCases {
        t.Run(test.name, func(t *testing.T) {
            got, err := test.ls.BalloonPayment()
            if err != nil {
                t.Errorf("error: %v", err)
            }
            if got != test.want {
                t.Errorf("got: %g, wanted: %g", got, test.want)
            }
        })
    }
}
