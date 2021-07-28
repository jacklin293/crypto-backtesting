package cryptodb

import (
	"time"

	"github.com/shopspring/decimal"
	"gorm.io/datatypes"
)

type Test struct {
	Id               int64
	StrategyId       int64
	StrategyType     string
	StrategyPair     string
	StrategyInterval string
	StrategyParams   datatypes.JSON
	Start            time.Time
	End              time.Time
	Cost             decimal.Decimal
	Revenue          decimal.Decimal
	Fee              decimal.Decimal
	Profit           decimal.Decimal
	ROI              decimal.Decimal
	TradeCount       int64
	Comment          string
	CreatedAt        time.Time
}
