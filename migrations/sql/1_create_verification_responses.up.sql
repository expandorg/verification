CREATE TABLE `verification_responses` (
  `id` int(10) unsigned AUTO_INCREMENT,
  `job_id` int unsigned NOT NULL,
  `task_id` int unsigned NOT NULL,
  `response_id` int unsigned NOT NULL,
  `worker_id` int unsigned NOT NULL,
  `verifier_id` int unsigned NULL,
  `accepted` tinyint(1) NOT NULL,
  `reason` text NULL,
  `created_at` TIMESTAMP NOT NULL DEFAULT NOW(),
  `updated_at` TIMESTAMP NOT NULL DEFAULT NOW() ON UPDATE NOW(),
  PRIMARY KEY (`id`),
  UNIQUE KEY `response_id` (`response_id`)
)
