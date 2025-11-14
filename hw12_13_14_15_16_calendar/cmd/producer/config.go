package main

import (
	"os"
	"strconv"

	"github.com/whatafunc/Golang_Otus_Labs/hw12_13_14_15_16_calendar/internal/config"

	"gopkg.in/yaml.v2"
)

type RabbitConfig struct {
	User         string `yaml:"user"`
	Password     string `yaml:"password"`
	Host         string `yaml:"host"`
	Port         string `yaml:"port"`
	Exchange     string `yaml:"exchange"`
	ExchangeType string `yaml:"exchange-type"`
	Queue        string `yaml:"queue"`
	Key          string `yaml:"key"`
	ConsumerTag  string `yaml:"consumerTag"`
	Lifetime     int    `yaml:"lifetime"` // in seconds (only relevant for consumers)
	Sync         bool   `yaml:"sync"`     // for synchronous confirm publishing
}

// ProducerConfig embeds both storage and RabbitMQ settings.
type ProducerConfig struct {
	config.Config `yaml:"config"` // reuse app config (includes storage, logger, etc.)
	Rabbit        RabbitConfig    `yaml:"rabbit"`
}

type StorageConfig struct {
	Type     string         `yaml:"type"`
	Postgres PostgresConfig `yaml:"postgres"`
}

type PostgresConfig struct {
	DSN string `yaml:"dsn"`
}

// LoadConfig loads from YAML file, then overrides with env vars if present.
func LoadConfig(path string) (ProducerConfig, error) {
	cfg := ProducerConfig{}

	// 1. Load from YAML file if available
	if f, err := os.Open(path); err == nil {
		defer f.Close()
		decoder := yaml.NewDecoder(f)
		_ = decoder.Decode(&cfg) // ignore if empty
	}

	// 2. Override with ENV vars if present
	if v := os.Getenv("RABBIT_USER"); v != "" {
		cfg.Rabbit.User = v
	}
	if v := os.Getenv("RABBIT_PASS"); v != "" {
		cfg.Rabbit.Password = v
	}
	if v := os.Getenv("RABBIT_HOST"); v != "" {
		cfg.Rabbit.Host = v
	}
	if v := os.Getenv("RABBIT_PORT"); v != "" {
		cfg.Rabbit.Port = v
	}
	if v := os.Getenv("RABBIT_EXCHANGE"); v != "" {
		cfg.Rabbit.Exchange = v
	}
	if v := os.Getenv("RABBIT_EXCHANGE_TYPE"); v != "" {
		cfg.Rabbit.ExchangeType = v
	}
	if v := os.Getenv("RABBIT_QUEUE"); v != "" {
		cfg.Rabbit.Queue = v
	}
	if v := os.Getenv("RABBIT_KEY"); v != "" {
		cfg.Rabbit.Key = v
	}
	if v := os.Getenv("RABBIT_TAG"); v != "" {
		cfg.Rabbit.ConsumerTag = v
	}
	if v := os.Getenv("RABBIT_LIFETIME"); v != "" {
		if n, err := strconv.Atoi(v); err == nil {
			cfg.Rabbit.Lifetime = n
		}
	}
	if v := os.Getenv("RABBIT_SYNC"); v != "" {
		if b, err := strconv.ParseBool(v); err == nil {
			cfg.Rabbit.Sync = b
		}
	}

	return cfg, nil
}
