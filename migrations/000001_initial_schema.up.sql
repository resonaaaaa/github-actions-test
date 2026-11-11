
-- Create users table
CREATE TABLE `users` (
    `id` BIGINT NOT NULL AUTO_INCREMENT,
    `name` VARCHAR(100) NOT NULL,
    `email` VARCHAR(255) DEFAULT NULL,
    PRIMARY KEY (`id`),
    UNIQUE KEY `email` (`email`)
) ENGINE=InnoDB
  DEFAULT CHARSET=utf8mb4
  COLLATE=utf8mb4_0900_ai_ci;


-- Create products table
CREATE TABLE `products` (
    `id` BIGINT NOT NULL AUTO_INCREMENT,
    `name` VARCHAR(100) NOT NULL,
    `price_cents` BIGINT NOT NULL,
    `stock` INT NOT NULL DEFAULT '0',
    PRIMARY KEY (`id`),
    CONSTRAINT `chk_product_price`
        CHECK (`price_cents` >= 0),
    CONSTRAINT `chk_product_stock`
        CHECK (`stock` >= 0)
) ENGINE=InnoDB
  DEFAULT CHARSET=utf8mb4
  COLLATE=utf8mb4_0900_ai_ci;


-- Create orders table
CREATE TABLE `orders` (
    `id` BIGINT NOT NULL AUTO_INCREMENT,
    `user_id` BIGINT NOT NULL,
    `product_id` BIGINT NOT NULL,
    `quantity` INT NOT NULL,
    `unit_price` BIGINT NOT NULL,
    `request_id` VARCHAR(64) NOT NULL,
    `created_at` DATETIME(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6),

    PRIMARY KEY (`id`),

    UNIQUE KEY `uq_order_request`
        (`user_id`, `request_id`),

    KEY `fk_orders_product` (`product_id`),

    CONSTRAINT `fk_orders_product`
        FOREIGN KEY (`product_id`)
        REFERENCES `products` (`id`),

    CONSTRAINT `fk_orders_user`
        FOREIGN KEY (`user_id`)
        REFERENCES `users` (`id`),

    CONSTRAINT `chk_order_quantity`
        CHECK (`quantity` > 0)

) ENGINE=InnoDB
  DEFAULT CHARSET=utf8mb4
  COLLATE=utf8mb4_0900_ai_ci;
