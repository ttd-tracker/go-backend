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

		want := Ruble(1000)
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

//t.Run("new income", func(t *testing.T) {
//	svr := FinanceServer{}
//
//	req, err := http.NewRequest(http.MethodPost, "/op/income", nil)
//	if err != nil {
//		t.Fatalf("error creating req: %v", err)
//	}
//	res := httptest.NewRecorder()
//	svr.ServeHTTP(res, req)
//
//	if res.Code != http.StatusCreated {
//		t.Errorf("got status %d want %d", res.Code, http.StatusCreated)
//	}
//
//	// в res.Body мне ничего не нужно, раз я получаю успешный код. как убедиться, что
//	// баланс пользователя обновлён вне стора? как получить к нему доступ?
//	// getBalance поможет в таком случае. я прекрасно знаю, что мне нужна сущность-хранилище.
//
//	// если есть getBalance, что я верну? нужный доход. если изначальный = 0, то конечный = начальный + значение,
//	//	которое мы ввели. давай уж тогда сделаем гет беленс. это не во всех имплементациях будет дб оп, так что
//	//	может быть полезен
//})
