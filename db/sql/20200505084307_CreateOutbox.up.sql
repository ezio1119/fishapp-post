CREATE TABLE `outbox`(
  `id` VARCHAR(255) PRIMARY KEY,
  `event_type` VARCHAR(255) NOT NULL,
  `event_data` JSON NOT NULL,
  `aggregate_id` VARCHAR(255),
  `aggregate_type` VARCHAR(255)
);