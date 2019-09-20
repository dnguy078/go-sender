package config

import (
	"fmt"

	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	SendGridAPIKey string `envconfig:"SENDGRID_KEY"`
	SparkPostKey   string `envconfig:"SPARKPOST_KEY"`
	RabbitUsername string `envconfig:"RABBIT_USERNAME"`
	RabbitPassword string `envconfig:"RABBIT_PASSWORD"`
}

func NewConfig() (Config, error) {
	var c Config
	if err := envconfig.Process("sender", &c); err != nil {
		return c, fmt.Errorf("unable to read env vars, err: %v", err)
	}

	return c, nil
}
