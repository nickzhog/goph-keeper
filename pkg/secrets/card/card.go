package secretcard

import (
	"encoding/json"
	"errors"
)

type SecretCard struct {
	Number     string `json:"number,omitempty"`
	Month      string `json:"month,omitempty"`
	Year       string `json:"year,omitempty"`
	CVV        string `json:"cvv,omitempty"`
	HolderName string `json:"holder_name,omitempty"`
	Note       string `json:"note,omitempty"`
}

func (card *SecretCard) IsValid() bool {
	if card.Number == "" && card.Month == "" &&
		card.Year == "" && card.CVV == "" &&
		card.HolderName == "" && card.Note == "" {
		return false
	}

	return true
}

func (card *SecretCard) Marshal() []byte {
	data, _ := json.Marshal(card)
	return data
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

func New(number, month, year, cvv, holder, note string) *SecretCard {
	return &SecretCard{
		Number:     number,
		Month:      month,
		Year:       year,
		CVV:        cvv,
		HolderName: holder,
		Note:       note,
	}
}
