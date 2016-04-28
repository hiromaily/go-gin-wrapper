CREATE TABLE `t_users` (
  `user_id` int(11) NOT NULL AUTO_INCREMENT COMMENT 'User ID',
  `first_name` VARCHAR(20) COLLATE utf8_unicode_ci NOT NULL COMMENT 'First Name',
  `last_name` VARCHAR(20) COLLATE utf8_unicode_ci COMMENT 'Last Name',
  `email` VARCHAR(50) COLLATE utf8_unicode_ci COMMENT 'E-Mail Address',
  `password` VARCHAR(50) COLLATE utf8_unicode_ci COMMENT 'Password',
  `delete_flg` char(1) COLLATE utf8_unicode_ci DEFAULT '0' COMMENT 'delete flg',
  `create_datetime` datetime DEFAULT CURRENT_TIMESTAMP COMMENT 'created date',
  `update_datetime` datetime DEFAULT CURRENT_TIMESTAMP COMMENT 'updated date',
  PRIMARY KEY (`user_id`)
) ENGINE=InnoDB AUTO_INCREMENT=10 DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci COMMENT='Users Table';


INSERT INTO `t_users` (`user_id`, `first_name`, `last_name`, `delete_flg`, `create_datetime`, `update_datetime`)
VALUES
	(1, 'taro', 'yamada', 'aaaa@aa.jp', '', '0', '2016-03-02 12:18:13', '2016-03-02 12:18:13');

