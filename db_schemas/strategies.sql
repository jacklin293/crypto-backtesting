CREATE TABLE `strategies` (
  `id` int(10) UNSIGNED NOT NULL,
  `title` varchar(100) NOT NULL COMMENT 'Title',
  `description` text NOT NULL COMMENT 'Description',
  `strategy_type` varchar(30) NOT NULL COMMENT 'Strategy type',
  `params` text CHARACTER SET utf8mb4 COLLATE utf8mb4_bin NOT NULL COMMENT 'Detailed params',
  `start` datetime NOT NULL COMMENT 'Period start time',
  `end` datetime NOT NULL COMMENT 'Period end time',
  `cost` decimal(18,8) UNSIGNED NOT NULL DEFAULT 1000.00000000 COMMENT 'Cost',
  `enabled` tinyint(4) NOT NULL DEFAULT 1 COMMENT '0:disabled 1:enabled',
  `created_at` datetime NOT NULL DEFAULT current_timestamp() COMMENT 'Create time'
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

ALTER TABLE `strategies`
  ADD PRIMARY KEY (`id`),
  ADD KEY `enabled_type` (`enabled`,`strategy_type`);

ALTER TABLE `strategies`
  MODIFY `id` int(10) UNSIGNED NOT NULL AUTO_INCREMENT;
