package config

import (
	"log"
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Server struct {
		Port    int    `yaml:"port"`
		Address string `yaml:"address"`
	} `yaml:"server"`

	Database struct {
		User     string `yaml:"user"`
		Password string `yaml:"password"`
		Name     string `yaml:"name"`
	} `yaml:"database"`

	Email struct {
		Address  string `yaml:"address"`
		Password string `yaml:"password"`
	} `yaml:"email"`
}

var AllConfig Config

func LoadConfig() {
	file, err := os.Open("config/config.yaml")
	if err != nil {
		log.Fatalf("Failed to load config file: %v", err)
	}
	defer file.Close()

	decoder := yaml.NewDecoder(file)
	if err := decoder.Decode(&AllConfig); err != nil {
		log.Fatalf("Failed to decode config file: %v", err)
	}
}
