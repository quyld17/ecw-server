CREATE TABLE `users` (
  `user_id` INT PRIMARY KEY NOT NULL AUTO_INCREMENT,
  `email` VARCHAR(255) UNIQUE NOT NULL,
  `password` VARCHAR(255) NOT NULL,
  `full_name` CHAR(30),
  `date_of_birth` DATETIME,
  `phone_number` CHAR(11),
  `gender` TINYINT,
  `role_id` INT NOT NULL DEFAULT 1,
  `created_at` TIMESTAMP NOT NULL DEFAULT (CURRENT_TIMESTAMP)
);

CREATE TABLE `roles` (
  `role_id` INT PRIMARY KEY NOT NULL AUTO_INCREMENT,
  `role_name` VARCHAR(255) NOT NULL
);

CREATE TABLE `addresses` (
  `address_id` INT PRIMARY KEY NOT NULL AUTO_INCREMENT,
  `user_id` INT NOT NULL,
  `name` VARCHAR(255) NOT NULL,
  `address` VARCHAR(255) NOT NULL,
  `is_default` TINYINT NOT NULL
);

CREATE TABLE `orders` (
  `order_id` INT PRIMARY KEY NOT NULL AUTO_INCREMENT,
  `user_id` INT NOT NULL,
  `total_price` DECIMAL(12,0) NOT NULL,
  `payment_method` VARCHAR(255) NOT NULL,
  `address` VARCHAR(255) NOT NULL,
  `status` VARCHAR(255) NOT NULL,
  `created_at` TIMESTAMP NOT NULL DEFAULT (CURRENT_TIMESTAMP)
);

CREATE TABLE `order_products` (
  `id` INT PRIMARY KEY NOT NULL AUTO_INCREMENT,
  `order_id` INT NOT NULL,
  `product_id` INT NOT NULL,
  `product_name` VARCHAR(255) NOT NULL,
  `quantity` INT NOT NULL,
  `price` DECIMAL(12,0) NOT NULL,
  `image_url` VARCHAR(255) NOT NULL,
  `size_name` VARCHAR(255) NOT NULL
);

CREATE TABLE `products` (
  `product_id` INT PRIMARY KEY NOT NULL AUTO_INCREMENT,
  `product_name` VARCHAR(255) NOT NULL,
  `price` DECIMAL(12,0) NOT NULL,
  `total_quantity` INT NOT NULL DEFAULT 0
);

CREATE TABLE `product_images` (
  `image_id` INT PRIMARY KEY NOT NULL AUTO_INCREMENT,
  `product_id` INT NOT NULL,
  `image_url` VARCHAR(255) NOT NULL,
  `is_thumbnail` TINYINT NOT NULL
);

CREATE TABLE `cart_products` (
  `id` INT PRIMARY KEY NOT NULL AUTO_INCREMENT,
  `user_id` INT NOT NULL,
  `product_id` INT NOT NULL,
  `quantity` INT NOT NULL,
  `size_id` INT NOT NULL,
  `selected` TINYINT NOT NULL DEFAULT 0
);

CREATE TABLE `sizes` (
  `size_id` INT PRIMARY KEY NOT NULL AUTO_INCREMENT,
  `size_name` VARCHAR(50) NOT NULL,
  `product_id` INT NOT NULL,
  `quantity` INT NOT NULL
);

ALTER TABLE `users` ADD FOREIGN KEY (`role_id`) REFERENCES `roles` (`role_id`);

ALTER TABLE `addresses` ADD FOREIGN KEY (`user_id`) REFERENCES `users` (`user_id`);

ALTER TABLE `orders` ADD FOREIGN KEY (`user_id`) REFERENCES `users` (`user_id`);

ALTER TABLE `order_products` ADD FOREIGN KEY (`order_id`) REFERENCES `orders` (`order_id`);

ALTER TABLE `cart_products` ADD FOREIGN KEY (`user_id`) REFERENCES `users` (`user_id`);

ALTER TABLE `cart_products` ADD FOREIGN KEY (`product_id`) REFERENCES `products` (`product_id`);

ALTER TABLE `product_images` ADD FOREIGN KEY (`product_id`) REFERENCES `products` (`product_id`);

ALTER TABLE `cart_products` ADD FOREIGN KEY (`size_id`) REFERENCES `sizes` (`size_id`);

ALTER TABLE `sizes` ADD FOREIGN KEY (`product_id`) REFERENCES `products` (`product_id`);
