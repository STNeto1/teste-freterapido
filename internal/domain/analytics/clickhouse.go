package analytics

import (
	"github.com/shopspring/decimal"
)

type ServiceMetrics struct {
	Carrier     string          `ch:"carrier" json:"carrier"`
	TotalQuotes uint64          `ch:"total_quotes" json:"total_quotes"`
	TotalPrice  decimal.Decimal `ch:"total_price" json:"total_price"`
	AvgPrice    decimal.Decimal `ch:"avg_price" json:"avg_price"`
	MinPrice    decimal.Decimal `ch:"min_price" json:"min_price"`
	MaxPrice    decimal.Decimal `ch:"max_price" json:"max_price"`
}
