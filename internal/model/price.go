package model

import "time"

type Price struct {
	Exchange  string
	Coin      string
	Price     float64
	Timestamp time.Time
}
