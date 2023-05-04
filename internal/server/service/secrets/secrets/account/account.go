package secretaccount

import (
	"bytes"
	"encoding/gob"
	"errors"
)

type SecretAccount struct {
	SiteDomain string
	Login      string
	Password   string
	KeyTOTP    string
	Note       string
}

func (a *SecretAccount) IsValid() bool {
	if a.SiteDomain == "" && a.Login == "" &&
		a.Password == "" && a.KeyTOTP == "" && a.Note == "" {
		return false
	}

	return true
}

func New(secretData []byte) (*SecretAccount, error) {
	account := new(SecretAccount)

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
