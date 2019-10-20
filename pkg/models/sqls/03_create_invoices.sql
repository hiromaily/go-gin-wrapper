CREATE TABLE `t_invoices` (
  `invoice_id` int(11) NOT NULL AUTO_INCREMENT COMMENT 'Invoice ID',
  `user_id` int(11) NOT NULL COMMENT 'User ID',
  `memo` text COLLATE utf8_unicode_ci COMMENT 'Memo',
  `delete_flg` char(1) COLLATE utf8_unicode_ci DEFAULT '0' COMMENT 'delete flg',
  `create_datetime` datetime DEFAULT CURRENT_TIMESTAMP COMMENT 'created date',
  `update_datetime` datetime DEFAULT CURRENT_TIMESTAMP COMMENT 'updated date',
  PRIMARY KEY (`invoice_id`),
  KEY `idx_user` (`user_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci COMMENT='Invoices Table';

