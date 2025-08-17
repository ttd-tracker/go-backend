package server

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
)

func TestExpense(t *testing.T) {
	store := &StubStore{database: map[int]Ruble{
		1:  NewRuble(1000),
		20: NewRuble(5000),
	}}
	svr := NewServer(store)

	res := httptest.NewRecorder()
	id := 20
	svr.ServeHTTP(res, newExpenseRequest(t, id, 150))
	assertStatus(t, res.Code, http.StatusCreated)
	assertContentType(t, res, "application/json")

	got, err := newBalanceDTOFromResponse(res.Body)
	assertNoErr(t, err)
	assertBalance(t, got.Money, 4850)
	assertBalance(t, store.database[20].Float64(), 4850)
}

func newExpenseRequest(t *testing.T, id int, money float64) *http.Request {
	req, err := http.NewRequest(http.MethodPost, fmt.Sprintf("/op/expense/%.2f", money), nil)
	assertNoErr(t, err)
	req.Header.Set("Authorization", strconv.Itoa(id))

	return req
}
