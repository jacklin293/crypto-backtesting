package backfill

import (
	"crypto-backtesting/cryptodb"
	"errors"
	"fmt"
	"log"
	"sync"
)

const (
	// NOTE: This must be bigger than MA length
	KLINE_BATCH_LIMIT = 2000
)

type Backfiller interface {
	backfill() error
}

type baseMA struct {
	db       *cryptodb.DB
	maType   string // MA type
	pair     string
	interval string
	length   int
}

func Start(db *cryptodb.DB, maType string, pair string, interval string) (err error) {
	if pair == "" || interval == "" {
		return errors.New("All pair, interval should be specified")
	}

	var wg sync.WaitGroup
	for length := 10; length <= 200; length += 10 {
		wg.Add(1)
		go func(wg *sync.WaitGroup, db *cryptodb.DB, pair string, interval string, length int) {
			defer wg.Done()

			ma, err := newMA(db, maType, pair, interval, length)
			if err != nil {
				return
			}

			switch maType {
			case "ema":
				if err = ma.backfill(); err != nil {
					log.Println(err)
				}
			case "sma":
				if err = ma.backfill(); err != nil {
					log.Println(err)
				}
			}
		}(&wg, db, pair, interval, length)
	}
	wg.Wait()

	return
}

func newMA(db *cryptodb.DB, maType string, pair string, interval string, length int) (b Backfiller, err error) {
	base := baseMA{
		db:       db,
		maType:   maType,
		pair:     pair,
		interval: interval,
		length:   length,
	}
	switch maType {
	case "ema":
		b = &Ema{baseMA: base}
	case "sma":
		b = &Sma{baseMA: base}
	default:
		err = fmt.Errorf("ma_type '%s' not supported", maType)
	}

	return
}
