package server

import (
	"encoding/json"
	"net/http"
)

type Ruble int

type User struct {
	Balance Ruble
}

type FinanceStore interface {
	GetUser(id string) User
}

type FinanceServer struct {
	http.Handler
	store FinanceStore
}

func NewServer(store FinanceStore) FinanceServer {
	svr := FinanceServer{store: store}
	mux := http.NewServeMux()
	mux.Handle("/balance", NewEnsureAuth(svr.ExtractBalance, store))
	mux.Handle("/op/income", NewEnsureAuth(svr.RecordIncome, store))
	svr.Handler = mux
	return svr
}

func (f *FinanceServer) ExtractBalance(w http.ResponseWriter, r *http.Request, u User) {
	w.Header().Set("Content-Type", "application/json")

	status := http.StatusOK
	err := json.NewEncoder(w).Encode(BalanceDTO{u.Balance})
	if err != nil {
		status = http.StatusInternalServerError
	}
	w.WriteHeader(status)
}

// RecordIncome should write income and return actual balance
func (f *FinanceServer) RecordIncome(w http.ResponseWriter, r *http.Request, u User) {
	w.Header().Set("Content-Type", "application/json")

	w.WriteHeader(http.StatusCreated)
}
