package server

import "net/http"

type FinanceServer struct{}

func (f *FinanceServer) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	res.WriteHeader(http.StatusCreated)
}
