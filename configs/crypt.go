package configs

import (
	"github.com/kelseyhightower/envconfig"
)

type CryptConfig struct {
	Secret string `default:"0123456789abcdef"`
}

func GetCryptConfig() (CryptConfig, error) {
	var cfg CryptConfig
	err := envconfig.Process("", &cfg)
	if err != nil {
		return CryptConfig{}, err
	}
	return cfg, nil
}
