package secrets

import (
	"testing"

	"github.com/nickzhog/goph-keeper/pkg/encryption"
	"github.com/stretchr/testify/assert"
)

func TestAbstractSecret_Decrypt(t *testing.T) {
	a := assert.New(t)

	secretData := "secret data"
	s := NewSecret("id", "userID", "title", TypeAccount, []byte(secretData))

	priv, pub := encryption.NewRandomKeys()

	err := s.Encrypt(pub)
	a.NoError(err)

	a.NotEqual(secretData, string(s.Data))

	err = s.Decrypt(priv)
	a.NoError(err)
	a.Equal(secretData, string(s.Data))
}
