package types

import (
	"bonex-middleware/dao/models"
	"time"
)

type Message struct {
	TxHash         string `json:"tx_hash"`
	SenderPubkey   string `json:"sender_pubkey,omitempty"`
	ReceiverPubkey string `json:"receiver_pubkey"`
	Body           string `json:"body"`
}

func (this *Message) Validate() error {
	return nil
}

type MessageFilters struct {
	FromTime    *time.Time
	ToTime      *time.Time
	Statuses    []models.MessageStatus
	NotStatuses []models.MessageStatus
}
