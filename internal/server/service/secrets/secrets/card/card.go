package secretcard

import (
	"bytes"
	"encoding/gob"
	"errors"
)

type SecretCard struct {
	Number     string
	Month      string
	Year       string
	CVV        string
	HolderName string
	Note       string
}

func (card *SecretCard) IsValid() bool {
	if card.Number == "" && card.Month == "" &&
		card.Year == "" && card.CVV == "" &&
		card.HolderName == "" && card.Note == "" {
		return false
	}

	return true
}

func New(secretData []byte) (*SecretCard, error) {
	account := new(SecretCard)

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
