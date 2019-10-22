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
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci COMMENT='Users Table';

INSERT INTO t_users (first_name, last_name, email, password)
 VALUES ('taro', 'yamada', 'aaaa@aa.jp', 'fa2195e7026ad9298f9246047b97a83a');

#INSERT INTO `t_users` (`user_id`, `first_name`, `last_name`, `email`, `password`, `delete_flg`, `create_datetime`, `update_datetime`)
#VALUES
#	(1, 'taro', 'yamada', 'aaaa@aa.jp', 'password', '0', '2016-04-29 12:18:13', '2016-04-29 12:18:13');


