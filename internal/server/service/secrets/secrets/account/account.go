package secretaccount

import (
	"bytes"
	"encoding/gob"

	"github.com/nickzhog/goph-keeper/internal/server/service/secrets"
)

type SecretAccount struct {
	ID     string
	UserID string
	Title  string

	SiteDomain string
	Login      string
	Password   string
	KeyTOTP    string
	Note       string
}

func NewFromAbstractSecret(secret secrets.AbstractSecret) (SecretAccount, error) {
	var account SecretAccount
	buf := bytes.NewBuffer(secret.Data)

	err := gob.NewDecoder(buf).Decode(&account)
	if err != nil {
		return SecretAccount{}, err
	}

	return account, nil
}

func (a *SecretAccount) ExportToAbstractSecret() (secrets.AbstractSecret, error) {
	var buf bytes.Buffer
	err := gob.NewEncoder(&buf).Encode(*a)
	if err != nil {
		return secrets.AbstractSecret{}, err
	}

	abstractSecret := secrets.NewSecret(a.ID, a.UserID, a.Title, secrets.TypeAccount, buf.Bytes())

	return *abstractSecret, nil
}
