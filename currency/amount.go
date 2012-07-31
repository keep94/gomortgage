// Package currency contains utilites for handling monetary amounts.
package currency

import (
    "math"
    "strconv"
)

// Amount represents an amount in a particular currency.
type Amount interface {
  // Amount as a float
  Float64() float64
  // Amount as an int. Use to add or subtract without round off error.
  Int64() int64
  // Returns a new Amount for x in the same currency.
  FromFloat64(x float64) Amount
  // Returns a new amount for x in the same currency. 
  FromInt64(x int64) Amount
  // Returns the currency e.g "USD", "EUR"
  Currency() string
  // Returns string representation of amount without currency suffix.
  String() string
}

// NewAmount returns 0 in given currency. New panics if currency is unknown.
// Right now only USD, EUR, and JPY are supported.
func NewAmount(currency string) Amount {
  switch currency {
    case usdValue.Currency():
      return usdValue
    case jpyValue.Currency():
      return jpyValue
    case eurValue.Currency():
      return eurValue
  }
  panic("Unknown currency " + currency)
}

type prec2 int64

func (p prec2) Int64() int64 {
  return int64(p)
}

func (p prec2) Float64() float64 {
  return float64(p) / 100.0
}

func (p prec2) String() string {
  return strconv.FormatFloat(p.Float64(), 'f', 2, 64)
}

func toPrec2(x float64) prec2 {
  return prec2(math.Floor(x * 100.0 + 0.5))
}

type prec0 int64

func (p prec0) Int64() int64 {
  return int64(p)
}

func (p prec0) Float64() float64 {
  return float64(p)
}

func (p prec0) String() string {
  return strconv.FormatFloat(p.Float64(), 'f', 0, 64)
}

func toPrec0(x float64) prec0 {
  return prec0(math.Floor(x + 0.5))
}

type usd struct {
  prec2
}

func (u usd) FromFloat64(x float64) Amount {
  return usd{toPrec2(x)}
} 

func (u usd) FromInt64(x int64) Amount {
  return usd{prec2(x)}
}

func (u usd) Currency() string {
  return "USD"
} 

type jpy struct {
  prec0
}

func (j jpy) FromFloat64(x float64) Amount {
  return jpy{toPrec0(x)}
} 

func (j jpy) FromInt64(x int64) Amount {
  return jpy{prec0(x)}
}

func (j jpy) Currency() string {
  return "JPY"
} 

type eur struct {
  prec2
}

func (e eur) FromFloat64(x float64) Amount {
  return eur{toPrec2(x)}
} 

func (e eur) FromInt64(x int64) Amount {
  return eur{prec2(x)}
}

func (e eur) Currency() string {
  return "EUR"
} 

var usdValue usd
var jpyValue jpy
var eurValue eur
