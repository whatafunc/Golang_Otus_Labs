package main

import (
	"os"

	"github.com/whatafunc/Golang_Otus_Labs/hw12_13_14_15_calendar/internal/config"
	"gopkg.in/yaml.v3"
)

// При желании конфигурацию можно вынести в internal/config.
// Организация конфига в main принуждает нас сужать API компонентов, использовать
// при их конструировании только необходимые параметры, а также уменьшает вероятность циклической зависимости.
type Config struct {
	Logger         LoggerConf    `yaml:"logger"`
	HTTP           HTTPConf      `yaml:"http"`
	Storage        StorageConfig `yaml:"storage"`
	MigrationsPath string        `yaml:"migrations_path"`
}

type LoggerConf struct {
	Level string `yaml:"level"`
}

type HTTPConf struct {
	Listen string `yaml:"listen"`
}

type StorageConfig struct {
	Type     string                `yaml:"type"`
	Redis    config.RedisConfig    `yaml:"redis"`
	Postgres config.PostgresConfig `yaml:"postgres"`
}

func NewConfig() Config {
	return Config{}
}

func LoadConfig(path string) (Config, error) {
	f, err := os.Open(path)
	if err != nil {
		return Config{}, err
	}
	defer f.Close()

	var cfg Config
	decoder := yaml.NewDecoder(f)
	if err := decoder.Decode(&cfg); err != nil {
		return Config{}, err
	}
	return cfg, nil
}
