package secretbinary

import (
	"bytes"
	"encoding/gob"
	"errors"
)

type SecretBinary struct {
	Data []byte
	Note string
}

func (b *SecretBinary) IsValid() bool {
	if len(b.Data) == 0 && b.Note == "" {
		return false
	}

	return true
}

func New(secretData []byte) (*SecretBinary, error) {
	account := new(SecretBinary)

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
