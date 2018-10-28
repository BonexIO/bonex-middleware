-- +migrate Up
CREATE TABLE `accounts` (
  `acc_id` int(10) UNSIGNED NOT NULL,
  `acc_pubkey` char(56) NOT NULL,
  `acc_created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP
) ENGINE=InnoDB DEFAULT CHARSET=latin1;

CREATE TABLE `subscriptions` (
  `acc_id` int(10) UNSIGNED NOT NULL,
  `mer_id` int(10) UNSIGNED NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=latin1;


ALTER TABLE `accounts`
  ADD PRIMARY KEY (`acc_id`),
  ADD UNIQUE KEY `acc_name` (`acc_pubkey`);

ALTER TABLE `subscriptions`
  ADD PRIMARY KEY (`acc_id`,`mer_id`),
  ADD KEY `mer_id` (`mer_id`);


ALTER TABLE `accounts`
  MODIFY `acc_id` int(10) UNSIGNED NOT NULL AUTO_INCREMENT, AUTO_INCREMENT=2;

ALTER TABLE `subscriptions`
  ADD CONSTRAINT `subscriptions_ibfk_1` FOREIGN KEY (`acc_id`) REFERENCES `accounts` (`acc_id`) ON UPDATE CASCADE,
  ADD CONSTRAINT `subscriptions_ibfk_2` FOREIGN KEY (`mer_id`) REFERENCES `merchants` (`mer_id`) ON UPDATE CASCADE;


-- +migrate Down
DROP TABLE subscriptions;
DROP TABLE accounts;