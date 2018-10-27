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
}
