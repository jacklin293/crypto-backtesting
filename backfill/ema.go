package backfill

import (
	"crypto-backtesting/cryptodb"
	"crypto-backtesting/utils"
	"fmt"
	"time"

	"github.com/shopspring/decimal"
)

func (ma *MA) HandleBackfillEma(batchLimit int) (err error) {
	// Today and yesterday's multiplier
	tdyMul := decimal.NewFromInt(int64(2)).Div(decimal.NewFromInt(int64(ma.Length + 1)))
	ytdMul := decimal.NewFromInt(1).Sub(tdyMul)
	lengthMins, err := utils.ConvertIntervalToMins(ma.Interval)
	if err != nil {
		return
	}

	for err == nil {
		latestEma, maCount, err := ma.Db.GetLastestMovingAverage(ma.MaType, ma.Pair, ma.Interval, ma.Length)
		if err != nil {
			return err
		}

		var klines *[]cryptodb.Kline
		var ema decimal.Decimal
		var klineCount int64

		// If there is no ma before, start from beginning
		if maCount == 0 {
			klines, klineCount, err = ma.Db.GetKlines(ma.Pair, ma.Interval, batchLimit, "ASC")
		} else {
			startTime := latestEma.OpenTime.Add(time.Minute * time.Duration(lengthMins))
			klines, klineCount, err = ma.Db.GetKlinesByOpenTime(ma.Pair, ma.Interval, batchLimit, startTime, "ASC")
		}
		if err != nil {
			return err
		}
		if klineCount == 0 {
			return fmt.Errorf("There is no more klines of %s-%s to backfill EMA", ma.Pair, ma.Interval)
		}

		if maCount == 0 {
			// If there is no history, use first close price as last ema
			ema = (*klines)[0].Close
		} else {
			// Use the latest EMA as last ema
			ema = latestEma.Value
		}

		// Start to calculate ema based on existing klines data and return MovingAverage struct
		emas := ma.calculateEma(tdyMul, ytdMul, ema, klines)
		maCount, err = ma.Db.BatchInsertMovingAverages(emas)
		if err != nil {
			return err
		}
		fmt.Printf("%d rows have been inserted into table 'moving_averages' successfully\n", maCount)
	}

	return
}

// EMA=Price(t)×k+EMA(y)×(1−k)
// t=today, y=yesterday, k=2÷(N+1), N=length
func (ma *MA) calculateEma(tdyMul decimal.Decimal, ytdMul decimal.Decimal, lastEma decimal.Decimal, klines *[]cryptodb.Kline) (emas []cryptodb.MovingAverage) {
	var maVal decimal.Decimal
	for i, kline := range *klines {
		if i == 0 {
			maVal = kline.Close.Mul(tdyMul).Add(lastEma.Mul(ytdMul))
		} else {
			maVal = kline.Close.Mul(tdyMul).Add(emas[i-1].Value.Mul(ytdMul))
		}
		maData := map[string]interface{}{
			"ma_type":   ma.MaType,
			"pair":      ma.Pair,
			"interval":  ma.Interval,
			"length":    ma.Length,
			"value":     maVal,
			"open_time": kline.OpenTime,
		}
		ema := cryptodb.NewMovingAverage(maData)
		emas = append(emas, ema)
	}
	return
}
