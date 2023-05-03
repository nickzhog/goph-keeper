package service

import (
	"context"
	"log"
	"time"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/nickzhog/goph-keeper/internal/server/service/account"
	accountdb "github.com/nickzhog/goph-keeper/internal/server/service/account/db"
	"github.com/nickzhog/goph-keeper/internal/server/service/secrets"
	secretsdb "github.com/nickzhog/goph-keeper/internal/server/service/secrets/db"
	"github.com/nickzhog/goph-keeper/pkg/logging"
)

type Storage struct {
	Account account.Repository
	Secrets secrets.Repository
}

func NewPostgresStorage(ctx context.Context, logger *logging.Logger, dsn string) Storage {
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	pool, err := pgxpool.Connect(ctx, dsn)
	if err != nil {
		log.Fatal(err)
	}

	return Storage{
		Account: accountdb.NewRepository(logger, pool),
		Secrets: secretsdb.NewRepository(logger, pool),
	}
}
