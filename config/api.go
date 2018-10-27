package config

import (
	"fmt"
)

type ApiConfig struct {
	Port           int
	AllowedOrigins []string
}

// Validate checks all MysqlConfig fields
func (this ApiConfig) Validate() error {
	if this.Port == 0 {
		return fmt.Errorf("Port")
	}

	if len(this.AllowedOrigins) == 0 {
		return fmt.Errorf("AllowedOrigins")
	}

	return nil
}
