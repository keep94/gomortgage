package mortgage

import (
    "github.com/keep94/gomortgage/currency"
    "testing"
)

func TestNeedsOneNilField(t *testing.T) {
  _, e := NewLoan(&LoanSpec{}, true)
  if e != ErrNeedsOneNilField {
    t.Errorf("Expected ErrNeedsOneNilField got %v", e)
  }
}

func TestNeedPositive(t *testing.T) {
  _, e := NewLoan(&LoanSpec{Amount: newUSD(-1.0), Rate: PtrFloat64(0.0), Length: PtrInt(1)}, true)
  if e != ErrNeedPositive {
    t.Errorf("Expected ErrNeedPositive got %v", e)
  }
}

func TestZeroLength(t *testing.T) {
  _, e := NewLoan(&LoanSpec{Amount: newUSD(1.0), Rate: PtrFloat64(0.0), Length: PtrInt(0)}, true)
  if e != ErrZeroLength {
    t.Errorf("Expected ErrZeroLength got %v", e)
  }
}

func TestRegular(t *testing.T) {
  l, e := NewLoan(&LoanSpec{Amount: newUSD(238000.0), Rate: PtrFloat64(.04 / 12.0), Length: PtrInt(360)}, true)
  if e != nil {
    t.Errorf("Expected no error but got %v", e)
  }
  if amount := l.Payment.String(); amount != "1136.25" {
    t.Errorf("Expected payment 1136.25 but got %v", amount)
  }
  verifyTerms(t, l.Terms, newUSD(238000.0))
}

func TestBigInterest(t *testing.T) {
  l, e := NewLoan(&LoanSpec{Amount: newUSD(238000.0), Rate: PtrFloat64(.05), Length: PtrInt(360)}, true)
  if e != nil {
    t.Errorf("Expected no error but got %v", e)
  }
  if amount := l.Payment.String(); amount != "11900.01" {
    t.Errorf("Expected payment 11900.01 but got %v", amount)
  }
  verifyTerms(t, l.Terms, newUSD(238000.0))
}

func TestNegativeInterest(t *testing.T) {
  l, e := NewLoan(&LoanSpec{Amount: newUSD(238000.0), Rate: PtrFloat64(-.05), Length: PtrInt(360)}, true)
  if e != nil {
    t.Errorf("Expected no error but got %v", e)
  }
  if amount := l.Payment.String(); amount != "0.01" {
    t.Errorf("Expected payment 0.01 but got %v", amount)
  }
  verifyTerms(t, l.Terms, newUSD(238000.0))
}

func verifyTerms(t *testing.T, terms []*Term, balance currency.Amount) {
  prevBalance := balance.Int64()
  for _, term := range terms {
    verifyTerm(t, term)
    balance := term.Balance().Int64()
    if prevBalance - term.Principal().Int64() != balance {
      t.Error("Balance wrong")
    }
    prevBalance = balance
  }
  if prevBalance != 0 {
    t.Error("Final balance not zero.")
  }
}
  
func verifyTerm(t *testing.T, term *Term) {
  if term.Payment().Int64() != term.Principal().Int64() + term.Interest().Int64() {
    t.Error("Term does not add up.")
  }
}

func newUSD(x float64) currency.Amount {
  return currency.NewAmount("USD").FromFloat64(x)
}
