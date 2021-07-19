package backtest

import (
	"crypto-backtesting/cryptodb"
	"crypto-backtesting/utils"
	"fmt"
	"time"

	"github.com/shopspring/decimal"
)

// FIXME
var status string = "buy"
var bidPrice, bidVolume decimal.Decimal

var seedMoney decimal.Decimal = decimal.NewFromFloat(1000)

var cost decimal.Decimal = seedMoney
var balance decimal.Decimal
var revenue, returnPercent decimal.Decimal

var tradeCount int

type emaLastKline struct {
	db          *cryptodb.DB
	maType      string
	pair        string
	interval    string
	length      int
	periodStart time.Time
	periodEnd   time.Time
}

func (p *Params) handleEmaLastKline(db *cryptodb.DB) (err error) {
	// FIXME
	maType := "ema"
	pair := "btcusdt"
	interval := "4h"
	length := 18
	dateStart := "2020-10-01"
	dateEnd := "2021-06-30"

	periodStart, err := time.Parse("2006-01-02 15:04:05", dateStart+" 00:00:00")
	if err != nil {
		return
	}
	periodEnd, err := time.Parse("2006-01-02 15:04:05", dateEnd+" 23:59:59")
	if err != nil {
		return
	}

	// FIXME
	elk := emaLastKline{
		db:          db,
		maType:      maType,
		pair:        pair,
		interval:    interval,
		length:      length,
		periodStart: periodStart,
		periodEnd:   periodEnd,
	}

	lengthMins, err := utils.ConvertIntervalToMins(elk.interval)
	if err != nil {
		return
	}

	if err = elk.initialiseTimeBlocks(lengthMins); err != nil {
		return
	}

	return nil
}

func (elk *emaLastKline) initialiseTimeBlocks(lengthMins int) (err error) {
	// Initial time block
	blockStart, blockEnd, err := elk.getTimeBlockByLength(elk.periodStart, lengthMins)
	if err != nil {
		return err
	}

	// Use 1m klines as prices changed over time
	var minLengthStart, minLengthEnd time.Time

	// Start to go over each time block
	for !blockStart.After(elk.periodEnd) {
		minLengthStart = blockStart
		if minLengthStart.Before(elk.periodStart) {
			minLengthStart = elk.periodStart
		}
		minLengthEnd = blockEnd.Add(time.Minute * time.Duration(-1))
		if minLengthEnd.After(elk.periodEnd) {
			minLengthEnd = elk.periodEnd
		}

		// Current block of kline and ema are ongoing until this block ends. Use previous kline and ema as baseline
		previousOpenTime := blockStart.Add(time.Minute * time.Duration(-lengthMins))

		// Get previous kline by length
		baselineKline, count, err := elk.db.GetKlineByOpenTime(elk.pair, elk.interval, previousOpenTime)
		if err != nil {
			return err
		}
		if count == 0 {
			return fmt.Errorf("No baseline-kline row found at %v", previousOpenTime)
		}

		// Get previous ema by length
		baselineEma, count, err := elk.db.GetMovingAveragesByOpenTime(elk.maType, elk.pair, elk.interval, elk.length, previousOpenTime)
		if err != nil {
			return err
		}
		if count == 0 {
			return fmt.Errorf("No baseline-EMA(%d) row found at %v", elk.length, previousOpenTime)
		}

		if err = elk.checkPriceDuringTimeBlock(minLengthStart, minLengthEnd, baselineKline, baselineEma); err != nil {
			return err
		}

		// Add lenght mins for nex time block
		blockStart = blockStart.Add(time.Minute * time.Duration(lengthMins))
		blockEnd = blockEnd.Add(time.Minute * time.Duration(lengthMins))
	}

	fmt.Printf("\n%s ~ %s result: %s -> %s (%s%%) tradeCount: %d\n\n", elk.periodStart.Format("2006-01-02"), elk.periodEnd.Format("2006-01-02"), seedMoney.StringFixed(0), balance.StringFixed(0), balance.Sub(seedMoney).Div(seedMoney).Mul(decimal.NewFromInt(100)).StringFixed(1), tradeCount)
	return
}

func (elk *emaLastKline) checkPriceDuringTimeBlock(tStart time.Time, tEnd time.Time, baselineKline *cryptodb.Kline, baselineEma *cryptodb.MovingAverage) (err error) {
	var count int64
	var minKlines *[]cryptodb.Kline
	nextStart := tStart

	// FIXME
	batchLimit := 200

	for {
		minKlines, count, err = elk.db.GetKlinesByPeriod(elk.pair, "1m", nextStart, tEnd, batchLimit, "ASC")
		if err != nil {
			return
		}
		if count == 0 {
			break
		}

		if err = elk.checkEmaAndLastKline(minKlines, baselineKline, baselineEma); err != nil {
			return
		}
		nextStart = (*minKlines)[count-1].OpenTime.Add(time.Minute * 1)

		if count < int64(batchLimit) {
			continue
		}
	}

	return
}

func (elk *emaLastKline) checkEmaAndLastKline(klines *[]cryptodb.Kline, baselineKline *cryptodb.Kline, baselineEma *cryptodb.MovingAverage) error {
	for _, kline := range *klines {
		switch status {
		case "buy":
			/*
				-1 if d <  d2
				 0 if d == d2
				+1 if d >  d2
			*/
			if kline.Close.Cmp(baselineKline.High) == 1 && kline.Close.Cmp(baselineEma.Value) == 1 {
				if balance.Cmp(decimal.NewFromInt(0)) == 0 {
					cost = seedMoney
				} else {
					cost = balance
				}
				status = "sell"
				bidPrice = kline.Close
				bidVolume = cost.Div(kline.Close)
				tradeCount++
				fmt.Printf("%s [B] %s at %s\n", kline.OpenTime.Format("2006-01-02 15:04"), cost.StringFixed(2), bidPrice.StringFixed(2))
			}
		case "sell":
			if kline.Close.Cmp(baselineKline.Low) == -1 && kline.Close.Cmp(baselineEma.Value) == -1 {
				status = "buy"
				revenue = kline.Close.Mul(bidVolume)
				returnPercent = revenue.Sub(cost).Div(cost)
				balance = revenue
				tradeCount++
				fmt.Printf("%s [S] %s at %s (%s%%)\n\n", kline.OpenTime.Format("2006-01-02 15:04"), revenue.StringFixed(2), kline.Close.StringFixed(2), returnPercent.Mul(decimal.NewFromInt(100)).StringFixed(1))
			}
		}
	}
	return nil
}

// Get time block based on length with time point given
// For example, if time pint is 10:00 and lenght is 4h, the time block would be between 08:00-12:00
func (elk *emaLastKline) getTimeBlockByLength(t time.Time, lengthMins int) (tStart time.Time, tEnd time.Time, err error) {
	// Get the minutes of time
	mins := t.Hour()*60 + t.Minute()
	timeBlocks := mins / lengthMins
	tStart = time.Date(t.Year(), t.Month(), t.Day(), 0, timeBlocks*lengthMins, 0, 0, time.UTC)
	tEnd = time.Date(t.Year(), t.Month(), t.Day(), 0, (timeBlocks+1)*lengthMins, 0, 0, time.UTC)
	return
}
