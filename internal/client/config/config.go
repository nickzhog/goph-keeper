package config

import (
	"flag"

	"github.com/caarlos0/env"
)

type Config struct {
	Settings struct {
		ConnAddress string `env:"CONNECTION_ADDRESS"`

		CertSSL string `env:"CERT_SSL_PATH"`
	}
}

func GetConfig() *Config {
	cfg := new(Config)
	flag.StringVar(&cfg.Settings.ConnAddress, "a", "localhost:3200", "connection address")
	flag.StringVar(&cfg.Settings.CertSSL, "c", "cert.crt", "path to cert for ssl")

	flag.Parse()

	env.Parse(&cfg.Settings)

	return cfg
}
