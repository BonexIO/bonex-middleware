package dao

import (
    "bonex-middleware/dao/models"
    "bonex-middleware/types"
    "github.com/wedancedalot/decimal"
)

type DAOTx interface {
    CommitTx() error
    RollbackTx() error
}

type DAO interface {
    DbDAO
    RedisDAO
}

type DbDAO interface {
    BeginTx() (DAOTx, error)

    // Merchants
    CreateMerchant(*models.Merchant, DAOTx) error
    GetMerchants() ([]*models.Merchant, error)
    GetMerchantByPubkeyOrNil(string, DAOTx) (*models.Merchant, error)

    CreateAccount(*models.Account) error
    GetAccountByPubkeyOrNil(string) (*models.Account, error)

    Subscribe(accountId, merchantId uint64) error
    Unsubscribe(accountId, merchantId uint64) error

    GetSubscriptions(accountId uint64) ([]*models.Merchant, error)
    GetSubscribers(merchantId uint64) ([]*models.Account, error)
}

type RedisDAO interface {
    SetAccountVolume(ipAddress string, volume decimal.Decimal, ttl int64) error
    GetAccountVolume(ipAddress string) (decimal.Decimal, error)

    AddToQueue(qi *types.QueueItem) error
    PopFromQueue() (*types.QueueItem, error)
}

func New(redisDAO RedisDAO, dbDAO DbDAO) DAO {
    return &struct {
        RedisDAO
        DbDAO
    }{
        redisDAO,
        dbDAO,
    }
}
