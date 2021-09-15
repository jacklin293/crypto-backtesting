package backtest

import (
	"crypto-backtesting/cryptodb"
	"crypto-backtesting/market/future"
	"crypto-backtesting/market/order"
	"fmt"
	"log"
	"time"

	"github.com/shopspring/decimal"
)

type strategyFutureContract struct {
	baseStrategy
	params strategyFutureContractParams
}

type strategyFutureContractParams struct {
	contract *future.Contract
}

type contractHook struct {
	trade *trade
	test  *test
}

func (s *strategyFutureContract) backtest() (err error) {
	// New contract hook
	ch := &contractHook{
		trade: &s.trade,
		test:  &s.test,
	}
	s.params.contract.SetHook(ch)

	var klineCount int64
	var minKlines *[]cryptodb.Kline
	var halted bool               // if it's true, end the test
	var lastPrice decimal.Decimal // For calculating the market value if the order is still opened

	nextStart := s.start // Set the start time
	startTime := time.Now()
	for {
		minKlines, klineCount, err = s.db.GetKlinesByPeriod(s.pair, "1m", nextStart, s.end, MIN_KLINE_BATCH_LIMIT, "ASC")
		if err != nil {
			return
		}
		if klineCount == 0 {
			break
		}

		for _, kline := range *minKlines {
			err, halted = s.params.contract.CheckPrice(kline.OpenTime, kline.Close)
			if err != nil {
				log.Println("checkPrices err:", err)
			}
			if halted {
				break
			}
		}
		if halted {
			break
		}

		nextStart = (*minKlines)[klineCount-1].OpenTime.Add(time.Minute * 1)
		lastPrice = (*minKlines)[klineCount-1].Close
	}

	// Calculate market value based on current price
	if s.params.contract.Status == future.OPENED {
		ch.test.marketValue = ch.trade.bidVolume.Mul(lastPrice)
	}

	roi := ch.test.marketValue.Sub(ch.test.cost).Div(ch.test.cost).Mul(decimal.NewFromInt(100))
	fmt.Println("-----------------------")
	fmt.Printf("'%s %s' $%s => $%s (%s%%) (%s ~ %s) %s\n", s.pair, s.strategyType, s.test.cost.StringFixed(0), s.test.marketValue.StringFixed(0), roi.StringFixed(1), s.start.Format("2006-01-02"), s.end.Format("2006-01-02"), time.Since(startTime))

	return nil
}

func (ch *contractHook) EntryTriggered(c *future.Contract, t time.Time, p decimal.Decimal) (decimal.Decimal, error, bool) {
	fmt.Printf("EntryTriggered             baseline: %s  %s\n", c.EntryOrder.(*order.Entry).BaselineTrigger.GetPrice(t).StringFixed(2), c.EntryOrder.(*order.Entry).BaselineTrigger)
	fmt.Printf("EntryTriggered                entry: %s  %s\n", c.EntryOrder.GetTrigger().GetPrice(t).StringFixed(2), c.EntryOrder.GetTrigger())
	fmt.Printf("EntryTriggered                  buy: %s  '%s'\n", p, t.Format("2006-01-02 15:04"))
	ch.trade.cost = ch.test.marketValue
	ch.trade.bidPrice = p
	ch.trade.bidVolume = ch.trade.cost.Div(p)
	return p, nil, false
}

func (ch *contractHook) StopLossTriggerCreated(c *future.Contract) (error, bool) {
	fmt.Printf("StopLossTriggerCreated    stop-loss: %s  %s\n", c.StopLossOrder.GetTrigger().GetPrice(time.Now()).StringFixed(2), c.StopLossOrder.GetTrigger().GetOperator())
	return nil, false
}

func (ch *contractHook) StopLossTriggered(c *future.Contract, t time.Time, p decimal.Decimal) error {
	ch.trade.revenue = p.Mul(ch.trade.bidVolume)
	ch.trade.returnPercent = ch.trade.revenue.Sub(ch.trade.cost).Div(ch.trade.cost)
	ch.test.marketValue = ch.trade.revenue

	fmt.Printf("! StopLossTriggered            sell: %s  '%s'  ($%s => $%s)\n", p, t.Format("2006-01-02 15:04"), ch.trade.cost.StringFixed(0), ch.trade.revenue.StringFixed(0))
	return nil
}

func (ch *contractHook) EntryBaselineTriggerUpdated(c *future.Contract) {
	fmt.Println("- EntryBaselineTriggerUpdated breakout:", c.BreakoutPeak)
	fmt.Println("- EntryBaselineTriggerUpdated baseline:", c.EntryOrder.(*order.Entry).BaselineTrigger)
	fmt.Println("- EntryBaselineTriggerUpdated    entry:", c.EntryOrder.GetTrigger())
	fmt.Println("-----------------------")
}

func (ch *contractHook) TakeProfitTriggered(c *future.Contract, t time.Time, p decimal.Decimal) error {
	ch.trade.revenue = p.Mul(ch.trade.bidVolume)
	ch.trade.returnPercent = ch.trade.revenue.Sub(ch.trade.cost).Div(ch.trade.cost)
	ch.test.marketValue = ch.trade.revenue

	fmt.Printf("! TakeProfitTriggered          sell: %s  '%s'  ($%s => $%s)\n", p, t.Format("2006-01-02 15:04"), ch.trade.cost.StringFixed(0), ch.trade.revenue.StringFixed(0))
	return nil
}

func (ch *contractHook) OrderTriggerUpdated(c *future.Contract) {
}

func (ch *contractHook) StatusChanged(c *future.Contract) {
	fmt.Println("current status:", c.Status)
}
