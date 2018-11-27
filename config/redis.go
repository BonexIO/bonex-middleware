package config

import "fmt"

// Redis config struct
type RedisConfig struct {
	Host string
	Port int
	DB   int //Default 0
}

// Validate checks Redis field
func (this RedisConfig) Validate() error {
	if this.Host == "" {
		return fmt.Errorf("Redis param: Host")
	}
	if this.Port == 0 {
		return fmt.Errorf("Redis param: Port")
	}

	return nil
}
