package cryptodb

import (
	"errors"
	"time"

	"github.com/shopspring/decimal"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type Strategy struct {
	Id           int64
	StrategyType string // strategy type
	Pair         string
	Interval     string
	Params       datatypes.JSONMap // JSON string to store ad-hoc params of strategy type e.g. {"ma_type":"ema","pair":"btcusdt","interval":"4h","length":18}
	Start        time.Time
	End          time.Time
	Cost         decimal.Decimal
	Enabled      int // 0: disabled 1: enabled
	CreatedAt    time.Time
}

func (db *DB) GetAllEnabledStrategies() (*[]Strategy, int64, error) {
	var strategies []Strategy
	result := db.GormDB.Where("enabled = 1").Order("strategy_type").Find(&strategies)
	if result.Error != nil && !errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return &strategies, 0, result.Error
	}
	return &strategies, result.RowsAffected, result.Error
}

func (db *DB) GetStrategyById(id int64) (*Strategy, int64, error) {
	var strategy Strategy
	result := db.GormDB.Where("id = ?", id).Find(&strategy)
	if result.Error != nil && !errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return &strategy, 0, result.Error
	}
	return &strategy, result.RowsAffected, result.Error
}
