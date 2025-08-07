package server

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
)

type StubStore struct {
	database map[int]Ruble
}

func (s StubStore) GetUser(id int) User {
	return User{id, s.database[id]}
}

func (s StubStore) RecordIncome(id int, income Ruble) Ruble {
	s.database[id] += income
	return s.database[id]
}

func TestFinanceServer(t *testing.T) {
	store := StubStore{map[int]Ruble{
		1:  1000,
		20: 5000,
	}}
	svr := NewServer(store)

	t.Run("get one's balance", func(t *testing.T) {
		req := newBalanceRequest(t, 1)
		res := httptest.NewRecorder()
		svr.ServeHTTP(res, req)
		assertStatus(t, res.Code, http.StatusOK)
		assertContentType(t, res, "application/json")

		got, err := NewBalanceDTO(res.Body)
		assertNoErr(t, err)
		assertBalance(t, got.Value, Ruble(1000))
	})

	t.Run("get another's balance", func(t *testing.T) {
		req := newBalanceRequest(t, 20)
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
	store := StubStore{map[int]Ruble{
		1:  1000,
		20: 5000,
	}}
	svr := NewServer(store)

	t.Run("income to user 1", func(t *testing.T) {
		req := newIncomeRequest(t, 1, 500)
		res := httptest.NewRecorder()
		svr.ServeHTTP(res, req)
		assertStatus(t, res.Code, http.StatusCreated)
		assertContentType(t, res, "application/json")

		got, err := NewBalanceDTO(res.Body)
		assertNoErr(t, err)
		assertBalance(t, got.Value, 1500)

		req = newIncomeRequest(t, 1, 500)
		res = httptest.NewRecorder()
		svr.ServeHTTP(res, req)

		got, err = NewBalanceDTO(res.Body)
		assertNoErr(t, err)
		assertBalance(t, got.Value, 2000)
	})
}

func newBalanceRequest(t *testing.T, id int) *http.Request {
	req, err := http.NewRequest(http.MethodGet, "/balance", nil)
	assertNoErr(t, err)
	req.Header.Set("Authorization", strconv.Itoa(id))
	return req
}

func newIncomeRequest(t *testing.T, id int, income Ruble) *http.Request {
	req, err := http.NewRequest(http.MethodPost, fmt.Sprintf("/op/income/%d", income), nil)
	assertNoErr(t, err)
	req.Header.Set("Authorization", strconv.Itoa(id))
	return req
}

func assertBalance(t testing.TB, got, want Ruble) {
	t.Helper()

	if got != want {
		t.Errorf("got balance %d want %d", got, want)
	}
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
