package dao

import (
    "bonex-middleware/dao/models"
)

type DAOTx interface {
    CommitTx() error
    RollbackTx() error
}

type DAO interface {
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
