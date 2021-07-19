package backfill

import (
	"crypto-backtesting/cryptodb"
)

type Params struct {
	Db       *cryptodb.DB
	MaType   string // MA type
	Pair     string
	Interval string
	Length   int
}
