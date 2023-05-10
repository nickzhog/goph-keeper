package account

import "context"

type Repository interface {
	Create(ctx context.Context, account Account) error
	FindForLogin(ctx context.Context, login string) (Account, error)
	FindForID(ctx context.Context, id string) (Account, error)
}
