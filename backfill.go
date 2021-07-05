package main

import (
	"crypto-backtesting/cryptodb"
	"flag"
	"fmt"
	"time"

	"github.com/shopspring/decimal"
)

type EMA struct {
	maType          string // MA type
	pair            string
	interval        string
	length          int
	todayMultiplier decimal.Decimal // Today's multiplier
	ydayMultiplier  decimal.Decimal // Yesterday's multiplier
}

func handleBackfillEma(db *cryptodb.DB, maType string, pair string, interval string, length int) (err error) {
	if maType == "" || pair == "" || interval == "" || length == 0 {
		flag.PrintDefaults()
		return fmt.Errorf("All of pair, interval and length need to be specified")
	}

	todayMul := decimal.NewFromInt(int64(2)).Div(decimal.NewFromInt(int64(length + 1)))
	ydayMul := decimal.NewFromInt(1).Sub(todayMul)
	ema := EMA{
		maType:          maType,
		pair:            pair,
		interval:        interval,
		length:          length,
		todayMultiplier: todayMul,
		ydayMultiplier:  ydayMul,
	}

	for err == nil {
		ma, count, err := db.GetLastestMovingAverage(maType, pair, interval, length)
		if err != nil {
			return err
		}

		var klines []cryptodb.Kline
		var lastEma decimal.Decimal

		// If there is no ma, backfill all
		if count == 0 {
			klines, err = db.GetKlines(pair, interval, DB_KLINES_BATCH_SELECT_NUMBER, "ASC")
		} else {
			// Start from next time
			startTime := ma.OpenTime.Add(time.Hour * time.Duration(4))
			klines, err = db.GetKlinesByOpenTime(pair, interval, DB_KLINES_BATCH_SELECT_NUMBER, startTime, "ASC")
		}
		if err != nil {
			return err
		}
		if len(klines) == 0 {
			return fmt.Errorf("There is no more klins of %s-%s to backfill EMA", pair, interval)
		}

		if count == 0 {
			// If there is no history, use first close price as last ema
			lastEma = klines[0].Close
		} else {
			// Use the latest EMA as last ema
			lastEma = ma.Value
		}

		// Start to calculate ema based on existing klines data and return MovingAverage struct
		emas := ema.calculateEma(lastEma, klines)
		count, err = db.BatchInsertMovingAverages(emas)
		if err != nil {
			return err
		}
		fmt.Printf("%d rows have been inserted into moving_averages successfully\n", count)
	}

	return
}

// EMA=Price(t)×k+EMA(y)×(1−k)
// t=today, y=yesterday, k=2÷(N+1), N=length
func (ema *EMA) calculateEma(lastEma decimal.Decimal, klines []cryptodb.Kline) (emas []cryptodb.MovingAverage) {
	var maVal decimal.Decimal
	for i, kline := range klines {
		if i == 0 {
			maVal = kline.Close.Mul(ema.todayMultiplier).Add(lastEma.Mul(ema.ydayMultiplier))
		} else {
			maVal = kline.Close.Mul(ema.todayMultiplier).Add(emas[i-1].Value.Mul(ema.ydayMultiplier))
		}
		ema := cryptodb.NewMovingAverage(ema.maType, ema.pair, ema.interval, ema.length, maVal, kline.OpenTime)
		emas = append(emas, ema)
	}
	return
}
