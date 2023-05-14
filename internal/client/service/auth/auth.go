package auth

import (
	"context"
	"strings"

	"github.com/nickzhog/goph-keeper/internal/client/api"
)

func Login(ctx context.Context, keeper api.KeeperClient, login, password string) error {
	login = strings.TrimSpace(login)
	password = strings.TrimSpace(password)

	toketString, err := keeper.LoginAccount(ctx, login, password)
	if err != nil {
		return err
	}

	keeper.ApplyToken(toketString)

	return nil
}

func Register(ctx context.Context, keeper api.KeeperClient, login, password string) error {
	login = strings.TrimSpace(login)
	password = strings.TrimSpace(password)

	err := keeper.CreateAccount(ctx, login, password)

	return err
}
