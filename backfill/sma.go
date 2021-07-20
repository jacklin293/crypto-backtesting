package backfill

import (
	"crypto-backtesting/cryptodb"
	"crypto-backtesting/utils"
	"fmt"
	"time"

	"github.com/shopspring/decimal"
)

type Sma struct {
	baseMA
}

func (ma *Sma) backfill() (err error) {
	lengthMins, err := utils.ConvertIntervalToMins(ma.interval)
	if err != nil {
		return
	}

	for err == nil {
		latestSma, maCount, err := ma.db.GetLastestMovingAverage(ma.maType, ma.pair, ma.interval, ma.length)
		if err != nil {
			return err
		}

		var klines *[]cryptodb.Kline
		var klineCount int64

		// If there is no ma before, start from beginning
		if maCount == 0 {
			klines, klineCount, err = ma.db.GetKlines(ma.pair, ma.interval, KLINE_BATCH_LIMIT, "ASC")
		} else {
			startTime := latestSma.OpenTime.Add(time.Minute * time.Duration(-lengthMins*(ma.length-2)))
			klines, klineCount, err = ma.db.GetKlinesByOpenTime(ma.pair, ma.interval, KLINE_BATCH_LIMIT, startTime, "ASC")
		}
		if err != nil {
			return err
		}
		if klineCount == 0 {
			return fmt.Errorf("There is no more klines of %s-%s to backfill SMA", ma.pair, ma.interval)
		}
		if len(*klines) < ma.length {
			return fmt.Errorf("There is no enough klines of %s-%s to backfill SMA with length %d", ma.pair, ma.interval, ma.length)
		}

		// Start to calculate sma based on existing klines data and return MovingAverage struct
		smas := ma.calculateSma(klines)
		maCount, err = ma.db.BatchInsertMovingAverages(smas)
		if err != nil {
			return err
		}
		fmt.Printf("%d rows have been inserted into table 'moving_averages' successfully\n", maCount)
	}

	return
}

func (ma *Sma) calculateSma(klines *[]cryptodb.Kline) (smas []cryptodb.MovingAverage) {
	for i := 0; i < len(*klines); i++ {
		// Add 'i != 0' into condition to prevent from breaking due to only 3 klines
		if i != 0 && i > len(*klines)-ma.length {
			break
		}

		l := i + ma.length
		var total decimal.Decimal
		for j := i; j < l; j++ {
			total = total.Add((*klines)[j].Close)
		}

		maData := map[string]interface{}{
			"ma_type":   ma.maType,
			"pair":      ma.pair,
			"interval":  ma.interval,
			"length":    ma.length,
			"value":     total.Div(decimal.NewFromInt(int64(ma.length))),
			"open_time": (*klines)[l-1].OpenTime,
		}
		sma := cryptodb.NewMovingAverage(maData)
		smas = append(smas, sma)
	}
	return
}
