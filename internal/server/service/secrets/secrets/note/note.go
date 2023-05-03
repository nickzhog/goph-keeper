package secretnote

import (
	"bytes"
	"encoding/gob"

	"github.com/nickzhog/goph-keeper/internal/server/service/secrets"
)

type SecretNote struct {
	ID     string
	UserID string
	Title  string

	Note string
}

func NewFromAbstractSecret(secret secrets.AbstractSecret) (SecretNote, error) {
	var note SecretNote
	buf := bytes.NewBuffer(secret.Data)

	err := gob.NewDecoder(buf).Decode(&note)
	if err != nil {
		return SecretNote{}, err
	}

	return note, nil
}

func (a *SecretNote) ExportToAbstractSecret() (secrets.AbstractSecret, error) {
	var buf bytes.Buffer
	err := gob.NewEncoder(&buf).Encode(*a)
	if err != nil {
		return secrets.AbstractSecret{}, err
	}

	abstractSecret := secrets.NewSecret(a.ID, a.UserID, a.Title, secrets.TypeAccount, buf.Bytes())

	return *abstractSecret, nil
}
