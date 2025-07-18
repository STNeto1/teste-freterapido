package quotes

import (
	"time"

	"github.com/gofrs/uuid"
	"github.com/shopspring/decimal"
)

type Quote struct {
	ID        uuid.UUID       `ch:"id" json:"-"`
	Name      string          `ch:"name" json:"name"`
	Service   string          `ch:"service" json:"service"`
	Deadline  uint8           `ch:"deadline" json:"deadline"`
	Price     decimal.Decimal `ch:"price" json:"price"`
	CreatedAt time.Time       `ch:"timestamp" json:"-"`
}
