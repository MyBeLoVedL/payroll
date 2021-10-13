CREATE TABLE `employees` (
  `id` bigint PRIMARY KEY AUTO_INCREMENT,
  `name` varchar(255) NOT NULL DEFAULT "guest",
  `password` varchar(255) NOT NULL DEFAULT "123456",
  `type` ENUM ('hour', 'salaried', 'commissioned'),
  `mail` varchar(255) NOT NULL,
  `social_security_number` varchar(255) NOT NULL,
  `standard_tax_deductions` decimal(4,2) NOT NULL,
  `other_deductions` decimal(10,2) NOT NULL,
  `phone_number` varchar(255) NOT NULL,
  `salary_rate` decimal(10,2) NOT NULL,
  `hour_limit` int DEFAULT 99999999,
  `payment_method` ENUM ('pick_up', 'mail', 'deposit') DEFAULT "pick_up",
  `deleted` tinyint DEFAULT 0,
  `root` tinyint DEFAULT 0
);

CREATE TABLE `employee_account` (
  `id` bigint PRIMARY KEY AUTO_INCREMENT,
  `bank_name` varchar(255) NOT NULL,
  `account_number` varchar(255) NOT NULL
);

CREATE TABLE `timecard` (
  `id` bigint PRIMARY KEY AUTO_INCREMENT,
  `emp_id` bigint NOT NULL,
  `start_date` timestamp DEFAULT now(),
  `committed` tinyint DEFAULT 0
);

CREATE TABLE `timecard_record` (
  `id` bigint PRIMARY KEY AUTO_INCREMENT,
  `charge_number` bigint NOT NULL,
  `card_id` bigint NOT NULL,
  `hours` smallint NOT NULL,
  `date` date NOT NULL
);

CREATE TABLE `order_info` (
  `order_id` bigint NOT NULL,
  `product_id` bigint NOT NULL,
  `amount` decimal(10,3) NOT NULL
);

CREATE TABLE `purchase_order` (
  `id` bigint PRIMARY KEY AUTO_INCREMENT,
  `emp_id` bigint NOT NULL,
  `customer_contact` varchar(255) NOT NULL,
  `customer_address` varchar(255) NOT NULL,
  `date` timestamp NOT NULL,
  `closed` smallint DEFAULT 0
);

CREATE TABLE `paycheck` (
  `id` bigint PRIMARY KEY AUTO_INCREMENT,
  `emp_id` bigint not null,
  `amount` decimal(10,2) not null,
  `start_date` timestamp not null,
  `end_date` timestamp not null
);

ALTER TABLE `timecard` ADD FOREIGN KEY (`emp_id`) REFERENCES `employees` (`id`);

ALTER TABLE `timecard_record` ADD FOREIGN KEY (`card_id`) REFERENCES `timecard` (`id`);

ALTER TABLE `order_info` ADD FOREIGN KEY (`order_id`) REFERENCES `purchase_order` (`id`);

ALTER TABLE `purchase_order` ADD FOREIGN KEY (`emp_id`) REFERENCES `employees` (`id`);

ALTER TABLE `paycheck` ADD FOREIGN KEY (`emp_id`) REFERENCES `employees` (`id`);

ALTER TABLE `employee_account` ADD FOREIGN KEY (`id`) REFERENCES `employees` (`id`);

CREATE INDEX `employees_index_0` ON `employees` (`id`, `password`);
