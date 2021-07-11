package cryptodb

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

type Kline struct {
	KlineKey  string // e.g. btcusdt_1h
	Open      decimal.Decimal
	High      decimal.Decimal
	Low       decimal.Decimal
	Close     decimal.Decimal
	Volume    decimal.Decimal
	OpenTime  time.Time
	CloseTime time.Time
}

func getKlineKey(pair string, interval string) string {
	return fmt.Sprintf("%s_%s", strings.ToLower(pair), interval)
}

func (db *DB) GetKlines(pair string, interval string, limit int, orderType string) (*[]Kline, int64, error) {
	var klines []Kline
	result := db.GormDB.Where("kline_key = ?", getKlineKey(pair, interval)).Limit(limit).Order(fmt.Sprintf("open_time %s", orderType)).Find(&klines)
	if result.Error != nil && !errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return &klines, 0, result.Error
	}
	return &klines, result.RowsAffected, result.Error
}

func (db *DB) GetKlinesByOpenTime(pair string, interval string, limit int, openTime time.Time, orderType string) (*[]Kline, int64, error) {
	var klines []Kline
	result := db.GormDB.Where("kline_key = ? AND open_time >= ?", getKlineKey(pair, interval), openTime).Limit(limit).Order(fmt.Sprintf("open_time %s", orderType)).Find(&klines)
	if result.Error != nil && !errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return &klines, 0, result.Error
	}
	return &klines, result.RowsAffected, result.Error
}

func (db *DB) GetKlineByOpenTime(pair string, interval string, openTime time.Time) (*Kline, int64, error) {
	var kline Kline
	result := db.GormDB.Where("kline_key = ? AND open_time >= ?", getKlineKey(pair, interval), openTime).Find(&kline)
	if result.Error != nil && !errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return &kline, 0, result.Error
	}
	return &kline, result.RowsAffected, nil
}

func (db *DB) GetKlinesByPeriod(pair string, interval string, periodStart time.Time, periodEnd time.Time, limit int, orderType string) (*[]Kline, int64, error) {
	var klines []Kline
	result := db.GormDB.Where("kline_key = ? AND open_time BETWEEN ? AND ?", getKlineKey(pair, interval), periodStart, periodEnd).Limit(limit).Order(fmt.Sprintf("open_time %s", orderType)).Find(&klines)
	if result.Error != nil && !errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return &klines, 0, result.Error
	}
	return &klines, result.RowsAffected, result.Error
}
