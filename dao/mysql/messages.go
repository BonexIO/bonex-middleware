package mysql

import (
	"bonex-middleware/dao"
	"bonex-middleware/dao/models"
	"bonex-middleware/types"
	"fmt"
	"github.com/satori/go.uuid"
	sq "github.com/wedancedalot/squirrel"
)

func (this *mysqlDAO) CreateMessage(m *models.Message, tx dao.DAOTx) error {
	var err error

	m.Id, err = this.mysql.Insert(sq.Insert(models.MessagesTable).SetMap(sq.Eq{
		"msg_uuid":            m.Uuid,
		"msg_sender_pubkey":   m.SenderPubkey,
		"msg_receiver_pubkey": m.ReceiverPubkey,
		"msg_tx_hash":         m.TxHash,
		"msg_status":          m.Status,
	}), this.daoTx2Sqlx(tx))

	return err
}

func (this *mysqlDAO) UpdateMessage(m *models.Message, tx dao.DAOTx) error {
	if m.Id == 0 {
		return fmt.Errorf("message was not properly loaded")
	}

	q, params, err := sq.Update(models.MessagesTable).
		SetMap(sq.Eq{
			"msg_status":      m.Status,
			"msg_send_counts": m.SendCounts,
			"msg_received_at": m.ReceivedAt,
		}).
		Where("msg_id = ?", m.Id).
		ToSql()

	_, err = this.mysql.Exec(q, params, err, this.daoTx2Sqlx(tx))

	return err
}

func (this *mysqlDAO) GetMessageById(id uint64, tx dao.DAOTx) (*models.Message, error) {
	if id == 0 {
		return nil, types.NewError(types.ErrBadParam, "id")
	}

	var m models.Message
	q := sq.Select(
		sq.Field(models.MessagesTable, "*"),
	).ForUpdate(tx != nil).
		From(models.MessagesTable).Where(sq.Eq{"msg_id": id})
	err := this.mysql.FindFirst(&m, q, this.daoTx2Sqlx(tx))

	return &m, err
}

func (this *mysqlDAO) GetMessageByUuid(mUuid uuid.UUID) (*models.Message, error) {
	if mUuid == uuid.Nil {
		return nil, types.NewError(types.ErrBadParam, "uuid")
	}

	var m models.Message
	q := sq.Select(
		sq.Field(models.MessagesTable, "*"),
	).
		From(models.MessagesTable).Where(sq.Eq{"msg_uuid": mUuid})
	err := this.mysql.FindFirst(&m, q)

	return &m, err
}
