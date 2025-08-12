package server

import "time"

type Op struct {
	Ruble
	Time time.Time
	Type OpType
}

type OpType int

const (
	OpIncome OpType = iota
)
