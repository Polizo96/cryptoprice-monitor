package storage

import (
	"cryptoprice-monitor/internal/model"
	"sync"
	"time"
)

//type PriceData struct {
//	Exchange string
//	Coin     string
//	Price    float64
//}

type Cache struct {
	mu     sync.RWMutex
	prices map[string]model.Price
}

func NewCache() *Cache {
	return &Cache{
		prices: make(map[string]model.Price),
	}
}

func (c *Cache) Set(exchange, coin string, price float64) {
	key := exchange + ":" + coin
	c.mu.Lock()
	defer c.mu.Unlock()
	c.prices[key] = model.Price{
		Exchange:  exchange,
		Coin:      coin,
		Price:     price,
		Timestamp: time.Now()}
}

func (c *Cache) Get(exchange, coin string) (model.Price, bool) {
	key := exchange + ":" + coin
	c.mu.RLock()
	defer c.mu.RUnlock()
	data, ok := c.prices[key]
	return data, ok
}
