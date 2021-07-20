package backtest

import (
	"crypto-backtesting/cryptodb"
	"errors"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/shopspring/decimal"
)

const (
	MIN_KLINE_BATCH_LIMIT = 1000 // 1-min kline batch limit

	// FIXME
	COST = 1000
)

type Backtester interface {
	backtest() error
}

type baseStrategy struct {
	db           *cryptodb.DB
	strategyType string
	maType       string
	pair         string
	interval     string
	length       int
	start        time.Time
	end          time.Time
	trade        trade
	test         test
}

type trade struct {
	status        string
	cost          decimal.Decimal
	bidPrice      decimal.Decimal
	bidVolume     decimal.Decimal
	revenue       decimal.Decimal
	returnPercent decimal.Decimal
}

type test struct {
	cost        decimal.Decimal
	marketValue decimal.Decimal
	tradeCount  int
}

func Start(db *cryptodb.DB, strategyId int64) {
	strategies := &[]cryptodb.Strategy{}
	var count int64
	var err error
	if strategyId == 0 {
		strategies, count, err = db.GetAllEnabledStrategies()
		if err != nil {
			return
		}
		if count == 0 {
			fmt.Println("There is no row in the table 'strategies'")
			return
		}
	} else {
		// For testing
		strategy, count, err := db.GetStrategyById(strategyId)
		if err != nil {
			return
		}
		if count == 0 {
			fmt.Println("Can't find strategy by id:", strategyId)
			return
		}
		*strategies = append(*strategies, *strategy)
	}

	var wg sync.WaitGroup
	for _, strategy := range *strategies {
		for length := 20; length <= 180; length += 20 {
			wg.Add(1)
			go func(wg *sync.WaitGroup, strategy cryptodb.Strategy, length int) {
				defer wg.Done()

				s, err := newStrategy(db, &strategy, length)
				if err != nil {
					log.Fatalf("Strategy id '%d' failed, err: %v\n", strategy.Id, err)
				} else {
					if err = s.backtest(); err != nil {
						fmt.Printf("Strategy id: %d failed, err: %v\n", strategy.Id, err)
					}
				}
			}(&wg, strategy, length)
		}
	}
	wg.Wait()
}

func newStrategy(db *cryptodb.DB, strategy *cryptodb.Strategy, length int) (s Backtester, err error) {
	base := baseStrategy{
		db:           db,
		strategyType: strategy.StrategyType,
		maType:       strategy.MaType,
		pair:         strategy.Pair,
		interval:     strategy.Interval,
		length:       length,
		start:        strategy.Start,
		end:          strategy.End.AddDate(0, 0, 1).Add(-time.Second),
		test: test{
			cost:        decimal.NewFromFloat(COST),
			marketValue: decimal.NewFromFloat(COST),
		},
		trade: trade{
			status: "waiting",
		},
	}

	switch strategy.StrategyType {
	case "ma_and_loss_tolerance":
		lossTolerance, ok := strategy.Params["loss_tolerance"].(float64)
		if !ok {
			err = errors.New("'loss_tolerance' is missing in params or not a float")
			return
		}
		s = &maAndLossTolerance{
			baseStrategy: base,
			params: maAndLossToleranceParams{
				lossTolerance: lossTolerance,
			},
		}
	case "ma_and_latest_kline":
		s = &maAndLastKline{
			baseStrategy: base,
		}
	default:
		err = fmt.Errorf("strategy_type '%s' not supported", strategy.StrategyType)
	}

	return
}
