package config

import (
	"log"
	"os"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Hostname    string `yaml:"hostname" env:"HOSTNAME" env-required:"true" env-default:"localhost"`
	Port        uint16 `yaml:"port" default:"8080"`
	DatabaseUrl string `yaml:"database_url" env:"DATABASE_URL" env-required:"true"`
}

func Load(configPath string) *Config {
	var cfg Config
	if configPath == "" {
		log.Fatalf("config path is required")
	}
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		log.Fatalf("config file does not exist: %v", err)
	}
	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		log.Fatalf("failed to read config: %v", err)
	}
	return &cfg
}
