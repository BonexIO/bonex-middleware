-- +migrate Up
CREATE TABLE `transactions` (
  `tx_id` int(10) UNSIGNED NOT NULL,
  `mer_pubkey` varchar(128) NOT NULL,
  `mer_asset_code` varchar(12) NOT NULL,
  `amount` int(10) UNSIGNED NOT NULL,
  `secret` varchar(60) NOT NULL,
  `tx_created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP
) ENGINE=InnoDB DEFAULT CHARSET=latin1;


ALTER TABLE `transactions`
  ADD PRIMARY KEY (`tx_id`),
  ADD UNIQUE KEY `secret` (`secret`);

ALTER TABLE `transactions`
  MODIFY `tx_id` int(10) UNSIGNED NOT NULL AUTO_INCREMENT;

-- +migrate Down
DROP TABLE transactions;