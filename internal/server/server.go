package server

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
)

type User struct {
	Id      int
	Balance Ruble
	History []Op
}

type FinanceStore interface {
	GetUser(id int) *User
	AddIncome(id int, income Ruble) Ruble
	AddExpense(id int, expense Ruble) Ruble
}

type FinanceServer struct {
	http.Handler
	store FinanceStore
}

func NewServer(store FinanceStore) *FinanceServer {
	svr := FinanceServer{store: store}

	mux := http.NewServeMux()
	mux.Handle("GET /balance", NewEnsureAuth(extractBalance, store))
	mux.Handle("POST /op/income/{cash}", NewEnsureAuth(svr.addIncome, store))
	mux.Handle("POST /op/expense/{cash}", NewEnsureAuth(svr.addExpense, store))

	svr.Handler = mux
	return &svr
}

func extractBalance(w http.ResponseWriter, r *http.Request, user *User) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(BalanceDTO{user.Balance.Float64()})
}

func (f *FinanceServer) addIncome(w http.ResponseWriter, r *http.Request, user *User) {
	w.Header().Set("Content-Type", "application/json")

	income, err := getCashPathParameter(r)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	balance := f.store.AddIncome(user.Id, income)
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(BalanceDTO{balance.Float64()})
}

// addExpense subtracts given amount from store balance, records new Op and returns BalanceDTO with user updated balance
func (f *FinanceServer) addExpense(w http.ResponseWriter, r *http.Request, user *User) {
	w.Header().Set("Content-Type", "application/json")

	expense, err := getCashPathParameter(r)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	balance := f.store.AddExpense(user.Id, expense)
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(BalanceDTO{balance.Float64()})
}

var ErrZeroCash = errors.New("cash is zero")

func getCashPathParameter(r *http.Request) (Ruble, error) {
	cash, err := strconv.ParseFloat(r.PathValue("cash"), 64)
	if err != nil {
		return NewRuble(0), err
	}
	if cash == 0 {
		return NewRuble(0), ErrZeroCash
	}

	return NewRuble(cash), nil
}
