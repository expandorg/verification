CREATE TABLE `settings` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT,
  `job_id` int(11) unsigned NOT NULL,
  `manual` tinyint(1) NOT NULL DEFAULT '1',
  `requester` tinyint(1) NOT NULL DEFAULT '0',
  `limit` int(11) NOT NULL DEFAULT 0,
  `whitelist` tinyint(1) NOT NULL DEFAULT '1',
  `agreement_count` int(11) NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `job_id` (`job_id`)
)
