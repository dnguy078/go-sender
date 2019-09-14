package config

import (
	"fmt"

	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	SendGridAPIKey string `envconfig:"SENDER_SENDGRID_KEY"`
	MailGunAPIKey  string `envconfig:"SENDER_MAILGUN_KEY"`
}

func NewConfig() (Config, error) {
	var c Config
	if err := envconfig.Process("sender", &c); err != nil {
		return c, fmt.Errorf("unable to read env vars, err: %v", err)
	}

	return c, nil
}
