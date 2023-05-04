package secretnote

import (
	"bytes"
	"encoding/gob"
	"errors"
)

type SecretNote struct {
	Note string
}

func (n *SecretNote) IsValid() bool {
	if n.Note == "" {
		return false
	}

	return true
}

func New(secretData []byte) (*SecretNote, error) {
	account := new(SecretNote)

	buf := bytes.NewBuffer(secretData)
	err := gob.NewDecoder(buf).Decode(&account)
	if err != nil {
		return nil, err
	}

	if !account.IsValid() {
		return nil, errors.New("not valid secret")
	}

	return account, nil
}
