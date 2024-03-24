package configs

import (
	"github.com/kelseyhightower/envconfig"
	"time"
)

type ServerConfig struct {
	Host            string        `default:"0.0.0.0"`
	Port            uint          `default:"8080"`
	ShutdownTimeout time.Duration `default:"10s"`
	Mode            string        `default:"debug"`
}

func GetServerConfig() (ServerConfig, error) {
	var cfg ServerConfig
	err := envconfig.Process("", &cfg)
	if err != nil {
		return ServerConfig{}, err
	}
	return cfg, nil
}
