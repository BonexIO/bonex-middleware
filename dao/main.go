package dao

import (
	"bonex-middleware/dao/models"
	"bonex-middleware/types"
	"github.com/stellar/go/clients/horizon"
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
		BlockchainDAO
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
		GetMessageByTxHash(txHash string) (*models.Message, error)
		GetMessages(f *types.MessageFilters) ([]*models.Message, error)
		DeleteMessage(m *models.Message, tx DAOTx) error
	}

	RedisDAO interface {
		SetAccountVolume(ipAddress string, volume decimal.Decimal, ttl int64) error
		GetAccountVolume(ipAddress string) (decimal.Decimal, error)

		AddToQueue(qi *types.QueueItem) error
		PopFromQueue() (*types.QueueItem, error)
	}

	FirebaseDAO interface {
		SendMessage(msg *models.Message) error
	}

	BlockchainDAO interface {
		BlockchainName() string
		SendMoney(string, decimal.Decimal) error
		SetPrivateKey(string) error
		ValidateAddress(string) error
		GetTransaction(txHash string) (*horizon.Transaction, error)
	}
)

func New(redisDAO RedisDAO, dbDAO DbDAO, fDAO FirebaseDAO, bDAO BlockchainDAO) DAO {
	return &struct {
		RedisDAO
		DbDAO
		FirebaseDAO
		BlockchainDAO
	}{
		RedisDAO:      redisDAO,
		DbDAO:         dbDAO,
		FirebaseDAO:   fDAO,
		BlockchainDAO: bDAO,
	}
}
