package main

import (
	"crypto-backtesting/cryptodb"
	"flag"
	"fmt"
	"log"
	"os"
)

const (
	DB_DSN                        = "root:root@tcp(127.0.0.1:3306)/crypto_db?charset=utf8mb4&parseTime=true"
	DB_KLINES_BATCH_SELECT_NUMBER = 2000
)

func main() {
	// Connect to DB
	db, err := cryptodb.NewDB(DB_DSN)
	if err != nil {
		log.Fatal(err)
	}

	task := flag.String("task", "", "e.g. 0: Run backtesting  1: Backfill EMA data")
	pair := flag.String("pair", "", "Moving Average param e.g. BTCUSDT")
	interval := flag.String("interval", "", "Moving Average param e.g. 1m, 1h 4h")
	length := flag.Int("length", 0, "Moving Average param e.g. 18")
	dateStart := flag.String("start", "", "e.g. 2020-10-01")
	dateEnd := flag.String("end", "", "e.g. 2021-06-30")
	flag.Parse()

	switch *task {
	// Run backtesting
	case "1":
		// TODO put all params into struct
		if err = handleBacktesting(db, *interval, *length, *dateStart, *dateEnd); err != nil {
			log.Fatal(err)
		}
	// Backfill EMA data
	case "2":
		if err = handleBackfillEma(db, "ema", *pair, *interval, *length); err != nil {
			log.Fatal(err)
		}
	default:
		fmt.Println("Please choose a task. Print usage with '-h'")
		os.Exit(0)
	}
}
