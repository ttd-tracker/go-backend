package server

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestFinanceServer(t *testing.T) {
	svr := FinanceServer{}

	req, err := http.NewRequest(http.MethodPost, "/op/income", nil)
	if err != nil {
		t.Fatalf("error creating req: %v", err)
	}

	res := httptest.NewRecorder()

	svr.ServeHTTP(res, req)

	if res.Code != http.StatusCreated {
		t.Errorf("got status %d want %d", res.Code, http.StatusCreated)
	}
}
