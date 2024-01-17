package configs

import (
	"github.com/kelseyhightower/envconfig"
	"time"
)

type DatabaseConfig struct {
	DbUser    string        `default:"user"`
	DbPass    string        `default:"password"`
	DbHost    string        `default:"localhost"`
	DbPort    uint          `default:"3306"`
	DbName    string        `default:"example"`
	DbTimeout time.Duration `default:"5s"`
}

func GetDatabaseConfig() (DatabaseConfig, error) {
	var cfg DatabaseConfig
	err := envconfig.Process("", &cfg)
	if err != nil {
		return DatabaseConfig{}, err
	}
	return cfg, nil
}
