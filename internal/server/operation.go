package server

import "time"

type Ruble int

type Op struct {
	Money Ruble
	Time  time.Time
	Type  OpType
}

type OpType int

const (
	OpIncome OpType = iota
)
