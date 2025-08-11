package server

import "time"

type Ruble int

type Op struct {
	Money Ruble
	Time  time.Time
	Type  opType // to be enum
}

type opType int

const (
	OpIncome opType = iota
)
