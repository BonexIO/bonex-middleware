-- +migrate Up
CREATE TABLE `merchants` (
  `mer_id` int(10) UNSIGNED NOT NULL,
  `mer_title` varchar(128) NOT NULL,
  `mer_pubkey` char(56) NOT NULL,
  `mer_asset_code` varchar(12) NOT NULL,
  `mer_logo` text NOT NULL,
  `mer_created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP
) ENGINE=InnoDB DEFAULT CHARSET=latin1;


ALTER TABLE `merchants`
  ADD PRIMARY KEY (`mer_id`),
  ADD UNIQUE KEY `mer_title` (`mer_title`),
  ADD UNIQUE KEY `mer_asset_code` (`mer_asset_code`);


ALTER TABLE `merchants`
  MODIFY `mer_id` int(10) UNSIGNED NOT NULL AUTO_INCREMENT;

-- +migrate Down
DROP TABLE merchants;