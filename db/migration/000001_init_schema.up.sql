CREATE TABLE `employees` (
  `id` bigint PRIMARY KEY AUTO_INCREMENT,
  `name` varchar(255) DEFAULT "employee 007",
  `password` varchar(255) DEFAULT "123456",
  `type` ENUM ('hour', 'salaried', 'commissioned'),
  `mail` varchar(255),
  `social_security_number` varchar(255),
  `standard_tax_deductions` varchar(255),
  `other_duductions` varchar(255),
  `phone_number` varchar(255),
  `rate` decimal,
  `hour_limit` int DEFAULT 99999999,
  `payment_method` ENUM ('pick_up', 'mail', 'deposit') DEFAULT "pick_up"
);
CREATE TABLE `employee_account` (
  `id` bigint AUTO_INCREMENT PRIMARY KEY,
  `bank_name` varchar(255),
  `account_number` varchar(255)
);
CREATE TABLE `timecard` (
  `id` bigint AUTO_INCREMENT PRIMARY KEY,
  `employees_id` bigint,
  `start_date` timestamp,
  `end_date` timestamp,
  `status` smallint DEFAULT 0
);
CREATE TABLE `timecard_record` (
  `id` bigint AUTO_INCREMENT PRIMARY KEY,
  `charge_number` bigint,
  `card_id` bigint,
  `hours` smallint,
  `date` timestamp
);
CREATE TABLE `order_info` (
  `order_id` bigint,
  `product_id` bigint,
  `amount` int,
  primary key (`order_id`, `product_id`)
);
CREATE TABLE `purchase_order` (
  `id` bigint AUTO_INCREMENT PRIMARY KEY,
  `emp_id` bigint,
  `customer_contact` varchar(255),
  `customer_address` varchar(255),
  `order_info_id` bigint,
  `date` timestamp,
  `status` smallint
);
CREATE TABLE `paycheck` (
  `id` bigint AUTO_INCREMENT PRIMARY KEY,
  `emp_id` bigint,
  `amount` decimal(10, 3),
  `start_date` timestamp,
  `end_date` timestamp
);
ALTER TABLE
  `timecard`
ADD
  FOREIGN KEY (`employees_id`) REFERENCES `employees` (`id`);
ALTER TABLE
  `timecard_record`
ADD
  FOREIGN KEY (`card_id`) REFERENCES `timecard` (`id`);
ALTER TABLE
  `order_info`
ADD
  FOREIGN KEY (`order_id`) REFERENCES `purchase_order` (`id`);
ALTER TABLE
  `purchase_order`
ADD
  FOREIGN KEY (`emp_id`) REFERENCES `employees` (`id`);
ALTER TABLE
  `paycheck`
ADD
  FOREIGN KEY (`emp_id`) REFERENCES `employees` (`id`);
ALTER TABLE
  `employee_account`
ADD
  FOREIGN KEY (`id`) REFERENCES `employees` (`id`);
CREATE INDEX `employees_index_0` ON `employees` (`id`, `password`);

-- v1.1 
-- change the type of columnns regarding tax from varchar to decimal in table employees
ALTER TABLE `payroll_system`.`employees` 
CHANGE COLUMN `standard_tax_deductions` `standard_tax_deductions` DECIMAL(3,2) NULL DEFAULT NULL ,
CHANGE COLUMN `other_duductions` `other_duductions` DECIMAL(10) NULL DEFAULT NULL ,
-- clarify the definition of rate which is used to represents the three attributes: 
-- hourly rate, salary and commission rate as a whole
CHANGE COLUMN `rate` `salary_rate` DECIMAL(10,0) NULL DEFAULT NULL;

-- unify the column name for the same attribute
ALTER TABLE `payroll_system`.`timecard` 
DROP FOREIGN KEY `timecard_ibfk_1`;
ALTER TABLE `payroll_system`.`timecard` 
CHANGE COLUMN `employees_id` `emp_id` BIGINT NULL DEFAULT NULL ;
ALTER TABLE `payroll_system`.`timecard` 
ADD CONSTRAINT `timecard_ibfk_1`
  FOREIGN KEY (`emp_id`)
  REFERENCES `payroll_system`.`employees` (`id`);


