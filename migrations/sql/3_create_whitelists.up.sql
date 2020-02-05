CREATE TABLE `whitelists` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT,
  `verifier_id` int(11) unsigned NOT NULL,
  `job_id` int(11) unsigned NOT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `job_verifier` (`job_id`,`verifier_id`)
)