CREATE TABLE `tests` (
  `id` int(10) UNSIGNED NOT NULL,
  `strategy_id` int(10) NOT NULL COMMENT 'Strategy Id',
  `strategy_title` varchar(100) NOT NULL COMMENT 'Strategy title',
  `strategy_desc` text NOT NULL,
  `strategy_type` varchar(30) NOT NULL COMMENT 'Strategy type',
  `strategy_params` text NOT NULL COMMENT 'Strategy params',
  `start` datetime NOT NULL COMMENT 'Period start time',
  `end` datetime NOT NULL COMMENT 'Period end time',
  `cost` decimal(18,8) UNSIGNED NOT NULL COMMENT 'Cost',
  `revenue` decimal(18,8) NOT NULL COMMENT 'Revenue',
  `fee` decimal(18,8) NOT NULL COMMENT 'Total trade fee',
  `profit` decimal(18,8) NOT NULL COMMENT 'Profit',
  `ROI` decimal(18,8) NOT NULL COMMENT 'Return on Investment',
  `trade_count` int(10) UNSIGNED NOT NULL COMMENT 'Total trade count',
  `comment` text NOT NULL COMMENT 'Comment',
  `created_at` datetime NOT NULL DEFAULT current_timestamp() COMMENT 'Create time'
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

ALTER TABLE `tests`
  ADD PRIMARY KEY (`id`),
  ADD KEY `strategy_id` (`strategy_id`);

ALTER TABLE `tests`
  MODIFY `id` int(10) UNSIGNED NOT NULL AUTO_INCREMENT;
