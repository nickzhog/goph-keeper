package secretcard

import (
	"encoding/json"
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

func Unmarshal(secretData []byte) (*SecretCard, error) {
	var card SecretCard

	err := json.Unmarshal(secretData, &card)
	if err != nil {
		return nil, err
	}

	if !card.IsValid() {
		return nil, errors.New("not valid secret")
	}

	return &card, nil
}
