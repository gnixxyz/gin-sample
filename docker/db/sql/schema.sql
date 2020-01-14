-- create database sandbox
DROP DATABASE IF EXISTS `sandbox`;
CREATE DATABASE IF NOT EXISTS `sandbox` DEFAULT CHARACTER SET utf8;

--
USE `sandbox`

-- drop tables
DROP TABLE IF EXISTS `user`;
DROP TABLE IF EXISTS `post`;

-- user
CREATE TABLE IF NOT EXISTS `user` (
    `id` int(11) unsigned NOT NULL AUTO_INCREMENT,
    `username` varchar(64) NOT NULL,
    `password` varchar(128) NOT NULL,
    `created_at` datetime NOT NULL,
    `updated_at` datetime NOT NULL,
    `deleted_at` datetime NOT NULL,
    PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

-- post
CREATE TABLE IF NOT EXISTS `post` (
    `id` int(11) unsigned NOT NULL AUTO_INCREMENT,
    `title` varchar(140) NOT NULL,
    `content` text,
    `created_at` datetime NOT NULL,
    `updated_at` datetime NOT NULL,
    `deleted_at` datetime NOT NULL,
    `user_id` int(11),
    PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;