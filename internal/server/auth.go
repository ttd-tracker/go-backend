package server

import (
	"net/http"
)

type AuthenticatedHandler func(http.ResponseWriter, *http.Request, *User)

type EnsureAuth struct {
	handler AuthenticatedHandler
}

func (e *EnsureAuth) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// authorization code
	e.handler(w, r, &User{0})
}

func NewEnsureAuth(handler AuthenticatedHandler) *EnsureAuth {
	return &EnsureAuth{handler}
}
