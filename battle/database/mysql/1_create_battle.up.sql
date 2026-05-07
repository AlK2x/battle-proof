CREATE TABLE `battle` (
     `id` VARCHAR(36) NOT NULL,
     `event_id` VARCHAR(36) NOT NULL,
     `dancer1` VARCHAR(36) NOT NULL,
     `dancer2` VARCHAR(36) NOT NULL,
     `winner` VARCHAR(36) DEFAULT NULL,
     `status` ENUM('pending', 'in_progress', 'finished') DEFAULT NULL,
     `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
     `updated_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
     PRIMARY KEY (`id`)
)