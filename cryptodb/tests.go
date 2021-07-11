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

	Cost    decimal.Decimal
	Revenue decimal.Decimal // TODO could be negative
	Fee     decimal.Decimal
	Profit  decimal.Decimal // TODO cound be negative, revenue - cost - fee
	Return  decimal.Decimal // TODO count be negative, percentage, profit / cost

	TradeCount int64 // count of trades
	Comment    string
	CreatedAt  time.Time
	UpdatedAt  time.Time
}
