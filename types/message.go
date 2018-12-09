package types

type Message struct {
	SenderPubkey   string `json:"sender_pubkey,omitempty"`
	ReceiverPubkey string `json:"receiver_pubkey"`
	Uuid           string `json:"uuid"`
	Body           string `json:"body"`
	TxHash         string `json:"tx_hash,omitempty"`
}

func (this *Message) Validate() error {
	return nil
}
