package models

import (
	"time"
)

type Transaction struct {
	Amount    float64   `json:"amount"`
	Timestamp time.Time `json:"timestamp"`
	Location  string    `json:"location"`
}
