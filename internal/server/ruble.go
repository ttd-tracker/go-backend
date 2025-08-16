package server

import "github.com/shopspring/decimal"

type Ruble struct {
	decimal decimal.Decimal
}

func NewRuble(value float64) Ruble {
	return Ruble{decimal.NewFromFloat(value)}
}

func (r Ruble) Float64() float64 {
	value, _ := r.decimal.Float64()
	return value
}

func (r Ruble) Add(r2 Ruble) Ruble {
	sum, _ := r.decimal.Add(r2.decimal).Float64()
	return NewRuble(sum)
}
