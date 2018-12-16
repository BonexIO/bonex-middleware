package models

import (
	"database/sql"
	"github.com/go-sql-driver/mysql"
)

type MessageStatus string

const (
	MessageStatusCreated          MessageStatus = "created"
	MessageStatusWaitingForTx     MessageStatus = "waiting_for_tx"
	MessageStatusNotificationSent MessageStatus = "notification_sent"
	MessageStatusReceived         MessageStatus = "received"
	MessageStatusError            MessageStatus = "error"
)

const MessagesTable = "messages"

type Message struct {
	Id             uint64         `db:"msg_id"`
	TxHash         string         `db:"msg_tx_hash"`
	Body           string         `db:"msg_body"`
	SenderPubkey   sql.NullString `db:"msg_sender_pubkey"`
	ReceiverPubkey string         `db:"msg_receiver_pubkey"`
	Status         MessageStatus  `db:"msg_status"`
	CreatedAt      mysql.NullTime `db:"msg_created_at"`
	UpdatedAt      mysql.NullTime `db:"msg_updated_at"`
	ReceivedAt     mysql.NullTime `db:"msg_received_at"`
}
