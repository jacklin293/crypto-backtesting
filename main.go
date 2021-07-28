package main

import (
	"crypto-backtesting/backfill"
	"crypto-backtesting/backtest"
	"crypto-backtesting/cryptodb"
	"flag"
	"fmt"
	"log"
	"os"
)

const (
	DB_DSN = "root:root@tcp(127.0.0.1:3306)/crypto_db?charset=utf8mb4&parseTime=true"
)

func main() {
	// Connect to DB
	db, err := cryptodb.NewDB(DB_DSN)
	if err != nil {
		log.Fatal(err)
	}

	// backtest
	task := flag.String("task", "", "e.g. 0: Run backtesting  1: Backfill EMA data")
	// backtest a specific strategy
	strategyid := flag.Int("strategyid", 0, "Strategy id for testing e.g. 5")

	// backfill
	maType := flag.String("ma_type", "", "e.g. ema sma")
	pair := flag.String("pair", "", "Moving Average param e.g. BTCUSDT")
	interval := flag.String("interval", "", "Moving Average param e.g. 1m 1h 4h")
	flag.Parse()

	switch *task {
	// Run backtesting
	case "1":
		backtest.Start(db, int64(*strategyid))
	// Backfill EMA data
	case "2":
		if err := backfill.Start(db, *maType, *pair, *interval); err != nil {
			log.Fatal(err)
		}
	default:
		fmt.Println("Please choose a task. Print usage with '-h'")
		os.Exit(0)
	}
}
