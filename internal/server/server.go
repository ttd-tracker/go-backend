package server

import (
	"encoding/json"
	"net/http"
)

type BalanceDTO struct {
	Value int
}

type FinanceServer struct{}

func (f *FinanceServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(BalanceDTO{1000})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
