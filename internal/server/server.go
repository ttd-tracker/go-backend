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
	return NewEnsureAuth(ServeFinances)
}

func ServeFinances(w http.ResponseWriter, r *http.Request, u *User) {
	w.Header().Set("Content-Type", "application/json")

	status := http.StatusOK
	err := json.NewEncoder(w).Encode(BalanceDTO{u.Balance})
	if err != nil {
		status = http.StatusInternalServerError
	}
	w.WriteHeader(status)
}
