package configs

import "github.com/kelseyhightower/envconfig"

type SmtpConfig struct {
	// Credentials have no defaults on purpose: they must come from the
	// environment, never be baked into the binary.
	Email    string `envconfig:"SMTP_EMAIL"`
	Password string `envconfig:"SMTP_PASSWORD"`
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
