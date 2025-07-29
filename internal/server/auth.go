package server

import (
	"net/http"
)

type AuthenticatedHandler func(http.ResponseWriter, *http.Request, *User)

type EnsureAuth struct {
	handler AuthenticatedHandler
}

func (e *EnsureAuth) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("Authorization") == "1" {
		e.handler(w, r, &User{1000})
		return
	}

	e.handler(w, r, &User{5000})
}

func NewEnsureAuth(handler AuthenticatedHandler) *EnsureAuth {
	return &EnsureAuth{handler}
}
