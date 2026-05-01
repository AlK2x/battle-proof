SET time_zone = '+00:00';

CREATE TABLE `user` (
     `id` VARCHAR(36) NOT NULL,
     `name` VARCHAR(255) NOT NULL,
     `email` VARCHAR(255) NOT NULL,
     `level` ENUM('beginner', 'amateur', 'pro') DEFAULT NULL,
     `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
     PRIMARY KEY (`id`),
     UNIQUE (email)
);