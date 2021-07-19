package cryptodb

import (
	"time"

	"github.com/shopspring/decimal"
	"gorm.io/datatypes"
)

type Strategy struct {
	Id           int64
	Title        string
	Description  string         // Description
	StrategyType string         // strategy type
	Params       datatypes.JSON // JSON string to store ad-hoc params of strategy type
	Start        time.Time
	End          time.Time
	Cost         decimal.Decimal
	Enabled      int // 0: disabled 1: enabled
	CreatedAt    time.Time
}
