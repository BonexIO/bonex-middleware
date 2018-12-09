package api

import (
	"bonex-middleware/dao/models"
	"bonex-middleware/log"
	"bonex-middleware/services/api/response"
	"bonex-middleware/types"
	"database/sql"
	"encoding/json"
	"github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"github.com/satori/go.uuid"
	"goutils"
	"net/http"
	"time"
)

func (this *api) parseMessageUuid(mUuid string) (uuid.UUID, error) {
	u, err := uuid.FromString(mUuid)
	if err != nil {
		// try to parse as packed uuid
		sUuid, err := goutils.Base64ToUuid(mUuid)
		if err != nil {
			return uuid.Nil, err
		}
		return uuid.FromStringOrNil(sUuid), nil
	}

	return u, nil
}

func (this *api) messageSend(w http.ResponseWriter, r *http.Request) {
	var m types.Message
	err := json.NewDecoder(r.Body).Decode(&m)
	if err != nil {
		log.Warn(err.Error())
		response.JsonError(w, types.NewError(types.ErrBadRequest, err.Error()))
		return
	}

	mUuid, err := this.parseMessageUuid(m.Uuid)
	if err != nil {
		log.Warn(err.Error())
		response.JsonError(w, types.NewError(types.ErrBadParam, "uuid"))
		return
	}

	msg := models.Message{
		Uuid:           mUuid,
		SenderPubkey:   goutils.ToNullString(m.SenderPubkey),
		ReceiverPubkey: m.ReceiverPubkey,
		TxHash:         goutils.ToNullString(m.TxHash),
		Status:         models.MessageStatusCreated,
	}

	err = this.dao.CreateMessage(&msg, nil)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Warn(err.Error())
			response.JsonError(w, types.NewError(types.ErrAlreadyExists, msg.Uuid.String()))
			return
		}
		log.Errorf("CreateMessage error: %s", err.Error())
		response.JsonError(w, types.NewError(types.ErrService))
		return
	}

	err = this.dao.SendMessage(m)
	if err != nil {
		log.Errorf("SendMessage error: %s", err.Error())
		response.JsonError(w, types.NewError(types.ErrService))
		return
	}

	response.Json(w, map[string]interface{}{
		"uuid":   msg.Uuid.String(),
		"status": msg.Status,
	})
	return
}

func (this *api) messageReceived(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	mUuid, err := this.parseMessageUuid(params["uuid"])
	if err != nil {
		log.Warn(err.Error())
		response.JsonError(w, types.NewError(types.ErrBadParam, "uuid"))
		return
	}

	//load message by uuid
	msg, err := this.dao.GetMessageByUuid(mUuid)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Warn(err.Error())
			response.JsonError(w, types.NewError(types.ErrNotFound, "uuid"))
			return
		}
		log.Errorf("GetMessageByUuid error: %s", err.Error())
		response.JsonError(w, types.NewError(types.ErrService))
		return
	}

	//check the current message status
	switch msg.Status {
	case models.MessageStatusError:
		fallthrough
	case models.MessageStatusReceived:
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
		"uuid":   msg.Uuid.String(),
		"status": msg.Status,
	})
	return
}
