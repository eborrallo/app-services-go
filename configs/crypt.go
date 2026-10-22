package configs

import (
	"github.com/kelseyhightower/envconfig"
)

type CryptConfig struct {
	// Secret is required: there is deliberately no default, so a missing SECRET
	// fails fast instead of silently running on a public, hard-coded key.
	Secret string `required:"true"`
}

func GetCryptConfig() (CryptConfig, error) {
	var cfg CryptConfig
	err := envconfig.Process("", &cfg)
	if err != nil {
		return CryptConfig{}, err
	}
	return cfg, nil
}
