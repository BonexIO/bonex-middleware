package models

const SubscriptionsTable = "subscriptions"

type Subscription struct {
    AccountId  uint64 `db:"acc_id"`
    MerchantId uint64 `db:"mer_id"`
}
