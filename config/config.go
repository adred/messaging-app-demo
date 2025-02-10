package config

import "github.com/kelseyhightower/envconfig"

// Config holds application configuration.
type Config struct {
	DBDSN         string `envconfig:"DB_DSN"`
	RabbitMQDSN   string `envconfig:"RABBITMQ_DSN"`
	RabbitMQQueue string `envconfig:"RABBITMQ_QUEUE"`
	HTTPPort      string `envconfig:"HTTP_PORT"`
}

// LoadConfig processes environment variables into a Config struct.
func LoadConfig() (*Config, error) {
	var cfg Config
	err := envconfig.Process("", &cfg)
	if err != nil {
		return nil, err
	}
	return &cfg, nil
}
