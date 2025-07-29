package server

import (
	"encoding/json"
	"io"
)

type BalanceDTO struct {
	Value Ruble
}

func NewBalanceDTO(rdr io.Reader) (BalanceDTO, error) {
	var result BalanceDTO
	err := json.NewDecoder(rdr).Decode(&result)
	return result, err
}
