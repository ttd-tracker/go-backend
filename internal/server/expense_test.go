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
	svr.ServeHTTP(res, newExpenseRequest(t, 20, 150))
	assertStatus(t, res.Code, http.StatusCreated)
	//assertBalance(t, store.database[20].Value(), 4850)
}

func newExpenseRequest(t *testing.T, id int, money float64) *http.Request {
	req, err := http.NewRequest(http.MethodPost, fmt.Sprintf("/op/expense/%.2f", money), nil)
	assertNoErr(t, err)
	req.Header.Set("Authorization", strconv.Itoa(id))

	return req
}
