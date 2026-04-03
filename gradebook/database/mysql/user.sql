SET NAMES utf8;
SET time_zone = '+00:00';

CREATE TABLE `user` (
     `id` VARCHAR(36) NOT NULL,
     `name` VARCHAR(255) NOT NULL,
     `password` VARBINARY(255) NOT NULL,
     `salt` VARCHAR(255) DEFAULT NULL,
     PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;