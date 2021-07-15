# crypto-backtesting
Backtesting tool for cryptocurrency market

# How to use

Run backtesting

```
go build && ./crypto-backtesting -task=1 -interval=1d -start=2020-10-01 -end=2021-06-30 -length=100
```

Backfill MA data

```
go build && ./crypto-backtesting -task=2 -pair=btcusdt -matype=ema -interval=4h -length=18
```
