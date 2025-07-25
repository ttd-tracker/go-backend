package middleware

import (
	"net/http"
)

type User struct{}

type AuthenticatedHandler func(http.ResponseWriter, *http.Request, *User)

type EnsureAuth struct {
	handler AuthenticatedHandler
}

func (e *EnsureAuth) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// authorization code
	e.handler(w, r, &User{})
}

func NewEnsureAuth(handler AuthenticatedHandler) *EnsureAuth {
	return &EnsureAuth{handler}
}
