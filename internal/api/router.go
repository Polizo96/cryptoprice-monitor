package api

import (
	"cryptoprice-monitor/internal/storage"
	"github.com/gorilla/mux"
	"net/http"
)

func NewRouter(cache *storage.Cache) http.Handler {
	r := mux.NewRouter()
	r.HandleFunc("/prices", PricesHandler(cache)).Methods("GET")

	return r
}
