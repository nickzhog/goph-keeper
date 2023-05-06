package secret

import (
	"context"

	"github.com/nickzhog/goph-keeper/pkg/secrets"
)

type Repository interface {
	Create(ctx context.Context, secret secrets.AbstractSecret) error
	FindForUser(ctx context.Context, usrID string) ([]secrets.AbstractSecret, error)
	FindByID(ctx context.Context, id string) (secrets.AbstractSecret, error)
	Update(ctx context.Context, secret secrets.AbstractSecret) error
	DeleteByID(ctx context.Context, id string) error
}
