package configs

import (
	"github.com/kelseyhightower/envconfig"
)

type BlokchainConfig struct {
	ChainId string `default:"1"`
	Chain   string `default:"mainnet"`
	Network string `default:"ethereum"`
	RpcUrl  string `default:"https://eth.public-rpc.com"`
}

func GetBlokchainConfig() (BlokchainConfig, error) {
	var cfg BlokchainConfig
	err := envconfig.Process("", &cfg)
	if err != nil {
		return BlokchainConfig{}, err
	}
	return cfg, nil
}
