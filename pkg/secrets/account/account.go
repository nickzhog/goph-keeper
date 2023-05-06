package secretaccount

import (
	"encoding/json"
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

func Unmarshal(secretData []byte) (*SecretAccount, error) {
	var account SecretAccount

	err := json.Unmarshal(secretData, &account)
	if err != nil {
		return nil, err
	}

	if !account.IsValid() {
		return nil, errors.New("not valid secret")
	}

	return &account, nil
}
