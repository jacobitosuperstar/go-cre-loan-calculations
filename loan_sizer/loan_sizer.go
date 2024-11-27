// We will create a LoanSizer struct and we will be adding things to the struct
// as is needed to show the different things needed or wanted from the loan
// sizing part.
// TODO: Things that should be returned to the user
// [X] maximum loan amount
// [X] requested loan amount
// [X] loan payment
// [X] loan io period payment
// [X] balloon payment
// [X] payments distribution

package loan_sizer

import (
    "fmt"
    "math"
    "sort"
    ff "github.com/jacobitosuperstar/go-cre-loan-calculations/internal/financial_formulas";
    utils "github.com/jacobitosuperstar/go-cre-loan-calculations/internal/utils";
)

// LoanSizer creates a struct that has all the information regarding the loan
// information.
type LoanSizer struct {
    MaxLTV              float64
    MinDSCR             float64
    Amortization        int
    Term                int
    IOPeriod            int
    Rate                float64
    PropertyValue       int
    NOI                 float64
    RequestedLoanAmount int
    LoanOriginationFees float64
}

// Constructor

// NewLoanSizer returns a LoanSizer struct if the values given are valid. If
// not, returns a default struct with the error.
func NewLoanSizer(
    maxLTV float64,
    minDSCR float64,
    amortization int,
    term int,
    ioPeriod int,
    interestRate float64,
    propertyValue int,
    noi float64,
    requestedLoanAmount int,
    loanOriginationFees float64,
) (
    LoanSizer,
    error,
) {
    // Data Validation
    if maxLTV < 0 || maxLTV > 1 {
        return LoanSizer{}, fmt.Errorf("The loan to value ratio must be between 0 and 1.")
    }
    if minDSCR < 1 {
        return LoanSizer{}, fmt.Errorf("The minDSCR cannot be lower than 1.")
    }
    if term > amortization {
        return LoanSizer{}, fmt.Errorf("The term of the loan cannot be greater than its amortization.")
    }
    if ioPeriod < 0 {
        return LoanSizer{}, fmt.Errorf("The ioPeriod cannot be than zero.")
    }
    if ioPeriod > term {
        return LoanSizer{}, fmt.Errorf("The ioPeriod cannot be greater than the term of the loan.")
    }
    if interestRate < 0 || interestRate > 1 {
        return LoanSizer{}, fmt.Errorf("The interest rate of a loan must be between 0 and 1.")
    }
    if propertyValue < 0 {
        return LoanSizer{}, fmt.Errorf("The property cannot have a value below zero.")
    }
    if requestedLoanAmount < 0 {
        return LoanSizer{}, fmt.Errorf("The requestedLoanAmount cannot have a value below zero.")
    }
    if requestedLoanAmount == 0 {
        requestedLoanAmount = propertyValue
    }
    // Struct Creation
    ls := LoanSizer{
        MaxLTV: maxLTV,
        MinDSCR: minDSCR,
        Amortization: amortization,
        Term: term,
        IOPeriod: ioPeriod,
        Rate: interestRate,
        PropertyValue: propertyValue,
        NOI: noi,
        RequestedLoanAmount: requestedLoanAmount,
        LoanOriginationFees: loanOriginationFees,
    }
    return ls, nil
}

// Calculation methods

// Internal

// max_ltv_loan_amount returns the maximum loan amount given the maximum loan
// to value ratio
func (ls LoanSizer) max_ltv_loan_amount () float64 {
    ltv_mla  := math.Floor(ls.MaxLTV * float64(ls.PropertyValue))
    return ltv_mla
}

// max_mindscr_loan_amount returns the maximum loan amount given the minimum
// dscr
func (ls LoanSizer) max_mindscr_loan_amount () (float64, error) {
    // monthly_rate := ls.Rate / 12
    // amoritzation_months := ls.Amortization * 12
    payment := - ls.NOI / ls.MinDSCR
    dscr_mla, err := ff.PresentValue(ls.Rate, ls.Amortization, payment, 0, 0)
    if err != nil {
        return 0.0, fmt.Errorf("PresentValue internal error: %v", err)
    }
    return math.Floor(dscr_mla), err
}

// External

// MaximumLoanAmount returns the maximum loan amount of a LoanSizer struct,
// with the ltv and dscr restrictions. If the RequestedLoanAmount is lower than
// the possibly higher loan amount, the RequestedLoanAmount will be returned.
func (ls *LoanSizer) MaximumLoanAmount () (float64, error) {
    max_mindscr_loan_amount, err := ls.max_mindscr_loan_amount()
    if err != nil {
        return 0, fmt.Errorf("max_mindscr_loan_amount internal error: %v", err)
    }
    loan_values := [3]float64{
        ls.max_ltv_loan_amount(),
        max_mindscr_loan_amount,
        float64(ls.RequestedLoanAmount),
    }
    sort.Float64s(loan_values[:])
    return loan_values[0], nil
}

// IOLoanPayment returns the loan payments during the IO periods for the
// maximum loan amount.
func (ls LoanSizer) IOLoanPayment () (float64, error){
    mla, err := ls.MaximumLoanAmount()
    if err != nil {
        return 0.0, fmt.Errorf("MaximumLoanAmount internal error: %v", err)
    }
    return ff.IOPayment(ls.Rate, mla), nil
}

// LoanPayment returns the periodic loan payments for the maximum loan amount.
func (ls LoanSizer) LoanPayment () (float64, error) {
    mla, err := ls.MaximumLoanAmount()
    if err != nil {
        return 0.0, fmt.Errorf("MaximumLoanAmount internal error: %v", err)
    }
    loan_payment, err := ff.Payment(ls.Rate, ls.Amortization, mla, 0, 0)
    if err != nil {
        return 0.0, fmt.Errorf("Payment internal error: %v", err)
    }
    return loan_payment, nil
}

// BallonPayment returns the balloon payment at the end of the term for the
// maximum loan amount.
func (ls *LoanSizer) EndofTermBalloonPayment () (float64, error) {
    mla, err := ls.MaximumLoanAmount()
    if err != nil {
        return 0.0, fmt.Errorf("MaximumLoanAmount internal error: %v", err)
    }
    principal_payments, err := ff.PrincipalPayments(ls.Rate, ls.Amortization, mla, 0, 0)
    if err != nil {
        return 0.0, fmt.Errorf("PrincipalPayments internal error: %v", err)
    }

    // here we create a 0s array and then append to it the principal payments
    // array that will represent the no principal payment while the IO period.
    if ls.IOPeriod > 0 {
        io_period_ppmt := make([]float64, ls.IOPeriod)
        principal_payments = append(io_period_ppmt, principal_payments...)
    }

    capital := mla
    // There is no need to create a new slice that will contain only the term
    // as the iteration will iterate till the term value. Genious move!!.
    for i:=0; i < ls.Term; i++ {
        capital += principal_payments[i]
    }
    return utils.Round2(capital), nil
}

// SaleYearBallonPayment returns the balloon payment at year that the property
// is being sold.
func (ls *LoanSizer) SaleYearBalloonPayment (saleYear int) (float64, error) {
    mla, err := ls.MaximumLoanAmount()
    if err != nil {
        return 0.0, fmt.Errorf("MaximumLoanAmount internal error: %v", err)
    }
    principal_payments, err := ff.PrincipalPayments(ls.Rate, ls.Amortization, mla, 0, 0)
    if err != nil {
        return 0.0, fmt.Errorf("PrincipalPayments internal error: %v", err)
    }

    // here we create a 0s array and then append to it the principal payments
    // array that will represent the no principal payment while the IO period.
    if ls.IOPeriod > 0 {
        io_period_ppmt := make([]float64, ls.IOPeriod)
        principal_payments = append(io_period_ppmt, principal_payments...)
    }

    capital := mla
    // There is no need to create a new slice that will contain only the term
    // as the iteration will iterate till the term value. Genious move!!.
    for i:=0; i < saleYear; i++ {
        capital += principal_payments[i]
    }
    return utils.Round2(capital), nil
}

// PaymentDistribution returns the slices of the different interest and
// principal payments of the loan within the duration of the term for the
// maximum loan amount.
func (ls LoanSizer) PaymentDistribution () (
    ppmt []float64,
    ipmt []float64,
    err error,
) {
    mla, err := ls.MaximumLoanAmount()
    if err != nil {
        return ppmt, ipmt, fmt.Errorf("MaximumLoanAmount internal error: %v", err)
    }
    // Principal Payments
    ppmt, err = ff.PrincipalPayments(ls.Rate, ls.Amortization, mla, 0, 0)
    if err != nil {
        return ppmt, ipmt, fmt.Errorf("PrincipalPayments internal error: %v", err)
    }
    // adding the IO period payments at the begining of the slice.
    if ls.IOPeriod > 0 {
        io_period_ppmt := make([]float64, ls.IOPeriod)
        ppmt = append(io_period_ppmt, ppmt...)
    }
    // taking the slice with the size of the term.
    ppmt = ppmt[:ls.Term]

    // Interest Payments
    ipmt, err = ff.InterestPayments(ls.Rate, ls.Amortization, mla, 0, 0)
    if err != nil {
        return ppmt, ipmt, fmt.Errorf("InterestPayments internal error: %v", err)
    }
    // adding the IO period payments at the begining of the slice.
    if ls.IOPeriod > 0 {
        io_pmt, err := ls.IOLoanPayment()
        if err != nil {
            return ppmt, ipmt, fmt.Errorf("IOLoanPayment internal error: %v", err)
        }
        io_period_ipmt := make([]float64, ls.IOPeriod)
        for i := 0; i < ls.IOPeriod; i++ {
            io_period_ipmt[i] = io_pmt
        }
        ipmt = append(io_period_ipmt, ipmt...)
    }
    // taking the slice with the size of the term.
    ipmt = ipmt[:ls.Term]
    // Returning everything
    return ppmt, ipmt, nil
}
