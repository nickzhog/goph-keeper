package server

import (
	"context"
	"crypto/rsa"

	"github.com/nickzhog/goph-keeper/internal/server/config"
	"github.com/nickzhog/goph-keeper/internal/server/service"
	"github.com/nickzhog/goph-keeper/internal/server/service/account"
	"github.com/nickzhog/goph-keeper/internal/server/service/secrets"
	"github.com/nickzhog/goph-keeper/pkg/encryption"
	"github.com/nickzhog/goph-keeper/pkg/logging"
)

type Server struct {
	Logger  *logging.Logger
	storage service.Storage
	cfg     *config.Config

	priv         *rsa.PrivateKey
	pub          *rsa.PublicKey
	jwtSecretKey []byte
}

func NewServer(logger *logging.Logger, cfg *config.Config, storage service.Storage) *Server {
	priv, err := encryption.NewPrivateKey(cfg.Settings.PrivateKey)
	if err != nil {
		logger.Fatal(err)
	}

	pub, err := encryption.NewPublicKey(cfg.Settings.PublicKey)
	if err != nil {
		logger.Fatal(err)
	}

	jwtKey := cfg.Settings.JWTkey

	return &Server{
		Logger:       logger,
		storage:      storage,
		cfg:          cfg,
		priv:         priv,
		pub:          pub,
		jwtSecretKey: []byte(jwtKey),
	}
}

func (s *Server) Register(ctx context.Context, login, password string) error {
	return account.Create(ctx, s.storage.Account, login, password)
}

func (s *Server) Login(ctx context.Context, login, password string) (string, error) {
	usr, err := account.Login(ctx, s.storage.Account, login, password)
	if err != nil {
		return "", err
	}

	return account.CreateJWT(usr, s.jwtSecretKey)
}

func (s *Server) FindSecretsForUser(ctx context.Context, usrID string) ([]secrets.AbstractSecret, error) {
	return s.storage.Secrets.FindForUser(ctx, usrID)
}

func (s *Server) FindSecretByID(ctx context.Context, id string) (secrets.AbstractSecret, error) {
	secret, err := s.storage.Secrets.FindByID(ctx, id)
	if err != nil {
		return secrets.AbstractSecret{}, err
	}

	err = secret.Decrypt(s.priv)
	if err != nil {
		return secrets.AbstractSecret{}, err
	}

	return secret, nil
}
