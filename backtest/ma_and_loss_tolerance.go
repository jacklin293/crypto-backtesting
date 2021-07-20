package backtest

import (
	"crypto-backtesting/cryptodb"
	"crypto-backtesting/utils"
	"fmt"
	"strings"
	"time"

	"github.com/shopspring/decimal"
)

type maAndLossTolerance struct {
	baseStrategy
	params maAndLossToleranceParams
}

type maAndLossToleranceParams struct {
	lossTolerance float64
}

func (s *maAndLossTolerance) backtest() (err error) {
	startTime := time.Now()

	lengthMins, err := utils.ConvertIntervalToMins(s.interval)
	if err != nil {
		return
	}

	// Initial the first time block
	blockStart, blockEnd, err := utils.GetTimeBlockByLength(s.start, lengthMins)
	if err != nil {
		return err
	}

	// Use 1m-kline to simulate prices changed over time
	var minKlineStart, minKlineEnd time.Time

	// Start to go over each time block
	for blockStart.Before(s.end) {
		minKlineStart = blockStart
		// The period start might be in the middle of time block, for more info, please have a look at `utils.GetTimeBlockByLength`
		if minKlineStart.Before(s.start) {
			// If block start is earlier than period start, min-kline start should be period start, because klines earlier than period start aren't needed
			// This will only be satisfied in the beginning
			minKlineStart = s.start
		}
		minKlineEnd = blockEnd.Add(time.Minute * time.Duration(-1))
		if minKlineEnd.After(s.end) {
			// If block end is later than period end, min-kline end should be period end, because klines later than period end aren't needed
			// This will only be satisfied in the beginning
			minKlineEnd = s.end
		}

		// Current length-kline and MA value are ongoing, use previous ones as baseline
		previousOpenTime := blockStart.Add(time.Minute * time.Duration(-lengthMins))

		// Get previous MA by length
		baselineMA, count, err := s.db.GetMovingAveragesByOpenTime(s.maType, s.pair, s.interval, s.length, previousOpenTime)
		if err != nil {
			return err
		}
		if count == 0 {
			return fmt.Errorf("baseline-%s(%d) not found at %v", strings.ToUpper(s.maType), s.length, previousOpenTime)
		}

		if err = s.checkPricesInTimeBlock(minKlineStart, minKlineEnd, baselineMA); err != nil {
			return err
		}

		// Add lenght mins for nex time block
		blockStart = blockStart.Add(time.Minute * time.Duration(lengthMins))
		blockEnd = blockEnd.Add(time.Minute * time.Duration(lengthMins))
	}

	roi := s.test.marketValue.Sub(s.test.cost).Div(s.test.cost).Mul(decimal.NewFromInt(100))
	roiPerTrade := decimal.NewFromInt(0)
	if s.test.tradeCount != 0 {
		roiPerTrade = roi.Div(decimal.NewFromInt(int64(s.test.tradeCount)))
	}
	fmt.Printf("'%s %s%d %s %s' $%s => $%s (%s%%) tradeCount: %d (%s%%/t) (%s ~ %s) %s\n", s.pair, s.maType, s.length, s.interval, s.strategyType, s.test.cost.StringFixed(0), s.test.marketValue.StringFixed(0), roi.StringFixed(1), s.test.tradeCount, roiPerTrade.StringFixed(1), s.start.Format("2006-01-02"), s.end.Format("2006-01-02"), time.Since(startTime))

	return
}

func (s *maAndLossTolerance) checkPricesInTimeBlock(tStart time.Time, tEnd time.Time, baselineMA *cryptodb.MovingAverage) (err error) {
	var klineCount int64
	var minKlines *[]cryptodb.Kline
	nextStart := tStart

	// Batch fetch min-kline until no more
	for {
		// FIXME bottleneck
		minKlines, klineCount, err = s.db.GetKlinesByPeriod(s.pair, "1m", nextStart, tEnd, MIN_KLINE_BATCH_LIMIT, "ASC")
		if err != nil {
			return
		}
		if klineCount == 0 {
			break
		}

		// Check min klines with baseline length-kline and baseline length-MA
		if err = s.checkPricesWithMinKlines(minKlines, baselineMA); err != nil {
			return
		}
		nextStart = (*minKlines)[klineCount-1].OpenTime.Add(time.Minute * 1)
	}

	return
}

func (s *maAndLossTolerance) checkPricesWithMinKlines(klines *[]cryptodb.Kline, baselineMA *cryptodb.MovingAverage) error {
	sellPercent := decimal.NewFromFloat(float64(1) - s.params.lossTolerance)
	for _, kline := range *klines {
		switch s.trade.status {
		case "short", "waiting":
			/*
				-1 if d <  d2
				 0 if d == d2
				+1 if d >  d2
			*/
			if kline.Close.Cmp(baselineMA.Value) >= 0 {
				// Buy
				s.trade.cost = s.test.marketValue
				s.trade.status = "long"
				s.trade.bidPrice = kline.Close
				s.trade.bidVolume = s.trade.cost.Div(kline.Close)
				s.test.tradeCount++
				// fmt.Printf("%s [B] %s at %s\n", kline.OpenTime.Format("2006-01-02 15:04"), s.trade.cost.StringFixed(2), s.trade.bidPrice.StringFixed(2))
			}
		case "long":
			if kline.Close.Cmp(baselineMA.Value.Mul(sellPercent)) <= 0 {
				// Sell
				s.trade.status = "short"
				s.trade.revenue = kline.Close.Mul(s.trade.bidVolume)
				s.trade.returnPercent = s.trade.revenue.Sub(s.trade.cost).Div(s.trade.cost)
				s.test.marketValue = s.trade.revenue
				s.test.tradeCount++
				// fmt.Printf("%s [S] %s at %s (%s%%)\n\n", kline.OpenTime.Format("2006-01-02 15:04"), s.trade.revenue.StringFixed(2), kline.Close.StringFixed(2), s.trade.returnPercent.Mul(decimal.NewFromInt(100)).StringFixed(1))
			}
		}
	}

	return nil
}
