package api

import (
	"net/http"
	"test3/package/repository"

	"github.com/gorilla/mux"
)

type api struct {
	r  *mux.Router
	db *repository.PGRepo
}

func New(router *mux.Router, db *repository.PGRepo) *api {
	return &api{r: router, db: db}
}

func (api *api) FillEndpoints() {
	api.r.HandleFunc("/api/auth/{guid}", api.auth)
	api.r.HandleFunc("/api/refresh/{refresh}", api.IpCheck(api.refresh))

}

func (api *api) ListenAndServe(addr string) error {
	return http.ListenAndServe(addr, api.r)
}
