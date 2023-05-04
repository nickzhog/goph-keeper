package account

import (
	"context"

	"golang.org/x/crypto/bcrypt"
)

type Account struct {
	ID           string
	Login        string
	PasswordHash string
}

func Create(ctx context.Context, rep Repository, login, password string) error {
	phash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	account := Account{Login: login, PasswordHash: string(phash)}

	return rep.Create(ctx, account)
}

func Login(ctx context.Context, rep Repository, login, password string) (Account, error) {
	usr, err := rep.FindForLogin(ctx, login)
	if err != nil {
		return Account{}, ErrWrongPassword
	}

	err = bcrypt.CompareHashAndPassword([]byte(usr.PasswordHash), []byte(password))
	if err != nil {
		return Account{}, ErrWrongPassword
	}

	return usr, nil
}
