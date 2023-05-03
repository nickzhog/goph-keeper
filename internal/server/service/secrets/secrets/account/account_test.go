package secretaccount

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewFromAbstractSecret(t *testing.T) {
	a := assert.New(t)

	expected := SecretAccount{
		ID:         "testid",
		UserID:     "userID",
		Title:      "title",
		SiteDomain: "vk.com",
		Login:      "login",
		Password:   "1234",
	}

	abstractSecret, err := expected.ExportToAbstractSecret()
	a.NoError(err)

	got, err := NewFromAbstractSecret(abstractSecret)
	a.NoError(err)

	a.Equal(expected, got)
}
