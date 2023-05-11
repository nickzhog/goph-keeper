package secretaccount

import (
	"encoding/json"
	"errors"
)

type SecretAccount struct {
	SiteDomain string `json:"site_domain,omitempty"`
	Login      string `json:"login,omitempty"`
	Password   string `json:"password,omitempty"`
	KeyTOTP    string `json:"key_totp,omitempty"`
	Note       string `json:"note,omitempty"`
}

func (a *SecretAccount) IsValid() bool {
	if a.SiteDomain == "" && a.Login == "" &&
		a.Password == "" && a.KeyTOTP == "" && a.Note == "" {
		return false
	}

	return true
}

func (a *SecretAccount) Marshal() []byte {
	data, _ := json.Marshal(a)
	return data
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

func New(site, login, password, keytotp, note string) *SecretAccount {
	return &SecretAccount{
		SiteDomain: site,
		Login:      login,
		Password:   password,
		KeyTOTP:    keytotp,
		Note:       note,
	}
}
