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
	AddExpense(id int, expense Ruble) Ruble
}

// Which methods must contain FinanceServer?
type FinanceServer struct {
	http.Handler
	store FinanceStore
}

func NewServer(store FinanceStore) *FinanceServer {
	svr := FinanceServer{store: store}

	mux := http.NewServeMux()
	mux.Handle("GET /balance", NewEnsureAuth(extractBalance, store))
	mux.Handle("POST /op/income/{ruble}", NewEnsureAuth(svr.addIncome, store))
	mux.Handle("POST /op/expense/{ruble}", NewEnsureAuth(svr.addExpense, store))

	svr.Handler = mux
	return &svr
}

func extractBalance(w http.ResponseWriter, r *http.Request, user User) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(BalanceDTO{user.Balance.Float64()})
}

func (f *FinanceServer) addIncome(w http.ResponseWriter, r *http.Request, user User) {
	w.Header().Set("Content-Type", "application/json")

	income, err := strconv.ParseFloat(r.PathValue("ruble"), 64)
	// tests are not providing this scenario
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	balance := f.store.AddIncome(user.Id, NewRuble(income))

	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(BalanceDTO{balance.Float64()})
}

func (f *FinanceServer) addExpense(w http.ResponseWriter, r *http.Request, user User) {
	w.Header().Set("Content-Type", "application/json")

	expense, _ := strconv.ParseFloat(r.PathValue("ruble"), 64)
	balance := f.store.AddExpense(user.Id, NewRuble(expense))
	//if err != nil {
	//	w.WriteHeader(http.StatusBadRequest)
	//	return
	//}

	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(BalanceDTO{balance.Float64()})
}
