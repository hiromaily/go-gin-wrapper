CREATE TABLE `t_news` (
  `news_id` int(11) NOT NULL AUTO_INCREMENT COMMENT 'News ID',
  `article` TEXT COLLATE utf8_unicode_ci NOT NULL COMMENT 'Article',
  `delete_flg` char(1) COLLATE utf8_unicode_ci DEFAULT '0' COMMENT 'delete flg',
  `create_datetime` datetime DEFAULT CURRENT_TIMESTAMP COMMENT 'created date',
  `update_datetime` datetime DEFAULT CURRENT_TIMESTAMP COMMENT 'updated date',
  PRIMARY KEY (`news_id`)
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci COMMENT='News Table';

INSERT INTO t_news ("article") VALUES ('article01');
INSERT INTO t_news ("article") VALUES ('article02');
INSERT INTO t_news ("article") VALUES ('article03');
INSERT INTO t_news ("article") VALUES ('article04');
INSERT INTO t_news ("article") VALUES ('article05');
INSERT INTO t_news ("article") VALUES ('article06');
INSERT INTO t_news ("article") VALUES ('article07');
INSERT INTO t_news ("article") VALUES ('article08');
INSERT INTO t_news ("article") VALUES ('article09');
INSERT INTO t_news ("article") VALUES ('article10');
