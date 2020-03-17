CREATE TABLE `assignments` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `job_id` int(10) unsigned NOT NULL,
  `task_id` int(10) unsigned NOT NULL,
  `verifier_id` int(10) unsigned NOT NULL,
  `response_id` int(10) unsigned,
  `active` tinyint(1) DEFAULT '1',
  `status` varchar(20) NOT NULL DEFAULT 'active', 
  `assigned_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `expires_at` timestamp NULL DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `verifier_id` (`verifier_id`,`job_id`, `active`, `response_id`),
  KEY `job_id` (`job_id`),
  KEY `task_id` (`task_id`)
)