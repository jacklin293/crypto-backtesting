package cryptodb

import (
	"time"

	"github.com/shopspring/decimal"
)

type Test struct {
	Id          int64
	StrategyId  int64
	PeriodStart time.Time
	PeriodEnd   time.Time
	Fee         decimal.Decimal // percent
	SeedMoney   decimal.Decimal
	Revenue     decimal.Decimal // TODO could be negative
	Return      decimal.Decimal // TODO percent, could be negative
	WinRate     decimal.Decimal // percent, could be negative
	Comment     string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
