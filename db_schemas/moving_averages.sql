-- MySQL dump 10.13  Distrib 8.0.19, for osx10.15 (x86_64)
--
-- Host: 127.0.0.1    Database: crypto_db
-- ------------------------------------------------------
-- Server version	5.5.5-10.6.2-MariaDB-1:10.6.2+maria~focal

/*!40101 SET @OLD_CHARACTER_SET_CLIENT=@@CHARACTER_SET_CLIENT */;
/*!40101 SET @OLD_CHARACTER_SET_RESULTS=@@CHARACTER_SET_RESULTS */;
/*!40101 SET @OLD_COLLATION_CONNECTION=@@COLLATION_CONNECTION */;
/*!50503 SET NAMES utf8mb4 */;
/*!40103 SET @OLD_TIME_ZONE=@@TIME_ZONE */;
/*!40103 SET TIME_ZONE='+00:00' */;
/*!40014 SET @OLD_UNIQUE_CHECKS=@@UNIQUE_CHECKS, UNIQUE_CHECKS=0 */;
/*!40014 SET @OLD_FOREIGN_KEY_CHECKS=@@FOREIGN_KEY_CHECKS, FOREIGN_KEY_CHECKS=0 */;
/*!40101 SET @OLD_SQL_MODE=@@SQL_MODE, SQL_MODE='NO_AUTO_VALUE_ON_ZERO' */;
/*!40111 SET @OLD_SQL_NOTES=@@SQL_NOTES, SQL_NOTES=0 */;

--
-- Table structure for table `moving_averages`
--

DROP TABLE IF EXISTS `moving_averages`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `moving_averages` (
  `ma_key` varchar(20) NOT NULL COMMENT 'ma type+pair+interval e.g. ema_btcusdt_1h',
  `length` smallint(5) unsigned NOT NULL COMMENT 'length of klines',
  `value` decimal(18,8) unsigned NOT NULL COMMENT 'Value',
  `open_time` datetime NOT NULL COMMENT 'Open time',
  UNIQUE KEY `makey_length_opentime` (`ma_key`,`length`,`open_time`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='Moving Averages';
/*!40101 SET character_set_client = @saved_cs_client */;
