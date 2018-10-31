package models

import (
	"github.com/go-sql-driver/mysql"
)

const TransactionsTable = "transactions"

type Transaction struct {
	Id        uint64         `db:"tx_id"`
	Pubkey    string         `db:"mer_pubkey"`
	AssetCode string         `db:"mer_asset_code"`
	Amount    uint64         `db:"amount"`
	Secret    string         `db:"secret"`
	CreatedAt mysql.NullTime `db:"tx_created_at"`
}
