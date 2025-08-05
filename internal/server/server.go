package server

import (
	"encoding/json"
	"net/http"
	"strconv"
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
	mux.Handle("GET /balance", NewEnsureAuth(svr.ExtractBalance, store))
	mux.Handle("POST /op/income/{ruble}", NewEnsureAuth(svr.RecordIncome, store))
	svr.Handler = mux
	return svr
}

func (f *FinanceServer) ExtractBalance(w http.ResponseWriter, r *http.Request, u User) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(BalanceDTO{u.Balance})
}

// RecordIncome should write income and return actual balance
func (f *FinanceServer) RecordIncome(w http.ResponseWriter, r *http.Request, u User) {
	w.Header().Set("Content-Type", "application/json")

	pathValue := r.PathValue("ruble")
	income, err := strconv.Atoi(pathValue)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(BalanceDTO{Ruble(income) + u.Balance})
}
