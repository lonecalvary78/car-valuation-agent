package model

import (
	"time"

	"github.com/shopspring/decimal"
)

type CarValuationResponse struct {
	Brand    string      `json:"brand"`
	Model    string      `json:"model"`
	Year     int         `json:"year"`
	Market   string      `json:"market"`
	Price    QuotedPrice `json:"quoted_price"`
	QuotedAt time.Time   `json:"quoted_at"`
}

type QuotedPrice struct {
	CurrencyCode string          `json:"currency"`
	LowerPrice   decimal.Decimal `json:"lower"`
	HigherPrice  decimal.Decimal `json:"higher"`
}
