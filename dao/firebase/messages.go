package firebase

import (
	"bonex-middleware/dao/models"
	"fmt"
)

func (this *firebaseDAO) SendMessage(msg *models.Message) error {
	data := map[string]string{
		"tx_hash": msg.TxHash,
		"body":    msg.Body,
		"sender":  msg.SenderPubkey.String,
	}

	this.client.NewFcmMsgTo(msg.ReceiverPubkey, data)
	status, err := this.client.Send()
	if err != nil {
		return err
	}
	if status == nil {
		return fmt.Errorf("no satus was recieved")
	}
	if !status.Ok {
		return fmt.Errorf("bad satus was recieved: %s", status.Err)
	}

	return nil
}
