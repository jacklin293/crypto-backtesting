package backfill

import (
	"crypto-backtesting/cryptodb"
	"crypto-backtesting/utils"
	"fmt"
	"time"

	"github.com/shopspring/decimal"
)

func (p *Params) HandleBackfillSma(batchLimit int) (err error) {
	lengthMins, err := utils.ConvertIntervalToMins(p.Interval)
	if err != nil {
		return
	}

	for err == nil {
		latestSma, maCount, err := p.Db.GetLastestMovingAverage(p.MaType, p.Pair, p.Interval, p.Length)
		if err != nil {
			return err
		}

		var klines *[]cryptodb.Kline
		var klineCount int64

		// If there is no ma before, start from beginning
		if maCount == 0 {
			klines, klineCount, err = p.Db.GetKlines(p.Pair, p.Interval, batchLimit, "ASC")
		} else {
			startTime := latestSma.OpenTime.Add(time.Minute * time.Duration(-lengthMins*(p.Length-2)))
			klines, klineCount, err = p.Db.GetKlinesByOpenTime(p.Pair, p.Interval, batchLimit, startTime, "ASC")
		}
		if err != nil {
			return err
		}
		if klineCount == 0 {
			return fmt.Errorf("There is no more klines of %s-%s to backfill SMA", p.Pair, p.Interval)
		}
		if len(*klines) < p.Length {
			return fmt.Errorf("There is no enough klines of %s-%s to backfill SMA with length %d", p.Pair, p.Interval, p.Length)
		}

		// Start to calculate sma based on existing klines data and return MovingAverage struct
		smas := p.calculateSma(klines)
		maCount, err = p.Db.BatchInsertMovingAverages(smas)
		if err != nil {
			return err
		}
		fmt.Printf("%d rows have been inserted into table 'moving_averages' successfully\n", maCount)
	}

	return
}

func (p *Params) calculateSma(klines *[]cryptodb.Kline) (smas []cryptodb.MovingAverage) {
	for i := 0; i < len(*klines); i++ {
		// Add 'i != 0' into condition to prevent from breaking due to only 3 klines
		if i != 0 && i > len(*klines)-p.Length {
			break
		}

		l := i + p.Length
		var total decimal.Decimal
		for j := i; j < l; j++ {
			total = total.Add((*klines)[j].Close)
		}

		maData := map[string]interface{}{
			"ma_type":   p.MaType,
			"pair":      p.Pair,
			"interval":  p.Interval,
			"length":    p.Length,
			"value":     total.Div(decimal.NewFromInt(int64(p.Length))),
			"open_time": (*klines)[l-1].OpenTime,
		}
		sma := cryptodb.NewMovingAverage(maData)
		smas = append(smas, sma)
	}
	return
}
