-- Host: 127.0.0.1    Database: gogin
-- ------------------------------------------------------
-- Server version	5.7.30

/*!40101 SET @OLD_CHARACTER_SET_CLIENT=@@CHARACTER_SET_CLIENT */;
/*!40101 SET @OLD_CHARACTER_SET_RESULTS=@@CHARACTER_SET_RESULTS */;
/*!40101 SET @OLD_COLLATION_CONNECTION=@@COLLATION_CONNECTION */;
/*!40101 SET NAMES utf8 */;
/*!40103 SET @OLD_TIME_ZONE=@@TIME_ZONE */;
/*!40103 SET TIME_ZONE='+00:00' */;
/*!40014 SET @OLD_UNIQUE_CHECKS=@@UNIQUE_CHECKS, UNIQUE_CHECKS=0 */;
/*!40014 SET @OLD_FOREIGN_KEY_CHECKS=@@FOREIGN_KEY_CHECKS, FOREIGN_KEY_CHECKS=0 */;
/*!40101 SET @OLD_SQL_MODE=@@SQL_MODE, SQL_MODE='NO_AUTO_VALUE_ON_ZERO' */;
/*!40111 SET @OLD_SQL_NOTES=@@SQL_NOTES, SQL_NOTES=0 */;


--
-- DATABASE gogin
--
DROP DATABASE IF EXISTS `gogin`;

CREATE DATABASE /*!32312 IF NOT EXISTS*/ `gogin` /*!40100 DEFAULT CHARACTER SET utf8 */;

USE `gogin`;


--
-- Table structure for table `t_user`
--

DROP TABLE IF EXISTS `t_user`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `t_user` (
  `id` int(11) NOT NULL AUTO_INCREMENT COMMENT 'user id',
  `first_name` varchar(20) COLLATE utf8_unicode_ci NOT NULL COMMENT 'first name',
  `last_name` varchar(20) COLLATE utf8_unicode_ci NOT NULL COMMENT 'last name',
  `email` varchar(50) COLLATE utf8_unicode_ci NOT NULL COMMENT 'mail address',
  `password` varchar(50) COLLATE utf8_unicode_ci NOT NULL COMMENT 'password',
  `oauth2_flg` char(1) COLLATE utf8_unicode_ci DEFAULT '0' COMMENT 'oauth flg',
  `delete_flg` char(1) COLLATE utf8_unicode_ci DEFAULT '0' COMMENT 'delete flg',
  `created_at` datetime DEFAULT CURRENT_TIMESTAMP COMMENT 'created date',
  `updated_at` datetime DEFAULT CURRENT_TIMESTAMP COMMENT 'updated date',
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_email` (`email`)
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci COMMENT='user table';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `t_user`
--

LOCK TABLES `t_user` WRITE;
/*!40000 ALTER TABLE `t_user` DISABLE KEYS */;
INSERT INTO `t_user` VALUES (1,'foo','bar','foobar@gogin.com','baa62a499e9b21940c2d763f58a25647','0','0','2021-01-10 21:43:15','2021-01-10 21:43:15'),
                             (2,'mark','harry','mark@gogin.com','d978eb967fbe04345371478a97f3c903','0','0','2021-01-11 20:20:28','2021-01-11 20:20:28');
/*!40000 ALTER TABLE `t_user` ENABLE KEYS */;
UNLOCK TABLES;
/*!40103 SET TIME_ZONE=@OLD_TIME_ZONE */;


--
-- Table structure for table `t_book`
--

DROP TABLE IF EXISTS `t_book`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `t_book` (
  `id` int(11) NOT NULL AUTO_INCREMENT COMMENT 'id',
  `title` varchar(50) COLLATE utf8_unicode_ci NOT NULL COMMENT 'title',
  `author` varchar(20) COLLATE utf8_unicode_ci NOT NULL COMMENT 'author',
  `category` varchar(10) COLLATE utf8_unicode_ci NOT NULL COMMENT 'category',
  `delete_flg` char(1) COLLATE utf8_unicode_ci DEFAULT '0' COMMENT 'delete flg',
  `created_at` datetime DEFAULT CURRENT_TIMESTAMP COMMENT 'created date',
  `updated_at` datetime DEFAULT CURRENT_TIMESTAMP COMMENT 'updated date',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci COMMENT='book table';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `t_book`
--

/*!40101 SET SQL_MODE=@OLD_SQL_MODE */;
/*!40014 SET FOREIGN_KEY_CHECKS=@OLD_FOREIGN_KEY_CHECKS */;
/*!40014 SET UNIQUE_CHECKS=@OLD_UNIQUE_CHECKS */;
/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40101 SET CHARACTER_SET_RESULTS=@OLD_CHARACTER_SET_RESULTS */;
/*!40101 SET COLLATION_CONNECTION=@OLD_COLLATION_CONNECTION */;
/*!40111 SET SQL_NOTES=@OLD_SQL_NOTES */;
