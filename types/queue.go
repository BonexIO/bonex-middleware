package types

import (
	"github.com/wedancedalot/decimal"
)

type QueueItem struct {
	Address string
	Amount  decimal.Decimal
}

func (this *QueueItem) Validate() error {
	if this.Address == "" {
		return NewError(ErrBadParam, "address")
	}

	if this.Amount.LessThanOrEqual(decimal.Zero) {
		return NewError(ErrBadParam, "amount")
	}

	return nil
}