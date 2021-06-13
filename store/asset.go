package store

import "time"

type Asset struct {
	Price float64
	Time time.Time
	Qty float64
}
