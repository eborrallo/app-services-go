package configs

import "github.com/kelseyhightower/envconfig"

type SmtpConfig struct {
	Email    string `default:"test@gmail.com"`
	Password string `default:"123123"`
	Host     string `default:"smtp.gmail.com"`
	Port     string `default:"587"`
}

func GetSmtpConfig() (SmtpConfig, error) {
	var cfg SmtpConfig
	err := envconfig.Process("", &cfg)
	if err != nil {
		return SmtpConfig{}, err
	}
	return cfg, nil
}
