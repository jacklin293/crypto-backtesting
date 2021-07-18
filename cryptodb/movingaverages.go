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
	MaKey    string // ma type+pair+interval
	Length   int
	Value    decimal.Decimal
	OpenTime time.Time
}

func getMaKey(maType string, pair string, interval string) string {
	return fmt.Sprintf("%s_%s_%s", maType, strings.ToLower(pair), interval)
}

func NewMovingAverage(data map[string]interface{}) MovingAverage {
	maType := data["ma_type"].(string)
	pair := data["pair"].(string)
	interval := data["interval"].(string)
	return MovingAverage{
		MaKey:    getMaKey(maType, pair, interval),
		Length:   data["length"].(int),
		Value:    data["value"].(decimal.Decimal),
		OpenTime: data["open_time"].(time.Time),
	}
}

func (db *DB) BatchInsertMovingAverages(mas []MovingAverage) (int64, error) {
	result := db.GormDB.Create(mas)
	return result.RowsAffected, result.Error
}

func (db *DB) GetLastestMovingAverage(maType string, pair string, interval string, length int) (*MovingAverage, int64, error) {
	var ma MovingAverage
	result := db.GormDB.Where("ma_key = ? AND length = ?", getMaKey(maType, pair, interval), length).Order("open_time DESC").First(&ma)
	if result.Error != nil && !errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return &ma, 0, result.Error
	}
	return &ma, result.RowsAffected, nil
}

func (db *DB) GetMovingAveragesByOpenTime(maType string, pair string, interval string, length int, openTime time.Time) (*MovingAverage, int64, error) {
	var ma MovingAverage
	result := db.GormDB.Where("ma_key = ? AND length = ? AND open_time = ?", getMaKey(maType, pair, interval), length, openTime).Find(&ma)
	if result.Error != nil && !errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return &ma, 0, result.Error
	}
	return &ma, result.RowsAffected, nil
}
