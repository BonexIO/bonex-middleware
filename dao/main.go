package dao

import (
	"bonex-middleware/dao/models"
	"bonex-middleware/types"
	"github.com/satori/go.uuid"
	"github.com/wedancedalot/decimal"
)

type DAOTx interface {
	CommitTx() error
	RollbackTx() error
}

type (
	DAO interface {
		DbDAO
		RedisDAO
		FirebaseDAO
	}

	DbDAO interface {
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

		CreateMessage(m *models.Message, tx DAOTx) error
		UpdateMessage(m *models.Message, tx DAOTx) error
		GetMessageById(id uint64, tx DAOTx) (*models.Message, error)
		GetMessageByUuid(mUuid uuid.UUID) (*models.Message, error)
	}

	RedisDAO interface {
		SetAccountVolume(ipAddress string, volume decimal.Decimal, ttl int64) error
		GetAccountVolume(ipAddress string) (decimal.Decimal, error)

		AddToQueue(qi *types.QueueItem) error
		PopFromQueue() (*types.QueueItem, error)
	}

	FirebaseDAO interface {
		SendMessage(msg types.Message) error
	}
)

func New(redisDAO RedisDAO, dbDAO DbDAO, fDAO FirebaseDAO) DAO {
	return &struct {
		RedisDAO
		DbDAO
		FirebaseDAO
	}{
		RedisDAO:    redisDAO,
		DbDAO:       dbDAO,
		FirebaseDAO: fDAO,
	}
}
