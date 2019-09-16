package config

import (
	"fmt"

	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	SendGridAPIKey string `envconfig:"SENDGRID_KEY"`
	MailGunAPIKey  string `envconfig:"MAILGUN_KEY"`
	SenderAPIPort  int    `envconfig:"API_PORT"`
}

func NewConfig() (Config, error) {
	var c Config
	if err := envconfig.Process("sender", &c); err != nil {
		return c, fmt.Errorf("unable to read env vars, err: %v", err)
	}

	return c, nil
}
