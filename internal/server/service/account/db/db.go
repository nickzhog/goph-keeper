package db

import (
	"context"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/nickzhog/goph-keeper/internal/server/service/account"
	"github.com/nickzhog/goph-keeper/pkg/logging"
)

var _ account.Repository = (*repository)(nil)

type repository struct {
	logger *logging.Logger
	client *pgxpool.Pool
}

func NewRepository(logger *logging.Logger, client *pgxpool.Pool) *repository {
	return &repository{
		logger: logger,
		client: client,
	}
}

func (r *repository) Create(ctx context.Context, usr account.Account) error {
	q := `
	INSERT INTO public.account 
		(login, password_hash) 
	VALUES 
		($1, $2)
	`
	_, err := r.client.Exec(ctx, q, usr.Login, usr.PasswordHash)
	return err
}

func (r *repository) FindForLogin(ctx context.Context, login string) (account.Account, error) {
	q := `
		SELECT 
			id, login, password_hash
		FROM
			public.account
		WHERE 
			login = $1
	`
	var usr account.Account
	err := r.client.QueryRow(ctx, q, login).Scan(&usr.ID, &usr.Login, &usr.PasswordHash)
	if err != nil {
		return account.Account{}, err
	}
	return usr, nil
}

func (r *repository) FindForID(ctx context.Context, id string) (account.Account, error) {
	q := `
		SELECT 
			id, login, password_hash
		FROM
			public.account
		WHERE 
			id = $1
	`
	var usr account.Account
	err := r.client.QueryRow(ctx, q, id).Scan(&usr.ID, &usr.Login, &usr.PasswordHash)
	if err != nil {
		return account.Account{}, err
	}
	return usr, nil
}
