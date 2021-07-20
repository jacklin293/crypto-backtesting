CREATE TABLE `strategies` (
  `id` int(10) UNSIGNED NOT NULL,
  `strategy_type` varchar(60) NOT NULL COMMENT 'Strategy type',
  `ma_type` varchar(10) NOT NULL COMMENT 'MA type',
  `pair` varchar(15) NOT NULL COMMENT 'Pair',
  `interval` varchar(3) NOT NULL COMMENT 'Interval',
  `params` longtext CHARACTER SET utf8mb4 COLLATE utf8mb4_bin NOT NULL COMMENT 'Detailed params' CHECK (json_valid(`params`)),
  `start` date NOT NULL COMMENT 'Period start time',
  `end` date NOT NULL COMMENT 'Period end time',
  `cost` decimal(18,8) UNSIGNED NOT NULL DEFAULT 1000.00000000 COMMENT 'Cost',
  `enabled` tinyint(4) NOT NULL DEFAULT 1 COMMENT '0:disabled 1:enabled',
  `created_at` datetime NOT NULL DEFAULT current_timestamp() COMMENT 'Create time',
  `updated_at` datetime NOT NULL DEFAULT current_timestamp() ON UPDATE current_timestamp()
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

ALTER TABLE `strategies`
  ADD PRIMARY KEY (`id`),
  ADD KEY `enabled_type` (`enabled`,`strategy_type`);

ALTER TABLE `strategies`
  MODIFY `id` int(10) UNSIGNED NOT NULL AUTO_INCREMENT;
