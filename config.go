package main

import (
	"os"

	"gopkg.in/yaml.v2"
)

func LoadConfig() (*Config, error) {
	var config Config
	configData, err := os.ReadFile("config.yml")
	if err != nil {
		return nil, err
	}
	err = yaml.Unmarshal(configData, &config)
	if err != nil {
		return nil, err
	}
	return &config, nil
}