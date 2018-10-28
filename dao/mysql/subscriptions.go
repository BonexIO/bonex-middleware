package mysql

import (
    "bonex-middleware/dao/models"
    "bonex-middleware/dao/mysql/driver"
    sq "github.com/wedancedalot/squirrel"
)

func (this *mysqlDAO) Subscribe(accountId, merchantId uint64) error {
    var err error

    _, err = this.mysql.Insert(sq.Insert(models.SubscriptionsTable).SetMap(sq.Eq{
        "acc_id": accountId,
        "mer_id": merchantId,
    }))

    if err == driver.ErrDuplicate{
        return nil
    }

    return err
}

func (this *mysqlDAO) Unsubscribe(accountId, merchantId uint64) error {
    var err error
    _, err = this.mysql.Exec(sq.Delete(models.SubscriptionsTable).Where(sq.Eq{
        "acc_id": accountId,
        "mer_id": merchantId,
    }).ToSql())
    return err
}

func (this *mysqlDAO) GetSubscriptions(accountId uint64) ([]*models.Merchant, error) {
    var res []*models.Merchant

    q := sq.Select(sq.Field(models.MerchantsTable, "*")).
        From(models.MerchantsTable).
        JoinTable(models.SubscriptionsTable, "mer_id", models.MerchantsTable).
        Where(sq.Eq{sq.Field(models.SubscriptionsTable, "acc_id"): accountId})

    err := this.mysql.Find(&res, q)
    return res, err
}

func (this *mysqlDAO) GetSubscribers(merchantId uint64) ([]*models.Account, error) {
    var res []*models.Account

    q := sq.Select(sq.Field(models.AccountsTable, "*")).
        From(models.AccountsTable).
        JoinTable(models.SubscriptionsTable, "acc_id", models.AccountsTable).
        Where(sq.Eq{sq.Field(models.SubscriptionsTable, "mer_id"): merchantId})

    err := this.mysql.Find(&res, q)
    return res, err
}
