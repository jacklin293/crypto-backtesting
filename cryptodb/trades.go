package cryptodb

import (
	"time"

	"github.com/shopspring/decimal"
)

type Trade struct {
	Id     int64
	TestId int64

	BidPrice  decimal.Decimal
	BidVolume decimal.Decimal
	BidFee    decimal.Decimal
	BoughtAt  time.Time

	AskPrice  decimal.Decimal
	AskVolume decimal.Decimal
	AskFee    decimal.Decimal
	SoldAt    time.Time

	Cost   decimal.Decimal
	Total  decimal.Decimal // TODO could be negative
	Return decimal.Decimal // TODO percent, could be negative

	Details string // JSON string to store trigger condition, etc.

	CreatedAt time.Time
	UpdateAt  time.Time
}
