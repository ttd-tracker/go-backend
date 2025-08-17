package server

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
	"time"
)

type StubStore struct {
	db map[int]*User
}

func (s *StubStore) GetUser(id int) *User {
	return s.db[id]
}

func (s *StubStore) AddIncome(id int, income Ruble) Ruble {
	s.db[id].Balance = s.db[id].Balance.Add(income)
	s.recordOp(id, Op{income, time.Now(), OpIncome})
	return s.db[id].Balance
}

func (s *StubStore) AddExpense(id int, expense Ruble) Ruble {
	s.db[id].Balance = s.db[id].Balance.Sub(expense)
	return s.db[id].Balance
}

func (s *StubStore) recordOp(id int, op Op) {
	s.db[id].History = append(s.db[id].History, op)
}

func TestExtractBalance(t *testing.T) {
	store := &StubStore{
		db: map[int]*User{
			1:  {1, NewRuble(1000), []Op{}},
			20: {20, NewRuble(5000), []Op{}},
		},
	}
	svr := NewServer(store)

	t.Run("get one's balance", func(t *testing.T) {
		res := httptest.NewRecorder()
		svr.ServeHTTP(res, newBalanceRequest(t, 1))
		assertStatus(t, res.Code, http.StatusOK)
		assertContentType(t, res, "application/json")

		got, err := newBalanceDTOFromResponse(res.Body)
		assertNoErr(t, err)
		assertBalance(t, got.Money, 1000)
	})

	t.Run("get another's balance", func(t *testing.T) {
		res := httptest.NewRecorder()
		svr.ServeHTTP(res, newBalanceRequest(t, 20))
		assertStatus(t, res.Code, http.StatusOK)
		assertContentType(t, res, "application/json")

		got, err := newBalanceDTOFromResponse(res.Body)
		assertNoErr(t, err)
		assertBalance(t, got.Money, 5000)
	})
}

func TestIncome(t *testing.T) {
	store := &StubStore{
		db: map[int]*User{
			1:  {1, NewRuble(1000), []Op{}},
			20: {20, NewRuble(5000), []Op{}},
		},
	}
	svr := NewServer(store)

	id := 1
	res := httptest.NewRecorder()
	svr.ServeHTTP(res, newIncomeRequest(t, id, 500))
	assertStatus(t, res.Code, http.StatusCreated)
	assertContentType(t, res, "application/json")
	assertBalance(t, store.db[id].Balance.Float64(), 1500)

	svr.ServeHTTP(httptest.NewRecorder(), newIncomeRequest(t, id, 500))
	assertBalance(t, store.db[id].Balance.Float64(), 2000)

	if len(store.db[id].History) == 0 {
		t.Fatalf("history is empty")
	}
	if time.Since(store.db[id].History[0].Time) > (5 * time.Second) {
		t.Errorf("history: op 1: since op %v passed too much time", store.db[id].History[0].Time)
	}
	if store.db[id].History[0].Ruble.Float64() != 500 {
		t.Errorf("history: op 1: got money %.2f want 500", store.db[id].History[0].Ruble.Float64())
	}
	if store.db[id].History[0].Type != OpIncome {
		t.Errorf("history: op 1: got Type %d want %q", store.db[id].History[0].Type, "income")
	}

	if len(store.db[20].History) != 0 {
		t.Errorf("user 20 op history is not empty")
	}
}

func newBalanceDTOFromResponse(rdr io.Reader) (BalanceDTO, error) {
	var result BalanceDTO
	err := json.NewDecoder(rdr).Decode(&result)
	return result, err
}

func newBalanceRequest(t *testing.T, id int) *http.Request {
	req, err := http.NewRequest(http.MethodGet, "/balance", nil)
	assertNoErr(t, err)
	req.Header.Set("Authorization", strconv.Itoa(id))
	return req
}

func newIncomeRequest(t *testing.T, id int, income float64) *http.Request {
	req, err := http.NewRequest(http.MethodPost, fmt.Sprintf("/op/income/%.2f", income), nil)
	assertNoErr(t, err)
	req.Header.Set("Authorization", strconv.Itoa(id))
	return req
}

func assertBalance(t testing.TB, got, want float64) {
	t.Helper()

	if got != want {
		t.Errorf("got balance %.2f want %.2f", got, want)
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
