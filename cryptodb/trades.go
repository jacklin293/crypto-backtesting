package cryptodb

import (
	"time"

	"github.com/shopspring/decimal"
)

type Trade struct {
	Id        int64
	TestId    int64
	Status    int64 // 0: undone 1: done
	BidPrice  decimal.Decimal
	BidVolume decimal.Decimal
	BidFee    decimal.Decimal
	BoughtAt  time.Time
	AskPrice  decimal.Decimal
	AskVolume decimal.Decimal
	AskFee    decimal.Decimal
	SoldAt    time.Time
	Cost      decimal.Decimal
	Revenue   decimal.Decimal
	Profit    decimal.Decimal
	ROI       decimal.Decimal
	CreatedAt time.Time
	UpdateAt  time.Time
}
