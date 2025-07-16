package quotes

import (
	"time"

	"github.com/gofrs/uuid"
	"github.com/shopspring/decimal"
)

type Quote struct {
	ID        uuid.UUID       `ch:"id"`
	Name      string          `ch:"name"`
	Service   string          `ch:"service"`
	Deadline  uint8           `ch:"deadline"`
	Price     decimal.Decimal `ch:"price"`
	CreatedAt time.Time       `ch:"timestamp"`
}
