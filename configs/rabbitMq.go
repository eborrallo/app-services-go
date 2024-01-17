package configs

import "github.com/kelseyhightower/envconfig"

type RabbitMqConfig struct {
	RabbitMQUser       string `default:"guest"`
	RabbitMQPassword   string `default:"guest"`
	RabbitMQVhost      string `default:"/"`
	RabbitMQSecure     bool   `default:"false"`
	RabbitMQHostname   string `default:"localhost"`
	RabbitMQPort       int    `default:"5672"`
	RabbitMQExchange   string `default:"domain_events"`
	RabbitMQMaxRetries int    `default:"3"`
	RabbitMQRetryTtl   int    `default:"1000"`
}

func GetRabbitMQConfig() (RabbitMqConfig, error) {
	var cfg RabbitMqConfig
	err := envconfig.Process("", &cfg)
	if err != nil {
		return RabbitMqConfig{}, err
	}
	return cfg, nil
}
