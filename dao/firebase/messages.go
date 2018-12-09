package firebase

import (
	"bonex-middleware/types"
	"fmt"
	"goutils"
)

func (this *firebaseDAO) SendMessage(msg types.Message) error {
	pUuid, err := goutils.UUIDToBase64(msg.Uuid)
	if err != nil {
		return err
	}

	data := map[string]string{
		"uuid":        msg.Uuid,
		"uuid_packed": pUuid,
		"body":        msg.Body,
		"sender":      msg.SenderPubkey,
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
