package cryptodb

import (
	"fmt"
	"time"

	"github.com/shopspring/decimal"
)

type Kline struct {
	PairInterval string // e.g. btcusdt_1h
	Open         decimal.Decimal
	High         decimal.Decimal
	Low          decimal.Decimal
	Close        decimal.Decimal
	Volume       decimal.Decimal
	OpenTime     time.Time
	CloseTime    time.Time
}

func (db *DB) GetKlines(pair string, interval string, limit int, orderType string) ([]Kline, error) {
	var klines []Kline
	result := db.GormDB.Where("pair_interval = ?", getPairInterval(pair, interval)).Limit(limit).Order(fmt.Sprintf("open_time %s", orderType)).Find(&klines)
	return klines, result.Error
}

func (db *DB) GetKlinesByOpenTime(pair string, interval string, limit int, openTime time.Time, orderType string) ([]Kline, error) {
	var klines []Kline
	result := db.GormDB.Where("pair_interval = ? AND open_time >= ?", getPairInterval(pair, interval), openTime).Limit(limit).Order(fmt.Sprintf("open_time %s", orderType)).Find(&klines)
	return klines, result.Error
}

// TODO havn't been tested yest
func (db *DB) GetKlinesByPeriod(pair string, interval string, periodStart time.Time, periodEnd time.Time, limit int) ([]Kline, error) {
	var klines []Kline
	result := db.GormDB.Where("pair_interval = ? AND open_time BETWEEN ? AND ?", getPairInterval(pair, interval), periodStart, periodEnd).Limit(limit).Find(&klines)
	return klines, result.Error
}
