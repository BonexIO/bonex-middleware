package models

import (
	"database/sql"
	"github.com/go-sql-driver/mysql"
	"github.com/satori/go.uuid"
)

type MessageStatus string

const (
	MessageStatusCreated      MessageStatus = "created"
	MessageStatusWaitingForTx MessageStatus = "waiting_for_tx"
	MessageStatusReceived     MessageStatus = "received"
	MessageStatusError        MessageStatus = "error"
)

const MessagesTable = "messages"

type Message struct {
	Id             uint64         `db:"msg_id"`
	Uuid           uuid.UUID      `db:"msg_uuid"`
	SenderPubkey   sql.NullString `db:"msg_sender_pubkey"`
	ReceiverPubkey string         `db:"msg_receiver_pubkey"`
	TxHash         sql.NullString `db:"msg_tx_hash"`
	Status         MessageStatus  `db:"msg_status"`
	SendCounts     uint64         `db:"msg_send_counts"`
	CreatedAt      mysql.NullTime `db:"msg_created_at"`
	UpdatedAt      mysql.NullTime `db:"msg_updated_at"`
	ReceivedAt     mysql.NullTime `db:"msg_received_at"`
}
