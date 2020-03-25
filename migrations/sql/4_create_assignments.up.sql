CREATE TABLE `assignments` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `job_id` int(10) unsigned NOT NULL,
  `task_id` int(10) unsigned NOT NULL,
  `verifier_id` int(10) unsigned DEFAULT NULL,
  `response_id` int(10) unsigned NOT NULL,
  `active` tinyint(1) DEFAULT NULL,
  `status` varchar(20) NOT NULL DEFAULT 'inactive', 
  `assigned_at` timestamp NULL DEFAULT NULL,
  `expires_at` timestamp NULL DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `verifier_id` (`verifier_id`,`job_id`, `active`, `response_id`),
  KEY `job_id` (`job_id`),
  KEY `task_id` (`task_id`)
)