# crypto-backtesting
Backtesting tool for cryptocurrency market

# How to use

Backfill EMA data

```
./crypto-backtesting -task=2 -pair=btcusdt -ma_type=ema -interval=4h
```

Backfill SMA data

```
./crypto-backtesting -task=2 -pair=btcusdt -ma_type=sma -interval=4h
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

Run backtesting for enabled strategies

```
./crypto-backtesting -task=1
```

Run backtesting for a specific strategy with `-strategyid`

```
./crypto-backtesting -task=1 -strategyid=3
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

Use the latest MA value and kline as baselines

* Buy when the price is above MA value and kline's high
* Sell when the price is below MA value and kline's low

Params example:

```
{
  "ma_type": "ema"
}
```


#### Strategy type - `future_contract`

Entry type: limit

* `flip_operator_enabled` (bool)
    * The situations that this param is recommended to enable are as follows:
        - buying the dip for long
        - selling at high for short
    * For example: the mark price is 48000 and you want to buy below 42000 and stop-loss below 41000. The problem is that when the price goes through 42000 and 41000 then fluctuates up and down around 41000. The entry and stop-loss orders will be triggerer constantly.
    * The solution for this is to enable this param that converts the `operator` of entry only once when it's first triggered and make entry `>= 42000` so it won't be triggered below 42000

```
{
  "position_type": "long",
  "entry_type": "limit",
  "entry_order": {
    "trigger": {
      "trigger_type": "limit",
      "operator": "<=",
      "price": 47200
    },
    "flip_operator_enabled": true
  },
  "stop_loss_order": {
    "trigger": {
      "trigger_type": "limit",
      "operator": "<=",
      "price": 47000
    }
  },
  "take_profit_order": {
    "trigger": {
      "trigger_type": "limit",
      "operator": ">=",
      "price": 50000
    }
  }
}
```

Entry type: baseline

* `baseline_offset_percent` (float) e.g. 0.01 denotes 1%
    * This param allows you to buy above or below the baseline
    * long
        - positive float: `above` baseline
        - negative float: `below` baseline
    * short
        - positive float: `below` baseline
        - negative float: `above` baseline
* `flip_operator_enabled` (bool)
* `loss_tolerance_percent` (float) e.g. 0.01 denotes 1%
    * The amount that you can accept for loss in each attempt
    * The baseline price is based on the real entry price as opposed to the price on entry trigger, for example, if the average cost of entry order is 30 and loss tolerance is 1%, then stop-loss will be 29.7 (30*0.99)
* `baseline_readjustment_enabled` (bool)
    * To re-adjust entry trigger to deal with the false breakout by updating `time_2` and `price_2` which are recorded for the highest price happended before stop-loss triggered. When the stop-loss triggered, it's treated as a false breakout. In order to avoid losing money, this param raises the threshold of entry.
    * For example, if the baseline is in a downtrend, the entry order is triggered when the price breaks out the line. Then the price goes down and trigger stop-loss order. (see the timeline below) The `time_2` and `price_2` of entry order will be updated with `10:05` and `130`.
        - 10:01: Entry order is triggered at $100
        - 10:05: The price reaches the highest point at $130
        - 10:07: Then the price goes down and stop-loss is triggered

```
{
  "position_type": "long",
  "entry_type": "baseline",
  "entry_order": {
    "baseline_trigger": {
      "trigger_type": "line",
      "operator": ">=",
      "time_1": "2021-08-18 18:00:00",
      "price_1": 46000,
      "time_2": "2021-08-19 01:45:00",
      "price_2": 45234.56
    },
    "baseline_offset_percent": 0.005
  },
  "stop_loss_order": {
    "loss_tolerance_percent": 0.005,
    "baseline_readjustment_enabled": true
  },
  "take_profit_order": {
    "trigger": {
      "trigger_type": "limit",
      "operator": ">=",
      "price": 46195
    }
  }
}
```
