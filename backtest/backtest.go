package backtest

import (
	"crypto-backtesting/cryptodb"
)

type Params struct {
	Db *cryptodb.DB
}

func (p *Params) HandleBacktesting(db *cryptodb.DB) (err error) {
	// TODO Read table strategies and test
	if err = p.handleEmaLastKline(db); err != nil {
		return
	}
	return nil
}
