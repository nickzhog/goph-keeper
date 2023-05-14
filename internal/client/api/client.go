package api

import (
	"context"
	"errors"

	"github.com/nickzhog/goph-keeper/pkg/secrets"
)

var ErrInvalidToken = errors.New("token invalid")

type KeeperClient interface {
	CreateAccount(ctx context.Context, login, password string) error
	LoginAccount(ctx context.Context, login, password string) (token string, err error)

	SecretsView(ctx context.Context) ([]secrets.AbstractSecret, error)

	CreateSecret(ctx context.Context, secret secrets.AbstractSecret) error
	GetSecretByID(ctx context.Context, secretID string) (secrets.AbstractSecret, error)
	UpdateSecretByID(ctx context.Context, secret secrets.AbstractSecret) error
	DeleteSecretByID(ctx context.Context, secretID string) error

	ApplyToken(token string)
	ResetToken()
}
