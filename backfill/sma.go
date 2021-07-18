package backfill

import (
	"crypto-backtesting/cryptodb"
	"crypto-backtesting/utils"
	"fmt"
	"time"

	"github.com/shopspring/decimal"
)

func (ma *MA) HandleBackfillSma(batchLimit int) (err error) {
	lengthMins, err := utils.ConvertIntervalToMins(ma.Interval)
	if err != nil {
		return
	}

	for err == nil {
		latestSma, maCount, err := ma.Db.GetLastestMovingAverage(ma.MaType, ma.Pair, ma.Interval, ma.Length)
		if err != nil {
			return err
		}

		var klines *[]cryptodb.Kline
		var klineCount int64

		// If there is no ma before, start from beginning
		if maCount == 0 {
			klines, klineCount, err = ma.Db.GetKlines(ma.Pair, ma.Interval, batchLimit, "ASC")
		} else {
			startTime := latestSma.OpenTime.Add(time.Minute * time.Duration(-lengthMins*(ma.Length-2)))
			klines, klineCount, err = ma.Db.GetKlinesByOpenTime(ma.Pair, ma.Interval, batchLimit, startTime, "ASC")
		}
		if err != nil {
			return err
		}
		if klineCount == 0 {
			return fmt.Errorf("There is no more klines of %s-%s to backfill SMA", ma.Pair, ma.Interval)
		}
		if len(*klines) < ma.Length {
			return fmt.Errorf("There is no enough klines of %s-%s to backfill SMA with length %d", ma.Pair, ma.Interval, ma.Length)
		}

		// Start to calculate sma based on existing klines data and return MovingAverage struct
		smas := ma.calculateSma(klines)
		maCount, err = ma.Db.BatchInsertMovingAverages(smas)
		if err != nil {
			return err
		}
		fmt.Printf("%d rows have been inserted into table 'moving_averages' successfully\n", maCount)
	}

	return
}

func (ma *MA) calculateSma(klines *[]cryptodb.Kline) (smas []cryptodb.MovingAverage) {
	for i := 0; i < len(*klines); i++ {
		// Add 'i != 0' into condition to prevent from breaking due to only 3 klines
		if i != 0 && i > len(*klines)-ma.Length {
			break
		}

		l := i + ma.Length
		var total decimal.Decimal
		for j := i; j < l; j++ {
			total = total.Add((*klines)[j].Close)
		}

		maData := map[string]interface{}{
			"ma_type":   ma.MaType,
			"pair":      ma.Pair,
			"interval":  ma.Interval,
			"length":    ma.Length,
			"value":     total.Div(decimal.NewFromInt(int64(ma.Length))),
			"open_time": (*klines)[l-1].OpenTime,
		}
		sma := cryptodb.NewMovingAverage(maData)
		smas = append(smas, sma)
	}
	return
}
