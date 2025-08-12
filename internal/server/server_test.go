package server

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
	"time"
)

type StubStore struct {
	database map[int]Ruble
	history  []Op
}

func (s *StubStore) GetUser(id int) User {
	return User{id, s.database[id]}
}

func (s *StubStore) AddIncome(id int, income Ruble) Ruble {
	s.database[id] += income
	s.recordOp(Op{income, time.Now(), OpIncome})
	return s.database[id]
}

func (s *StubStore) recordOp(op Op) {
	s.history = append(s.history, op)
}

func TestFinanceServer(t *testing.T) {
	store := &StubStore{database: map[int]Ruble{
		1:  1000,
		20: 5000,
	}}
	svr := NewServer(store)

	t.Run("get one's balance", func(t *testing.T) {
		res := httptest.NewRecorder()
		svr.ServeHTTP(res, newBalanceRequest(t, 1))
		assertStatus(t, res.Code, http.StatusOK)
		assertContentType(t, res, "application/json")

		got, err := NewBalanceDTO(res.Body)
		assertNoErr(t, err)
		assertBalance(t, got.Value, Ruble(1000))
	})

	t.Run("get another's balance", func(t *testing.T) {
		res := httptest.NewRecorder()
		svr.ServeHTTP(res, newBalanceRequest(t, 20))
		assertStatus(t, res.Code, http.StatusOK)
		assertContentType(t, res, "application/json")

		got, err := NewBalanceDTO(res.Body)
		assertNoErr(t, err)
		assertBalance(t, got.Value, Ruble(5000))
	})
}

func TestIncome(t *testing.T) {
	store := &StubStore{database: map[int]Ruble{
		1:  1000,
		20: 5000,
	}}
	svr := NewServer(store)

	t.Run("income to user 1", func(t *testing.T) {
		id := 1
		res := httptest.NewRecorder()
		svr.ServeHTTP(res, newIncomeRequest(t, id, 500))
		assertStatus(t, res.Code, http.StatusCreated)
		assertContentType(t, res, "application/json")
		assertBalance(t, store.database[id], 1500)

		svr.ServeHTTP(httptest.NewRecorder(), newIncomeRequest(t, id, 500))
		assertBalance(t, store.database[id], 2000)

		if len(store.history) == 0 {
			t.Fatalf("history is empty")
		}

		if time.Since(store.history[0].Time) > (5 * time.Second) {
			t.Errorf("history: op 1: since op %v passed too much time", store.history[0].Time)
		}

		if store.history[0].Money != 500 {
			t.Errorf("history: op 1: got money %d want 500", store.history[0].Money)
		}

		if store.history[0].Type != OpIncome {
			t.Errorf("history: op 1: got Type %d want %q", store.history[0].Type, "income")
		}
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
