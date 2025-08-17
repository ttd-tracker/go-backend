package server

import (
	"net/http"
	"strconv"
)

type AuthenticatedHandler func(http.ResponseWriter, *http.Request, *User)

type EnsureAuth struct {
	handler AuthenticatedHandler
	store   FinanceStore
}

func (e *EnsureAuth) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.Header.Get("Authorization"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	// no check whether a user exists
	user := e.store.GetUser(id)
	e.handler(w, r, user)
}

func NewEnsureAuth(handler AuthenticatedHandler, store FinanceStore) *EnsureAuth {
	return &EnsureAuth{handler, store}
}
