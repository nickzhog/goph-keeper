package secrets

import "context"

type Repository interface {
	Create(ctx context.Context, secret AbstractSecret) error
	FindForUser(ctx context.Context, usrID string) ([]AbstractSecret, error)
	FindByID(ctx context.Context, id string) (AbstractSecret, error)
	Update(ctx context.Context, secret AbstractSecret) error
	DeleteByID(ctx context.Context, id string) error
}
