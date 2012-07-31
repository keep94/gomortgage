package currency

import (
    "fmt"
    "testing"
)

func Add(a Amount, b... Amount) Amount {
  sum := a.Int64()
  for _, x := range b {
    sum += x.Int64()
  }
  return a.FromInt64(sum)
}

func PlusTax(a Amount) Amount {
  return a.FromFloat64(a.Float64() * 1.0825)
}

func TestUSD(t *testing.T) {
  u := NewAmount("USD")
  sumString := Add(u.FromFloat64(-3.006), u.FromFloat64(-4.006)).String()
  if sumString != "-7.02" {
    t.Errorf("Wanted -7.02 got %v", sumString)
  }
  sumString = Add(u.FromFloat64(3.006), u.FromFloat64(4.006)).String()
  if sumString != "7.02" {
    t.Errorf("Wanted 7.02 got %v", sumString)
  }
  sumString = Add(u.FromFloat64(-3.004), u.FromFloat64(-4.004)).String()
  if sumString != "-7.00" {
    t.Errorf("Wanted -7.00 got %v", sumString)
  }
  sumString = Add(u.FromFloat64(3.004), u.FromFloat64(4.004)).String()
  if sumString != "7.00" {
    t.Errorf("Wanted 7.00 got %v", sumString)
  }
}

func TestUSDFloat(t *testing.T) {
  u := NewAmount("USD")
  sumString := PlusTax(u.FromFloat64(-3.0)).String()
  if sumString != "-3.25" {
    t.Errorf("Wanted -3.25 got %v", sumString)
  }
  sumString = PlusTax(u.FromFloat64(3.0)).String()
  if sumString != "3.25" {
    t.Errorf("Wanted 3.25 got %v", sumString)
  }
}

func TestJPY(t *testing.T) {
  j := NewAmount("JPY")
  sumString := Add(j.FromFloat64(-300.6), j.FromFloat64(-400.6)).String()
  if sumString != "-702" {
    t.Errorf("Wanted -702 got %v", sumString)
  }
  sumString = Add(j.FromFloat64(300.6), j.FromFloat64(400.6)).String()
  if sumString != "702" {
    t.Errorf("Wanted 702 got %v", sumString)
  }
  sumString = Add(j.FromFloat64(-300.4), j.FromFloat64(-400.4)).String()
  if sumString != "-700" {
    t.Errorf("Wanted -700 got %v", sumString)
  }
  sumString = Add(j.FromFloat64(300.4), j.FromFloat64(400.4)).String()
  if sumString != "700" {
    t.Errorf("Wanted 700 got %v", sumString)
  }
}

func TestJPYFloat(t *testing.T) {
  j := NewAmount("JPY")
  sumString := PlusTax(j.FromFloat64(-300.0)).String()
  if sumString != "-325" {
    t.Errorf("Wanted -325 got %v", sumString)
  }
  sumString = PlusTax(j.FromFloat64(300.0)).String()
  if sumString != "325" {
    t.Errorf("Wanted 325 got %v", sumString)
  }
}

func TestCurrency(t *testing.T) {
  c := NewAmount("USD").Currency()
  if c != "USD" {
    t.Errorf("Wanted USD got %v", c)
  }
  c = NewAmount("JPY").Currency()
  if c != "JPY" {
    t.Errorf("Wanted JPY got %v", c)
  }
  c = NewAmount("EUR").Currency()
  if c != "EUR" {
    t.Errorf("Wanted EUR got %v", c)
  }
}

func TestUnknownCurrency(t *testing.T) {
  defer func() {
    recover()
  }()
  NewAmount("ABC")
  t.Error("Failed to panic for currency ABC")
}

func ExampleAmount() {
  u := NewAmount("USD")
  a := u.FromFloat64(3.53)
  b := u.FromFloat64(6.48)
  sum := u.FromInt64(a.Int64() + b.Int64())
  fmt.Println(sum)
  // Output: 10.01
}
