package secretbinary

import (
	"bytes"
	"encoding/gob"

	"github.com/nickzhog/goph-keeper/internal/server/service/secrets"
)

type SecretBinary struct {
	ID     string
	UserID string
	Title  string

	Data []byte
	Note string
}

func NewFromAbstractSecret(secret secrets.AbstractSecret) (SecretBinary, error) {
	var data SecretBinary
	buf := bytes.NewBuffer(secret.Data)

	err := gob.NewDecoder(buf).Decode(&data)
	if err != nil {
		return SecretBinary{}, err
	}

	return data, nil
}

func (a *SecretBinary) ExportToAbstractSecret() (secrets.AbstractSecret, error) {
	var buf bytes.Buffer
	err := gob.NewEncoder(&buf).Encode(*a)
	if err != nil {
		return secrets.AbstractSecret{}, err
	}

	abstractSecret := secrets.NewSecret(a.ID, a.UserID, a.Title, secrets.TypeAccount, buf.Bytes())

	return *abstractSecret, nil
}
