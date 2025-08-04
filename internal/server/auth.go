package server

import (
	"net/http"
)

type AuthenticatedHandler func(http.ResponseWriter, *http.Request, User)

type EnsureAuth struct {
	handler AuthenticatedHandler
	store   FinanceStore
}

func (e *EnsureAuth) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	user := e.store.GetUser(r.Header.Get("Authorization"))
	e.handler(w, r, user)
}

func NewEnsureAuth(handler AuthenticatedHandler, store FinanceStore) *EnsureAuth {
	return &EnsureAuth{handler, store}
}
