# crypto-backtesting
Backtesting tool for cryptocurrency market

# How to use

Backfill EMA data

```
go build && ./crypto-backtesting -task=2 -pair=btcusdt -matype=ema -interval=4h -length=18
```

Backfill SMA data

```
go build && ./crypto-backtesting -task=2 -pair=btcusdt -matype=sma -interval=4h -length=18
```

Create a new strategy in the table `strategies`, e.g.

// TODO

Run backtesting

```
go build && ./crypto-backtesting -task=1
```
