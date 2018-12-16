-- +migrate Up
CREATE TABLE `messages` (
  `msg_id` int(10) UNSIGNED NOT NULL,
  `msg_tx_hash` varchar(255) NOT NULL,
  `msg_body` TEXT,
  `msg_sender_pubkey` char(56) NULL DEFAULT NULL,
  `msg_receiver_pubkey` char(56) NOT NULL,
  `msg_status` enum('created','waiting_for_tx','notification_sent','received','error') COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT 'created',
  `msg_created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `msg_updated_at` timestamp NULL DEFAULT NULL ON UPDATE CURRENT_TIMESTAMP,
  `msg_received_at` timestamp NULL DEFAULT NULL,
) ENGINE=InnoDB DEFAULT CHARSET=latin1;

ALTER TABLE `messages`
  ADD PRIMARY KEY (`msg_id`),
  ADD UNIQUE KEY `msg_tx_hash_uk` (`msg_tx_hash`);

ALTER TABLE `messages`
  MODIFY `msg_id` int(10) UNSIGNED NOT NULL AUTO_INCREMENT, AUTO_INCREMENT=1;

-- +migrate Down
DROP TABLE messages;