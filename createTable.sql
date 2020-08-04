/*
SQLyog 
MySQL - 10.3.22-MariaDB-0+deb10u1 : Database - happyblog
*********************************************************************
*/

/*!40101 SET NAMES utf8 */;

/*!40101 SET SQL_MODE=''*/;

/*!40014 SET @OLD_UNIQUE_CHECKS=@@UNIQUE_CHECKS, UNIQUE_CHECKS=0 */;
/*!40101 SET @OLD_SQL_MODE=@@SQL_MODE, SQL_MODE='NO_AUTO_VALUE_ON_ZERO' */;
/*!40111 SET @OLD_SQL_NOTES=@@SQL_NOTES, SQL_NOTES=0 */;
/*Table structure for table `happyblog_settings` */

CREATE TABLE `happyblog_settings` (
  `id` bigint(32) NOT NULL AUTO_INCREMENT COMMENT 'index',
  `configGroup` varchar(20) NOT NULL,
  `configName` varchar(20) NOT NULL,
  `configValue` varchar(100) NOT NULL,
  `configOrder` int(10) NOT NULL DEFAULT 1,
  `configType` enum('string','image') NOT NULL DEFAULT 'string',
  PRIMARY KEY (`id`),
  UNIQUE KEY `group_configKey` (`configGroup`,`configName`)
) ENGINE=InnoDB AUTO_INCREMENT=10 DEFAULT CHARSET=utf8;

/*Table structure for table `happyblog_tblAlbum` */

CREATE TABLE `happyblog_tblAlbum` (
  `id` int(10) NOT NULL AUTO_INCREMENT COMMENT 'the index key',
  `albumName` varchar(20) NOT NULL COMMENT 'the name of Album',
  `isPublic` tinyint(1) DEFAULT 1 COMMENT 'wheather public',
  `articleTotal` int(10) NOT NULL DEFAULT 0 COMMENT 'total articles',
  `authorId` int(10) NOT NULL COMMENT 'the Album author',
  `createTime` datetime DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=14 DEFAULT CHARSET=utf8 COMMENT='album list';

/*Table structure for table `happyblog_tblAlbumRe` */

CREATE TABLE `happyblog_tblAlbumRe` (
  `id` bigint(32) NOT NULL AUTO_INCREMENT COMMENT 'index id',
  `articleId` int(10) NOT NULL COMMENT 'articleId',
  `albumId` int(10) NOT NULL COMMENT 'albumId',
  PRIMARY KEY (`id`),
  KEY `ablum_article` (`albumId`,`articleId`)
) ENGINE=InnoDB AUTO_INCREMENT=8 DEFAULT CHARSET=utf8 COMMENT='relation article and album';

/*Table structure for table `happyblog_tblArticle` */

CREATE TABLE `happyblog_tblArticle` (
  `id` bigint(32) NOT NULL AUTO_INCREMENT COMMENT 'the index key',
  `title` varchar(100) NOT NULL COMMENT 'the article title',
  `content` longtext DEFAULT NULL COMMENT 'content',
  `pubStatus` tinyint(1) NOT NULL DEFAULT 0 COMMENT 'publish status',
  `createTime` datetime DEFAULT NULL COMMENT 'create time',
  `updateTime` datetime DEFAULT NULL COMMENT 'update tim',
  `authorId` int(10) DEFAULT NULL COMMENT 'author id',
  `independPage` tinyint(1) DEFAULT 2 COMMENT 'wheather page independ',
  `brief` varchar(100) DEFAULT NULL COMMENT 'article describe info',
  `keywords` varchar(50) DEFAULT NULL COMMENT 'article key words',
  `headimage` varchar(200) DEFAULT NULL COMMENT 'article head image',
  `uri` varchar(100) DEFAULT NULL COMMENT 'user define uri',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=60 DEFAULT CHARSET=utf8;

/*Table structure for table `happyblog_tblTag` */

CREATE TABLE `happyblog_tblTag` (
  `id` int(10) NOT NULL AUTO_INCREMENT COMMENT 'the index key',
  `tagName` varchar(20) NOT NULL COMMENT 'the tag name',
  `createTime` datetime NOT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `tagname` (`tagName`)
) ENGINE=InnoDB AUTO_INCREMENT=23 DEFAULT CHARSET=utf8;

/*Table structure for table `happyblog_tblUser` */

CREATE TABLE `happyblog_tblUser` (
  `id` int(10) NOT NULL AUTO_INCREMENT COMMENT 'the index key',
  `accountEmail` varchar(100) NOT NULL COMMENT 'the account email address',
  `accountPassword` char(32) NOT NULL COMMENT 'the password whth solt',
  `nickName` varchar(50) NOT NULL COMMENT 'the nickName',
  `createTime` datetime DEFAULT NULL,
  `updateTime` datetime DEFAULT NULL,
  `lastLogin` datetime DEFAULT NULL,
  `headImageUri` varchar(200) NOT NULL COMMENT 'the headimage uri',
  `emailVerify` tinyint(1) DEFAULT 0 COMMENT 'wheather the email valid',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=3 DEFAULT CHARSET=utf8;

/*Table structure for table `happylblog_tblTagRe` */

CREATE TABLE `happylblog_tblTagRe` (
  `id` bigint(32) NOT NULL AUTO_INCREMENT COMMENT 'index id',
  `articleId` int(10) NOT NULL COMMENT 'articleId',
  `tagId` int(10) DEFAULT NULL COMMENT 'tagId',
  PRIMARY KEY (`id`),
  UNIQUE KEY `tags_articleid` (`tagId`,`articleId`),
  UNIQUE KEY `articleId_tags` (`articleId`,`tagId`)
) ENGINE=InnoDB AUTO_INCREMENT=1505 DEFAULT CHARSET=utf8 COMMENT='relation tag and article';

/*!40101 SET SQL_MODE=@OLD_SQL_MODE */;
/*!40014 SET UNIQUE_CHECKS=@OLD_UNIQUE_CHECKS */;
/*!40111 SET SQL_NOTES=@OLD_SQL_NOTES */;
