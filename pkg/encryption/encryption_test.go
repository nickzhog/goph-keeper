package encryption

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDecryptData(t *testing.T) {
	a := assert.New(t)
	priv, pub := NewRandomKeys()

	expected := "test"

	encrypted, err := EncryptData([]byte(expected), pub)
	a.NoError(err)

	a.NotEqual(expected, encrypted)

	decrypted, err := DecryptData(encrypted, priv)
	a.NoError(err)

	a.Equal(expected, string(decrypted))
}
