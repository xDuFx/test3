package api

import (
	// "bytes"
	"net/http"
	"net/http/httptest"

	"test3/package/repository"
	"testing"

	"github.com/gorilla/mux"
)

func Test_auth(t *testing.T) {
	db, err := repository.New("postgres://postgres:123@localhost:5432/postgres")
	if err != nil {
		t.Fatal(err.Error())
	}
	api := New(mux.NewRouter(), db)

	req, err := http.NewRequest("GET", "https://localhost:8090/auth/1452asd145", nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()


	router := mux.NewRouter()
	router.HandleFunc("/auth/{guid}", api.auth)
	router.ServeHTTP(rr, req)	
	
	if rr.Code != http.StatusOK  {
		t.Errorf("handler should have failed on routeVariable : got %v want %v",
			 rr.Code, http.StatusOK)
	}
	token :=req.FormValue("Authorization")
	req.ParseForm()
	t.Log(token, rr.HeaderMap.Values("Authorization") )
}
