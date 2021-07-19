CREATE TABLE `trades` (
  `id` int(11) NOT NULL,
  `test_id` int(11) NOT NULL COMMENT 'test id',
  `status` tinyint(4) NOT NULL DEFAULT 0 COMMENT '0: undone 1: done',
  `bid_price` decimal(18,8) UNSIGNED NOT NULL COMMENT 'Bid price',
  `bid_volume` decimal(18,8) UNSIGNED NOT NULL COMMENT 'Bid volume',
  `bid_fee` decimal(18,8) NOT NULL COMMENT 'Bid fee',
  `bought at` datetime NOT NULL COMMENT 'Bought time',
  `ask_price` decimal(18,8) UNSIGNED NOT NULL COMMENT 'Ask price',
  `ask_volume` decimal(18,8) UNSIGNED NOT NULL COMMENT 'Ask volume',
  `ask_fee` decimal(18,8) NOT NULL COMMENT 'Ask fee',
  `sold_at` datetime NOT NULL COMMENT 'Sold time',
  `cost` decimal(18,8) UNSIGNED NOT NULL COMMENT 'Cost',
  `revenue` decimal(18,8) NOT NULL COMMENT 'Revenue',
  `profit` decimal(18,8) NOT NULL COMMENT 'Profit',
  `ROI` decimal(18,8) NOT NULL COMMENT 'Return on investment',
  `created_at` datetime NOT NULL DEFAULT current_timestamp() COMMENT 'Create time',
  `updated_at` datetime NOT NULL DEFAULT current_timestamp() ON UPDATE current_timestamp() COMMENT 'Update time'
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

ALTER TABLE `trades`
  ADD PRIMARY KEY (`id`),
  ADD KEY `test_id` (`test_id`);

ALTER TABLE `trades`
  MODIFY `id` int(11) NOT NULL AUTO_INCREMENT;
