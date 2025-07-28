package server

import (
	"encoding/json"
	"net/http"
)

type Ruble int

type User struct {
	Balance Ruble
}

func NewServer() *EnsureAuth {
	return &EnsureAuth{ServeFinances}
}

func ServeFinances(w http.ResponseWriter, r *http.Request, u *User) {
	w.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(BalanceDTO{u.Balance})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
