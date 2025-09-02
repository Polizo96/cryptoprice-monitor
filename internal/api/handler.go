package api

import (
	"cryptoprice-monitor/internal/storage"
	"encoding/json"
	"net/http"
)

type PriceResponse struct {
	Exchange string  `json:"exchange"`
	Coin     string  `json:"coin"`
	Price    float64 `json:"price"`
}

func PricesHandler(cache *storage.Cache) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		coin := r.URL.Query().Get("coin")
		exchange := r.URL.Query().Get("exchange")

		if coin == "" || exchange == "" {
			http.Error(w, "query params 'coin' and 'exchange' are required", http.StatusBadRequest)
			return
		}

		priceData, found := cache.Get(exchange, coin)
		if !found {
			http.Error(w, "price data not found", http.StatusNotFound)
			return
		}

		response := PriceResponse{
			Exchange: priceData.Exchange,
			Coin:     priceData.Coin,
			Price:    priceData.Price,
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}
}
