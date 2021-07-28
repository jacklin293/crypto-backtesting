-- phpMyAdmin SQL Dump
-- version 5.1.1
-- https://www.phpmyadmin.net/
--
-- Host: db
-- Generation Time: Jul 28, 2021 at 02:56 PM
-- Server version: 10.6.2-MariaDB-1:10.6.2+maria~focal
-- PHP Version: 7.4.20

SET SQL_MODE = "NO_AUTO_VALUE_ON_ZERO";
START TRANSACTION;
SET time_zone = "+00:00";

--
-- Database: `crypto_db`
--

-- --------------------------------------------------------

--
-- Table structure for table `strategies`
--

CREATE TABLE `strategies` (
  `id` int(10) UNSIGNED NOT NULL,
  `strategy_type` varchar(60) NOT NULL COMMENT 'Strategy type',
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

--
-- Indexes for dumped tables
--

--
-- Indexes for table `strategies`
--
ALTER TABLE `strategies`
  ADD PRIMARY KEY (`id`),
  ADD KEY `enabled_type` (`enabled`,`strategy_type`);

--
-- AUTO_INCREMENT for dumped tables
--

--
-- AUTO_INCREMENT for table `strategies`
--
ALTER TABLE `strategies`
  MODIFY `id` int(10) UNSIGNED NOT NULL AUTO_INCREMENT;
COMMIT;

