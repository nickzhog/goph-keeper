package account

import (
	"context"
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go/v4"
)

const (
	ClaimsKeyUserID = "user_id"
	expireDuration  = time.Hour
)

type Claims struct {
	jwt.StandardClaims
	UserID string `json:"user_id"`
}

func CreateJWT(usr Account, secretKey []byte) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &Claims{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: jwt.At(time.Now().Add(expireDuration)),
			IssuedAt:  jwt.At(time.Now()),
		},
		UserID: usr.ID,
	})

	return token.SignedString(secretKey)
}

func ValidateJWT(tokenStr string, secretKey []byte) (userID string, err error) {
	token, err := jwt.ParseWithClaims(tokenStr, &Claims{}, func(t *jwt.Token) (interface{}, error) {
		_, ok := t.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			return "", fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}

		return secretKey, nil
	})

	if err != nil {
		return "", err
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims.UserID, nil
	}

	return "", ErrInvalidAccessToken
}

// context storage

type ContextUserKey string

const contextUserKey ContextUserKey = "user"

func WriteUserToContext(usr Account, ctx context.Context) context.Context {
	return context.WithValue(ctx, contextUserKey, usr)
}

func ReadUserFromContext(ctx context.Context) Account {
	return ctx.Value(contextUserKey).(Account)
}
