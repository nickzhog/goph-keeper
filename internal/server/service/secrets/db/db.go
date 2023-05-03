package db

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/nickzhog/goph-keeper/internal/server/service/secrets"
	"github.com/nickzhog/goph-keeper/pkg/logging"
)

var _ secrets.Repository = (*repository)(nil)

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

func (r *repository) Create(ctx context.Context, secret secrets.AbstractSecret) error {
	q := `
	INSERT INTO public.secret 
		(stype, title, data_encrypted, account_id)
	VALUES
		($1, $2, $3, $4)
	`
	_, err := r.client.Exec(ctx, q,
		secret.SType, secret.Title, secret.Data, secret.UserID)
	return err
}

func (r *repository) FindByID(ctx context.Context, id string) (secrets.AbstractSecret, error) {
	q := `
	SELECT 
		id, stype, title, data_encrypted, account_id
	FROM 
		public.secret
	WHERE 
		id = $1
	`
	var secret secrets.AbstractSecret
	secret.IsEncrypted = true
	err := r.client.QueryRow(ctx, q, id).
		Scan(&secret.ID, &secret.SType, &secret.Title, &secret.Data, &secret.UserID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return secrets.AbstractSecret{}, secrets.ErrNotFound

		}
		return secrets.AbstractSecret{}, err
	}

	return secret, nil
}

func (r *repository) FindForUser(ctx context.Context, usrID string) ([]secrets.AbstractSecret, error) {
	q := `
	SELECT 
		id, stype, title, account_id
	FROM 
		public.secret
	WHERE 
		id = $1
	`

	rows, err := r.client.Query(ctx, q)
	if err != nil {
		return nil, err
	}

	var answer []secrets.AbstractSecret
	for rows.Next() {
		var s secrets.AbstractSecret
		err := rows.Scan(&s.ID, &s.SType, &s.Title, &s.UserID)
		if err != nil {
			return nil, err
		}

		answer = append(answer, s)
	}

	return answer, nil
}

func (r *repository) Update(ctx context.Context, secret secrets.AbstractSecret) error {
	q := `
	UPDATE 
		public.secret 
	SET
		stype = $1, 
		title = $2,
	 	data_encrypted = $3,
		account_id = $4
	WHERE
		id = $5
	`
	_, err := r.client.Exec(ctx, q,
		secret.SType, secret.Title, secret.Data, secret.UserID, secret.ID)
	return err
}

func (r *repository) DeleteByID(ctx context.Context, id string) error {
	q := `
	DELETE FROM
		public.secret 
	WHERE
		id = $1
	`
	_, err := r.client.Exec(ctx, q, id)
	return err
}
