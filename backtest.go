package main

import (
	"crypto-backtesting/cryptodb"
	"crypto-backtesting/strategies"
)

func handleBacktesting(db *cryptodb.DB, interval string, length int, dateStart string, dateEnd string) (err error) {
	// TODO Read table strategies and test
	if err = strategies.HandleEmaLastKline(db, interval, length, dateStart, dateEnd); err != nil {
		return
	}
	return nil
}
