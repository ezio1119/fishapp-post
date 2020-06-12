CREATE TABLE `images`(
  `id` INT(11) PRIMARY KEY AUTO_INCREMENT,
  `name` VARCHAR(255) NOT NULL,
  `post_id` INT(11) NOT NULL,
  `created_at` DATETIME NOT NULL,
  `updated_at` DATETIME NOT NULL,
  UNIQUE (`name`, `post_id`),
  FOREIGN KEY (`post_id`)
    REFERENCES posts(`id`)
    ON DELETE CASCADE
);