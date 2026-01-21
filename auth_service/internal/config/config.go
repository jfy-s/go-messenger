package config

import (
	"log"
	"os"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Hostname       string `yaml:"hostname" env:"HOSTNAME" required:"true" env-default:"localhost"`
	Port           uint16 `yaml:"port" default:"52525"`
	DatabaseUrl    string `yaml:"database_url" env:"DATABASE_URL" env-required:"true"`
	PrivateKeyPath string `yaml:"private_key_path" env:"PRIVATE_KEY_PATH" env-required:"true"`
	PublicKeyPath  string `yaml:"public_key_path" env:"PUBLIC_KEY_PATH" env-required:"true"`
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
	privatekey, err := os.ReadFile(cfg.PrivateKeyPath)
	if err != nil {
		log.Fatalf("failed to read private key: %v", err)
	}
	publicKey, err := os.ReadFile(cfg.PublicKeyPath)
	if err != nil {
		log.Fatalf("failed to read public key: %v", err)
	}
	os.Setenv("AUTH_SERVICE_PRIVATE_KEY", string(privatekey))
	os.Setenv("AUTH_SERVICE_PUBLIC_KEY", string(publicKey))
	return &cfg
}
