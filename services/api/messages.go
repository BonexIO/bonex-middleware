package api

import (
	"bonex-middleware/dao/models"
	mysqld "bonex-middleware/dao/mysql/driver"
	"bonex-middleware/log"
	"bonex-middleware/services/api/response"
	"bonex-middleware/types"
	"database/sql"
	"encoding/json"
	"github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"github.com/Nargott/goutils"
	"net/http"
	"time"
)

func (this *api) messageSend(w http.ResponseWriter, r *http.Request) {
	var m types.Message
	err := json.NewDecoder(r.Body).Decode(&m)
	if err != nil {
		log.Warn(err.Error())
		response.JsonError(w, types.NewError(types.ErrBadRequest, err.Error()))
		return
	}

	if m.TxHash == "" {
		response.JsonError(w, types.NewError(types.ErrBadParam, "tx_hash"))
		return
	}

	if m.ReceiverPubkey == "" {
		response.JsonError(w, types.NewError(types.ErrBadParam, "receiver_pubkey"))
		return
	}

	if m.Body == "" {
		response.JsonError(w, types.NewError(types.ErrBadParam, "body"))
		return
	}

	msg := models.Message{
		TxHash:         m.TxHash,
		Body:           m.Body,
		SenderPubkey:   goutils.ToNullString(m.SenderPubkey),
		ReceiverPubkey: m.ReceiverPubkey,
		Status:         models.MessageStatusCreated,
	}

	err = this.dao.CreateMessage(&msg, nil)
	if err != nil {
		if err == mysqld.ErrDuplicate {
			log.Warn(err.Error())
			response.JsonError(w, types.NewError(types.ErrAlreadyExists, msg.TxHash))
			return
		}
		log.Errorf("CreateMessage error: %s", err.Error())
		response.JsonError(w, types.NewError(types.ErrService))
		return
	}

	//err = this.dao.SendMessage(m)
	//if err != nil {
	//	log.Errorf("SendMessage error: %s", err.Error())
	//	response.JsonError(w, types.NewError(types.ErrService))
	//	return
	//}

	response.Json(w, map[string]interface{}{
		"tx_hash": msg.TxHash,
		"status":  msg.Status,
	})
	return
}

func (this *api) getMessage(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	txHash := params["tx_hash"]

	//load message by uuid
	msg, err := this.dao.GetMessageByTxHash(txHash)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Warn(err.Error())
			response.JsonError(w, types.NewError(types.ErrNotFound, "tx_hash"))
			return
		}
		log.Errorf("GetMessageByTxHash error: %s", err.Error())
		response.JsonError(w, types.NewError(types.ErrService))
		return
	}

	//check the current message status
	switch msg.Status {
	case models.MessageStatusError:
		fallthrough
	case models.MessageStatusReceived: //TODO may be return this message?
		response.JsonError(w, types.NewError(types.ErrBadStatus, string(msg.Status)))
		return
	}

	//load message by id for update
	dbTx, err := this.dao.BeginTx()
	if err != nil {
		log.Errorf("BeginTx error: %s", err.Error())
		response.JsonError(w, types.NewError(types.ErrService))
		return
	}

	msg, err = this.dao.GetMessageById(msg.Id, dbTx)
	if err != nil {
		log.Errorf("GetMessageById error: %s", err.Error())
		response.JsonError(w, types.NewError(types.ErrService))
		return
	}

	msg.Status = models.MessageStatusReceived
	msg.ReceivedAt = mysql.NullTime{
		Time:  time.Now(),
		Valid: true,
	}

	err = this.dao.UpdateMessage(msg, dbTx)
	if err != nil {
		log.Errorf("UpdateMessage error: %s", err.Error())
		response.JsonError(w, types.NewError(types.ErrService))
		return
	}

	err = dbTx.CommitTx()
	if err != nil {
		log.Errorf("CommitTx error: %s", err.Error())
		response.JsonError(w, types.NewError(types.ErrService))
		return
	}

	response.Json(w, map[string]interface{}{
		"tx_hash": msg.TxHash,
		"body":    msg.Body,
		"status":  msg.Status,
	})
	return
}
