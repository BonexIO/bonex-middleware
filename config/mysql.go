package config

import "fmt"

// MysqlConfig is mysql config struct
type MysqlConfig struct {
	Host         string
	Port         int
	User         string
	Password     string
	Database     string
	DebugMode    bool
	MaxOpenConns int
}

// Validate checks all MysqlConfig fields
func (this MysqlConfig) Validate() error {
	if this.Port == 0 {
		return fmt.Errorf("Port")
	}

	if this.Host == "" {
		return fmt.Errorf("Host")
	}

	if this.User == "" {
		return fmt.Errorf("User")
	}

	if this.Database == "" {
		return fmt.Errorf("Database")
	}

	return nil
}
