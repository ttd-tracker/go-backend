package server

import (
	"encoding/json"
	"net/http"
	"strconv"
)

type User struct {
	Id      int
	Balance Ruble
}

type FinanceStore interface {
	GetUser(id int) User
	AddIncome(id int, income Ruble) Ruble
}

type FinanceServer struct {
	http.Handler
	store FinanceStore
}

func NewServer(store FinanceStore) *FinanceServer {
	svr := FinanceServer{store: store}

	mux := http.NewServeMux()
	mux.Handle("GET /balance", NewEnsureAuth(svr.ExtractBalance, store))
	mux.Handle("POST /op/income/{ruble}", NewEnsureAuth(svr.AddIncome, store))

	svr.Handler = mux
	return &svr
}

func (f *FinanceServer) ExtractBalance(w http.ResponseWriter, r *http.Request, user User) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(BalanceDTO{user.Balance.Value()})
}

func (f *FinanceServer) AddIncome(w http.ResponseWriter, r *http.Request, user User) {
	w.Header().Set("Content-Type", "application/json")

	income, err := strconv.ParseFloat(r.PathValue("ruble"), 64)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	balance := f.store.AddIncome(user.Id, NewRuble(income))

	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(BalanceDTO{balance.Value()})
}
