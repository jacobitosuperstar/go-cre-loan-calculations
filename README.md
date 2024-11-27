# go-cre-loan-calculations [![Go Reference](https://pkg.go.dev/badge/github.com/jacobitosuperstar/go-cre-loan-calculations.svg)](https://pkg.go.dev/github.com/jacobitosuperstar/go-cre-loan-calculations)

Opinionated package for loan sizing and loan analysis for commercial real
state.

## Functionalities

### Loan Sizing

Given the loan constrains (maximum loan to value ratio and minimum dscr) and
the loan information (interest rate, interest only period, term, amortization)
one should be able to calculate the following:

* Maximum Loan Amount: How much money the bank can lend me given their
  constrains (LTV and DSCR). If the RequestedLoanAmount is lower than the
  possibly higher loan amount, the RequestedLoanAmount will be returned.

* Loan Payment: Periodic loan payments for the MaximumLoanAmount. If the values
  represent years, the periodic loan payments will represent yearly payments.

* Interes Only Period Payment: Periodic loan payments for the period on which
  only interest will be payed.

* Balloon Payment: Payment that has the done to the bank at the end of the
  term.

* Sale Year Balloon Payment: Payment that has to be done to the bank at the
  time the property sales. The year on which the property will be sold, cannot
  be greater than the term.

* Payment Distribution: How much money goes to the interest and the principal
  during the term, taking into account the interest only period.

## Investment Analysis

Given the Loan constrains, the information of the deal, the tax assumptions and
the projected sale of the property, how does our investment looks like, within
the realm of a simple projection?. Thats the main question that is being solved
within this module.

* Net Cash Flow Projection: How do the payments unfold during the duration of
  the term, and how the money is distributed during the time up to the sale of
  the property.

## TO BE ADDED IN THE FUTURE

* Return of investment metrics, like IRR, Average cash on cash return, Equity
  multiple.

* Objective search: Maximum purchase price of the property given a set of
  return of investment metrics.

* Objective search: Minimum NOI required for the loan to go through, given the
  purchase price of the property.
