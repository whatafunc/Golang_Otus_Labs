package main

import (
	"os"
	"strconv"

	"gopkg.in/yaml.v2"
)

type ConsumerConfig struct {
	User         string `yaml:"user"`
	Password     string `yaml:"password"`
	Host         string `yaml:"host"`
	Port         string `yaml:"port"`
	Exchange     string `yaml:"exchange"`
	ExchangeType string `yaml:"exchangeType"`
	Queue        string `yaml:"queue"`
	Key          string `yaml:"key"`
	ConsumerTag  string `yaml:"consumerTag"`
	Lifetime     int    `yaml:"lifetime"` // in seconds
}

// LoadConfig loads from YAML file, then overrides with env vars if present.
func LoadConfig(path string) (ConsumerConfig, error) {
	cfg := ConsumerConfig{}

	// 1. Load from YAML file if available
	if f, err := os.Open(path); err == nil {
		defer f.Close()
		decoder := yaml.NewDecoder(f)
		_ = decoder.Decode(&cfg) // ignore if empty
	}

	// 2. Override with ENV vars if present
	if v := os.Getenv("RABBIT_USER"); v != "" {
		cfg.User = v
	}
	if v := os.Getenv("RABBIT_PASS"); v != "" {
		cfg.Password = v
	}
	if v := os.Getenv("RABBIT_HOST"); v != "" {
		cfg.Host = v
	}
	if v := os.Getenv("RABBIT_PORT"); v != "" {
		cfg.Port = v
	}
	if v := os.Getenv("RABBIT_EXCHANGE"); v != "" {
		cfg.Exchange = v
	}
	if v := os.Getenv("RABBIT_EXCHANGE_TYPE"); v != "" {
		cfg.ExchangeType = v
	}
	if v := os.Getenv("RABBIT_QUEUE"); v != "" {
		cfg.Queue = v
	}
	if v := os.Getenv("RABBIT_KEY"); v != "" {
		cfg.Key = v
	}
	if v := os.Getenv("RABBIT_TAG"); v != "" {
		cfg.ConsumerTag = v
	}
	if v := os.Getenv("RABBIT_LIFETIME"); v != "" {
		if n, err := strconv.Atoi(v); err == nil {
			cfg.Lifetime = n
		}
	}

	return cfg, nil
}
