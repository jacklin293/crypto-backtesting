package cryptodb

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

type MovingAverage struct {
	Type         string
	PairInterval string
	Length       int
	Value        decimal.Decimal
	OpenTime     time.Time
}

func getPairInterval(pair string, interval string) string {
	return fmt.Sprintf("%s_%s", strings.ToLower(pair), interval)
}

func NewMovingAverage(maType string, pair string, interval string, length int, value decimal.Decimal, openTime time.Time) MovingAverage {
	return MovingAverage{
		Type:         maType,
		PairInterval: getPairInterval(pair, interval),
		Length:       length,
		Value:        value,
		OpenTime:     openTime,
	}
}

func (db *DB) BatchInsertMovingAverages(mas []MovingAverage) (int64, error) {
	result := db.GormDB.Create(mas)
	return result.RowsAffected, result.Error
}

// TODO
func (db *DB) GetMovingAverages(maType string, pair string, interval string, length int, openTime time.Time) ([]MovingAverage, error) {
	return []MovingAverage{}, nil
}

func (db *DB) GetLastestMovingAverage(maType string, pair string, interval string, length int) (*MovingAverage, int64, error) {
	var ma MovingAverage
	result := db.GormDB.Where("type = ? AND pair_interval = ? AND length = ?", maType, getPairInterval(pair, interval), length).Order("open_time DESC").First(&ma)
	if result.Error != nil && !errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return &MovingAverage{}, 0, result.Error
	}
	return &ma, result.RowsAffected, nil
}
