package config

import (
	"fmt"

	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	SenderAPIPort int `envconfig:"API_PORT"`

	SendGridAPIKey string `envconfig:"SENDGRID_KEY"`
	MailGunAPIKey  string `envconfig:"MAILGUN_KEY"`
	SparkPostKey   string `envconfig:"SPARKPOST_KEY"`
}

func NewConfig() (Config, error) {
	var c Config
	if err := envconfig.Process("sender", &c); err != nil {
		return c, fmt.Errorf("unable to read env vars, err: %v", err)
	}
	fmt.Printf("%+v", c)

	return c, nil
}
