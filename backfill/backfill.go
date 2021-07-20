package backfill

import (
	"crypto-backtesting/cryptodb"
	"errors"
	"fmt"
	"log"
)

const (
	// TODO
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

func Start(db *cryptodb.DB, maType string, pair string, interval string, length int) (err error) {
	if pair == "" || interval == "" || length == 0 {
		return errors.New("All pair, interval and length should be specified")
	}

	if length > KLINE_BATCH_LIMIT {
		// TODO const name
		fmt.Errorf("KLINE_BATCH_LIMIT(%d) must be bigger than length(%d)", KLINE_BATCH_LIMIT, length)
	}

	ma, err := newMA(db, maType, pair, interval, length)
	if err != nil {
		return
	}

	switch maType {
	case "ema":
		if err = ma.backfill(); err != nil {
			log.Fatal(err)
		}
	case "sma":
		if err = ma.backfill(); err != nil {
			log.Fatal(err)
		}
	}

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
		err = fmt.Errorf("matype '%s' not supported", maType)
	}

	return
}
