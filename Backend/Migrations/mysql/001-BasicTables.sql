-- StoicMigration Up
CREATE TABLE IF NOT EXISTS `User` (
    `user_id` INT AUTO_INCREMENT NOT NULL,
    `user_name` NVARCHAR(256) NOT NULL UNIQUE,
    `first_name` NVARCHAR(256) NOT NULL,
    `last_name` NVARCHAR(256) NOT NULL,
    `email` NVARCHAR(256) NOT NULL,
    `user_status` NVARCHAR(256) NOT NULL,
    `department` NVARCHAR(256) NULL,
    PRIMARY KEY (`user_id`)
);

-- StoicMigration Down
DROP TABLE IF EXISTS `User`;