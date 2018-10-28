package models

import (
    "github.com/go-sql-driver/mysql"
)

const AccountsTable = "accounts"

type Account struct {
    Id        uint64         `db:"acc_id"`
    Pubkey    string         `db:"acc_pubkey"`
    CreatedAt mysql.NullTime `db:"acc_created_at"`
}
