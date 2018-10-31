package mysql

import (
	"bonex-middleware/dao"
	"bonex-middleware/dao/models"
	"bonex-middleware/dao/mysql/driver"
	sq "github.com/wedancedalot/squirrel"
	"strings"
)

func (this *mysqlDAO) CreateTransaction(tx *models.Transaction, dao dao.DAOTx) error {
	var err error

	tx.Pubkey = strings.ToUpper(tx.Pubkey)
	tx.Id, err = this.mysql.Insert(sq.Insert(models.TransactionsTable).SetMap(sq.Eq{
		"mer_pubkey":     tx.Pubkey,
		"mer_asset_code": tx.AssetCode,
	}), this.daoTx2Sqlx(dao))

	return err
}

func (this *mysqlDAO) RedeemTransaction(issuer, assetCode, secret string, dao dao.DAOTx) (*models.Transaction, error) {
	var res models.Transaction

	q := sq.Select("*").
		ForUpdate(dao != nil).
		From(models.TransactionsTable).
		Where(sq.Eq{
			"mer_pubkey":     strings.ToUpper(strings.ToUpper(issuer)),
			"mer_asset_code": assetCode,
			"secret":         secret,
		})

	err := this.mysql.Find(&res, q, this.daoTx2Sqlx(dao))

	if err == driver.ErrNoRows {
		return nil, nil
	}

	_, err = this.mysql.Exec(sq.Delete(models.TransactionsTable).Where(sq.Eq{
		"tx_id": res.Id,
	}).ToSql())

	return &res, err
}
