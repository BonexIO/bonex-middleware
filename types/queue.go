package types

import (
	"fmt"
	"github.com/wedancedalot/decimal"
)

type QueueItem struct {
	Address string `json:"address"`
	Amount  string `json:"amount"`
}

func (this *QueueItem) Validate() error {
	if this.Address == "" {
		return fmt.Errorf("no address provided")
	}

	if this.Amount == "" {
		return fmt.Errorf("no amount provided")
	}

	amount, err := decimal.NewFromString(this.Amount)
	if err != nil {
		return fmt.Errorf("bad amount format (not a decimal value): %s", err.Error())
	}

	if amount.LessThanOrEqual(decimal.Zero) {
		return fmt.Errorf("too low amount %s", amount.String())
	}

	return nil
}