package config

import (
	"flag"

	"github.com/caarlos0/env"
)

type Config struct {
	PostgresStorage struct {
		DatabaseDSN string `env:"DATABASE_DSN"`
	}

	GRPC struct {
		AddressGRPC   string `env:"ADDRESS_GRPC"`
		CertSSL       string `env:"CERT_SSL"`
		PrivateKeySSL string `env:"PRIVATE_KEY_SSL"`
	}

	Settings struct {
		PrivateKey string `env:"PRIVATE_KEY"`
		JWTkey     string `env:"JWT_SECRET_KEY"`
	}
}

func GetConfig() *Config {
	cfg := new(Config)

	flag.StringVar(&cfg.GRPC.AddressGRPC, "g", ":3200", "grpc port")
	flag.StringVar(&cfg.GRPC.CertSSL, "c", "cert.crt", "cert for ssl path")
	flag.StringVar(&cfg.GRPC.PrivateKeySSL, "private_key_ssl", "private.key", "private key for ssl path")

	flag.StringVar(&cfg.PostgresStorage.DatabaseDSN, "d", "", "database dsn")
	flag.StringVar(&cfg.Settings.PrivateKey, "private_key", "private.key", "private key path")
	flag.StringVar(&cfg.Settings.JWTkey, "j", "", "jwt secret key")

	flag.Parse()

	env.Parse(&cfg.Settings)
	env.Parse(&cfg.GRPC)
	env.Parse(&cfg.PostgresStorage)

	return cfg
}
