package main

type Config struct {
	Database struct {
		URL string `yaml:"url"`
	} `yaml:"database"`
}