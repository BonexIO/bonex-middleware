package mysql

import (
    "bonex-middleware/dao/models"
    "bonex-middleware/dao/mysql/driver"
    sq "github.com/wedancedalot/squirrel"
    "strings"
)

func (this *mysqlDAO) CreateAccount(m *models.Account) error {
    var err error

    m.Pubkey = strings.ToUpper(m.Pubkey)
    m.Id, err = this.mysql.Insert(sq.Insert(models.AccountsTable).SetMap(sq.Eq{
        "acc_pubkey": m.Pubkey,
    }))

    return err
}

func (this *mysqlDAO) GetAccountByPubkeyOrNil(pubkey string) (*models.Account, error) {
    var res models.Account

    q := sq.Select("*").
        From(models.AccountsTable).
        Where(sq.Eq{"acc_pubkey": strings.ToUpper(pubkey)})

    err := this.mysql.FindFirst(&res, q)
    if err == driver.ErrNoRows {
        return nil, nil
    }

    return &res, err
}
