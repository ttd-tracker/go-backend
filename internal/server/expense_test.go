package server

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
)

func TestExpense(t *testing.T) {
	store := &StubStore{
		db: map[int]*User{
			1:  {1, NewRuble(1000), []Op{}},
			20: {20, NewRuble(5000), []Op{}},
		},
	}
	svr := NewServer(store)

	res := httptest.NewRecorder()
	id := 20
	svr.ServeHTTP(res, newExpenseRequest(t, id, 150))
	assertStatus(t, res.Code, http.StatusCreated)
	assertContentType(t, res, "application/json")

	got, err := newBalanceDTOFromResponse(res.Body)
	assertNoErr(t, err)
	assertBalance(t, got.Cash, 4850)
	assertBalance(t, store.db[20].Balance.Float64(), 4850)

	assertHistoryOp(t, store.db[id].History, 150, OpExpense)
}

func newExpenseRequest(t *testing.T, id int, cash float64) *http.Request {
	req, err := http.NewRequest(http.MethodPost, fmt.Sprintf("/op/expense/%.2f", cash), nil)
	assertNoErr(t, err)
	req.Header.Set("Authorization", strconv.Itoa(id))

	return req
}
