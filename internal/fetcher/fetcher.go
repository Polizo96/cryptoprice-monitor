package fetcher

import (
	"cryptoprice-monitor/configs"
	"cryptoprice-monitor/internal/storage"
	"cryptoprice-monitor/internal/util"
)

func StartFetcher(config *configs.Config, cache *storage.Cache) *Scheduler {
	util.Init()

	scheduler := NewScheduler(config.Fetcher.IntervalSeconds)
	scheduler.Start(func() {
		fetchPrices(config, cache)
	})

	return scheduler
}

func fetchPrices(config *configs.Config, cache *storage.Cache) {
	for _, ex := range config.Fetcher.Exchanges {
		for _, coin := range config.Coins {
			var price float64
			var err error

			switch ex.Name {
			case "binance":
				price, err = FetchBinancePrice(ex.ApiUrl, coin)
			case "coinbase":
				price, err = FetchCoinbasePrice(ex.ApiUrl, coin)
			case "kraken":
				price, err = FetchKrakenPrice(ex.ApiUrl, coin)
			// add more brokers here
			default:
				util.Info.Printf("Exchange %s: fetch client not supported\n", ex.Name)
			}
			if err != nil {
				util.Error.Printf("Failed to fetch %s price from %s: %v", coin, ex.Name, err)
				continue
			}
			cache.Set(ex.Name, coin, price)
			util.Info.Printf("Updated price: %s %s = %f", ex.Name, coin, price)
		}
	}
}
