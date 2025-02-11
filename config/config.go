package config

import "github.com/kelseyhightower/envconfig"

// Config holds application configuration.
type Config struct {
	RabbitMQDSN   string `envconfig:"RABBITMQ_DSN"`
	RabbitMQQueue string `envconfig:"RABBITMQ_QUEUE"`
	HTTPPort      string `envconfig:"HTTP_PORT"`
	AuthUsername  string `envconfig:"AUTH_USERNAME"`
	AuthPassword  string `envconfig:"AUTH_PASSWORD"`
	RateLimit     int    `envconfig:"RATE_LIMIT"`
	ReadTimeout   int    `envconfig:"READ_TIMEOUT"`
	WriteTimeout  int    `envconfig:"WRITE_TIMEOUT"`
	IdleTimeout   int    `envconfig:"IDLE_TIMEOUT"`
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
