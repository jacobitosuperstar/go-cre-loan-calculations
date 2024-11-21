// Opinionated Financial Formulas. Every money value rounded down to 2 decimal
// places.
// TODO: Formulas needed
// [X] IOPayment (monthly)
// [X] Payment (monthly)
// [X] PrincipalPayment (monthly)
// [X] InterestPayment (monthly)
// [X] PresentValue
// If everything is already in years, this is not needed
// [X] YearlyIOPayment
// [X] YearlyPayment

package financial_formulas

import (
    "log";
    "fmt";
	"math";
    utils "github.com/jacobitosuperstar/go-cre-loan-calculations/internal/utils";
)

const (
    PayEnd = 0
    PayBegin = 1
)


// IOPayment returns the interest only payment for a cash flow with a constant
// interest rate.
func IOPayment(
    rate float64,
    pv float64,
) (
    pmt float64,
) {
    pmt = - pv * rate
    return utils.Round2(pmt)
}


// Payment returns the constant payment for a cash flow with a constant
// interest rate.
func Payment(
    rate float64,
    numPeriods int,
    pv float64,
    fv float64,
    paymentType int,
) (
    pmt float64,
    err error,
) {
    if numPeriods <= 0 {
        return 0.0, &ValidationError{"numPeriods", numPeriods, "The value must be greater than 0"}
	}
	if paymentType != PayEnd && paymentType != PayBegin {
		return 0.0, &ValidationError{"paymentType", paymentType, "The value must be 0 (PayEnd) or 1 (PayBegin)"}
	}
	if rate != 0 {
		pmt = ((-fv - pv*math.Pow(1+rate, float64(numPeriods))) * rate) /
        ((1 + rate*float64(paymentType)) * (math.Pow(1+rate, float64(numPeriods)) - 1))
	} else {
		pmt = (-pv - fv) / float64(numPeriods)
	}
	return utils.Round2(pmt), nil
}


// InterestAndPrincipalPayment returns an slice of interest payments, an array
// of principal payments and an error.
func interest_and_principal_payments(
    rate float64,
    numPeriods int,
    pv float64,
    fv float64,
    paymentType int,
) (
    ipmt []float64,
    ppmt []float64,
    err error,
) {
    pmt, err := Payment(rate, numPeriods, pv, fv, paymentType)
    if err != nil {
        return ipmt, ppmt, fmt.Errorf("interest_and_principal_payments internal error: %v", err)
    }

    capital := pv

    // Because of the commulative rounding error, we will always have
    // unmatching values at the end of the `numPeriods` between the principal
    // payment and the capital. To solve this, the last principal, and interest
    // payment, will be caltulated directly from the projected remaining
    // capital.
    var interest_payment, principal_payment float64
    for i := 1; i < numPeriods; i++ {
        if paymentType == PayBegin && i == 1 {
            interest_payment = 0.0
        } else {
            interest_payment = utils.Round2(- capital * rate)
        }
        ipmt = append(ipmt, interest_payment)
        principal_payment = utils.Round2(pmt - interest_payment)
        ppmt = append(ppmt, principal_payment)
        capital = capital + principal_payment
        capital = utils.Round2(capital)
    }

    // TODO: Look out for the cases where the `fv` is different from 0, there
    // could be more inconsistencies hidden.

    interest_payment = utils.Round2(- capital * rate)
    ipmt = append(ipmt, interest_payment)
    principal_payment = utils.Round2(- capital + fv)
    ppmt = append(ppmt, principal_payment)
    capital = utils.Round2(capital + principal_payment)

    if capital != fv{
        log.Printf("capital: %v  fv: %v", capital, fv)
        return ipmt, ppmt, &ValidationError{"pv", capital, "pv doesnt' match fv at the end of the numPeriods"}
    }
    return ipmt, ppmt, nil
}

// PrincipalPayments return an array and an error of all the principal payments
// during the number of periods
func PrincipalPayments(
    rate float64,
    numPeriods int,
    pv float64,
    fv float64,
    paymentType int,
) (
    []float64,
    error,
) {
    _, ppmt, err := interest_and_principal_payments(rate, numPeriods, pv, fv, paymentType)

    if err != nil {
        return ppmt, fmt.Errorf("interest_and_principal_payments internal error: %v", err)
    }
    return ppmt, nil
}

// InterestPayments returns an array and an error of all the interest payments
// during the number of periods
func InterestPayments(
    rate float64,
    numPeriods int,
    pv float64,
    fv float64,
    paymentType int,
) (
    []float64,
    error,
) {
    ipmt, _, err := interest_and_principal_payments(rate, numPeriods, pv, fv, paymentType)

    if err != nil {
        return ipmt, fmt.Errorf("interest_and_principal_payments internal error: %v", err)
    }
    return ipmt, nil
}

// PresentValue return the present value of a cashflow with constant interest
// rate and payments
func PresentValue(
    rate float64,
    numPeriods int,
    pmt float64,
    fv float64,
    paymentType int,
) (
    pv float64,
    err error,
) {
    if numPeriods <= 0 {
        return 0.0, &ValidationError{"numPeriods", numPeriods, "The value must be greater than 0"}
    }
    if paymentType != PayEnd && paymentType != PayBegin {
		return 0.0, &ValidationError{"paymentType", paymentType, "The value must be 0 (PayEnd) or 1 (PayBegin)"}
    }
    if rate != 0 {
        pv = (-pmt*(1+rate*float64(paymentType))*((math.Pow(1+rate, float64(numPeriods))-1)/rate) - fv) / math.Pow(1+rate, float64(numPeriods))
    } else {
        pv = -fv - pmt*float64(numPeriods)
    }
    return utils.Round2(pv), nil
}
