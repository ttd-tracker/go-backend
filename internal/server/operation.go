package server

import "time"

type OpType int

const (
	OpIncome OpType = iota
	OpExpense
)

type Op struct {
	Cash Ruble
	Time time.Time
	Type OpType
}

func NewOp(cash Ruble, opType OpType) Op {
	return Op{cash, time.Now(), opType}
}
