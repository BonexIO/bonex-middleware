package mysql

import (
	"bonex-middleware/dao"
	"bonex-middleware/dao/models"
	"bonex-middleware/types"
	"fmt"
	sq "github.com/wedancedalot/squirrel"
)

func (this *mysqlDAO) CreateMessage(m *models.Message, tx dao.DAOTx) error {
	var err error

	m.Id, err = this.mysql.Insert(sq.Insert(models.MessagesTable).SetMap(sq.Eq{
		"msg_tx_hash":         m.TxHash,
		"msg_body":            m.Body,
		"msg_sender_pubkey":   m.SenderPubkey,
		"msg_receiver_pubkey": m.ReceiverPubkey,
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

func (this *mysqlDAO) GetMessageByTxHash(txHash string) (*models.Message, error) {
	if txHash == "" {
		return nil, types.NewError(types.ErrBadParam, "txHash")
	}

	var m models.Message
	q := sq.Select(
		sq.Field(models.MessagesTable, "*"),
	).From(models.MessagesTable).Where(sq.Eq{"msg_tx_hash": txHash})
	err := this.mysql.FindFirst(&m, q)

	return &m, err
}

func (this *mysqlDAO) GetMessages(f *types.MessageFilters) ([]*models.Message, error) {

	var mLists []*models.Message
	q := sq.Select(
		sq.Field(models.MessagesTable, "*"),
	).From(models.MessagesTable)
	if f != nil {
		if f.FromTime != nil {
			q = q.Where("msg_created_at > ?", f.FromTime)
		}

		if f.ToTime != nil {
			q = q.Where("msg_created_at < ?", f.ToTime)
		}

		if len(f.Statuses) > 0 {
			q = q.Where(sq.Eq{"msg_status": f.Statuses})
		}

		if len(f.NotStatuses) > 0 {
			q = q.Where(sq.NotEq{"msg_status": f.NotStatuses})
		}
	}

	err := this.mysql.Find(&mLists, q)

	return mLists, err
}

func (this *mysqlDAO) DeleteMessage(m *models.Message, tx dao.DAOTx) error {
	if m.Id == 0 {
		return fmt.Errorf("message was not properly loaded")
	}

	q, params, err := sq.Delete(models.MessagesTable).
		Where("msg_id = ?", m.Id).
		ToSql()

	_, err = this.mysql.Exec(q, params, err, this.daoTx2Sqlx(tx))

	return err
}
