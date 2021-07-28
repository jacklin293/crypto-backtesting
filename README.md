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

Create a strategy in the table `strategies`, e.g.

* strategy_type: `ma_and_loss_tolerance`
* pair: `btcusdt`
* interval: `1d`
* params: `{"ma_type":"ema","loss_tolerance":0.03}`
* start: `2021-06-01`
* end: `2021-06-30`
* cost: `1000`
* enabled: `1`

Run backtesting

```
go build && ./crypto-backtesting -task=1
```

Run backtesting for a specific strategy with `-strategyid`

```
go build && ./crypto-backtesting -task=1 -strategyid=3
```


# Strategy Types and Params

#### Strategy type - `ma_and_loss_tolerance`

Use the latest MA as baseline

* Buy when the price is above MA value
* Sell when the price is below MA value by percent of loss_tolerance

Params example:

```
{
  "ma_type": "ema",
  "loss_tolerance": 0.03
}
```
> For above example, buy the coin at 1% above MA and sell at 1% below MA

#### Strategy type - `ma_and_latest_kline`

Use the latest MA value and kline as baselines.

* Buy when the price is above MA value and kline's high.
* Sell when the price is below MA value and kline's low

```
{
  "ma_type": "ema"
}
```


