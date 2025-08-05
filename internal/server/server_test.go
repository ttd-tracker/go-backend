package server

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

type StubStore struct {
	database map[string]Ruble
}

func (s StubStore) GetUser(id string) User {
	return User{s.database[id]}
}

func TestFinanceServer(t *testing.T) {
	store := StubStore{map[string]Ruble{
		"1":  1000,
		"20": 5000,
	}}
	svr := NewServer(store)

	t.Run("get one's balance", func(t *testing.T) {
		req := newBalanceRequest(t, "1")
		res := httptest.NewRecorder()

		svr.ServeHTTP(res, req)
		assertStatus(t, res.Code, http.StatusOK)
		assertContentType(t, res, "application/json")

		got, err := NewBalanceDTO(res.Body)
		assertNoErr(t, err)
		assertBalance(t, got.Value, Ruble(1000))
	})

	t.Run("get another's balance", func(t *testing.T) {
		req := newBalanceRequest(t, "20")
		res := httptest.NewRecorder()

		svr.ServeHTTP(res, req)
		assertStatus(t, res.Code, http.StatusOK)
		assertContentType(t, res, "application/json")

		got, err := NewBalanceDTO(res.Body)
		assertNoErr(t, err)
		assertBalance(t, got.Value, Ruble(5000))
	})
}

func TestIncome(t *testing.T) {
	store := StubStore{map[string]Ruble{
		"1":  1000,
		"20": 5000,
	}}
	svr := NewServer(store)

	t.Run("income to user 1", func(t *testing.T) {
		req, err := http.NewRequest(http.MethodPost, "/op/income/500", nil)
		assertNoErr(t, err)
		req.Header.Set("Authorization", "1")

		res := httptest.NewRecorder()
		svr.ServeHTTP(res, req)

		assertStatus(t, res.Code, http.StatusCreated)
		assertContentType(t, res, "application/json")

		got, err := NewBalanceDTO(res.Body)
		assertNoErr(t, err)

		assertBalance(t, got.Value, 1500)
	})
}

func assertBalance(t testing.TB, got, want Ruble) {
	t.Helper()

	if got != want {
		t.Errorf("got balance %d want %d", got, want)
	}
}

func newBalanceRequest(t testing.TB, id string) *http.Request {
	req, err := http.NewRequest(http.MethodGet, "/balance", nil)
	if err != nil {
		t.Fatalf("error creating req: %v", err)
	}

	req.Header.Set("Authorization", id)

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

func assertNoErr(t testing.TB, err error) {
	t.Helper()

	if err != nil {
		t.Errorf("unexpected error: %s", err)
	}
}
