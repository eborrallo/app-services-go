package configs

import "github.com/kelseyhightower/envconfig"

type RedisConfig struct {
	RedisHost string `default:"localhost"`
	RedisPort int    `default:"6379"`
}

func GetRedisConfig() (RedisConfig, error) {
	var cfg RedisConfig
	err := envconfig.Process("", &cfg)
	if err != nil {
		return RedisConfig{}, err
	}
	return cfg, nil
}
