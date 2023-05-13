package server

import (
	"context"
	"crypto/rsa"

	"github.com/nickzhog/goph-keeper/internal/server/config"
	"github.com/nickzhog/goph-keeper/internal/server/service/account"
	"github.com/nickzhog/goph-keeper/internal/server/storage"
	"github.com/nickzhog/goph-keeper/pkg/encryption"
	"github.com/nickzhog/goph-keeper/pkg/logging"
	"github.com/nickzhog/goph-keeper/pkg/secrets"
)

type Server struct {
	Logger  *logging.Logger
	storage storage.Storage
	cfg     *config.Config

	priv         *rsa.PrivateKey
	pub          *rsa.PublicKey
	jwtSecretKey []byte
}

func NewServer(logger *logging.Logger, cfg *config.Config, storage storage.Storage) *Server {
	priv, err := encryption.GetPrivateKeyFromFile(cfg.Settings.PrivateKey)
	if err != nil {
		logger.Fatal(err)
	}

	pub := &priv.PublicKey

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

func (s *Server) FindSecretByID(ctx context.Context, secretID, userID string) (secrets.AbstractSecret, error) {
	secret, err := s.storage.Secrets.FindByID(ctx, secretID)
	if err != nil {
		return secrets.AbstractSecret{}, err
	}

	if secret.UserID != userID {
		return secrets.AbstractSecret{}, secrets.ErrNotFound

	}

	err = secret.Decrypt(s.priv)
	if err != nil {
		return secrets.AbstractSecret{}, err
	}

	return secret, nil
}

func (s *Server) CreateSecret(ctx context.Context, userID string, secret secrets.AbstractSecret) error {
	if err := secret.Validate(); err != nil {
		return secrets.ErrInvalid
	}

	err := secret.Encrypt(s.pub)
	if err != nil {
		return err
	}

	return s.storage.Secrets.Create(ctx, userID, secret)
}

func (s *Server) DeleteSecret(ctx context.Context, secretID string, userID string) error {
	secret, err := s.storage.Secrets.FindByID(ctx, secretID)
	if err != nil {
		return err
	}

	if secret.UserID != userID {
		return secrets.ErrNotFound
	}

	return s.storage.Secrets.DeleteByID(ctx, secretID)
}

func (s *Server) UpdateSecret(ctx context.Context, secret secrets.AbstractSecret, userID string) error {
	if err := secret.Validate(); err != nil {
		return secrets.ErrInvalid
	}

	secretOld, err := s.storage.Secrets.FindByID(ctx, secret.ID)
	if err != nil {
		return err
	}
	if secretOld.UserID != userID {
		return secrets.ErrNotFound
	}

	if err := secret.Encrypt(s.pub); err != nil {
		return secrets.ErrInvalid
	}

	return s.storage.Secrets.Update(ctx, secret)
}
