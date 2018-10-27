package mysql

import (
    "bonex-middleware/dao"
    "bonex-middleware/dao/models"
    "bonex-middleware/dao/mysql/driver"
    sq "github.com/wedancedalot/squirrel"
    "strings"
)

func (this *mysqlDAO) CreateMerchant(m *models.Merchant, tx dao.DAOTx) error {
    var err error

    m.Pubkey = strings.ToUpper(m.Pubkey)
    m.Id, err = this.mysql.Insert(sq.Insert(models.MerchantsTable).SetMap(sq.Eq{
        "mer_title":      m.Title,
        "mer_pubkey":     m.Pubkey,
        "mer_asset_code": m.AssetCode,
        "mer_logo":       m.Logo,
    }), this.daoTx2Sqlx(tx))

    return err
}

func (this *mysqlDAO) GetMerchantByPubkeyOrNil(pubkey string, tx dao.DAOTx) (*models.Merchant, error) {
    var res models.Merchant

    q := sq.Select("*").
        ForUpdate(tx != nil).
        From(models.MerchantsTable).
        Where(sq.Eq{"mer_pubkey": strings.ToUpper(pubkey)})

    err := this.mysql.FindFirst(&res, q, this.daoTx2Sqlx(tx))
    if err == driver.ErrNoRows {
        return nil, nil
    }

    return &res, err
}

func (this *mysqlDAO) GetMerchants() ([]*models.Merchant, error) {
    var res []*models.Merchant

    q := sq.Select("*").From(models.MerchantsTable)
    err := this.mysql.Find(&res, q)
    if err == driver.ErrNoRows {
        return nil, nil
    }

    return res, err
}
