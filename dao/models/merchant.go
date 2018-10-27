package models

import (
    "github.com/go-sql-driver/mysql"
)

const MerchantsTable = "merchants"

type Merchant struct {
    Id        uint64         `db:"mer_id"`
    Title     string         `db:"mer_title"`
    Pubkey    string         `db:"mer_pubkey"`
    AssetCode string         `db:"mer_asset_code"`
    Logo      string         `db:"mer_logo"`
    CreatedAt mysql.NullTime `db:"mer_created_at"`
}
