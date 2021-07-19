CREATE TABLE `moving_averages` (
  `ma_key` varchar(20) NOT NULL COMMENT 'ma type+pair+interval e.g. ema_btcusdt_1h',
  `length` smallint(5) UNSIGNED NOT NULL COMMENT 'length of klines',
  `value` decimal(18,8) UNSIGNED NOT NULL COMMENT 'Value',
  `open_time` datetime NOT NULL COMMENT 'Open time'
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='Moving Averages';

ALTER TABLE `moving_averages`
  ADD UNIQUE KEY `ma_key_length_opentime` (`ma_key`,`length`,`open_time`);
