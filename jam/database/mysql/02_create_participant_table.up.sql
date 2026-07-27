CREATE TABLE `participant` (
     `jam_id` VARCHAR(36) NOT NULL,
     `user_id` VARCHAR(255) NOT NULL,
     `role` ENUM('dancer', 'judge', 'media') NOT NULL,
     `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
     PRIMARY KEY (`jam_id`, `user_id`)
)