package account

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidateJWT(t *testing.T) {
	a := assert.New(t)
	usr := Account{
		ID: "123",
	}

	secretKey := []byte("secret-key")

	tokenStr, err := CreateJWT(usr, secretKey)
	a.NoError(err)

	_, err = ValidateJWT(tokenStr, []byte("wrong key"))
	a.Error(err)

	t.Log(tokenStr)
	usrID, err := ValidateJWT(tokenStr, secretKey)
	a.NoError(err)
	a.Equal(usr.ID, usrID)
}

func TestWriteUserToContext(t *testing.T) {
	ctx := context.Background()
	user := Account{
		ID: "user id",
	}

	ctx = WriteUserToContext(ctx, user)

	userFromContext := ReadUserFromContext(ctx)

	assert.Equal(t, user, userFromContext)
}
