package config

import (
	"flag"

	"github.com/caarlos0/env"
)

type Config struct {
	PostgresStorage struct {
		DatabaseDSN string `env:"DATABASE_DSN"`
	}

	Settings struct {
		AddressGRPC string `env:"ADDRESS_GRPC" json:"ADDRESS_GRPC,omitempty"`
		PrivateKey  string `env:"PRIVATE_KEY"`
		PublicKey   string `env:"PUBLIC_KEY"`
		JWTkey      string `env:"JWT_SECRET_KEY"`
	}
}

func GetConfig() *Config {
	cfg := new(Config)
	flag.StringVar(&cfg.Settings.AddressGRPC, "g", ":3200", "grpc port")
	flag.StringVar(&cfg.PostgresStorage.DatabaseDSN, "d", "", "database dsn")
	flag.StringVar(&cfg.Settings.PrivateKey, "private_key", "", "private key path")
	flag.StringVar(&cfg.Settings.PublicKey, "public_key", "", "public key path")
	flag.StringVar(&cfg.Settings.JWTkey, "j", "", "jwt secret key")

	flag.Parse()

	env.Parse(&cfg.Settings)
	env.Parse(&cfg.PostgresStorage)

	return cfg
}
