// Package mortgage calculates Amortization schedules
package mortgage

import (
    "errors"
    "github.com/keep94/gomortgage/currency"
    "math"
)

var (
  ErrNeedsOneNilField = errors.New("LoanSpec must have exactly one nil field.")
  ErrNeedPositive = errors.New("Amount and Payment must be positive when specified.")
  ErrNotSupported = errors.New("This loan cannot be computed.")
  ErrZeroLength = errors.New("Length must be positive.")
)

func PtrFloat64(x float64) *float64 {
  return &x
}

func PtrInt(x int) *int {
  return &x
}

type LoanSpec struct {
  Amount currency.Amount
  Rate *float64
  Length *int
  Payment currency.Amount
}

func (s *LoanSpec) checkFields() error {
  var foundNil bool
  if s.Amount != nil {
    if s.Amount.Int64() <= 0 {
      return ErrNeedPositive
    }
  } else {
    foundNil = true
  }
  if s.Length != nil {
    if *s.Length <= 0 {
      return ErrZeroLength
    }
  } else if foundNil {
    return ErrNeedsOneNilField
  } else {
   foundNil = true
  }
  if s.Payment != nil {
    if s.Payment.Int64() <= 0 {
      return ErrNeedPositive
    }
  } else if foundNil {
    return ErrNeedsOneNilField
  } else {
    foundNil = true
  }
  if !foundNil {
    return ErrNeedsOneNilField
  }
  return nil
}

// Loan represents a loan.
type Loan struct {
  Amount currency.Amount
  Rate float64
  Length int
  Payment currency.Amount
  Terms []*Term
}

func NewLoan(spec *LoanSpec, computeTerms bool) (*Loan, error) {
  if err := spec.checkFields(); err != nil {
    return nil, err
  }
  if spec.Payment == nil {
    return fromSpecNeedingPayment(spec, computeTerms), nil
  }
  return nil, ErrNotSupported
}

func fromSpecNeedingPayment(spec *LoanSpec, computeTerms bool) *Loan {
  payment := solveForPayment(spec.Amount, *spec.Rate, *spec.Length)
  result := &Loan{Amount: spec.Amount, Rate: *spec.Rate, Length: *spec.Length, Payment: payment}
  if computeTerms {
    result.computeTerms()
  }
  return result
}

func solveForPayment(
    amount currency.Amount, rate float64, length int) currency.Amount {
  amountF := amount.Float64()
  if rate == 0.0 {
    return amount.FromFloat64(amountF / float64(length))
  }
  result := amount.FromFloat64(amountF * rate * (1.0 + 1.0 / (
      math.Pow((1.0 + rate), float64(length)) - 1.0)))

  resultI := result.Int64()
  if resultI <= 0 {
    resultI = 1
  }
  interestOnly := amount.FromFloat64(amountF * rate).Int64()
  if resultI <= interestOnly {
    resultI = interestOnly + 1
  }
  return result.FromInt64(resultI)
}

func (l *Loan) computeTerms() {
  balance := l.Amount
  for balance.Int64() > 0 {
    t := new(Term)
    t.interest = balance.FromFloat64(balance.Float64() * l.Rate)
    balanceI := balance.Int64() + t.interest.Int64()
    if l.Payment.Int64() > balanceI {
      t.payment = balance.FromInt64(balanceI)
    } else {
      t.payment = l.Payment
    }
    t.balance = balance.FromInt64(balanceI - t.payment.Int64())
    l.Terms = append(l.Terms, t)
    balance = t.balance
  }
}

type Term struct {
  payment currency.Amount
  interest currency.Amount
  balance currency.Amount
}

// Balance returns the balance left on the loan.
func (t *Term) Balance() currency.Amount {
  return t.balance
}

// Interest returns the amount going toward interest.
func (t *Term) Interest() currency.Amount {
  return t.interest
}

// Payment returns the payment due for this term.
func (t *Term) Payment() currency.Amount {
  return t.payment
}

// Principal returns the amount going toward principal.
func (t *Term) Principal() currency.Amount {
  return t.payment.FromInt64(t.payment.Int64() - t.interest.Int64())
}
