package main

import (
	"os"

	"github.com/whatafunc/Golang_Otus_Labs/hw12_13_14_15_calendar/internal/config" //nolint:depguard
	"gopkg.in/yaml.v3"                                                             //nolint:depguard
)

func LoadConfig(path string) (config.Config, error) {
	f, err := os.Open(path)
	if err != nil {
		return config.Config{}, err
	}
	defer f.Close()

	var cfg config.Config
	decoder := yaml.NewDecoder(f)
	if err := decoder.Decode(&cfg); err != nil {
		return config.Config{}, err
	}
	return cfg, nil
}
