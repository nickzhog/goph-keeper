package account

import (
	"context"
	"errors"
	"time"

	"github.com/dgrijalva/jwt-go/v4"
	"golang.org/x/crypto/bcrypt"
)

var ErrNotFound = errors.New("user not found")
var ErrWrongPassword = errors.New("wrong password")

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

const ClaimsKeyUserID = "user_id"

func CreateJWT(usr Account, secretKey []byte) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)

	claims["exp"] = time.Now().Add(time.Hour).Unix()
	claims[ClaimsKeyUserID] = usr.ID

	tokenStr, err := token.SignedString(secretKey)
	if err != nil {
		return "", err
	}

	return tokenStr, nil
}

type ContextUserKey string

const contextUserKey ContextUserKey = "user"

func WriteUserToContext(usr Account, ctx context.Context) context.Context {
	return context.WithValue(ctx, contextUserKey, usr)
}

func ReadUserFromContext(ctx context.Context) Account {
	return ctx.Value(contextUserKey).(Account)
}
