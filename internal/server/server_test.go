package server

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestFinanceServer(t *testing.T) {
	svr := FinanceServer{}

	t.Run("get balance", func(t *testing.T) {
		req := newBalanceRequest(t)
		res := httptest.NewRecorder()
		svr.ServeHTTP(res, req)

		assertStatus(t, res.Code, http.StatusOK)
		assertContentType(t, res, "application/json")

		var got BalanceDTO
		err := json.NewDecoder(res.Body).Decode(&got)
		if err != nil {
			t.Fatalf("error decoding res: %v", err)
		}

		want := 1000
		if got.Value != want {
			t.Errorf("got balance %q want %q", got, want)
		}
	})
}

func newBalanceRequest(t testing.TB) *http.Request {
	req, err := http.NewRequest(http.MethodGet, "/balance", nil)
	if err != nil {
		t.Fatalf("error creating req: %v", err)
	}

	return req
}

func assertStatus(t testing.TB, got, want int) {
	t.Helper()

	if got != want {
		t.Errorf("got status %d want %d", got, want)
	}
}

func assertContentType(t testing.TB, res http.ResponseWriter, want string) {
	contentType := res.Header().Get("Content-Type")
	if contentType != want {
		t.Fatalf("got content-type %q want %q", contentType, want)
	}
}
