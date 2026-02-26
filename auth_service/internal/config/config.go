package config

import (
	"fmt"
	"log"
	"os"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Hostname       string `yaml:"hostname" required:"true" env-default:"localhost"`
	Port           uint16 `yaml:"port" default:"52525"`
	DatabaseUrl    string `yaml:"database_url" env:"DATABASE_URL" env-required:"true"`
	PrivateKeyPath string `yaml:"private_key_path" env:"PRIVATE_KEY_PATH" required:"true"`
	PublicKeyPath  string `yaml:"public_key_path" env:"PUBLIC_KEY_PATH" required:"true"`
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

	entries, err := os.ReadDir("/app/auth_service")
	if err != nil {
		log.Fatalf("ошибка чтения директории: %v", err)
	}

	for _, e := range entries {
		info, err := e.Info()
		if err != nil {
			fmt.Printf("%s (ошибка: %v)\n", e.Name(), err)
			continue
		}

		mode := info.Mode()
		var kind string
		switch {
		case mode.IsDir():
			kind = "DIR"
		case mode.IsRegular():
			kind = "FILE"
		default:
			kind = "OTHER"
		}

		fmt.Printf("%s\t%s\t%d bytes\n", kind, e.Name(), info.Size())
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
