package config

import (
	"encoding/json"
	"fmt"
	"bonex-middleware/log"
	"os"
	"reflect"
)

type (
	Config struct {
		Api      ApiConfig
		LogLevel log.LogLevel
		Mysql    MysqlConfig
		Redis    RedisConfig
		Faucet   FaucetConfig

		HorizonClientURL string
	}

	baseConfig interface {
		Validate() error
	}
)

func NewFromFile(filename *string) (*Config, error) {
	var cfg = &Config{}

	//check if file with config exists
	if _, err := os.Stat(*filename); os.IsNotExist(err) {
		//expected file with config not exists
		//check special /.secrets directory (DevOps special)
		developmentConfigPath := "/.secrets/config.json"
		if _, err := os.Stat(developmentConfigPath); os.IsNotExist(err) {
			return nil, err
		}

		filename = &developmentConfigPath
	}

	file, err := os.Open(*filename)
	if err != nil {
		return nil, err
	}

	decoder := json.NewDecoder(file)
	err = decoder.Decode(&cfg)
	if err != nil {
		return nil, err
	}

	return cfg, cfg.Validate()
}

//validates all Config fields
func (this *Config) Validate() error {
	return this.validateBaseConfigStructs()
}

// ValidateBaseConfigStructs validates additional structures (which implements BaseConfig)
func (this *Config) validateBaseConfigStructs() (err error) {
	v := reflect.ValueOf(this).Elem()
	baseConfigType := reflect.TypeOf((*baseConfig)(nil)).Elem()

	for i := 0; i < v.NumField(); i++ {
		if v.Type().Field(i).Type.Implements(baseConfigType) {
			err = v.Field(i).Interface().(baseConfig).Validate()
			if err != nil {
				return fmt.Errorf("invalid param '%s.%s'", v.Field(i).Type().Name(), err.Error())
			}
		}
	}

	return
}
