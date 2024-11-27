// Testing of the Financial Formulas
// TODO: Tests needed
// [X] IOPayment (monthly)
// [X] YearlyIOPayment
// [X] Payment (monthly)
// [X] YearlyPayment
// [X] PrincipalPayment (monthly)
// [X] InterestPayment (monthly)
// [X] PresentValue

package financial_formulas
import (
    "testing";
    utils "github.com/jacobitosuperstar/go-cre-loan-calculations/internal/utils";
)


// tolerance
const TOL = 0.1

func TestIOPayment(t *testing.T){
    var testCases = []struct {
        rate float64
        pv float64
        want float64
    }{
        {0.00, 0.00, 0.00},
        {100, 0.00375, - 0.38},
        {200, 0.00375, - 0.75},
    }

    for _, test := range testCases {
        t.Run("Testing IO Payments", func(t *testing.T) {
            if got := IOPayment(test.rate, test.pv); !utils.Tolerance(got, test.want, TOL) {
                t.Errorf("got: %g, wanted: %g", got, test.want)
            }
        })
    }
}


func TestPayment(t *testing.T){
    var testCases = []struct {
        name string
        rate float64
        numPeriods int
        pv float64
        fv float64
        paymentType int
        want float64
    }{
        {
            name: "Everything Zero",
            rate: 0,
            numPeriods: 0,
            pv: 0,
            fv: 0,
            paymentType: 0,
            want: 0,
        },
        {
            name: "Valid numPeriods",
            rate: 0,
            numPeriods: 12,
            pv: 0,
            fv: 0,
            paymentType: 0,
            want: 0,
        },
        {
            name: "Invalid paymentType",
            rate: 0.00375,
            numPeriods: 12,
            pv: 100,
            fv: 0,
            paymentType: 5,
            want: 0,
        },
        {
            name: "Totally valid case",
            rate: 0.00375,
            numPeriods: 12,
            pv: 100,
            fv: 0,
            paymentType: 0,
            want: -8.54,
        },
        {
            name: "Totally valid case, fv != 0",
            rate: 0.00375,
            numPeriods: 1,
            pv: 100,
            fv: -100,
            paymentType: 0,
            want: -0.38,
        },
    }

    for _, test := range testCases {
        t.Run(test.name, func(t *testing.T) {
            got, _ := Payment(test.rate, test.numPeriods, test.pv, test.fv, test.paymentType)
            if !utils.Tolerance(got, test.want, TOL) {
                t.Errorf("got: %g, wanted: %g", got, test.want)
            }
        })
    }
}


func TestPrincipalPayments(t *testing.T){
    var testCases = []struct {
        name string
        rate float64
        numPeriods int
        pv float64
        fv float64
        paymentType int
        want []float64
    }{
        {
            name: "Everything Zero",
            rate: 0,
            numPeriods: 1,
            pv: 0,
            fv: 0,
            paymentType: 0,
            want: []float64{0},
        },
        {
            name: "Totally valid case",
            rate: 0.0375,
            numPeriods: 2,
            pv: 100,
            fv: 0,
            paymentType: 0,
            want: []float64{-49.08, -50.92},
        },
        {
            name: "Totally valid case, fv != 0",
            rate: 0.0375,
            numPeriods: 2,
            pv: 100,
            fv: -20,
            paymentType: 0,
            want: []float64{-39.26, -80.74},
        },
    }

    for _, test := range testCases {
        t.Run(test.name, func(t *testing.T) {
            got, err := PrincipalPayments(test.rate, test.numPeriods, test.pv, test.fv, test.paymentType)

            if err != nil {
                t.Errorf("PrincipalPayments internal error: %v", err)
                return
            }

            if len(got) != len(test.want) {
                t.Errorf("got: %g, wanted: %g", got, test.want)
                return
            } else {
                for i := range got {
                    if !utils.Tolerance(got[i], test.want[i], TOL) {
                        t.Errorf("got: %g, wanted: %g", got[i], test.want[i])
                    }
                }
            }
        })
    }
}


func TestInterestPayments(t *testing.T){
    var testCases = []struct {
        name string
        rate float64
        numPeriods int
        pv float64
        fv float64
        paymentType int
        want []float64
    }{
        {
            name: "Everything Zero, Valid numPeriods",
            rate: 0,
            numPeriods: 1,
            pv: 0,
            fv: 0,
            paymentType: 0,
            want: []float64{0},
        },
        {
            name: "Totally valid case",
            rate: 0.00375,
            numPeriods: 2,
            pv: 100,
            fv: 0,
            paymentType: 0,
            want: []float64{-0.38, -0.19},
        },
        {
            name: "Totally valid case, fv != 0",
            rate: 0.00375,
            numPeriods: 2,
            pv: 100,
            fv: 20,
            paymentType: 0,
            want: []float64{-0.38, -0.15},
        },
    }

    for _, test := range testCases {
        t.Run(test.name, func(t *testing.T) {
            got, err := InterestPayments(test.rate, test.numPeriods, test.pv, test.fv, test.paymentType)

            if err != nil {
                t.Errorf("PrincipalPayments internal error: %v", err)
                return
            }

            if len(got) != len(test.want) {
                t.Errorf("got: %g, wanted: %g", got, test.want)
                return
            } else {
                for i := range got {
                    if !utils.Tolerance(got[i], test.want[i], TOL) {
                        t.Errorf("got: %g, wanted: %g", got[i], test.want[i])
                    }
                }
            }
        })
    }
}


func TestPresentValue(t *testing.T){
    var testCases = []struct {
        name string
        rate float64
        numPeriods int
        pmt float64
        fv float64
        paymentType int
        want float64
    }{
        {
            name: "Everything Zero, Valid numPeriods",
            rate: 0,
            numPeriods: 1,
            pmt: 0,
            fv: 0,
            paymentType: 0,
            want: 0,
        },
        {
            name: "Totally valid case",
            rate: 0.08,
            numPeriods: 20,
            pmt: 500,
            fv: 0,
            paymentType: 0,
            want: -4909.07,
        },
    }

    for _, test := range testCases {
        t.Run(test.name, func(t *testing.T) {
            got, err := PresentValue(test.rate, test.numPeriods, test.pmt, test.fv, test.paymentType)

            if err != nil {
                t.Errorf("PresentValue internal error: %v", err)
                return
            }

            if !utils.Tolerance(got, test.want, TOL) {
                t.Errorf("got: %g, wanted: %g", got, test.want)
            }
        })
    }
}
