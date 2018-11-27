package config

import (
	"fmt"
	"github.com/wedancedalot/decimal"
)

// Redis config struct
type FaucetConfig struct {
	MaxAllowed24HoursValue decimal.Decimal
	RunEverySeconds uint64
}

// Validate checks Redis field
func (this FaucetConfig) Validate() error {
	if this.MaxAllowed24HoursValue == decimal.Zero {
		return fmt.Errorf("faucet param: MaxAllowed24HoursValue")
	}

	if this.RunEverySeconds == 0 {
		return fmt.Errorf("faucet param: RunEverySeconds")
	}

	return nil
}
